#set -o errexit

echo "Bringing up redis-master .... "
cluster/kubecfg.sh -c examples/guestbook/redis-master.json create pods

echo "Cheking to see if redis-master is running ... "
cluster/kubecfg.sh list pods

echo "Bringing up redis-master-service ... "
cluster/kubecfg.sh -c examples/guestbook/redis-master-service.json create services

echo "bringing up replication slaves... "
cluster/kubecfg.sh -c examples/guestbook/redis-slave-controller.json create replicationControllers

echo "Checking to see if all pods are running... "
cluster/kubecfg.sh list pods

echo "Bringing up redis-slave-service ... "
cluster/kubecfg.sh -c examples/guestbook/redis-slave-service.json create services

echo "Bringing up frontend pods ... "
cluster/kubecfg.sh -c examples/guestbook/frontend-controller.json create replicationControllers

echo "Listing all pods ... "
cluster/kubecfg.sh list pods

echo "Done."

