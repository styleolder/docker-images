#!/usr/bin/bash

#this script must be run from inside the cluster/vagrant directory
set -e

MASTER_IP=$1
MINION_IPS_STR=$2
JDK_PKG_VER=java-1.7.0-openjdk-1.7.0.65-2.5.2.5.fc20
JDK_DEVEL_PKG_VER=java-1.7.0-openjdk-devel-1.7.0.65-2.5.2.5.fc20
JAVA_HOME=/usr/lib/jvm/${JDK_PKG_VER}.x86_64
HADOOP_INSTALL_DIR=/home/vagrant/hadoop/install
HADOOP_HOME=${HADOOP_INSTALL_DIR}/hadoop-2.6.0-SNAPSHOT
ENV_CONFIG=${HADOOP_HOME}/env.sh
HADOOP_ARCHIVE_DIR=./hadoop
HADOOP_ARCHIVE=${HADOOP_ARCHIVE_DIR}/hadoop-2.6.0-SNAPSHOT.tgz
YARN_SITE=${HADOOP_HOME}/etc/hadoop/yarn-site.xml
SLAVES=${HADOOP_HOME}/etc/hadoop/slaves
MR_SCRIPT=test-pi-yarn.sh

echo "creating hadoop environment script... "
cat <<EOF > ${ENV_CONFIG}
export JAVA_HOME=$JAVA_HOME
export HADOOP_MAPRED_HOME=${HADOOP_HOME}
export HADOOP_COMMON_HOME=${HADOOP_HOME}
export HADOOP_HDFS_HOME=${HADOOP_HOME}
export YARN_HOME=${HADOOP_HOME}
export HADOOP_CONF_DIR=${HADOOP_HOME}/etc/hadoop
export YARN_CONF_DIR=${HADOOP_HOME}/etc/hadoop
EOF

echo "setting up yarn-site.xml .. "
sed -i "s/__YARN_RESOURCE_MANAGER_IP__/${MASTER_IP}/g" ${YARN_SITE}

echo "setting up slaves file .. "
echo "list of slaves : ${MINION_IPS_STR}"
echo ${MINION_IPS_STR} | sed "s/,/\n/g" > ${SLAVES}

echo "copying test map reduce script .. "
cp ${MR_SCRIPT} ${HADOOP_HOME}

echo "done installing hadoop environment"
