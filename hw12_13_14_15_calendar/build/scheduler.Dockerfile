# Собираем в базовом образе
FROM otus-golang/base:develop as build

COPY . .

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "${LDFLAGS}" \
        -o /opt/calendar/calendar-scheduler cmd/scheduler/*

# На выходе тонкий образ
FROM alpine:3.16

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="calendar_scheduler"
LABEL MAINTAINERS="student@otus.ru"

COPY --from=build /opt/calendar/calendar-scheduler /usr/local/bin/calendar-scheduler
COPY ./configs/scheduler_config.toml /etc/calendar/config.toml

CMD [ "calendar-scheduler", "--config", "/etc/calendar/config.toml" ]
