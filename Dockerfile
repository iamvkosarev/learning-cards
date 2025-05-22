# ==== Build ====
FROM golang:1.23 AS builder

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o cards ./cmd/cards

# ==== Runtime ====
FROM scratch

WORKDIR /

COPY --from=builder /app/cards /cards

EXPOSE ${GRPC_PORT} ${REST_PORT}

ENTRYPOINT ["/cards"]