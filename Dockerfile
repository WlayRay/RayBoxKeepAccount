FROM golang:latest AS builder

WORKDIR /app

COPY . .

RUN go mod download && \
    go mod tidy && \
    go build -o main . \
    chmod +x main

FROM alpine:latest

WORKDIR /ray

COPY --from=builder /app/main .

EXPOSE 8000

CMD ["./main"]

