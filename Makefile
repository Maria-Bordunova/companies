.SILENT:

OPENAPI_FILE=./openapi.yaml
TARGET:='development'

build:
	go mod download && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./.bin/app ./cmd/main.go

run: build
	docker-compose up --remove-orphans --build app

down:
	docker-compose down --volumes --remove-orphans

test:
	docker-compose exec -T app sh -c "go test -v -race ./..."

client-go:
	oapi-codegen -package oapi openapi.yaml > pkg/gen/oapi/oapi.go

generate:
	docker-compose exec app sh -c "CGO_ENABLED=0 go generate ./..."
