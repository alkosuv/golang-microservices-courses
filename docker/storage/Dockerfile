# Используем официальный образ Go как базовый
FROM golang:1.23-alpine AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /storage

# Копируем файлы go.mod и go.sum в рабочую директорию
COPY go.mod go.sum ./

# Скачиваем все зависимости
RUN go mod download

# Копируем остальные файлы проекта в рабочую директорию
COPY . .

# Собираем бинарный файл
RUN go build -o storage ./cmd/storage/main.go

# Указываем команду для запуска приложения
CMD ./storage -config=./config/config.docker.yml

