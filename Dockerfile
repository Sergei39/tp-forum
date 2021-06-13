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

RUN echo "synchronous_commit = off\nfsync = off\nshared_buffers = 256MB\n" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "wal_buffers = 2MB\nwal_writer_delay = 50ms\nrandom_page_cost = 1.0\nmax_connections = 100\nwork_mem = 10MB\nmaintenance_work_mem = 128MB\ncpu_tuple_cost = 0.0030\ncpu_index_tuple_cost = 0.0010\ncpu_operator_cost = 0.0005" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "full_page_writes = off" >> /etc/postgresql/$PGVER/main/postgresql.conf
# RUN echo "min_wal_size = 1GB\nmax_wal_size = 4GB\nmax_worker_processes = 2\nmax_parallel_workers_per_gather = 1\nmax_parallel_workers = 2\nmax_parallel_maintenance_workers = 1"

EXPOSE 5432

VOLUME ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

USER root

WORKDIR /usr/src/app

COPY . .
COPY --from=build /app/main .

EXPOSE 5000

ENV PGPASSWORD 1111

CMD service postgresql start && psql -h localhost -d forums -U sergei -p 5432 -a -q -f ./tabels.sql && ./main