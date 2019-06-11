FROM postgres:latest as base

ENV POSTGRES_DB 'trash_hunt'



FROM base as dev

ENV POSTGRES_USER 'dev'
ENV POSTGRES_DB 'dev_secret'

# Copies all seed files to image for initialization
COPY ./env/container/pg/seed_db/* /docker-entrypoint-initdb.d/
