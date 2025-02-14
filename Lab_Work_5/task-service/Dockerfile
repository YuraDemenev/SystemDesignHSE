# Используем официальный образ Go
FROM golang:1.22 AS build-stage

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum для загрузки зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Компилируем бинарный файл
RUN go build -o /task main.go

# Новый минимальный образ
FROM debian:bookworm  

# Устанавливаем рабочую директорию
WORKDIR /

COPY --from=build-stage /task /task

EXPOSE 8001

# Указываем команду для запуска сервиса
ENTRYPOINT ["/task"]
