package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"calculator/backend/api/controllers"
	"calculator/backend/api/repositories"
	"calculator/backend/api/routes"
	"calculator/backend/api/services"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	service := services.NewCalculatorService(repositories.NewMemoryHistory)
	controller := controllers.NewCalculatorController(service)
	router := routes.NewRouter(controller)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Printf("Listening on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
