package orders

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

type endpoints struct {
	create, get, changeStatus endpoint.Endpoint
}

type statusRequest struct {
	ID     string `json:"-"`
	Status string `json:"status"`
}

func makeEndpoints(s Service) endpoints {
	return endpoints{
		create: func(ctx context.Context, request any) (any, error) {
			id, err := s.Create(ctx, request.(Order))
			return map[string]string{"id": id}, err
		},
		get: func(ctx context.Context, request any) (any, error) {
			return s.GetByID(ctx, request.(string))
		},
		changeStatus: func(ctx context.Context, request any) (any, error) {
			r := request.(statusRequest)
			return map[string]string{"status": r.Status}, s.ChangeStatus(ctx, r.ID, r.Status)
		},
	}
}

func NewHTTPHandler(s Service) http.Handler {
	e := makeEndpoints(s)
	r := mux.NewRouter()
	opt := kithttp.ServerErrorEncoder(encodeError)
	r.Handle("/orders", kithttp.NewServer(e.create, decodeCreate, encodeJSON, opt)).Methods(http.MethodPost)
	r.Handle("/orders/{id}", kithttp.NewServer(e.get, decodeID, encodeJSON, opt)).Methods(http.MethodGet)
	r.Handle("/orders/{id}/status", kithttp.NewServer(e.changeStatus, decodeStatus, encodeJSON, opt)).Methods(http.MethodPatch)
	return r
}

func decodeCreate(_ context.Context, r *http.Request) (any, error) {
	var o Order
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		return nil, ErrInvalid
	}
	return o, nil
}

func decodeID(_ context.Context, r *http.Request) (any, error) {
	return mux.Vars(r)["id"], nil
}

func decodeStatus(_ context.Context, r *http.Request) (any, error) {
	req := statusRequest{ID: mux.Vars(r)["id"]}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, ErrInvalid
	}
	return req, nil
}

func encodeJSON(_ context.Context, w http.ResponseWriter, v any) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func encodeError(ctx context.Context, err error, w http.ResponseWriter) {
	code := http.StatusInternalServerError
	switch {
	case errors.Is(err, ErrInvalid):
		code = http.StatusBadRequest
	case errors.Is(err, ErrNotFound):
		code = http.StatusNotFound
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = encodeJSON(ctx, w, map[string]string{"error": err.Error()})
}
