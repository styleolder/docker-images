#!/usr/bin/bash

HADOOP_INSTALL_DIR=/home/vagrant/hadoop/install
HADOOP_HOME=${HADOOP_INSTALL_DIR}/hadoop-2.6.0-SNAPSHOT

source ${HADOOP_HOME}/env.sh

${HADOOP_HOME}/sbin/yarn-daemon.sh stop nodemanager
${HADOOP_HOME}/sbin/yarn-daemon.sh start nodemanager
${HADOOP_HOME}/sbin/hadoop-daemon.sh stop datanode
${HADOOP_HOME}/sbin/hadoop-daemon.sh start datanode
