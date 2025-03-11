package middleware

import (
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Chain(f http.HandlerFunc, middleware ...Middleware) http.HandlerFunc {
	for _, m := range middleware {
		f = m(f)
	}
	return f
}

type HandlerMiddleware func(handler http.Handler) http.Handler

func ChainHandlerFunc(h http.Handler, middleware ...HandlerMiddleware) http.Handler {
	for _, m := range middleware {
		h = m(h)
	}
	return h
}
