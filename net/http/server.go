package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	types "github.com/kilnfi/go-utils/common/types"
	kilnnet "github.com/kilnfi/go-utils/net"
	"github.com/sirupsen/logrus"
)

var (
	statusRunning = "running"
	statusStopped = "stopped"
)

type ConnStateCallbackTypeWrapper func(net.Conn, http.ConnState)

func (c ConnStateCallbackTypeWrapper) MarshalJSON() ([]byte, error) {
	// When the callback is overloaded.
	if c == nil {
		return []byte("\"onConnStateChangeOverloaded\""), nil
	}

	// Otherwise..
	return []byte("\"onConnStateChangeDefault\""), nil
}

type ServerConfig struct {
	Entrypoint *kilnnet.EntrypointConfig

	ReadTimeout       *types.Duration
	ReadHeaderTimeout *types.Duration

	WriteTimeout *types.Duration

	IdleTimeout *types.Duration

	MaxHeaderBytes *int

	ConnStateCallback ConnStateCallbackTypeWrapper
}

func (cfg *ServerConfig) SetDefault() *ServerConfig {
	if cfg.Entrypoint == nil {
		cfg.Entrypoint = &kilnnet.EntrypointConfig{}
	}
	cfg.Entrypoint.SetDefault()

	if cfg.ReadTimeout == nil {
		cfg.ReadTimeout = &types.Duration{Duration: 30 * time.Second}
	}

	if cfg.ReadHeaderTimeout == nil {
		cfg.ReadHeaderTimeout = &types.Duration{Duration: 30 * time.Second}
	}

	if cfg.WriteTimeout == nil {
		cfg.WriteTimeout = &types.Duration{Duration: 90 * time.Second}
	}

	if cfg.IdleTimeout == nil {
		cfg.IdleTimeout = &types.Duration{Duration: 90 * time.Second}
	}

	return cfg
}

type Server struct {
	cfg *ServerConfig

	logger logrus.FieldLogger

	entrypoint *kilnnet.Entrypoint

	server *http.Server

	mux    sync.Mutex
	status string

	startOnce sync.Once
	startErr  error

	done   chan struct{}
	srvErr error

	stopOnce sync.Once
	stopErr  error
}

func NewServer(cfg *ServerConfig) (*Server, error) {
	entrypoint, err := kilnnet.NewEntrypoint(cfg.Entrypoint)
	if err != nil {
		return nil, err
	}

	server := &Server{
		cfg: cfg,
		server: &http.Server{
			ReadTimeout:       cfg.ReadTimeout.Duration,
			ReadHeaderTimeout: cfg.ReadHeaderTimeout.Duration,
			WriteTimeout:      cfg.WriteTimeout.Duration,
			IdleTimeout:       cfg.IdleTimeout.Duration,
			ConnState:         cfg.ConnStateCallback,
		},
		entrypoint: entrypoint,
	}

	return server, nil
}

func (s *Server) SetLogger(logger logrus.FieldLogger) {
	s.logger = logger.WithField("component", "server")
	s.entrypoint.SetLogger(logger)
}

func (s *Server) Logger() logrus.FieldLogger {
	if s.logger == nil {
		s.SetLogger(logrus.StandardLogger())
	}
	return s.logger
}

func (s *Server) SetHandler(h http.Handler) *Server {
	s.server.Handler = h
	return s
}

func (s *Server) Start(ctx context.Context) error {
	s.startOnce.Do(func() {
		s.start(ctx)
	})

	return s.startErr
}

func (s *Server) start(ctx context.Context) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.status == statusStopped {
		s.startErr = fmt.Errorf("server already stopped")
		return
	}

	// Open connection and return possibly error
	l, err := s.entrypoint.Listen(ctx)
	if err != nil {
		s.startErr = err
		return
	}

	s.status = statusRunning

	s.Logger().Infof("start serving HTTP request")
	s.done = make(chan struct{})

	// Start serving in a separate go-routine
	go func() {
		s.srvErr = s.server.Serve(l)
		close(s.done)
		if s.srvErr != nil && s.srvErr != http.ErrServerClosed {
			s.Logger().WithError(s.srvErr).Infof("error while serving HTTP request")
		} else {
			s.Logger().Infof("stopped serving HTTP request")
		}
	}()
}

func (s *Server) Stop(ctx context.Context) error {
	s.stopOnce.Do(func() {
		s.stop(ctx)
	})

	return s.stopErr
}

func (s *Server) stop(ctx context.Context) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.status == statusRunning {
		s.Logger().Infof("stop server...")

		// Gracefully shutdown server
		err := s.server.Shutdown(ctx)
		if err != nil {
			_ = s.server.Close()
		}

		// Wait for Serve(...) to be done
		<-s.done

		// Return possible error from Serve(...)
		if err == nil && s.srvErr != nil && s.srvErr != http.ErrServerClosed {
			err = s.srvErr
		}

		s.stopErr = err
		if err != nil {
			s.Logger().WithError(err).Infof("error stoping server")
		}
		s.Logger().Infof("server successfully stopped")
	} else {
		s.Logger().Infof("server not started nothing to stop")
	}

	s.status = statusStopped
}

func (s *Server) Done() chan struct{} {
	return s.done
}

func (s *Server) Error() error {
	return s.srvErr
}
