package main

import (
	"context"
	"log"
	"net/http"

	"github.com/necrophonic/mocking-mongo/pkg/api"
)

func main() {
	ctx := context.Background()

	a := &api.API{}
	if err := a.ConnectDB(ctx, "mongodb://localhost:27017"); err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/thing", a.ThingHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
