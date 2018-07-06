FROM alpine:3.7
MAINTAINER 若虚 <slpcat@qq.com>

# Container variables
ENV \
    TERM="xterm" \
    LANG="en_US.UTF-8" \
    LANGUAGE="en_US.UTF-8" \
    LC_ALL="en_US.UTF-8" \
    TIMEZONE="Asia/Shanghai"

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/' /etc/apk/repositories

# Set timezone and locales
RUN set -ex \
    && apk update \
    && apk upgrade \
    && apk add \
           bash \
           tzdata \
           vim \
           tini \
           su-exec \
           gzip \
           tar \
           wget \
           curl \
    && echo "${TIMEZONE}" > /etc/TZ \
    && ln -sf /usr/share/zoneinfo/${TIMEZONE} /etc/localtime \
    # Network fix
    && echo 'hosts: files mdns4_minimal [NOTFOUND=return] dns mdns4' >> /etc/nsswitch.conf

ARG VERSION=${VERSION:-master}

RUN apk add --no-cache --virtual .build-deps \
      curl gcc g++ make cmake autoconf libc-dev libevent-dev linux-headers perl tar \
    && mkdir -p /qconf/tmp \
    && curl -Lk "https://github.com/Qihoo360/QConf/archive/${VERSION}.tar.gz" | \
       tar -xz -C /qconf/tmp --strip-components=1 \
    && cd /qconf/tmp \
    && mkdir build \
    && cd build \
    && cmake .. -DCMAKE_INSTALL_PREFIX=/qconf \
    && make -j$(getconf _NPROCESSORS_ONLN) \
    && make install PREFIX=/qconf \
    && rm -rf /qconf/tmp \
    && apk add --virtual .rundeps libstdc++ \
    && apk del .build-deps

EXPOSE 8888
#VOLUME /ssdb/var
WORKDIR /qconf

CMD ["bin/agent-cmd.sh", "start"]
