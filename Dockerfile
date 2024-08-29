# Build the frontend
FROM node:14 AS frontend-builder
WORKDIR /app/web
COPY web/package*.json ./
RUN npm install
COPY web .
RUN npm run build

# Build the backend
FROM golang:1.16 AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend-builder /app/web/build ./web/build
RUN CGO_ENABLED=0 GOOS=linux go build -o /registry-service ./cmd/server

# Final stage
FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=backend-builder /registry-service .
EXPOSE 8080
CMD ["./registry-service"]