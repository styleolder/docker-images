FROM centos:7

MAINTAINER 若虚 <slpcat@qq.com>

# Container variables
ENV \
    TERM="xterm" \
    LANG="zh_CN.UTF-8" \
    LANGUAGE="zh_CN.UTF-8" \
    TIMEZONE="Asia/Shanghai"

COPY epel.repo /etc/yum.repos.d/
COPY java.sh /etc/profile.d/

# set timezone
RUN ln -sf /usr/share/zoneinfo/${TIMEZONE} /etc/localtime

RUN yum update -y
RUN yum install -y wget && wget --no-check-certificate --no-cookies --header "Cookie: oraclelicense=accept-securebackup-cookie" http://download.oracle.com/otn-pub/java/jdk/8u161-b12/jdk-8u161-linux-x64.rpm
RUN rpm -ivh jdk-8u161-linux-x64.rpm && rm jdk-8u161-linux-x64.rpm
