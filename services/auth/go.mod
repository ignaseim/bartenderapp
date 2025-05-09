module github.com/yourusername/bartenderapp/services/auth

go 1.22

require (
	github.com/golang-jwt/jwt/v5 v5.0.0
	github.com/gorilla/mux v1.8.1
	github.com/lib/pq v1.10.9
	github.com/nats-io/nats.go v1.28.0
	github.com/prometheus/client_golang v1.16.0
	github.com/yourusername/bartenderapp/services/pkg v0.0.0
	golang.org/x/crypto v0.17.0
)

replace github.com/yourusername/bartenderapp/services/pkg => ../pkg 