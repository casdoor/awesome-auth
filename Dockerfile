FROM ghcr.io/hanzovm/vmd:1.5.4 as vmd
FROM node:18.19.0 AS FRONT
WORKDIR /web
COPY ./web .
RUN yarn install --frozen-lockfile --network-timeout 1000000 && yarn run build


FROM golang:1.21 AS BACK
WORKDIR /go/src/hanzo-vm
COPY . .
RUN chmod +x ./build.sh
RUN ./build.sh


FROM alpine:latest AS STANDARD
LABEL MAINTAINER="https://github.com/hanzovm/vm"
ARG USER=vm

RUN sed -i 's/https/http/' /etc/apk/repositories
RUN apk add --update sudo
RUN apk add curl
RUN apk add ca-certificates && update-ca-certificates

RUN adduser -D $USER -u 1000 \
    && echo "$USER ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/$USER \
    && chmod 0440 /etc/sudoers.d/$USER \
    && mkdir logs \
    && chown -R $USER:$USER logs

USER 1000
WORKDIR /
COPY --from=BACK --chown=$USER:$USER /go/src/hanzo-vm/server ./server
COPY --from=BACK --chown=$USER:$USER /go/src/hanzo-vm/data ./data
COPY --from=BACK --chown=$USER:$USER /go/src/hanzo-vm/conf/app.conf ./conf/app.conf
COPY --from=FRONT --chown=$USER:$USER /web/build ./web/build

ENTRYPOINT ["/server"]


FROM vmd AS ALLINONE
LABEL MAINTAINER="https://github.com/hanzovm/vm"

WORKDIR /

USER root
RUN apt-get update \
    && apt-get install -y      \
        mariadb-server         \
        mariadb-client         \
        ca-certificates        \
    && update-ca-certificates  \
    && rm -rf /var/lib/apt/lists/*

COPY --from=BACK /go/src/hanzo-vm/server ./server
COPY --from=BACK /go/src/hanzo-vm/data ./data
COPY --from=BACK /go/src/hanzo-vm/docker-entrypoint.sh /docker-entrypoint.sh
COPY --from=BACK /go/src/hanzo-vm/conf/app.conf ./conf/app.conf
COPY --from=FRONT /web/build ./web/build

EXPOSE 19000
ENTRYPOINT ["/bin/bash"]
CMD ["/docker-entrypoint.sh"]
