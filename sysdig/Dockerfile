#upstream https://github.com/draios/sysdig/tree/dev/docker/stable
FROM debian:unstable

MAINTAINER 若虚 <slpcat@qq.com>

# Container variables
ENV \ 
    TERM="xterm" \
    LANG="en_US.UTF-8" \ 
    LANGUAGE="en_US:en" \
    LC_ALL="en_US.UTF-8" \
    TIMEZONE="Asia/Shanghai"

RUN sed -i 's/deb.debian.org/mirrors.aliyun.com/' /etc/apt/sources.list \
    && sed -i 's/security.debian.org/mirrors.aliyun.com\/debian-security/' /etc/apt/sources.list

# Set timezone and locales
RUN \
    echo "${TIMEZONE}" > /etc/timezone \
    && echo "$LANG UTF-8" > /etc/locale.gen \
    && apt-get update -q \
    && DEBIAN_FRONTEND=noninteractive apt-get install -yq apt-utils dialog vim-tiny locales \ 
    && ln -sf /usr/share/zoneinfo/${TIMEZONE} /etc/localtime \
    && update-locale LANG=$LANG \
    && locale-gen $LANG \
    && DEBIAN_FRONTEND=noninteractive dpkg-reconfigure locales

# Install required packages
RUN \
    apt-get dist-upgrade -y

ENV SYSDIG_REPOSITORY stable

LABEL RUN="docker run -i -t --privileged -v /var/run/docker.sock:/host/var/run/docker.sock -v /dev:/host/dev -v /proc:/host/proc:ro -v /boot:/host/boot:ro -v /lib/modules:/host/lib/modules:ro -v /usr:/host/usr:ro --name NAME IMAGE"

ENV SYSDIG_HOST_ROOT /host

ENV HOME /root

RUN cp /etc/skel/.bashrc /root && cp /etc/skel/.profile /root

ADD http://download.draios.com/apt-draios-priority /etc/apt/preferences.d/

RUN echo "deb http://httpredir.debian.org/debian jessie main" > /etc/apt/sources.list.d/jessie.list \
 && apt-get update \
 && apt-get install -y --no-install-recommends \
	bash-completion \
	curl \
	gnupg2 \
	ca-certificates \
	gcc \
	gcc-5 \
	gcc-4.9 \
	libelf1 \
	less \
	procps \
 && rm -rf /var/lib/apt/lists/*

# Since our base Debian image ships with GCC 5.0 which breaks older kernels, revert the
# default to gcc-4.9. Also, since some customers use some very old distributions whose kernel
# makefile is hardcoded for gcc-4.6 or so (e.g. Debian Wheezy), we pretend to have gcc 4.6/4.7
# by symlinking it to 4.9

RUN rm -rf /usr/bin/gcc \
 && ln -s /usr/bin/gcc-4.9 /usr/bin/gcc \
 && ln -s /usr/bin/gcc-4.9 /usr/bin/gcc-6 \
 && ln -s /usr/bin/gcc-4.9 /usr/bin/gcc-4.8 \
 && ln -s /usr/bin/gcc-4.9 /usr/bin/gcc-4.7 \
 && ln -s /usr/bin/gcc-4.9 /usr/bin/gcc-4.6

RUN curl -s https://s3.amazonaws.com/download.draios.com/DRAIOS-GPG-KEY.public | apt-key add - \
 && curl -s -o /etc/apt/sources.list.d/draios.list http://download.draios.com/$SYSDIG_REPOSITORY/deb/draios.list \
 && apt-get update \
 && apt-get install -y --no-install-recommends sysdig \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*

RUN ln -s $SYSDIG_HOST_ROOT/lib/modules /lib/modules

COPY ./docker-entrypoint.sh /

ENTRYPOINT ["/docker-entrypoint.sh"]

CMD ["bash"]
