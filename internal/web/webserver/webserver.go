package webserver

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]http.HandlerFunc
	WebServerPort string
}

func NewWebServer(webServerPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]http.HandlerFunc),
		WebServerPort: webServerPort,
	}
}

func (s *WebServer) AddHandler(path string, handler http.HandlerFunc) {
	s.Handlers[path] = handler
}

func (s *WebServer) Start() {
	// Configurar o middleware de logger
	logger := middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log.New(os.Stdout, "", log.LstdFlags)})
	s.Router.Use(logger)

	// Adicionar handlers
	for path, handler := range s.Handlers {
		s.Router.Post(path, handler)
	}

	// Iniciar servidor HTTP
	log.Printf("Servidor da web iniciado na porta %s\n", s.WebServerPort)
	if err := http.ListenAndServe(s.WebServerPort, s.Router); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %s", err)
	}
}
