FROM golang:1.21.0 AS builder
WORKDIR /usr/local/src
COPY ./ ./

EXPOSE 8080

RUN go mod download
RUN go build ./cmd/server/main.go

CMD ["./main"]
