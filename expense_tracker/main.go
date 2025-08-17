package main

import (
	mainController "expense_tracker/controller"

	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// This is handler for the request
	mux := http.NewServeMux()
	mux.HandleFunc("/expense", RequestHandler)

	port := os.Getenv("PORT")
	if port == "" {
		log.Println("main: PORT not set, defaulting to 8080")
		port = "8080"
	} else {
		log.Printf("main: PORT found: %s", port)
	}

	// Create and configure the HTTP server with sane timeouts and logging middleware.
	log.Println("main: constructing http.Server instance")
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  12 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// Start the server (blocking call).
	log.Printf("Server listening on http://localhost:%s", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("main: server error: %v", err)
	}
	log.Println("main: server exited")
}

func RequestHandler(w http.ResponseWriter, req *http.Request) {
	mainController.MainController(w, req)
	// resp := map[string]string{"body": "Hello, World! from request handler"}
	// _ = json.NewEncoder(w).Encode(resp)
}
