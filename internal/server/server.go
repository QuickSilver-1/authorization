package server

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"auth/internal/logger"

	"github.com/gorilla/mux"
)

type Server struct {
	server *http.Server
}

// NewServer создает новый HTTP сервер с заданными параметрами
func NewServer(port, r, w int) *Server {

	mux := mux.NewRouter()
	mux.Use(LogMiddleware)
	mux.HandleFunc("/auth/login", Login).Methods("GET")
	mux.HandleFunc("/auth/refresh", Refresh).Methods("GET")

	s := &http.Server{
		Addr:         ":" + strconv.Itoa(port),
		Handler:      mux,
		ReadTimeout:  time.Duration(r) * time.Second,
		WriteTimeout: time.Duration(w) * time.Second,
	}

	return &Server{
		server: s,
	}
}

// StartServer запускает HTTP сервер и обрабатывает возможные ошибки
func (s *Server) StartServer() error {
	logger.Log.Debug(fmt.Sprintf("New server creating on port: %s", s.server.Addr))

	err := s.server.ListenAndServe()
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Ошибка при запуске сервера: %v", err))
	}

	return err
}
