# ==== Build ====
FROM golang:1.23 AS builder

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ../.. .

RUN go build -o cards ./cmd/cards
RUN go build -o reviews ./cmd/reviews

# ==== Runtime (cards) ====
FROM scratch AS cards

WORKDIR /

COPY --from=builder /app/cards /cards

EXPOSE ${CARDS_GRPC_PORT} ${CARDS_REST_PORT}
ENTRYPOINT ["/cards"]

# ==== Runtime (reviews) ====
FROM scratch AS reviews

WORKDIR /

COPY --from=builder /app/reviews /reviews

EXPOSE ${REVIEWS_GRPC_PORT} ${REVIEWS_REST_PORT}
ENTRYPOINT ["/reviews"]