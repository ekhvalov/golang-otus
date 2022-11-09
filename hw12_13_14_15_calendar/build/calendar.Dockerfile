# Собираем в базовом образе
FROM otus-golang/base:develop as build

COPY . .

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "${LDFLAGS}" \
        -o /opt/calendar/calendar-app cmd/calendar/*
RUN CGO_ENABLED=0 go build \
        -ldflags "${LDFLAGS}" \
        -o /opt/calendar/calendar-migrate cmd/migrate/*

# На выходе тонкий образ
FROM alpine:3.16

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="calendar"
LABEL MAINTAINERS="student@otus.ru"

ENV CALENDAR_HTTP_ADDRESS="0.0.0.0"
ENV CALENDAR_HTTP_PORT="8080"
ENV CALENDAR_GRPC_ADDRESS="0.0.0.0"
ENV CALENDAR_GRPC_PORT="8081"
EXPOSE 8080
EXPOSE 8081

COPY --from=build /opt/calendar/calendar-app /usr/local/bin/calendar
COPY --from=build /opt/calendar/calendar-migrate /usr/local/bin/calendar-migrate
COPY ./configs/config.toml /etc/calendar/config.toml

CMD [ "calendar", "--config", "/etc/calendar/config.toml" ]
