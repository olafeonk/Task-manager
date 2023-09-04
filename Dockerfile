FROM golang:1.21.0-alpine AS builer
WORKDIR /usr/local/src
COPY ./ ./


RUN go mod download
RUN go build ./cmd/main.go

FROM alpine
COPY --from=builer /usr/local/src/main /
CMD ["./main"]