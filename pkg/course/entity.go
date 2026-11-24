package course

import "fmt"

// Course — сущность курса
type Course struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Instructor string `json:"instructor"`
	Duration   int    `json:"duration"`
	Published  bool   `json:"published"`
	CreatedAt  string `json:"created_at"`
}

// CreateRequest — запрос на создание курса
type CreateRequest struct {
	Title      string `json:"title"`
	Instructor string `json:"instructor"`
	Duration   int    `json:"duration"`
}

// CreateResponse — ответ на создание курса
type CreateResponse struct {
	ID  string `json:"id"`
	Err error  `json:"-"`
}

// GetByIDRequest — запрос на получение курса по ID
type GetByIDRequest struct {
	ID string `json:"id"`
}

// GetByIDResponse — ответ на получение курса
type GetByIDResponse struct {
	Course Course `json:"course"`
	Err    error  `json:"-"`
}

// PublishRequest — запрос на публикацию курса
type PublishRequest struct {
	ID string `json:"id"`
}

// PublishResponse — ответ на публикацию курса
type PublishResponse struct {
	Err error `json:"-"`
}

// Validate проверяет данные курса
func (r CreateRequest) Validate() error {
	if r.Title == "" {
		return fmt.Errorf("title is required")
	}
	if r.Instructor == "" {
		return fmt.Errorf("instructor is required")
	}
	if r.Duration <= 0 {
		return fmt.Errorf("duration must be positive")
	}
	return nil
}
