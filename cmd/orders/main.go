package main

import (
	"log"
	"net/http"
	"time"

	"DevOps/pkg/orders"
)

func main() {
	srv := &http.Server{
		Addr:              "127.0.0.1:8080",
		Handler:           orders.NewHTTPHandler(orders.NewService()),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("Orders HTTP on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
