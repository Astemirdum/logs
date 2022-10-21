SERVICE_NAME = log-app
ENV = .env

.PHONY: run
run:
	docker-compose -f ./docker-compose.yaml --env-file $(ENV) up

.PHONY: stop
stop:
	docker-compose -f ./docker-compose.yaml --env-file $(ENV) down

.PHONY: build-go
build-go:
	go build -o log-app cmd/main.go

.PHONY: rm
rm: build-go
	rm log-app

.PHONY: run-go
run-go:
	go run cmd/main.go -c "configs/config.yml"


.PHONY: lint
lint:
	golangci-lint run --fix

.PHONY: test
test:
	 go test -v ./...

.PHONY: mocks
mocks:
	cd ./internal/handler/; go generate;

.PHONY: migrate
migrate:
	goose -dir ./migrations  \
      postgres "user=postgres password=postgres host=localhost port=5432 database=postgres sslmode=disable" \
      up

.PHONY: container
container: build
    docker build -t $(CONTAINER_IMAGE):$(RELEASE) .

.PHONY: push
push: container
    docker push $(CONTAINER_IMAGE):$(RELEASE)


