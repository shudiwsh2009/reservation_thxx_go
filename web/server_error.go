package web

import (
	"net/http"
)

func (s *Server) RenderError404(w http.ResponseWriter) {
	s.render.Html(w, http.StatusNotFound, "page_not_found_error", nil)
}

func (s *Server) RenderError500(w http.ResponseWriter, err error) {
	params := make(map[string]interface{})
	s.render.Html(w, http.StatusInternalServerError, "internal_server_error", params)
}
