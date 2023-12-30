FROM golang:latest-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o /bin/nats-cache ./cmd/nats-cache

FROM alpine as certificates
RUN apk add --no-cache ca-certificates

FROM scratch
COPY --from=certificates /etc/ssl/certs/ /etc/ssl/certs/
COPY --from=builder /bin/nats-cache /bin/nats-cache
ENTRYPOINT ["/bin/nats-cache"]
