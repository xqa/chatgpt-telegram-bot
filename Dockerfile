FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o chatgpt-telegram-bot .


FROM alpine
WORKDIR /app
COPY --from=builder /app/ .

CMD ["/app/chatgpt-telegram-bot"]