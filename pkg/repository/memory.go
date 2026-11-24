package repository

import (
	"context"
	"fmt"
	"sync"

	"DevOps/pkg/course"
)

// MemoryRepository — in-memory хранилище курсов
type MemoryRepository struct {
	mu      sync.RWMutex
	courses map[string]course.Course
}

// NewMemoryRepository создает новое in-memory хранилище
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		courses: make(map[string]course.Course),
	}
}

// Create сохраняет новый курс
func (r *MemoryRepository) Create(ctx context.Context, c course.Course) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.courses[c.ID]; exists {
		return fmt.Errorf("course with id %s already exists", c.ID)
	}

	r.courses[c.ID] = c
	return nil
}

// GetByID находит курс по ID
func (r *MemoryRepository) GetByID(ctx context.Context, id string) (course.Course, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	c, exists := r.courses[id]
	if !exists {
		return course.Course{}, fmt.Errorf("course with id %s not found", id)
	}
	return c, nil
}

// Publish меняет статус публикации курса
func (r *MemoryRepository) Publish(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	c, exists := r.courses[id]
	if !exists {
		return fmt.Errorf("course with id %s not found", id)
	}

	c.Published = true
	r.courses[id] = c
	return nil
}

// List возвращает все курсы
func (r *MemoryRepository) List(ctx context.Context) ([]course.Course, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]course.Course, 0, len(r.courses))
	for _, c := range r.courses {
		result = append(result, c)
	}
	return result, nil
}
