FROM golang:1.22-bookworm AS builder

RUN apt-get update && apt-get install -y protobuf-compiler && rm -rf /var/lib/apt/lists/*

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1

WORKDIR /app

COPY . .

RUN protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    internal/infra/grpc/pb/order.proto

RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/ordersystem/

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/server .
COPY migrations ./migrations

EXPOSE 8000 50051 8080

CMD ["./server"]
