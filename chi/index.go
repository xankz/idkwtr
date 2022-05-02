package chi

import (
	"net/http"
)

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if err := s.tmpl.execute(w, "templates/index.tmpl", nil); err != nil {
		// TODO: Implement custom error handler.
		panic(err)
	}
}
