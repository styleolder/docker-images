#!/bin/bash

# Copyright 2014 Google Inc. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# exit on any error
set -e

# Setup hosts file to support ping by hostname to master
if [ ! "$(cat /etc/hosts | grep $MASTER_NAME)" ]; then
  echo "Adding $MASTER_NAME to hosts file"
  echo "$MASTER_IP $MASTER_NAME" >> /etc/hosts
fi

# Setup hosts file to support ping by hostname to each minion in the cluster
for (( i=0; i<${#MINION_NAMES[@]}; i++)); do
  minion=${MINION_NAMES[$i]}
  ip=${MINION_IPS[$i]}
  if [ ! "$(cat /etc/hosts | grep $minion)" ]; then
    echo "Adding $minion to hosts file"
    echo "$ip $minion" >> /etc/hosts
  else
    host_entry=$(cat /etc/hosts | grep $minion)
    ip_in_file=$(echo $host_entry | awk '{print $1}')
    echo "existing host entry is \"$host_entry\""
    echo "ip is \"$ip_in_file\""
    if [ "$ip_in_file" == "127.0.0.1" ]; then
      echo "$minion has a 127.0.0.1 entry - fixing." 
      sed -i "s/127\.0\.0\.1.*/127.0.0.1 localhost/g" /etc/hosts
      echo "Adding $minion to hosts file"
      echo "$ip $minion" >> /etc/hosts
    fi
  fi
done

# Let the minion know who its master is
mkdir -p /etc/salt/minion.d
cat <<EOF >/etc/salt/minion.d/master.conf
master: '$(echo "$MASTER_NAME" | sed -e "s/'/''/g")'
EOF

cat <<EOF >/etc/salt/minion.d/log-level-debug.conf
log_level: debug
log_level_logfile: debug
EOF

# Our minions will have a pool role to distinguish them from the master.
cat <<EOF >/etc/salt/minion.d/grains.conf
grains:
  network_mode: openvswitch
  node_ip: '$(echo "$MINION_IP" | sed -e "s/'/''/g")'
  etcd_servers: '$(echo "$MASTER_IP" | sed -e "s/'/''/g")'
  api_servers: '$(echo "$MASTER_IP" | sed -e "s/'/''/g")'
  networkInterfaceName: eth1
  apiservers: '$(echo "$MASTER_IP" | sed -e "s/'/''/g")'
  roles:
    - kubernetes-pool
    - kubernetes-pool-vagrant
  cbr-cidr: '$(echo "$CONTAINER_SUBNET" | sed -e "s/'/''/g")'
  minion_ip: '$(echo "$MINION_IP" | sed -e "s/'/''/g")'
EOF

#Install hadoop before installing kubernetes
MINION_IPS_STR=${MINION_IPS[@]}
MINION_IPS_STR=${MINION_IPS_STR// /,}

echo "Installing hadoop ..."
pushd /vagrant/cluster/vagrant
./provision-hadoop-existing-hadoop.sh $MASTER_IP $MINION_IPS_STR
./restart-hadoop-slave-daemons.sh
popd

#enable/stop/start salt-minion
systemctl enable salt-minion.service
systemctl stop salt-minion.service
systemctl start salt-minion.service

# run the networking setup
"${KUBE_ROOT}/cluster/vagrant/provision-network.sh" $@
