include local.env

LOCAL_BIN:=$(CURDIR)/bin

LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=$(POSTGRES_HOST) port=$(POSTGRES_PORT) dbname=$(POSTGRES_DBNAME) user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) sslmode=disable"

install-deps:
# Скачать бинарники локально (фиксация версий обязательна)
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.18.0

generate:
	make generate-note-api

generate-note-api:
	mkdir -p pkg/note/v1
	protoc --proto_path api/note/v1 \
	--go_out=pkg/note/v1 \
	--go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/note/v1 \
	--go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/note/v1/note.proto

# Собрать приложение под конкретную архитектуру и ОС
build:
	GOOS=linux GOARCH=amd64 go build -o ./bin/ ./cmd/server/main.go

# Скопировать бинарный файл приложения на сервер
copy-binary-to-server:
	scp ./bin/main root@5.188.50.4:

docker-build-and-push:
	docker buildx build --no-cache --platform linux/amd64 -t 5.188.50.4:5000/test-server:v0.0.1 .
	docker push 5.188.50.4:5000/test-server:v0.0.1

local-migration-status:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migration-down:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v
