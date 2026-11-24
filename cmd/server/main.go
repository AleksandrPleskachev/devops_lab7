package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	kitlog "github.com/go-kit/kit/log"
	"github.com/gorilla/mux"

	"DevOps/pkg/course"
	"DevOps/pkg/grpc"
	"DevOps/pkg/repository"
)

func main() {
	// Логгер
	logger := kitlog.NewLogfmtLogger(os.Stdout)

	// Репозиторий (in-memory)
	repo := repository.NewMemoryRepository()

	// Сервис
	svc := course.NewCourseService(repo)

	// Эндпоинты
	endpoints := course.Endpoints{
		CreateEndpoint:  course.CreateEndpoint(svc),
		GetByIDEndpoint: course.GetByIDEndpoint(svc),
		PublishEndpoint: course.PublishEndpoint(svc),
	}

	// ============================================
	// HTTP транспорт (порт 8080)
	// ============================================
	r := mux.NewRouter()
	r.Methods("POST").Path("/courses").Handler(course.MakeCreateHandler(endpoints.CreateEndpoint, logger))
	r.Methods("GET").Path("/courses/{id}").Handler(course.MakeGetByIDHandler(endpoints.GetByIDEndpoint, logger))
	r.Methods("POST").Path("/courses/{id}/publish").Handler(course.MakePublishHandler(endpoints.PublishEndpoint, logger))

	httpSrv := &http.Server{Addr: address("HTTP_ADDR", "127.0.0.1:8080"), ReadHeaderTimeout: 5 * time.Second, Handler: r}

	go func() {
		fmt.Println("HTTP server running on http://localhost:8080")
		fmt.Println("  POST   /courses              — создать курс")
		fmt.Println("  GET    /courses/{id}         — получить курс по ID")
		fmt.Println("  POST   /courses/{id}/publish — опубликовать курс")
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// ============================================
	// gRPC транспорт (порт 9090)
	// ============================================
	grpcSrv := grpc.NewServer(endpoints, logger)

	go func() {
		if err := grpcSrv.Run(address("GRPC_ADDR", "127.0.0.1:9090")); err != nil {
			log.Fatalf("gRPC server failed: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nShutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	httpSrv.Shutdown(ctx)
	grpcSrv.Close()

	fmt.Println("All servers stopped")
}

func address(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
