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

if [ -d "${HADOOP_HOME}" ]
then
  echo "hadoop home : ${HADOOP_HOME} already exists. terminating installation"
  exit 0
fi

if [ ! -f "${HADOOP_ARCHIVE}" ]
then
  echo -n "hadoop archive not found. downloading archive ( ~140MB ). this could take a while ...  "
  mkdir -p $HADOOP_ARCHIVE_DIR
  curl -s -S -L http://bit.ly/kubernetes-yarn-hadoop-snapshot -o $HADOOP_ARCHIVE
  echo "done."
fi 

echo "installing JDK: $JDK_PKG_VER $JDK_DEVEL_PKG_VER ... "
yum install -y $JDK_PKG_VER $JDK_DEVEL_PKG_VER

echo "creating hadoop install dir: ${HADOOP_INSTALL_DIR} "
mkdir -p ${HADOOP_INSTALL_DIR}

echo "extracting hadoop installation ... "
tar -zxvf ${HADOOP_ARCHIVE} -C ${HADOOP_INSTALL_DIR}

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
echo "list of slaves: ${MINION_IPS_STR}"
echo ${MINION_IPS_STR} | sed "s/,/\n/g" > ${SLAVES}

echo "done installing hadoop environment"

#echo "changing owner for hadoop installation.. "
#echo "chown -R vagrant:vagrant ${HADOOP_INSTALL_DIR}
