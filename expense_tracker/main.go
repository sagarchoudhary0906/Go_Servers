package main

import (
	// "encoding/json"
	mainController "expense_tracker/controller"

	"database/sql"

	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {

	db, err := openPostgres() // Initialize database pool (verifies DB with Ping inside)
	if err != nil {
		log.Fatalf("main: failed to open postgres: %v", err) // Fail fast if DB is unreachable/misconfigured
	}
	defer db.Close() // Ensure pool is closed on program exit (graceful resource cleanup)

	// This is handler for the request
	mux := http.NewServeMux()
	mux.HandleFunc("/expense", RequestHandler)

	// // Check DB connectivity
	// mux.HandleFunc("/db/ping", func(w http.ResponseWriter, r *http.Request) {
	// 	var one int
	// 	if err := db.QueryRow("SELECT 1").Scan(&one); err != nil {
	// 		http.Error(w, "db not reachable", http.StatusInternalServerError)
	// 		return
	// 	}
	// 	w.Header().Set("Content-Type", "application/json")
	// 	_ = json.NewEncoder(w).Encode(map[string]any{"db": "ok", "value": one})
	// })

	// // Read user row: SELECT id FROM users WHERE login_token=?
	// mux.HandleFunc("/users/get", func(w http.ResponseWriter, r *http.Request) {
	// 	loginToken := r.URL.Query().Get("login_token")
	// 	if loginToken == "" {
	// 		http.Error(w, "login_token required", http.StatusBadRequest)
	// 		return
	// 	}
	// 	var id string
	// 	err := db.QueryRow(`SELECT id FROM users WHERE login_token=$1`, loginToken).Scan(&id)
	// 	if err == sql.ErrNoRows {
	// 		http.Error(w, "not found", http.StatusNotFound)
	// 		return
	// 	}
	// 	if err != nil {
	// 		http.Error(w, "query failed: "+err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// 	w.Header().Set("Content-Type", "application/json")
	// 	_ = json.NewEncoder(w).Encode(map[string]string{"login_token": loginToken, "id": id})
	// })

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

func openPostgres() (*sql.DB, error) {
	// Read the Postgres connection string (DSN) from the environment.
	// Example DSN: postgres://USER:PASS@HOST:PORT/DB?sslmode=disable
	dsn := os.Getenv("POSTGRES_URL")

	// If no env var is set, fall back to a sensible local default that matches your Docker setup.
	if dsn == "" {
		dsn = "postgres://postgres:dev@localhost:5432/expense_tracker?sslmode=disable"
	}

	// Initialize a database handle using the pgx driver. This does NOT open connections yet;
	// it just prepares the pool with the given DSN.
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// Pool settings:
	// Limit the total number of open connections to the DB (in-use + idle).
	db.SetMaxOpenConns(5)

	// Limit how many idle (pre-warmed) connections to keep ready for new requests.
	db.SetMaxIdleConns(5)

	// Set a maximum lifetime for each connection. Forces periodic rotation to avoid
	// stale connections and play nice with proxies/load balancers.
	db.SetConnMaxLifetime(30 * time.Minute)

	// Actually verify connectivity and credentials by pinging the DB.
	// This dials the database and returns an error if unreachable.
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Return the ready-to-use pooled DB handle to callers.
	return db, nil
}
