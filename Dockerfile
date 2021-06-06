FROM golang:1.13 AS build

ADD . /app
WORKDIR /app
RUN go build ./app/cmd/main.go

FROM ubuntu:20.04

RUN apt-get -y update && apt-get install -y tzdata

ENV PGVER 12
RUN apt-get -y update && apt-get install -y postgresql-$PGVER

USER postgres

RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER sergei WITH SUPERUSER PASSWORD '1111';" &&\
    createdb -E UTF8 -O sergei forums &&\
    /etc/init.d/postgresql stop

EXPOSE 5432

VOLUME ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

USER root

WORKDIR /usr/src/app

COPY . .
COPY --from=build /app/main .

EXPOSE 5000

ENV PGPASSWORD 1111

CMD service postgresql start && psql -h localhost -d forums -U sergei -p 5432 -a -q -f ./tables.sql && ./main