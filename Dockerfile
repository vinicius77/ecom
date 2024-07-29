# Build the application from source
FROM golang:1.18-alpine3.14 AS build-stage 

  WORKDIR /app

  COPY go.mod go.sum ./

  RUN go mod download

  COPY . .

  RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/main.go


# Run the tests in the container
FROM build-stage AS run-test-stage
  
  RUN go test -v ./...

# Run the db migrations in the container
FROM run-test-stage AS run-migrations-stage
  
  RUN make migration add-user-table \
  make migration add-product-table \
  make migration add-order-table \
  make migration add-order-items-table \
  migrate-up

# Deploy the application binary into a lean image
FROM scratch AS build-release-stage
  
  WORKDIR /

  COPY --from=build-stage /api /api

  EXPOSE 8080

  ENTRYPOINT ["/api"]