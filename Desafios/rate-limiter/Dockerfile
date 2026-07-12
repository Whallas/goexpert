# --- build stage ---
FROM golang:1.22-alpine AS build

WORKDIR /app

# Cache deps separately from source for faster rebuilds.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Static binary so the final image needs no Go toolchain or libc.
RUN CGO_ENABLED=0 GOOS=linux go build -o /rate-limiter ./cmd/server

# --- run stage ---
FROM alpine:3.20

WORKDIR /app
COPY --from=build /rate-limiter /rate-limiter

# Create .env from the example template. The app loads it on startup via
# godotenv; values are overridden by any env vars set in docker-compose.
COPY .env.example ./.env.example
RUN cp .env.example .env

EXPOSE 8080
ENTRYPOINT ["/rate-limiter"]
