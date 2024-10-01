# Используем базовый образ Go
FROM golang:1.18-alpine

WORKDIR /app

# Копируем все файлы в контейнер
COPY . .

# Устанавливаем зависимости
RUN go mod tidy

# Команда по умолчанию для запуска приложения
CMD ["go", "run", "main.go"]
