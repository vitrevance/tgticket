FROM golang:1.25.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /tgticket ./cmd/tgticket

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

COPY --from=builder /tgticket /tgticket

ENTRYPOINT ["/tgticket"]
CMD ["-config", "/config.yaml"]
