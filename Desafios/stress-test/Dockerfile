# ---- build stage ----
FROM golang:1.22-alpine AS builder

WORKDIR /src

# Cache module downloads (no external deps, but keeps layers clean).
COPY go.mod ./
RUN go mod download

COPY . .

# Static binary so it runs on scratch.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /stress-test .

# ---- runtime stage ----
FROM scratch

# CA certificates so HTTPS targets work.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /stress-test /stress-test

ENTRYPOINT ["/stress-test"]
