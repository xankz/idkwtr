// Package chi exposes a web server that interfaces with the nibo application. The main dependency
// is go-chi/chi for the router, along with html/template for rendering.
package chi

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	"github.com/elijk/nibo"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	//go:embed templates/*.tmpl templates/**/*.tmpl
	templateFiles embed.FS

	templateFuncMap template.FuncMap = make(template.FuncMap)
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
	tmpl     *template.Template
	services Services
	config   Configuration
}

// ServerOpts represents a set of options for creating a new server.
type ServerOpts struct {
	Services      Services
	Configuration Configuration
}

// NewServerFromChi returns a new server, given an existing Chi router. Useful for consumers
// desiring to load custom middleware.
func NewServerFromChi(opts ServerOpts, r chi.Router) (*Server, error) {
	s := &Server{
		r:        r,
		services: opts.Services,
		config:   opts.Configuration,
	}
	if err := setupServer(s); err != nil {
		return nil, err
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
	if err := setupServer(s); err != nil {
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

func setupServer(s *Server) error {
	// Setup templates.
	t, err := parseTemplates(templateFiles, templateFuncMap)
	if err != nil {
		return fmt.Errorf("parsing templates: %w", err)
	}
	s.tmpl = t

	// Setup server middleware.
	s.r.Use(middleware.Logger)

	// Setup Chi routes.
	s.r.Get("/", s.handleIndex)

	return nil
}

// parseTemplates walks through and parses a filesystem of templates, returning a root template
// having each view loaded as a sub-template. This strategy avoids name collisions unlike the
// built-in template.ParseFiles.
func parseTemplates(
	dir embed.FS,
	funcMap template.FuncMap,
) (*template.Template, error) {
	root := template.New("")

	if err := fs.WalkDir(dir, ".", func(path string, info fs.DirEntry, err error) error {
		// Skip directories and non-HTML templates.
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".tmpl") {
			return err
		}

		fmt.Printf("path: %v\n", path)

		// Read file data and parse into root sub-template.
		body, err := fs.ReadFile(dir, path)
		if err != nil {
			return err
		}
		tmpl := root.New(path).Funcs(funcMap)
		if _, err := tmpl.Parse(string(body)); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return root, nil
}
