FROM golang:1.21.0-alpine AS builder
WORKDIR /usr/local/src
COPY ./ ./


RUN go mod download
RUN go build ./cmd/server/main.go
RUN go build -o createuser ./cmd/command/main.go

FROM alpine
COPY --from=builder /usr/local/src/main /
COPY --from=builder /usr/local/src/schema /schema
COPY --from=builder /usr/local/src/.env /
COPY --from=builder /usr/local/src/createuser /
CMD ["./main"]