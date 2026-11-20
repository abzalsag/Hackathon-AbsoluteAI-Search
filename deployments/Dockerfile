FROM golang:1.26-alpine AS builder

WORKDIR /app

# 1. Тек go.mod және go.sum көшіреміз
COPY go.mod go.sum ./

# 2. Пакеттерді жүктейміз (бұл файлдарды өзгертпейді, тек жүктейді)
RUN go mod download

# 3. Кодты көшіріп, жинаймыз
COPY . .
RUN go build -o main cmd/api/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/db ./db

EXPOSE 8080

CMD ["./main"]