#FROM  maven:3.5.3
FROM java:8

# Rocketmq version
ENV ROCKETMQ_VERSION="4.1.0"
# Rocketmq home
ENV ROCKETMQ_HOME="/opt/rocketmq"

ENV JAVA_OPT=" -Duser.home=/opt/rocketmq"

RUN groupadd rocketmq && \
    useradd -g rocketmq -s /bin/bash -c RocketMQ rocketmq  && \
    mkdir /home/rocketmq && \
    chown -R rocketmq:rocketmq /home/rocketmq/ && \
    mkdir -p /opt/rocketmq && \

# install from source
#RUN cd /opt && \
#    git clone https://github.com/apache/incubator-rocketmq.git
#WORKDIR /home/admin/incubator-rocketmq
#    mvn -Prelease-all -DskipTests clean install -U
#cd distribution/target/apache-rocketmq
# install from binary
RUN cd /opt && \
    wget https://mirrors.tuna.tsinghua.edu.cn/apache/rocketmq/${ROCKETMQ_VERSION}-incubating/rocketmq-all-${ROCKETMQ_VERSION}-incubating-bin-release.zip && \
    unzip rocketmq-all-${ROCKETMQ_VERSION}-incubating-bin-release.zip && \
    rm -rf rocketmq-all-${ROCKETMQ_VERSION}-incubating-bin-release.zip && \
    mv rocketmq-all-${ROCKETMQ_VERSION}-incubating rocketmq && \
    chown -R rocketmq:rocketmq /opt/rocketmq

USER rocketmq
WORKDIR /opt/rocketmq
EXPOSE 9876 10909 10911
CMD  ["sh", "-c", ". ./bin/play.sh; while sleep 10; do echo RocketMQ, GO ROCK; done"]
