package webserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type WebServer struct {
	Router        *chi.Mux
	WebServerPort string
	handlers      []routeEntry
}

type routeEntry struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

func NewWebServer(port string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		WebServerPort: port,
	}
}

func (ws *WebServer) AddHandler(method, path string, h http.HandlerFunc) {
	ws.handlers = append(ws.handlers, routeEntry{Method: method, Path: path, Handler: h})
}

func (ws *WebServer) Start() {
	ws.Router.Use(middleware.Logger)
	ws.Router.Use(middleware.Recoverer)
	for _, e := range ws.handlers {
		switch e.Method {
		case http.MethodGet:
			ws.Router.Get(e.Path, e.Handler)
		case http.MethodPost:
			ws.Router.Post(e.Path, e.Handler)
		case http.MethodPut:
			ws.Router.Put(e.Path, e.Handler)
		case http.MethodDelete:
			ws.Router.Delete(e.Path, e.Handler)
		}
	}
	http.ListenAndServe(":"+ws.WebServerPort, ws.Router)
}
