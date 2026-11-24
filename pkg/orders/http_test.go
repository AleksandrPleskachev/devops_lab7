package orders

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOrderHTTPWorkflow(t *testing.T) {
	handler := NewHTTPHandler(NewService())
	call := func(method, path, body string, expected int) *httptest.ResponseRecorder {
		t.Helper()
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, httptest.NewRequest(method, path, bytes.NewBufferString(body)))
		if response.Code != expected {
			t.Fatalf("%s %s: status %d, want %d; body %s", method, path, response.Code, expected, response.Body)
		}
		return response
	}
	create := call(http.MethodPost, "/orders", `{"customer_id":"c","restaurant_id":"r","order_items":[{"quantity":1,"unit_price":10}]}`, http.StatusOK)
	var id struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(create.Body.Bytes(), &id); err != nil || id.ID == "" {
		t.Fatalf("invalid create response: %s", create.Body)
	}
	path := "/orders/" + id.ID
	var order Order
	read := call(http.MethodGet, path, "", http.StatusOK)
	if err := json.Unmarshal(read.Body.Bytes(), &order); err != nil || order.Status != "Pending" {
		t.Fatalf("expected pending order: %s", read.Body)
	}
	call(http.MethodPatch, path+"/status", `{"status":"Completed"}`, http.StatusOK)
	read = call(http.MethodGet, path, "", http.StatusOK)
	if err := json.Unmarshal(read.Body.Bytes(), &order); err != nil || order.Status != "Completed" || len(order.Items) != 1 {
		t.Fatalf("expected completed order with one item: %s", read.Body)
	}
	call(http.MethodGet, "/orders/missing", "", http.StatusNotFound)
}
