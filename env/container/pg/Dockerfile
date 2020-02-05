FROM postgres:latest as base

ENV POSTGRES_DB 'trash_hunt'



FROM base as dev

ENV POSTGRES_USER 'dev'
ENV POSTGRES_PASSWORD 'dev_secret'

RUN useradd -ms /bin/bash user
USER user
