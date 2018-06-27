FROM centos:7

MAINTAINER 若虚 <slpcat@qq.com>

# Container variables
ENV \
    TERM="xterm" \
    LANG="en_US.UTF-8" \
    LANGUAGE="en_US.UTF-8" \
    TIMEZONE="Asia/Shanghai"

COPY epel.repo /etc/yum.repos.d/
COPY java.sh /etc/profile.d/

# set timezone
RUN ln -sf /usr/share/zoneinfo/${TIMEZONE} /etc/localtime

RUN \
    yum update -y && \
    yum install -y wget && \
    wget --no-check-certificate --no-cookies --header "Cookie: oraclelicense=accept-securebackup-cookie" http://download.oracle.com/otn-pub/java/jdk/8u171-b11/512cd62ec5174c3487ac17c61aaa89e8/jdk-8u171-linux-x64.rpm && \
    rpm -ivh jdk-8u171-linux-x64.rpm && \
    rm jdk-8u171-linux-x64.rpm

ENV HBASE_VERSION 1.2.6.1
ENV JAVA_HOME /usr/java/jdk1.8.0_171
ENV HBASE_INSTALL_DIR /opt/hbase

RUN mkdir -p ${HBASE_INSTALL_DIR} && \
    curl -L http://mirrors.aliyun.com/apache/hbase//${HBASE_VERSION}/hbase-${HBASE_VERSION}-bin.tar.gz | tar -xz --strip=1 -C ${HBASE_INSTALL_DIR}

ADD hbase-site.xml /opt/hbase/conf/hbase-site.xml
ADD start-k8s-hbase.sh /opt/hbase/bin/start-k8s-hbase.sh
RUN chmod +x /opt/hbase/bin/start-k8s-hbase.sh

WORKDIR ${HBASE_INSTALL_DIR}
RUN echo "export HBASE_JMX_BASE=\"-Dcom.sun.management.jmxremote.ssl=false -Dcom.sun.management.jmxremote.authenticate=false\"" >> conf/hbase-env.sh && \
    echo "export HBASE_MASTER_OPTS=\"\$HBASE_MASTER_OPTS \$HBASE_JMX_BASE -Dcom.sun.management.jmxremote.port=10101\"" >> conf/hbase-env.sh && \
    echo "export HBASE_REGIONSERVER_OPTS=\"\$HBASE_REGIONSERVER_OPTS \$HBASE_JMX_BASE -Dcom.sun.management.jmxremote.port=10102\"" >> conf/hbase-env.sh && \
    echo "export HBASE_THRIFT_OPTS=\"\$HBASE_THRIFT_OPTS \$HBASE_JMX_BASE -Dcom.sun.management.jmxremote.port=10103\"" >> conf/hbase-env.sh && \
    echo "export HBASE_ZOOKEEPER_OPTS=\"\$HBASE_ZOOKEEPER_OPTS \$HBASE_JMX_BASE -Dcom.sun.management.jmxremote.port=10104\"" >> conf/hbase-env.sh && \
    echo "export HBASE_REST_OPTS=\"\$HBASE_REST_OPTS \$HBASE_JMX_BASE -Dcom.sun.management.jmxremote.port=10105\"" >> conf/hbase-env.sh

ENV PATH=$PATH:/opt/hbase/bin

CMD /opt/hbase/bin/start-k8s-hbase.sh
