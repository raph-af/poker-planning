FROM postgres:12.2

ENV POSTGRES_USER postgres
ENV POSTGRES_PASSWORD postgres
ENV POSTGRES_DB postgres

COPY ./db /docker-entrypoint-initdb.d
