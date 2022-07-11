package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/hellofresh/health-go/v4"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	kilnlog "github.com/kilnfi/go-utils/log"
	kilnhttp "github.com/kilnfi/go-utils/net/http"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var (
	statusInitializing = "initializing"
	statusInitErr      = "initErr"
	statusStarting     = "starting"
	statusRunning      = "running"
	statusStopping     = "stopping"
	statusStopped      = "stopped"
)

// An App is a modular application with HTTP server. It allows users to register custom Services
// and automatically deals with the lifecycle of those services on call to Run()
type App struct {
	cfg *Config

	middlewares alice.Chain
	mux         *httprouter.Router
	server      *kilnhttp.Server

	healthMux    *httprouter.Router
	healthServer *kilnhttp.Server

	liveness  *health.Health
	readiness *health.Health

	prometheus *prometheus.Registry

	logger *logrus.Logger

	services []interface{}
	toStop   []Runnable

	statusMux sync.Mutex
	status    string

	done chan os.Signal
}

// New creates a new App
func New(cfg *Config) (*App, error) {
	logger, err := kilnlog.New(cfg.Logger)
	if err != nil {
		return nil, err
	}

	server, err := kilnhttp.NewServer(cfg.Server)
	if err != nil {
		return nil, err
	}
	server.SetLogger(logger)

	healthServer, err := kilnhttp.NewServer(cfg.Healthz)
	if err != nil {
		return nil, err
	}
	healthServer.SetLogger(logger)

	liveness, _ := health.New()
	readiness, _ := health.New()

	return &App{
		cfg:          cfg,
		mux:          httprouter.New(),
		middlewares:  alice.New(),
		server:       server,
		healthMux:    httprouter.New(),
		healthServer: healthServer,
		liveness:     liveness,
		readiness:    readiness,
		prometheus:   prometheus.NewRegistry(),
		logger:       logger,
		done:         make(chan os.Signal, 1),
	}, nil
}

func (app *App) setStatus(status string) {
	app.statusMux.Lock()
	app.status = status
	app.statusMux.Unlock()
}

func (app *App) isStatus(status string) bool {
	app.statusMux.Lock()
	defer app.statusMux.Unlock()

	return app.status == status
}

func (app *App) registerBaseChecks() error {
	if err := app.liveness.Register(health.Config{
		Name:    "app",
		Timeout: time.Second,
		Check:   app.livecheck,
	}); err != nil {
		return err
	}

	if err := app.readiness.Register(health.Config{
		Name:    "app",
		Timeout: time.Second,
		Check:   app.readycheck,
	}); err != nil {
		return err
	}

	return nil
}

func (app *App) registerBaseMetrics() {
	app.prometheus.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	app.prometheus.MustRegister(collectors.NewGoCollector())
}

func (app *App) livecheck(ctx context.Context) error {
	if app.isStatus(statusInitErr) {
		return fmt.Errorf("app initilalization failed")
	}

	if app.isStatus(statusStopped) {
		return fmt.Errorf("app stopped")
	}

	return nil
}

func (app *App) readycheck(ctx context.Context) error {
	if app.isStatus("") {
		return fmt.Errorf("app has not yet been started")
	}

	if app.isStatus(statusInitializing) {
		return fmt.Errorf("app is initializing")
	}

	if app.isStatus(statusStarting) {
		return fmt.Errorf("app is starting")
	}

	if app.isStatus(statusStopping) {
		return fmt.Errorf("app is stopping")
	}

	return nil
}

// RegisterService register a Service on App

// RegisterService MUST be called before App.Run() is called to successfully register the service.

// Service is an object which lifecycle can managed by an App.

// A service can optionally
// - maintain long live tasks
// - register HTTP routes on App multiplexer to get incoming HTTP requests routed to to the service

// Service lifecycle

// 0. Service Construction: `svc := New(...)`
//    Constructor
//       - MUST not perform any side effects (like API calls, disk io, logging, accessing db,  etc.)
//       - MUST not start any long live tasks (like starting a server, etc)
//       - MAY perform some config validation checks
//       - SHOULD create the Service object as fast as possible or return an error

// 1. Service Registration: `app.RegisterService(svc)`
//   App calls
//      - `svc.SetLogger(...)` to feed a logger to the Service
//      - `svc.RegisterHandler(...)` allowing the Service to register HTTP routes on the application. If multiple
//         services register the same route this will fail
//      - `svc.RegisterMiddleware(...)` allowing the Service to register a middleware on the application

// 2. Service Initialization: `svc.Init(context.Context) error`
//   Service
//      - MAY perform side effects (like API calls, disk io, logging, accessing db,  etc.)
//      - MUST return an error in case initialization failed and nil if succeeded
//      - MUST stop initialization and return ASAP if the context is canceled
//      - MUST not start any long live tasks (like starting a server, etc)

// 3. Service Start: `svc.Init(context.Context) error`
//  Service
//      - MAY start long live tasks
//      - MUST return an error in case starting failed and nil if succeeded
//      - MUST stop starting and return ASAP if the context is canceled

// 4. Service Stop: `svc.Stop(context.Context) error`
//  Service
//      - MUST gracefully close all long running tasks
//      - MUST return an error in case stopping failed and nil if succeeded
//      - MUST stop stopping and return ASAP if the context is canceled

// 4. Service Close: `svc.Close() error`
// Service
//      - MAY clean its state
func (app *App) RegisterService(svc interface{}) {
	// Register Handler and SetLogger on service registration
	if loggable, ok := svc.(Loggable); ok {
		loggable.SetLogger(app.logger)
	}

	if api, ok := svc.(API); ok {
		api.RegisterHandler(app.mux)
	}

	if mid, ok := svc.(Middleware); ok {
		app.middlewares = mid.RegisterMiddleware(app.middlewares)
	}

	app.services = append(app.services, svc)
}

// Run application and all registered Service

// Run first starts HTTP server and accepts connection before starting services

// Run can be interupted by sending a SIGTERM, SIGINT signal. In which, case Run will
// attempt to gracefully stop all services and HTTP server
func (app *App) Run() error {
	return app.run()
}

func (app *App) Logger() *logrus.Logger {
	return app.logger
}

func (app *App) initServices(ctx context.Context) error {
	app.setStatus(statusInitializing)
	app.logger.Infof("initialize services...")

	initCtx, cancelInit := context.WithCancel(ctx)

	wg := &sync.WaitGroup{}
	errors := make(chan error)

	wg.Add(len(app.services))
	for _, svc := range app.services {
		go func(svc interface{}) {
			defer wg.Done()

			// Register checks before starting initialization
			if check, ok := svc.(Checkable); ok {
				if err := check.RegisterCheck(app.readiness); err != nil {
					errors <- err
					return
				}
			}

			// Register metrics before starting initialization
			if measure, ok := svc.(Measurable); ok {
				if err := measure.RegisterMetrics(app.prometheus); err != nil {
					errors <- err
					return
				}
			}

			if init, ok := svc.(Initializable); ok {
				err := init.Init(initCtx)
				if err != nil {
					errors <- err
					return
				}
			}
		}(svc)
	}

	var initErr error

	go func() {
		initErr = <-errors
		if initErr != nil {
			cancelInit()
			for range errors {
				// drain errors
			}
		}
	}()

	wg.Wait()
	close(errors)

	if initErr != nil {
		app.setStatus(statusInitErr)
		app.logger.WithError(initErr).Errorf("error initializing services")
	}

	return initErr
}

func (app *App) startServices(ctx context.Context) (err error) {
	app.logger.Infof("start services...")
	app.setStatus(statusStarting)

	for _, svc := range app.services {
		if run, ok := svc.(Runnable); ok {
			app.toStop = append(app.toStop, run)
			err = run.Start(ctx)
			if err != nil {
				app.logger.WithError(err).Errorf("error starting services")

				// before returning we stop already started services
				_ = app.stopServices(ctx)

				return err
			}
		}
	}

	app.logger.Infof("services successfully started")
	app.setStatus(statusRunning)

	return
}

func (app *App) stopServices(ctx context.Context) error {
	app.logger.Infof("stop services...")
	app.setStatus(statusStopping)
	var rErr error
	for i := range app.toStop {
		if err := app.toStop[len(app.toStop)-i-1].Stop(ctx); err != nil && rErr == nil {
			rErr = err
		}
	}

	if rErr != nil {
		app.logger.WithError(rErr).Errorf("error stopping service")
	}

	app.logger.Infof("services stopped...")
	app.setStatus(statusStopped)

	return rErr
}

func (app *App) run() error {
	if raw, err := json.Marshal(app.cfg); err != nil {
		app.logger.WithError(err).Errorf("invalid config")
		return err
	} else {
		app.logger.WithField("config", string(raw)).Infof("run app...")
	}

	startCtx, cancel := context.WithTimeout(context.Background(), app.cfg.StartTimeout.Duration)
	defer cancel()

	// first thing we start listen to signals and open server connection so we listen to incoming request
	if err := app.startSignalsAndServers(startCtx); err != nil {
		return err
	}

	// initialize services
	if err := app.initServices(startCtx); err != nil {
		return err
	}

	// start services
	if err := app.startServices(startCtx); err != nil {
		_ = app.stopSignalsAndServers(startCtx)
		return err
	}

	select {
	case sig := <-app.done:
		app.logger.WithField("signal", sig.String()).Errorf("received termination signal")
	case <-app.server.Done():
		app.logger.WithError(app.server.Error()).Errorf("server error")
	}

	// we received a stop signal so we stop
	stopCtx, cancel := context.WithTimeout(context.Background(), app.cfg.StopTimeout.Duration)
	defer cancel()

	if err := app.stopServices(stopCtx); err == nil {
		return app.stopSignalsAndServers(stopCtx)
	} else {
		_ = app.stopSignalsAndServers(stopCtx)
		return err
	}
}

func (app *App) listenSignals() {
	signal.Notify(app.done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
}

func (app *App) stopListeningSignals() {
	signal.Stop(app.done)
}

func (app *App) setHandlers() error {
	if err := app.setHealthHandler(); err != nil {
		return err
	}
	app.setHandler()
	return nil
}

func (app *App) setHealthHandler() error {
	if err := app.registerChecksHandler(); err != nil {
		return err
	}

	app.registerMetricsHandler()

	app.healthServer.SetHandler(app.healthMux)

	return nil
}

func (app *App) registerChecksHandler() error {
	// register base checks
	if err := app.registerBaseChecks(); err != nil {
		return err
	}

	app.healthMux.HandlerFunc(http.MethodGet, "/live", app.liveness.HandlerFunc)
	app.healthMux.HandlerFunc(http.MethodGet, "/ready", app.readiness.HandlerFunc)

	return nil
}

func (app *App) registerMetricsHandler() {
	app.registerBaseMetrics()

	app.healthMux.Handler(
		http.MethodGet,
		"/metrics",
		promhttp.HandlerFor(app.prometheus, promhttp.HandlerOpts{
			ErrorLog: app.logger,
		}),
	)
}

func (app *App) setHandler() {
	h := app.instrumentMiddleware().Extend(app.middlewares).Then(app.mux)
	app.server.SetHandler(h)
}

func (app *App) instrumentMiddleware() alice.Chain {
	return alice.New(
		app.loggerMiddleware,
		app.requestMetricsMiddleware,
	)
}

func (app *App) loggerMiddleware(h http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(app.logger.Out, h)
}

func (app *App) requestMetricsMiddleware(h http.Handler) http.Handler {
	return promhttp.InstrumentMetricHandler(app.prometheus, h)
}

func (app *App) startSignalsAndServers(ctx context.Context) error {
	app.listenSignals()

	if err := app.setHandlers(); err != nil {
		return err
	}

	if err := app.healthServer.Start(ctx); err != nil {
		_ = app.stopSignalsAndServers(ctx)
		return err
	}

	if err := app.server.Start(ctx); err != nil {
		_ = app.stopSignalsAndServers(ctx)
		return err
	}

	return nil
}

func (app *App) stopSignalsAndServers(ctx context.Context) error {
	app.stopListeningSignals()
	sErr := app.server.Stop(ctx)
	hErr := app.healthServer.Stop(ctx)
	if hErr != nil {
		return hErr
	}

	return sErr
}
