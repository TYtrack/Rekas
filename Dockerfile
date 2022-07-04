FROM mirrors.tencent.com/tlinux/tlinux3.2:latest

USER root

ENV APP_HOME /app


WORKDIR $APP_HOME

COPY ./config.toml ./config.toml
COPY ./server ./server

CMD  ./server