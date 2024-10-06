# Этап сборки
FROM golang:1.22.4 AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./

# Устанавливаем зависимости
RUN go mod download

# Копируем весь исходный код проекта
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bot ./main.go

# Этап выполнения
FROM alpine:latest

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /root/

# Копируем бинарный файл из предыдущего этапа
COPY --from=builder /app/bot .

# Копируем конфигурационный файл, если необходимо
COPY .env ./

# Запускаем приложение
CMD ["./bot"]
