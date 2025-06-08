# ==== Build stage ====
FROM golang:1.24-bullseye AS build

RUN apt-get update && \
    apt-get install -y mecab libmecab-dev mecab-ipadic pkg-config && \
    rm -rf /var/lib/apt/lists/*

ENV CGO_ENABLED=1
ENV CGO_CFLAGS="-I/usr/include"
ENV CGO_LDFLAGS="-lmecab -lstdc++"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o cards ./cmd/cards

# ==== Runtime stage ====
FROM debian:bullseye AS cards

RUN apt-get update && \
    apt-get install -y mecab libmecab2 mecab-ipadic-utf8 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=build /app/cards /cards

ENV MECAB_DICDIR=/usr/lib/mecab/dic/ipadic-utf8

EXPOSE ${GRPC_PORT} ${REST_PORT}
ENTRYPOINT ["/cards"]