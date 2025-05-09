package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/yourusername/bartenderapp/services/pkg/auth"
)

// RequestLogger logs information about each HTTP request
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Wrap the response writer to capture status code
		lrw := NewLoggingResponseWriter(w)
		
		// Call the next handler
		next.ServeHTTP(lrw, r)
		
		// Log the request details
		log.Printf(
			"method=%s path=%s status=%d duration=%s ip=%s user_agent=%s",
			r.Method,
			r.URL.Path,
			lrw.statusCode,
			time.Since(start),
			r.RemoteAddr,
			r.UserAgent(),
		)
	})
}

// LoggingResponseWriter wraps http.ResponseWriter to capture status code
type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// NewLoggingResponseWriter creates a new LoggingResponseWriter
func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{w, http.StatusOK}
}

// WriteHeader captures the status code before writing the header
func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Authenticate validates the JWT token and adds user info to the request context
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from request
		token, err := auth.ExtractTokenFromRequest(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		
		// Validate token
		claims, err := auth.ValidateToken(token)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}
		
		// Add claims to request context
		ctx := context.WithValue(r.Context(), "claims", claims)
		
		// Call the next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole checks if the user has the required role
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get claims from context
			claims, ok := r.Context().Value("claims").(*auth.Claims)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			
			// Check if user has one of the required roles
			if !auth.HasRole(claims, roles...) {
				http.Error(w, "Forbidden - Insufficient permissions", http.StatusForbidden)
				return
			}
			
			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// CORS middleware adds CORS headers to the response
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// JSONContentType sets the Content-Type header to application/json
func JSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// RespondWithJSON sends a JSON response
func RespondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Printf("Error encoding JSON response: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// RespondWithError sends an error response
func RespondWithError(w http.ResponseWriter, status int, message string) {
	RespondWithJSON(w, status, map[string]string{"error": message})
} 