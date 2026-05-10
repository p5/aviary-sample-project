package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/p5/aviary-sample-project/pkg/api"
	"github.com/p5/aviary-sample-project/pkg/auth"
	"github.com/p5/aviary-sample-project/pkg/store"
)

func main() {
	db := store.NewMemoryStore()
	authenticator := auth.New(os.Getenv("JWT_SECRET"))

	handler := api.NewHandler(db, authenticator)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Starting server on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
