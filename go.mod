module github.com/Cdaprod/registry-service

go 1.19

require (
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/rs/cors v1.11.1
	go.uber.org/zap v1.27.0
)

require go.uber.org/multierr v1.10.0 // indirect

replace github.com/Cdaprod/repocate => ../repocate
