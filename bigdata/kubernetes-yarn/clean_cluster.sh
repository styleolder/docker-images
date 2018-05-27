CMD=cluster/kubecfg.sh

echo "deleting services ... "
SERVICES=`$CMD list services | grep name= | cut -d " " -f1`
for SERVICE in $SERVICES
do
        echo "deleting $SERVICE ..."
        $CMD delete services/$SERVICE
done

echo "deleting replication controllers .. "
REPLICATION_CONTROLLERS=`$CMD list replicationControllers | grep name= | cut -d " " -f1`
for CONTROLLER in $REPLICATION_CONTROLLERS
do
        echo "stopping $CONTROLLER ..."
        $CMD stop $CONTROLLER
        echo "deleting $CONTROLLER ..."
        $CMD delete replicationControllers/$CONTROLLER
done

echo "deleting pods ... "
PODS=`$CMD list pods | grep "name=" | cut -d " " -f1`
for POD in $PODS
do
        echo "deleting $POD ... "
        $CMD delete pods/$POD
done

echo "done."
