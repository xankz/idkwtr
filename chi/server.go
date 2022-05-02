// Package chi exposes a web server that interfaces with the nibo application. The main dependency
// is go-chi/chi for the router, along with html/template for rendering.
package chi

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/elijk/nibo"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Services is a group of domain services to facilitate processing client actions.
type Services struct {
	UserRepository nibo.UserRepository
}

// Configuration represents the server-specific configuration of a server.
type Configuration struct {
	// Host represents the host address that the server will listen on.
	Host string

	// Port represents the port that the server will listen on.
	Port int
}

// Server is a Chi web server that interfaces with the main application.
type Server struct {
	r        chi.Router
	tmpl     map[string]*template.Template
	services Services
	config   Configuration
}

// ServerOpts represents a set of options for creating a new server.
type ServerOpts struct {
	Services      Services
	Configuration Configuration

	// StaticDir is a path that points to a directory of publicly-available files.
	StaticDir string
}

// NewServerFromChi returns a new server, given an existing Chi router. Useful for consumers
// desiring to load custom middleware.
func NewServerFromChi(opts ServerOpts, r chi.Router) (*Server, error) {
	s := &Server{
		r:        r,
		services: opts.Services,
		config:   opts.Configuration,
	}
	if err := setupServer(s, opts.StaticDir); err != nil {
		return nil, fmt.Errorf("creating server: %w", err)
	}
	return s, nil
}

// NewServer returns a new Server.
func NewServer(opts ServerOpts) (*Server, error) {
	s := &Server{
		r:        chi.NewRouter(),
		services: opts.Services,
		config:   opts.Configuration,
	}
	if err := setupServer(s, opts.StaticDir); err != nil {
		return nil, fmt.Errorf("creating server: %w", err)
	}
	return s, nil
}

// Start starts the server with http.ListenAndServe. This blocks until an error is returned from
// halting.
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	return http.ListenAndServe(addr, s.r)
}

func setupServer(s *Server, sdir string) error {
	// Setup templates.
	t, err := parseTemplates(templatePartials, templateViews, templateFuncMap)
	if err != nil {
		return fmt.Errorf("parsing templates: %w", err)
	}
	s.tmpl = t

	// Setup server middleware.
	s.r.Use(middleware.Logger)

	// Setup static filesystem access.
	if sdir != "" {
		s.r.Handle("/*", http.FileServer(http.Dir(sdir)))
	}

	// Setup Chi routes.
	s.r.Get("/", s.handleIndex)

	return nil
}
