package api

import (
	"net/http"
)

func (s *Server) routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("GET /champion-mastery", s.GetChampionMastery)

	v1 := http.NewServeMux()
	v1.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

	chain := MiddlewareChain(
		RequestLoggerMiddleware,
	)
	return chain(v1)
}
