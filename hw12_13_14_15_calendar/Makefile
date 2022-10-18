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

run: build
	$(BIN) --config ./configs/config.toml

.PHONY: build-img-base
build-img-base:
	docker build \
		--tag $(DOCKER_IMG_BASE) \
		--file build/base.Dockerfile .

.PHONY: build-img-calendar
build-img-calendar: build-img-base
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		--tag $(DOCKER_IMG_CALENDAR) \
		--file build/calendar.Dockerfile .
	docker image rm $$(docker image ls -q -f label=otusgolang-cache=true)

.PHONY: build-img-scheduler
build-img-scheduler: build-img-base
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		--tag $(DOCKER_IMG_SCHEDULER) \
		--file build/scheduler.Dockerfile .
	docker image rm $$(docker image ls -q -f label=otusgolang-cache=true)

.PHONY: build-img-sender
build-img-sender: build-img-base
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		--tag $(DOCKER_IMG_SENDER) \
		--file build/sender.Dockerfile .
	docker image rm $$(docker image ls -q -f label=otusgolang-cache=true)

build-img: build-img-calendar build-img-scheduler build-img-sender

run-img: build-img
	docker run $(DOCKER_IMG_CALENDAR)

version: build
	$(BIN) version

test:
	go test -race ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.50.0

lint: install-lint-deps
	golangci-lint run ./...

generate:
	go generate github.com/ekhvalov/hw12_13_14_15_calendar/internal/app
	go generate github.com/ekhvalov/hw12_13_14_15_calendar/internal/app/event/command
	go generate github.com/ekhvalov/hw12_13_14_15_calendar/internal/app/notification/queue
	go generate github.com/ekhvalov/hw12_13_14_15_calendar/internal/domain/event
	go generate github.com/ekhvalov/hw12_13_14_15_calendar/internal/server/http

up:
	cd deployments && docker-compose up -d

.PHONY: build run build-img run-img version test lint generate up
