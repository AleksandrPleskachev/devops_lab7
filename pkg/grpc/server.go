package grpc

import (
	"context"
	"fmt"
	"net"

	"DevOps/pkg/course"

	pb "DevOps/proto"
	kitlog "github.com/go-kit/kit/log"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	grpc "google.golang.org/grpc"
)

// Server — gRPC сервер
type Server struct {
	srv *grpc.Server
	lis net.Listener
}

// NewServer создает gRPC сервер с обработчиками
func NewServer(endpoints course.Endpoints, logger kitlog.Logger) *Server {
	s := grpc.NewServer()

	pb.RegisterCourseServiceServer(s, &courseGRPCServer{
		createHandler:  MakeCreateCourseHandler(endpoints.CreateEndpoint, logger),
		getHandler:     MakeGetCourseHandler(endpoints.GetByIDEndpoint, logger),
		publishHandler: MakePublishCourseHandler(endpoints.PublishEndpoint, logger),
	})

	return &Server{srv: s}
}

// courseGRPCServer — реализация gRPC сервиса
type courseGRPCServer struct {
	pb.UnimplementedCourseServiceServer
	createHandler  kitgrpc.Handler
	getHandler     kitgrpc.Handler
	publishHandler kitgrpc.Handler
}

func (s *courseGRPCServer) CreateCourse(ctx context.Context, req *pb.CreateCourseRequest) (*pb.CreateCourseResponse, error) {
	_, res, err := s.createHandler.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*pb.CreateCourseResponse), nil
}

func (s *courseGRPCServer) GetCourse(ctx context.Context, req *pb.GetCourseRequest) (*pb.GetCourseResponse, error) {
	_, res, err := s.getHandler.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*pb.GetCourseResponse), nil
}

func (s *courseGRPCServer) PublishCourse(ctx context.Context, req *pb.PublishCourseRequest) (*pb.PublishCourseResponse, error) {
	_, res, err := s.publishHandler.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*pb.PublishCourseResponse), nil
}

// Run запускает gRPC сервер
func (s *Server) Run(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.lis = lis
	fmt.Printf("gRPC server running on %s\n", addr)
	return s.srv.Serve(lis)
}

// Close останавливает gRPC сервер
func (s *Server) Close() error {
	if s.srv != nil {
		s.srv.GracefulStop()
	}
	if s.lis != nil {
		s.lis.Close()
	}
	return nil
}
