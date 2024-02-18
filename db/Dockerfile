FROM postgres:latest
LABEL authors="buguzei"

ENV POSTGRES_PASSWORD=123456
ENV POSTGRES_USER=buguzei
ENV POSTGRES_DB=sellerbot

COPY 20240124173955_psql_migrations.sql /docker-entrypoint-initdb.d/