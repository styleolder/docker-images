FROM  maven:3.5.3
MAINTAINER 若虚 <slpcat@qq.com>

# Container variables
ENV \
    TERM="xterm" \
    LANG="en_US.UTF-8" \
    LANGUAGE="en_US:en" \
    LC_ALL="en_US.UTF-8" \
    TIMEZONE="Asia/Shanghai"

RUN \
    sed -i 's/deb.debian.org/mirrors.aliyun.com/' /etc/apt/sources.list \
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
    apt-get upgrade -y

ENV JAVA_OPTS="-Drocketmq.namesrv.addr=127.0.0.1:9876 -Dcom.rocketmq.sendMessageWithVIPChannel=false"

RUN groupadd rocketmq && \
    useradd -g rocketmq -s /bin/bash -c RocketMQ rocketmq  && \
    mkdir /home/rocketmq && \
    chown -R rocketmq:rocketmq /home/rocketmq && \
    mkdir -p /opt/rocketmq && \
    chown -R rocketmq:rocketmq /opt/rocketmq

# install from source
USER rocketmq
RUN cd /opt/rocketmq && \
    git clone https://github.com/apache/rocketmq-externals.git

WORKDIR /opt/rocketmq/rocketmq-externals/rocketmq-console
RUN \
    mvn clean package -Dmaven.test.skip=true && \
    mv target/rocketmq-console-ng-1.0.0.jar /opt/rocketmq && \
    rm -rf /opt/rocketmq-externals && \
    rm -rf ~/.m2 

WORKDIR /opt/rocketmq
EXPOSE 8080
CMD  ["sh" "-c" "java" "$JAVA_OPTS" "-jar" "rocketmq-console-ng-1.0.0.jar"]
