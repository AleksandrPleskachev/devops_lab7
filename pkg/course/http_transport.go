package course

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// MakeCreateHandler создает HTTP handler для Create
func MakeCreateHandler(ep endpoint.Endpoint, logger kitlog.Logger) http.Handler {
	return kithttp.NewServer(
		ep,
		func(_ context.Context, r *http.Request) (interface{}, error) {
			var req CreateRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				return nil, err
			}
			return req, nil
		},
		func(_ context.Context, w http.ResponseWriter, resp interface{}) error {
			response := resp.(CreateResponse)
			if response.Err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return json.NewEncoder(w).Encode(map[string]string{"error": response.Err.Error()})
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			return json.NewEncoder(w).Encode(map[string]string{"id": response.ID})
		},
		kithttp.ServerErrorLogger(logger),
	)
}

// MakeGetByIDHandler создает HTTP handler для GetByID
func MakeGetByIDHandler(ep endpoint.Endpoint, logger kitlog.Logger) http.Handler {
	return kithttp.NewServer(
		ep,
		func(_ context.Context, r *http.Request) (interface{}, error) {
			vars := mux.Vars(r)
			id := vars["id"]
			return GetByIDRequest{ID: id}, nil
		},
		func(_ context.Context, w http.ResponseWriter, resp interface{}) error {
			response := resp.(GetByIDResponse)
			if response.Err != nil {
				w.WriteHeader(http.StatusNotFound)
				return json.NewEncoder(w).Encode(map[string]string{"error": response.Err.Error()})
			}
			w.Header().Set("Content-Type", "application/json")
			return json.NewEncoder(w).Encode(response.Course)
		},
		kithttp.ServerErrorLogger(logger),
	)
}

// MakePublishHandler создает HTTP handler для Publish
func MakePublishHandler(ep endpoint.Endpoint, logger kitlog.Logger) http.Handler {
	return kithttp.NewServer(
		ep,
		func(_ context.Context, r *http.Request) (interface{}, error) {
			vars := mux.Vars(r)
			id := vars["id"]
			return PublishRequest{ID: id}, nil
		},
		func(_ context.Context, w http.ResponseWriter, resp interface{}) error {
			response := resp.(PublishResponse)
			if response.Err != nil {
				w.WriteHeader(http.StatusNotFound)
				return json.NewEncoder(w).Encode(map[string]string{"error": response.Err.Error()})
			}
			w.Header().Set("Content-Type", "application/json")
			return json.NewEncoder(w).Encode(map[string]string{"message": "course published"})
		},
		kithttp.ServerErrorLogger(logger),
	)
}
