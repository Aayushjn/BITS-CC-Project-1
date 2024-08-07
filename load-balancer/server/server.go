package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aayushjn/load-balancer/balancer"
	"github.com/aayushjn/load-balancer/balancer/backend"
	"github.com/go-chi/chi/v5"
)

func apiRouter(r chi.Router, lb *balancer.LoadBalancer) chi.Router {
	return r.Route("/_api/backend", func(r chi.Router) {
		r.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			backends := lb.Backends()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(backends)
		}))
		r.Post("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			var params map[string]any
			err := json.NewDecoder(r.Body).Decode(&params)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			backendUrl := params["url"].(string)

			// TODO: Make POST idempotent for the same backendUrl

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
			var params map[string]any
			err := json.NewDecoder(r.Body).Decode(&params)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			backendUrl := params["url"].(string)
			err = lb.Unregister(backendUrl)
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
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/_api/backend/" {
			if r.Method == http.MethodGet {
				backends := lb.Backends()
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(backends)
			} else if r.Method == http.MethodPost {
				defer r.Body.Close()
				var params map[string]any
				err := json.NewDecoder(r.Body).Decode(&params)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				backendUrl := params["url"].(string)

				// TODO: Make POST idempotent for the same backendUrl

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
			} else if r.Method == http.MethodDelete {
				defer r.Body.Close()
				var params map[string]any
				err := json.NewDecoder(r.Body).Decode(&params)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				backendUrl := params["url"].(string)
				err = lb.Unregister(backendUrl)
				if err == nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(err.Error()))
					return
				}
				w.WriteHeader(http.StatusNoContent)
			} else {
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
			return
		}

		lb.RequestHandler()(w, r)
	})

	server := &http.Server{
		Addr:              fmt.Sprintf("0.0.0.0:%d", port),
		Handler:           handler,
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}
	return server
}
