# Build the manager binary
FROM golang:1.16 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY cmd/ cmd/
COPY webhook/ webhook/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o storageclass-accessor ./cmd/main.go

ENTRYPOINT ["/storageclass-accessor"]