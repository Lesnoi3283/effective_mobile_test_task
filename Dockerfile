FROM golang:1.23.1-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o musiclib_bin ./cmd/musiclib/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/musiclib_bin .

CMD ["./musiclib_bin"]