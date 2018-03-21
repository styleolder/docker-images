# Build a minimal distribution container
# upstream https://github.com/gliderlabs/docker-alpine
FROM alpine:3.7
MAINTAINER 若虚 <slpcat@qq.com>

# Container variables
ENV \
    TERM="xterm" \
    LANG="en_US.UTF-8" \ 
    LANGUAGE="en_US.UTF-8" \
    LC_ALL="en_US.UTF-8" \
    TIMEZONE="Asia/Shanghai"

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
    echo 'hosts: files mdns4_minimal [NOTFOUND=return] dns mdns4' >> /etc/nsswitch.conf \
    && rm -f /var/cache/apk/*.tar.gz

CMD ["/bin/bash"]
