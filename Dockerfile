FROM golang:1.23

WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd/
COPY configs ./configs/
COPY internal ./internal/
COPY pkg ./pkg/

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /my-app

# Run
CMD ["/ccli"]
