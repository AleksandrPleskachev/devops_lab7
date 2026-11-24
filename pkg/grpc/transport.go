package grpc

import (
	"context"

	pb "DevOps/proto"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kitgrpc "github.com/go-kit/kit/transport/grpc"

	"DevOps/pkg/course"
)

// MakeCreateCourseHandler создает gRPC handler для Create
func MakeCreateCourseHandler(ep endpoint.Endpoint, logger kitlog.Logger) kitgrpc.Handler {
	return kitgrpc.NewServer(
		ep,
		func(ctx context.Context, req interface{}) (interface{}, error) {
			r := req.(*pb.CreateCourseRequest)
			courseReq := course.CreateRequest{
				Title:      r.GetTitle(),
				Instructor: r.GetInstructor(),
				Duration:   int(r.GetDuration()),
			}
			return courseReq, nil
		},
		func(ctx context.Context, resp interface{}) (interface{}, error) {
			r := resp.(course.CreateResponse)
			if r.Err != nil {
				return &pb.CreateCourseResponse{Error: r.Err.Error()}, r.Err
			}
			return &pb.CreateCourseResponse{Id: r.ID}, nil
		},
		kitgrpc.ServerErrorLogger(logger),
	)
}

// MakeGetCourseHandler создает gRPC handler для GetByID
func MakeGetCourseHandler(ep endpoint.Endpoint, logger kitlog.Logger) kitgrpc.Handler {
	return kitgrpc.NewServer(
		ep,
		func(ctx context.Context, req interface{}) (interface{}, error) {
			r := req.(*pb.GetCourseRequest)
			return course.GetByIDRequest{ID: r.GetId()}, nil
		},
		func(ctx context.Context, resp interface{}) (interface{}, error) {
			r := resp.(course.GetByIDResponse)
			if r.Err != nil {
				return &pb.GetCourseResponse{Error: r.Err.Error()}, r.Err
			}
			c := r.Course
			return &pb.GetCourseResponse{
				Course: &pb.CourseMessage{
					Id:         c.ID,
					Title:      c.Title,
					Instructor: c.Instructor,
					Duration:   int32(c.Duration),
					Published:  c.Published,
					CreatedAt:  c.CreatedAt,
				},
			}, nil
		},
		kitgrpc.ServerErrorLogger(logger),
	)
}

// MakePublishCourseHandler создает gRPC handler для Publish
func MakePublishCourseHandler(ep endpoint.Endpoint, logger kitlog.Logger) kitgrpc.Handler {
	return kitgrpc.NewServer(
		ep,
		func(ctx context.Context, req interface{}) (interface{}, error) {
			r := req.(*pb.PublishCourseRequest)
			return course.PublishRequest{ID: r.GetId()}, nil
		},
		func(ctx context.Context, resp interface{}) (interface{}, error) {
			r := resp.(course.PublishResponse)
			if r.Err != nil {
				return &pb.PublishCourseResponse{Error: r.Err.Error()}, r.Err
			}
			return &pb.PublishCourseResponse{Message: "course published"}, nil
		},
		kitgrpc.ServerErrorLogger(logger),
	)
}
