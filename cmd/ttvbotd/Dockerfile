FROM golang:1.21-alpine AS builder
RUN apk update && apk add --no-cache git gcc g++
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
ADD . .
RUN GOOS=linux CGO_ENABLED=0 GOPROXY=direct go build -ldflags="-w -s" ./cmd/ttvbotd
ENTRYPOINT ["/app/ttvbotd"]

FROM scratch
COPY --from=builder /app/ttvbotd /ttvbotd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["/ttvbotd"]
