package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aayushjn/load-balancer/balancer"
	"github.com/aayushjn/load-balancer/balancer/backend"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func apiRouter(r chi.Router, lb *balancer.LoadBalancer) chi.Router {
	return r.Route("/_api/backend", func(r chi.Router) {
		r.Post("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			var params map[string]any
			err := json.NewDecoder(r.Body).Decode(&params)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			backendUrl := params["url"].(string)

			b, err := backend.NewBackend(backendUrl)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			err = lb.Register(b, params)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		}))
		r.Delete("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			body, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			err = lb.Unregister(string(body))
			if err == nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			w.WriteHeader(http.StatusNoContent)
		}))
	})
}

func NewServer(lb *balancer.LoadBalancer, port int) *http.Server {
	r := chi.NewRouter()
	r.Use(middleware.Logger, middleware.Recoverer)

	apiRouter(r, lb)
	r.Get("/", lb.RequestHandler())

	server := &http.Server{
		Addr:              fmt.Sprintf("0.0.0.0:%d", port),
		Handler:           r,
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}
	return server
}
