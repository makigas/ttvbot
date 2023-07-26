FROM golang:1.20
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
ADD . .
RUN go build ./cmd/ttvbotd
ENTRYPOINT ["/app/ttvbotd"]