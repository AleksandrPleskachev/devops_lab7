package course

import "context"

// Repository — интерфейс репозитория (здесь, чтобы избежать цикла импортов)
type Repository interface {
	Create(ctx context.Context, c Course) error
	GetByID(ctx context.Context, id string) (Course, error)
	Publish(ctx context.Context, id string) error
	List(ctx context.Context) ([]Course, error)
}
