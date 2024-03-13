include local.env

LOCAL_BIN:=$(CURDIR)/bin

LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=$(POSTGRES_HOST) port=$(POSTGRES_PORT) dbname=$(POSTGRES_DBNAME) user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) sslmode=disable"

# Скачать бинарники локально (фиксация версий обязательна)
install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.18.0
	GOBIN=$(LOCAL_BIN) go install github.com/gojuno/minimock/v3/cmd/minimock@v3.3.6
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.19.1
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.19.1
	GOBIN=$(LOCAL_BIN) go install github.com/rakyll/statik@v0.1.7
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v0.10.1

generate:
	mkdir -p pkg/swagger
	make generate-note-api
	$(LOCAL_BIN)/statik -f -src=pkg/swagger/ -include='*.json,*.png,*.html,*.css,*.js'

# Сгенерировать код из proto-файлов
generate-note-api:
	mkdir -p pkg/note/v1
	protoc --proto_path api/note/v1 --proto_path vendor.protogen \
	--go_out=pkg/note/v1 \
	--go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/note/v1 \
	--go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	--grpc-gateway_out=pkg/note/v1 \
	--grpc-gateway_opt=paths=source_relative \
	--plugin=protoc-gen-grpc-gateway=bin/protoc-gen-grpc-gateway \
	--openapiv2_out=allow_merge=true,merge_file_name=api:pkg/swagger \
	--plugin=protoc-gen-openapiv2=bin/protoc-gen-openapiv2 \
	--validate_out lang=go:pkg/note/v1 \
	--validate_opt=paths=source_relative \
	--plugin=protoc-gen-validate=bin/protoc-gen-validate \
	api/note/v1/note.proto

# Собрать приложение под конкретную архитектуру и ОС
build:
	GOOS=linux GOARCH=amd64 go build -o ./bin/ ./cmd/server/main.go

# Скопировать бинарный файл приложения на сервер
copy-binary-to-server:
	scp ./bin/main root@5.188.50.4:

# Собрать докера образ и запушить его на тестовый сервер
docker-build-and-push:
	docker buildx build --no-cache --platform linux/amd64 -t 5.188.50.4:5000/test-server:v0.0.1 .
	docker push 5.188.50.4:5000/test-server:v0.0.1

# Манипулировать миграциями бд
local-migration-status:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migration-down:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

# Сгенерировать моки для тестов
generate-mocks:
	make generate-repo-mocks generate-service-mocks

generate-repo-mocks:
	go generate ./internal/repository

generate-service-mocks:
	go generate ./internal/service

# Запустить тесты с просмотром покрытия
test:
	go clean -testcache
	go test ./... -covermode count -coverpkg=github.com/igorakimy/bigtech_microservices/internal/service/...,github.com/igorakimy/bigtech_microservices/internal/api/... -count 5

test-coverage:
	go clean -testcache
	go test ./... -covermode count -coverpkg=github.com/igorakimy/bigtech_microservices/internal/service/...,github.com/igorakimy/bigtech_microservices/internal/api/... -count 5
	grep -v 'mocks\|config' coverage.tmp.out > coverage.out
	rm coverage.tmp.out
	go tool cover -html=coverage.out;
	go tool cover -func=./coverage.out | grep "total";
	grep -sqFx "/coverage.out" .gitignore || echo "/coverage.out" >> .gitignore

# Скачать недостающие сторонние proto-файлы
vendor-proto:
	@if [ ! -d vendor.protogen/validate ]; then \
  		mkdir -p vendor.protogen/validate &&\
  		git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/protoc-gen-validate &&\
  		mv vendor.protogen/protoc-gen-validate/validate/*.proto vendor.protogen/validate && \
  		rm -rf vendor.protogen/protoc-gen-validate ;\
  	fi
	@if [ ! -d vendor.protogen/google ]; then \
  		git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
  		mkdir -p vendor.protogen/google/ &&\
  		mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
  		rm -rf vendor.protogen/googleapis ;\
  	fi
	@if [ ! -d vendor.protogen/protoc-gen-openapiv2 ]; then \
		mkdir -p vendor.protogen/protoc-gen-openapiv2/options &&\
		git clone https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/openapiv2 &&\
		mv vendor.protogen/openapiv2/protoc-gen-openapiv2/options/*.proto vendor.protogen/protoc-gen-openapiv2/options &&\
		rm -rf vendor.protogen/openapiv2 ;\
	fi