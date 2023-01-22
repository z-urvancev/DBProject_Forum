FROM golang:1.17 AS builder

WORKDIR /server
COPY . .
RUN go mod tidy
RUN go build -o main cmd/main.go

FROM ubuntu:20.04

RUN apt-get -y update && apt-get install -y tzdata
ENV TZ=Russia/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENV PGVER 12
RUN apt-get -y update && apt-get install -y postgresql-$PGVER
USER postgres

RUN /etc/init.d/postgresql start &&\
    psql --command "create user forum with superuser password 'forum';" &&\
    createdb -O forum forum &&\
    /etc/init.d/postgresql stop

EXPOSE 5432
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]
USER root

WORKDIR /server
COPY . .
COPY --from=builder /server/main /server/main

EXPOSE 5000
ENV PGPASSWORD forum
CMD service postgresql start && psql -h localhost -d forum -U forum -p 5432 -a -q -f ./db/db.sql && ./main
