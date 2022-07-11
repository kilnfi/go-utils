package app

import (
	"context"

	"github.com/hellofresh/health-go/v4"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// Loggaable is a service that holds a logger
type Loggable interface {
	SetLogger(logger logrus.FieldLogger)
}

// API is a service that exposes API routes
type API interface {
	RegisterHandler(mux *httprouter.Router)
}

// Middleware is a service that exposes a middleware to be set on an App
type Middleware interface {
	RegisterMiddleware(chain alice.Chain) alice.Chain
}

// Initializable is a service that can initialize
type Initializable interface {
	Init(context.Context) error
}

// Runnable is a service that maintains long living task(s)
type Runnable interface {
	Start(context.Context) error
	Stop(context.Context) error
}

// Closable is a service that needs to clean its state at the end of its execution
type Closable interface {
	Close() error
}

// Checkable is a service that can expose its health status
type Checkable interface {
	RegisterCheck(h *health.Health) error
}

// Checkable is a service that can expose metrics
type Measurable interface {
	RegisterMetrics(prometheus.Registerer) error
}
