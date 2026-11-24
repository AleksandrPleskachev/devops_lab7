package course

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints — набор эндпоинтов go-kit
type Endpoints struct {
	CreateEndpoint  endpoint.Endpoint
	GetByIDEndpoint endpoint.Endpoint
	PublishEndpoint endpoint.Endpoint
}

// CreateEndpoint — эндпоинт для создания курса
func CreateEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)
		id, err := svc.Create(ctx, req)
		return CreateResponse{ID: id, Err: err}, nil
	}
}

// GetByIDEndpoint — эндпоинт для получения курса по ID
func GetByIDEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetByIDRequest)
		course, err := svc.GetByID(ctx, req.ID)
		return GetByIDResponse{Course: course, Err: err}, nil
	}
}

// PublishEndpoint — эндпоинт для публикации курса
func PublishEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(PublishRequest)
		err := svc.Publish(ctx, req.ID)
		return PublishResponse{Err: err}, nil
	}
}
