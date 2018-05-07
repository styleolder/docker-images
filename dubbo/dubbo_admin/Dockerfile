#upstream https://github.com/carlossg/docker-maven/tree/master/jdk-8
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

ARG MAVEN_VERSION=3.5.3
ARG USER_HOME_DIR="/root"
ARG SHA=b52956373fab1dd4277926507ab189fb797b3bc51a2a267a193c931fffad8408
#ARG BASE_URL=https://apache.osuosl.org/maven/maven-3/${MAVEN_VERSION}/binaries
ARG BASE_URL=http://mirrors.tuna.tsinghua.edu.cn/apache/maven/maven-3/${MAVEN_VERSION}/binaries

RUN mkdir -p /usr/share/maven /usr/share/maven/ref \
  && curl -fsSL -o /tmp/apache-maven.tar.gz ${BASE_URL}/apache-maven-${MAVEN_VERSION}-bin.tar.gz \
  && echo "${SHA}  /tmp/apache-maven.tar.gz" | sha256sum -c - \
  && tar -xzf /tmp/apache-maven.tar.gz -C /usr/share/maven --strip-components=1 \
  && rm -f /tmp/apache-maven.tar.gz \
  && ln -s /usr/share/maven/bin/mvn /usr/bin/mvn

ENV MAVEN_HOME /usr/share/maven
ENV MAVEN_CONFIG "$USER_HOME_DIR/.m2"

COPY mvn-entrypoint.sh /usr/local/bin/mvn-entrypoint.sh
COPY settings-docker.xml /usr/share/maven/ref/

# Add deps & patch file
ADD alibaba-m2-deps.tar.gz /root/.m2/repository/com/alibaba/
COPY patch.diff /opt/dubbo/patch.diff

WORKDIR /opt/dubbo

# Install from Dubbo source code
RUN \
    wget -qO- https://github.com/alibaba/dubbo/archive/dubbo-2.5.3.tar.gz | tar -xzf - -C /opt \
    && mv /opt/dubbo-dubbo-2.5.3/* /opt/dubbo/ \
    #&& patch -p1 < patch.diff \
    && mvn package -Dmaven.test.skip=true \
    && mkdir /opt/dubbo-dist \
    && mv /opt/dubbo/dubbo/target/dubbo-2.5.3.jar /opt/dubbo-dist/dubbo-2.5.3.jar \
    && mv /opt/dubbo/dubbo-simple/dubbo-monitor-simple/target/dubbo-monitor-simple-2.5.3-assembly.tar.gz /opt/dubbo-dist/dubbo-monitor-simple-2.5.3-assembly.tar.gz \
    && mv /opt/dubbo/dubbo-admin/target/dubbo-admin-2.5.3.war /opt/dubbo-dist/dubbo-admin-2.5.3.war \
    && mv /opt/dubbo/dubbo-demo/dubbo-demo-provider/target/dubbo-demo-provider-2.5.3-assembly.tar.gz /opt/dubbo-dist/dubbo-demo-provider-2.5.3-assembly.tar.gz \
    && rm -rf /opt/maven \
    && rm -rf /opt/dubbo \
    && rm -rf /root/.m2


ENTRYPOINT ["/usr/local/bin/mvn-entrypoint.sh"]
CMD ["mvn"]
