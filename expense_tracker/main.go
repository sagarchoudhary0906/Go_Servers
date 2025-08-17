package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

// requestCounter increments for every request to help generate simple request IDs.
var requestCounter uint64

// loggingResponseWriter wraps http.ResponseWriter to capture status code and bytes written.
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	bytesSent  int
}

// WriteHeader captures the status code before delegating.
func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

// Write captures how many bytes were written. If WriteHeader wasn't called, default to 200.
func (lrw *loggingResponseWriter) Write(body []byte) (int, error) {
	if lrw.statusCode == 0 {
		lrw.statusCode = http.StatusOK
	}
	n, err := lrw.ResponseWriter.Write(body)
	lrw.bytesSent += n
	return n, err
}

// newRequestID returns a lightweight unique-ish ID for a request.
func newRequestID() string {
	c := atomic.AddUint64(&requestCounter, 1)
	return fmt.Sprintf("req-%d-%d", time.Now().UnixNano(), c)
}

func main() {
	log.Println("main: application starting")

	// Create a new ServeMux to route incoming requests.
	mux := http.NewServeMux()
	log.Println("main: ServeMux created")

	// Root route: provides basic API info.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("handler /: start - preparing response")
		w.Header().Set("Content-Type", "application/json")
		log.Println("handler /: set Content-Type header to application/json")
		w.WriteHeader(http.StatusOK)
		log.Println("handler /: wrote status 200 OK")
		_, err := w.Write([]byte(`{"message":"Expense Tracker API","status":"ok"}`))
		if err != nil {
			log.Printf("handler /: error writing response body: %v", err)
		} else {
			log.Println("handler /: wrote response body successfully")
		}
		log.Println("handler /: end")
	})

	// Health check route: used by monitors or load balancers to verify liveness.
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		log.Println("handler /healthz: start - preparing response")
		w.Header().Set("Content-Type", "application/json")
		log.Println("handler /healthz: set Content-Type header to application/json")
		w.WriteHeader(http.StatusOK)
		log.Println("handler /healthz: wrote status 200 OK")
		_, err := w.Write([]byte(`{"status":"healthy"}`))
		if err != nil {
			log.Printf("handler /healthz: error writing response body: %v", err)
		} else {
			log.Println("handler /healthz: wrote response body successfully")
		}
		log.Println("handler /healthz: end")
	})

	// Determine the port to listen on (PORT env var, default 8080).
	log.Println("main: reading PORT environment variable")
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
		Addr:              ":" + port,
		Handler:           loggingMiddleware(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("main: server configured to listen on :%s", port)

	// Start the server (blocking call).
	log.Printf("Server listening on http://localhost:%s", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("main: server error: %v", err)
	}
	log.Println("main: server exited")
}

// loggingMiddleware adds structured request logging around each handler invocation.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a request-scoped ID for tracing.
		requestID := newRequestID()
		start := time.Now()
		log.Printf("req %s: start %s %s from %s", requestID, r.Method, r.URL.RequestURI(), r.RemoteAddr)

		// Wrap ResponseWriter to capture status and size.
		lrw := &loggingResponseWriter{ResponseWriter: w}

		// Ensure we log completion details and recover from panics.
		defer func() {
			duration := time.Since(start)
			if rec := recover(); rec != nil {
				// If a panic occurred, ensure a 500 is recorded and reported to the client if not already.
				if lrw.statusCode == 0 {
					lrw.WriteHeader(http.StatusInternalServerError)
				}
				log.Printf("req %s: panic recovered: %v", requestID, rec)
				http.Error(lrw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			if lrw.statusCode == 0 {
				// No explicit status was written; default is 200 OK per net/http behavior when body is written.
				lrw.statusCode = http.StatusOK
			}
			log.Printf(
				"req %s: complete %s %s -> %d (%d bytes) in %s UA=%q",
				requestID,
				r.Method,
				r.URL.RequestURI(),
				lrw.statusCode,
				lrw.bytesSent,
				duration,
				r.UserAgent(),
			)
		}()

		// Call the next handler in the chain.
		next.ServeHTTP(lrw, r)
	})
}
