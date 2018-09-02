#!/usr/bin/env bash
cp /opt/rocketmq/conf/broker.conf.sample /opt/rocketmq/conf/broker.conf

HOST=`hostname -s`

if [[ $HOST =~ (.*)-([0-9]+)$ ]]; then
    CLUSTER_NAME=${BASH_REMATCH[1]}
    BROKER_NAME=${BASH_REMATCH[1]}-${BASH_REMATCH[2]}
    BROKER_ROLE=${BASH_REMATCH[3]}
else
    echo "Failed to extract ordinal from hostname $HOST"
    exit 1
fi

sed -i s/^brokerClusterName=.*$/brokerClusterName=$CLUSTER_NAME/ /opt/rocketmq/conf/broker.conf
sed -i s/^brokerName=.*$/brokerName=$BROKER_NAME/ /opt/rocketmq/conf/broker.conf

if [[ $BROKER_ROLE == slave ]]; then 
sed -i s/^brokerId=.*$/brokerId=1/ /opt/rocketmq/conf/broker.conf
fi
