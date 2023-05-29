package api

import (
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"palo-alto/conf"
	logger "palo-alto/logging"
	"palo-alto/metric"
	"palo-alto/service"
)

// API serves the end users requests.
type API struct {
	cfg     conf.Config
	m       *metric.Metric
	Router  *httprouter.Router
	server  *http.Server
	service service.DictionaryHandler
}

// New return new API instance
func New(cfg conf.Config, m *metric.Metric, service service.DictionaryHandler) *API {
	return &API{
		cfg:     cfg,
		m:       m,
		service: service,
	}
}

// Title returns the title.
func (api *API) Title() string {
	return "API"
}

// Start starts the http server and binds the handlers.
func (api *API) Start() {
	api.Initialize()
	api.startServer()
}

// Stop stops server
func (api *API) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), api.cfg.GracefulShutdownSec)
	defer cancel()

	api.server.SetKeepAlivesEnabled(false)

	err := api.server.Shutdown(ctx)
	if err != nil {
		logger.Log.Errorf("api shutdown error: %s" + err.Error())
	}
}

// Initialize init api
func (api *API) Initialize() {
	api.Router = httprouter.New()

	logMiddleware := []func(next httprouter.Handle, name string) httprouter.Handle{
		api.RequestLogger,
	}

	api.registerRoutes("GET", "/api/v1/similar", api.getSimilar, logMiddleware...)
	api.registerRoutes("GET", "/api/v1/stats", api.getStats, logMiddleware...)

	api.Router.GET("/health", api.Health)

	api.server = &http.Server{
		Addr:         api.cfg.Port,
		Handler:      api.Router,
		ReadTimeout:  api.cfg.ServerReadTimeoutSec,
		WriteTimeout: api.cfg.ServerWriteTimeoutSec,
		IdleTimeout:  api.cfg.ServerIdleTimeoutSec,
	}
}

func (api *API) registerRoutes(method, path string, handler httprouter.Handle, mws ...func(next httprouter.Handle, name string) httprouter.Handle) {
	for _, mw := range mws {
		handler = mw(handler, path)
	}

	api.Router.Handle(method, path, handler)
}

func (api *API) startServer() {
	log.Printf("Listening on port %s", api.cfg.Port)
	if err := api.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Log.Fatal("Error can't launch the server on port: " + api.cfg.Port)
	}
}

// JSON writes to ResponseWriter a single JSON-object
func JSON(w http.ResponseWriter, data interface{}) {
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		logger.Log.Error(err)
	}
}
