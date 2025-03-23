# ==== Build ====
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

# ==== Runtime ====
FROM scratch

WORKDIR /

COPY --from=builder /app/server /server

# REST
EXPOSE 8080
# gRPC
EXPOSE 50051

ENTRYPOINT ["/server"]