package orders

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("order not found")
var ErrInvalid = errors.New("invalid order request")

type Item struct {
	ProductCode string  `json:"product_code"`
	Name        string  `json:"name"`
	UnitPrice   float64 `json:"unit_price"`
	Quantity    int     `json:"quantity"`
}

type Order struct {
	ID           string `json:"id"`
	CustomerID   string `json:"customer_id"`
	RestaurantID string `json:"restaurant_id"`
	Status       string `json:"status"`
	CreatedOn    int64  `json:"created_on"`
	Items        []Item `json:"order_items"`
}

type Service interface {
	Create(context.Context, Order) (string, error)
	GetByID(context.Context, string) (Order, error)
	ChangeStatus(context.Context, string, string) error
}

type MemoryRepository struct {
	mu     sync.RWMutex
	orders map[string]Order
}

type service struct{ repo *MemoryRepository }

func NewService() Service {
	return &service{repo: &MemoryRepository{orders: make(map[string]Order)}}
}

func (s *service) Create(_ context.Context, o Order) (string, error) {
	if strings.TrimSpace(o.CustomerID) == "" || strings.TrimSpace(o.RestaurantID) == "" || len(o.Items) == 0 {
		return "", ErrInvalid
	}
	for _, item := range o.Items {
		if item.Quantity <= 0 || item.UnitPrice < 0 {
			return "", ErrInvalid
		}
	}
	o.ID = uuid.NewString()
	o.Status = "Pending"
	o.CreatedOn = time.Now().Unix()
	o.Items = append([]Item(nil), o.Items...)
	s.repo.mu.Lock()
	defer s.repo.mu.Unlock()
	s.repo.orders[o.ID] = o
	return o.ID, nil
}

func (s *service) GetByID(_ context.Context, id string) (Order, error) {
	s.repo.mu.RLock()
	defer s.repo.mu.RUnlock()
	o, ok := s.repo.orders[id]
	if !ok {
		return Order{}, ErrNotFound
	}
	o.Items = append([]Item(nil), o.Items...)
	return o, nil
}

func (s *service) ChangeStatus(_ context.Context, id, status string) error {
	if strings.TrimSpace(status) == "" {
		return ErrInvalid
	}
	s.repo.mu.Lock()
	defer s.repo.mu.Unlock()
	o, ok := s.repo.orders[id]
	if !ok {
		return ErrNotFound
	}
	o.Status = status
	s.repo.orders[id] = o
	return nil
}
