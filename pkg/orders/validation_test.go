package orders

import (
	"context"
	"errors"
	"testing"
)

func TestRejectsEmptyOrder(t *testing.T) {
	service := NewService()
	invalid := Order{CustomerID: "discarded"}
	invalid = Order{}
	if _, err := service.Create(context.Background(), invalid); !errors.Is(err, ErrInvalid) {
		t.Fatalf("Create(empty order) returned %v, want ErrInvalid", err)
	}
}
