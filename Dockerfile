ARG ALPINE_BASE_IMAGE=3.18

FROM guacamole/guacd:1.5.4 as guacd
FROM node:18.19.0 AS FRONT
WORKDIR /web
COPY ./web .
RUN yarn install --frozen-lockfile --network-timeout 1000000 && yarn run build


FROM golang:1.20.12 AS BACK
WORKDIR /go/src/casvisor
COPY . .
RUN chmod +x ./build.sh
RUN ./build.sh


FROM alpine:${ALPINE_BASE_IMAGE} AS STANDARD
LABEL MAINTAINER="https://casvisor.org/"
ARG USER=casvisor

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
COPY --from=BACK --chown=$USER:$USER /go/src/casvisor/server ./server
COPY --from=BACK --chown=$USER:$USER /go/src/casvisor/data ./data
COPY --from=BACK --chown=$USER:$USER /go/src/casvisor/conf/app.conf ./conf/app.conf
COPY --from=FRONT --chown=$USER:$USER /web/build ./web/build

ENTRYPOINT ["/server"]


FROM alpine:${ALPINE_BASE_IMAGE} AS ALLINONE
LABEL MAINTAINER="https://casvisor.org/"

WORKDIR /
ARG PREFIX_DIR=/opt/guacamole

ENV LD_LIBRARY_PATH=${PREFIX_DIR}/lib

COPY --from=guacd ${PREFIX_DIR} ${PREFIX_DIR}

# Bring runtime environment up to date and install runtime dependencies
RUN apk add --no-cache                \
        openrc                        \
        mariadb                       \
        mariadb-client                \
        ca-certificates               \
        font-noto-cjk                 \
        ghostscript                   \
        netcat-openbsd                \
        shadow                        \
        terminus-font                 \
        ttf-dejavu                    \
        ttf-liberation                \
        util-linux-login && \
    xargs apk add --no-cache < ${PREFIX_DIR}/DEPENDENCIES && \
    apk cache clean && \
    rm -rf /var/cache/apk/*

COPY --from=BACK /go/src/casvisor/server ./server
COPY --from=BACK /go/src/casvisor/data ./data
COPY --from=BACK /go/src/casvisor/docker-entrypoint.sh /docker-entrypoint.sh
COPY --from=BACK /go/src/casvisor/conf/app.conf ./conf/app.conf
COPY --from=FRONT /web/build ./web/build

ENTRYPOINT ["/bin/bash"]
CMD ["/docker-entrypoint.sh"]
