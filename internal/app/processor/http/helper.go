package rprocessor

import (
	"net/http"

	"github.com/RonIT-401/catalog-service/internal/pkg/http/httph"
	"github.com/gorilla/mux"
)

func reg(r *mux.Router, method, path string, handler http.Handler) {
	r.Methods(method).Path(path).Handler(handler)
}

func middlewaresToGorilla(middlewares []httph.Middleware) []mux.MiddlewareFunc {
	out := make([]mux.MiddlewareFunc, 0, len(middlewares))
	for _, mw := range middlewares {
		out = append(out, mux.MiddlewareFunc(mw))
	}
	return out
}
