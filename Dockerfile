# Этап сборки
FROM golang:1.23 AS builder

WORKDIR /app

# Сначала копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Затем копируем весь остальной код
COPY . .

# Собираем статический бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/main ./cmd/main.go

# Финальный образ
FROM alpine:3.19

WORKDIR /app

# Копируем бинарник из этапа сборки
COPY --from=builder /app/main /app/main

# Устанавливаем зависимости для Alpine
RUN apk --no-cache add ca-certificates

# Делаем бинарник исполняемым
RUN chmod +x /app/main

# Открываем порт
EXPOSE 8080

# Запускаем приложение
ENTRYPOINT ["/app/main"]