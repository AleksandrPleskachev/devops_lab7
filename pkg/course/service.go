package course

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Service — интерфейс сервисного слоя
type Service interface {
	Create(ctx context.Context, req CreateRequest) (string, error)
	GetByID(ctx context.Context, id string) (Course, error)
	Publish(ctx context.Context, id string) error
}

// CourseService — реализация сервисного слоя
type CourseService struct {
	repo Repository
}

// NewCourseService создает новый сервис курсов
func NewCourseService(repo Repository) Service {
	return &CourseService{repo: repo}
}

// Create создает новый курс
func (s *CourseService) Create(ctx context.Context, req CreateRequest) (string, error) {
	if err := req.Validate(); err != nil {
		return "", err
	}

	id := uuid.New().String()
	c := Course{
		ID:         id,
		Title:      req.Title,
		Instructor: req.Instructor,
		Duration:   req.Duration,
		Published:  false,
		CreatedAt:  time.Now().Format(time.RFC3339),
	}

	if err := s.repo.Create(ctx, c); err != nil {
		return "", err
	}

	return id, nil
}

// GetByID получает курс по ID
func (s *CourseService) GetByID(ctx context.Context, id string) (Course, error) {
	return s.repo.GetByID(ctx, id)
}

// Publish публикует курс
func (s *CourseService) Publish(ctx context.Context, id string) error {
	return s.repo.Publish(ctx, id)
}
