package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/briheet/nozarashi/examples/simple_docker/backend/internal/config"
	"github.com/briheet/nozarashi/examples/simple_docker/backend/internal/db"
	"github.com/briheet/nozarashi/examples/simple_docker/backend/internal/logger"
	"github.com/briheet/nozarashi/examples/simple_docker/backend/internal/redis"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

// API owns the HTTP transport layer and the application services exposed through it.
type API struct {
	config   *config.Config
	logger   *logger.Logger
	validate *validator.Validate
	redis    *redis.Client
}

// NewAPI is the HTTP composition root. It constructs repositories and services once,
// then keeps handlers dependent on business interfaces rather than infrastructure details.
func NewAPI(
	config *config.Config,
	logger *logger.Logger,
	dbClient *db.Client,
	redisClient *redis.Client,
) *API {
	// One validator instance enforces the request and domain-event contracts.
	validate := validator.New(validator.WithRequiredStructEnabled())

	return &API{
		config:   config,
		logger:   logger,
		validate: validate,
		redis:    redisClient,
	}
}

// Server applies configured network timeouts around the API router.
func (a *API) Server(port int) *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           a.Routes(),
		ReadHeaderTimeout: time.Duration(a.config.API.ReadHeaderTimeout) * time.Second,
		ReadTimeout:       time.Duration(a.config.API.ReadTimeout) * time.Second,
		WriteTimeout:      time.Duration(a.config.API.WriteTimeout) * time.Second,
		IdleTimeout:       time.Duration(a.config.API.IdleTimeout) * time.Second,
	}
}

// Routes assembles middleware, operational endpoints, and the versioned public API.
func (a *API) Routes() *mux.Router {
	r := mux.NewRouter()

	// Business endpoints are versioned independently from operational metrics.
	sub := r.PathPrefix("/api/v1").Subrouter()

	a.registerHealthHandlers(sub)

	r.PathPrefix("/").Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	return r
}
