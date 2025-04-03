# Сборка приложения
FROM golang:1.23 AS builder

# Установка рабочей директории
WORKDIR /app

# Копирование зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN go build -o main ./cmd/main.go

# Итоговый образ
FROM alpine:latest

# Установка рабочей директории
WORKDIR /app

# Копирование бинарного файла из builder-стадии
COPY --from=builder /app/main /app/main

# Добавление прав на выполнение
RUN chmod +x /app/main

# Открытие порта
EXPOSE 8080

# Запуск приложения
CMD ["./main"]