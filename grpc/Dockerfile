# Указание образа и название этапа сборки
FROM golang:1.21.6-alpine AS builder

# Скопировать содержимое в папку
COPY . /github.com/igorakimy/bigtech_microservices/grpc/source/

# Пометить папку, как рабочую
WORKDIR /github.com/igorakimy/bigtech_microservices/grpc/source/

# Запустить команды
RUN apk add git
RUN go mod tidy
RUN go build -o ./bin/crud cmd/server/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /github.com/igorakimy/bigtech_microservices/grpc/source/bin/crud .

# Выполнить команду
CMD ["./crud"]
