# Собираем в базовом образе
FROM otus-golang/base:develop as build
LABEL otusgolang-cache="true"

COPY . .

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "${LDFLAGS}" \
        -o /opt/calendar/calendar-sender cmd/sender/*

# На выходе тонкий образ
FROM alpine:3.16

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="calendar_sender"
LABEL MAINTAINERS="student@otus.ru"

COPY --from=build /opt/calendar/calendar-sender /usr/local/bin/calendar-sender
COPY ./configs/sender_config.toml /etc/calendar/config.toml

CMD [ "calendar-sender", "--config", "/etc/calendar/config.toml" ]
