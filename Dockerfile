FROM golang:alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o /bin/api-server ./cmd/api-server
RUN CGO_ENABLED=1 GOOS=linux go build -o /bin/tcp-server ./cmd/tcp-server
RUN CGO_ENABLED=1 GOOS=linux go build -o /bin/udp-server ./cmd/udp-server
RUN CGO_ENABLED=1 GOOS=linux go build -o /bin/grpc-server ./cmd/grpc-server

FROM alpine:latest
WORKDIR /app

COPY --from=builder /bin/api-server /app/
COPY --from=builder /bin/tcp-server /app/
COPY --from=builder /bin/udp-server /app/
COPY --from=builder /bin/grpc-server /app/

COPY ./data /app/data/
