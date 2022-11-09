BIN := "./bin/calendar"
DOCKER_IMG_BASE="otus-golang/base:develop"
DOCKER_IMG_CALENDAR="otus-golang/calendar:develop"
DOCKER_IMG_SCHEDULER="otus-golang/calendar-scheduler:develop"
DOCKER_IMG_SENDER="otus-golang/calendar-sender:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o ./bin/calendar_scheduler -ldflags "$(LDFLAGS)" ./cmd/scheduler
	go build -v -o ./bin/calendar_sender -ldflags "$(LDFLAGS)" ./cmd/sender
	go build -v -o ./bin/migrate -ldflags "$(LDFLAGS)" ./cmd/migrate

run: build
	$(BIN) --config ./configs/config.toml

build-img-base:
	docker build \
		--tag $(DOCKER_IMG_BASE) \
		--file build/base.Dockerfile .

build-img-calendar: build-img-base
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		--tag $(DOCKER_IMG_CALENDAR) \
		--force-rm=true \
		--file build/calendar.Dockerfile .

build-img-scheduler: build-img-base
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		--tag $(DOCKER_IMG_SCHEDULER) \
		--force-rm=true \
		--file build/scheduler.Dockerfile .

build-img-sender: build-img-base
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		--tag $(DOCKER_IMG_SENDER) \
		--force-rm=true \
		--file build/sender.Dockerfile .

build-img: build-img-calendar build-img-scheduler build-img-sender

run-img: build-img
	docker run $(DOCKER_IMG_CALENDAR)

version: build
	$(BIN) version

test:
	go test -race ./internal/... ./pkg/...

test-integration: build-img
	cd deployments/integration-tests && docker-compose up --exit-code-from tester && docker-compose down -v

test-full: lint build build-img test-integration

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.50.0

lint: install-lint-deps
	golangci-lint run ./...

generate:
	go generate github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/app
	go generate github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/app/event/command
	go generate github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/app/notification/queue
	go generate github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/domain/event
	go generate github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/server/grpc
	go generate github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/server/http

up: build-img
	cd deployments && docker-compose up -d

down:
	cd deployments && docker-compose down --remove-orphans

.PHONY: build build-img build-img-base build-img-calendar build-img-scheduler build-img-sender down generate lint run \
	run-img test test-full test-integration up version
