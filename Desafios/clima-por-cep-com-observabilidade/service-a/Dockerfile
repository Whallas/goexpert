FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /service-a ./service-a

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /service-a /service-a
EXPOSE 8080
ENTRYPOINT ["/service-a"]
