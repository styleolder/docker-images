#!/usr/bin/bash

HADOOP_INSTALL_DIR=/home/vagrant/hadoop/install
HADOOP_HOME=${HADOOP_INSTALL_DIR}/hadoop-2.6.0-SNAPSHOT

source ${HADOOP_HOME}/env.sh

${HADOOP_HOME}/sbin/yarn-daemon.sh stop resourcemanager
${HADOOP_HOME}/sbin/yarn-daemon.sh start resourcemanager

${HADOOP_HOME}/sbin/mr-jobhistory-daemon.sh stop historyserver
${HADOOP_HOME}/sbin/mr-jobhistory-daemon.sh start historyserver

#fails if namenode has already been formatted
${HADOOP_HOME}/sbin/hadoop-daemon.sh stop namenode
${HADOOP_HOME}/bin/hadoop namenode -format -nonInteractive
${HADOOP_HOME}/sbin/hadoop-daemon.sh start namenode
