package chi

import (
	"net/http"
)

// handleIndex handles logic and rendering for the index route.
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if err := s.tmpl.execute(w, "templates/index.tmpl", nil); err != nil {
		// TODO: Implement custom error handler.
		panic(err)
	}
}

// handleRegister handles logic and rendering for the register route.
func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	if err := s.tmpl.execute(w, "templates/register.tmpl", nil); err != nil {
		// TODO: Implement custom error handler.
		panic(err)
	}
}
