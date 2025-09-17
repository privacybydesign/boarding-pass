package main

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type Server struct {
	server *http.Server
	config *ServerConfig
}

type ServerState struct {
	irmaServerURL    string
	tokenStorage     TokenStorage
	credentialConfig CredentialConfig
	ticketStore      *TicketStore
	sessionTracker   *SessionTracker
}

type SpaHandler struct {
	staticPath string
	indexPath  string
}

type routeRegistrar func(router *mux.Router, state *ServerState)

func NewServer(state *ServerState, config *ServerConfig) *Server {
	router := mux.NewRouter()

	registerRoutes := []routeRegistrar{
		registerTicketRoutes,
		registerSessionRoutes,
		registerCallbackRoute,
	}
	for _, register := range registerRoutes {
		register(router, state)
	}

	spa := SpaHandler{staticPath: "../frontend/dist", indexPath: "index.html"}
	router.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler:      router,
		Addr:         config.Host + ":" + strconv.Itoa(config.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return &Server{
		server: srv,
		config: config,
	}
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func (h SpaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(h.staticPath, r.URL.Path)
	fi, err := os.Stat(path)
	if os.IsNotExist(err) || fi.IsDir() {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func registerTicketRoutes(router *mux.Router, state *ServerState) {
	router.HandleFunc("/api/tickets", func(w http.ResponseWriter, r *http.Request) {
		handleCreateTicket(w, r, state)
	}).Methods(http.MethodPost)

	router.HandleFunc("/api/tickets/{ticketId}", func(w http.ResponseWriter, r *http.Request) {
		handleGetTicket(w, r, state)
	}).Methods(http.MethodGet)
}

func registerSessionRoutes(router *mux.Router, state *ServerState) {
	router.HandleFunc("/api/start", func(w http.ResponseWriter, r *http.Request) {
		handleStartIRMASession(w, r, state)
	}).Methods(http.MethodPost)

	router.HandleFunc("/api/result", func(w http.ResponseWriter, r *http.Request) {
		handleResultIRMASession(w, r, state)
	}).Methods(http.MethodGet)
}

func registerCallbackRoute(router *mux.Router, state *ServerState) {
	router.HandleFunc("/api/irma/callback", func(w http.ResponseWriter, r *http.Request) {
		handleIRMAServerCallback(w, r, state)
	}).Methods(http.MethodPost)
}
