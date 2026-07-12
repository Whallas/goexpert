FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /service-b ./service-b

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /service-b /service-b
EXPOSE 8081
ENTRYPOINT ["/service-b"]
