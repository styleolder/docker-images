Ceph on Kubernetes
ceph作为Kubernetes持久存储服务，以pod形式运行,提供rbd存储类，支持自动创建持久卷供其他pod使用
使用限制与要求
kubernetes 节点内核版本 >= 4.5
public and cluster networks 必须相同，并且是kubernetes的集群内部网络，本例使用kubespray安装脚本的默认网络10.233.0.0/16If the storage class user id is not admin, you will have to manually create the user in your Ceph cluster and create its secret in Kubernetes
ceph-mgr can only run with 1 replica
rbd块设备，cephfs仅支持集群内访问
rgw对象存储支持集群内部和外部同时访问
因为主机名检测问题，暂时不兼容istio

生产环境安装过程
规划
奇数个宿主机节点作为mon，直接使用宿主机目录作为持久存储
osd数据盘使用宿主机无分区裸盘
宿主机节点安装ceph客户端
宿主机节点dns解析使用集群dns服务器
/etc/resolv.conf
domain <EXISTING_DOMAIN>
search <EXISTING_DOMAIN>

#Your kubernetes cluster ip domain
search ceph.svc.cluster.local svc.cluster.local cluster.local

nameserver 10.233.0.3     #The cluster IP of skyDNS
nameserver <EXISTING_RESOLVER_IP>

客户端（控制台）要求
In addition to kubectl, jinja2 or sigil is required for template handling and must be installed in your system PATH. Instructions can be found here for jinja2 https://github.com/mattrobenolt/jinja2-cli or here for sigil https://github.com/gliderlabs/sigil.

覆盖默认的网络设置
export osd_cluster_network=10.233.0.0/16
export osd_public_network=10.233.0.0/16

生成keys和ceph配置

cd generator
./generate_secrets.sh all `./generate_secrets.sh fsid`

kubectl create namespace ceph

kubectl create secret generic ceph-conf-combined --from-file=ceph.conf --from-file=ceph.client.admin.keyring --from-file=ceph.mon.keyring --namespace=ceph
kubectl create secret generic ceph-bootstrap-rgw-keyring --from-file=ceph.keyring=ceph.rgw.keyring --namespace=ceph
kubectl create secret generic ceph-bootstrap-mds-keyring --from-file=ceph.keyring=ceph.mds.keyring --namespace=ceph
kubectl create secret generic ceph-bootstrap-osd-keyring --from-file=ceph.keyring=ceph.osd.keyring --namespace=ceph
kubectl create secret generic ceph-bootstrap-rbd-keyring --from-file=ceph.keyring=ceph.rbd.keyring --namespace=ceph
kubectl create secret generic ceph-client-key --from-file=ceph-client-key --namespace=ceph

cd ..

生产环境部署ceph组件
生成授权
kubectl create -f ceph-rbac.yaml

mds
ceph-mds-v1-dp.yaml

mgr
ceph-mgr-v1-dp.yaml
ceph-mgr-dashboard-v1-svc.yam
ceph-mgr-prometheus-v1-svc.yaml

mon
2n+1台 生产5台以上，宿主机节点需要标签  ceph-mon: enabled
ceph-mon-v1-svc.yaml
ceph-mon-v1-ds.yaml
ceph-mon-check-v1-dp.yam

osd 使用持久存储
每个盘对应一个pod,本示例使用/dev/vdb
先初始化磁盘

kubectl create -f ceph-osd-prepare-v1-ds.yaml --namespace=ceph
初始化成功后删除pod
kubectl delete -f ceph-osd-prepare-v1-ds.yaml --namespace=ceph

初始化完毕，激活磁盘
kubectl create -f ceph-osd-activate-v1-ds.yaml --namespace=ceph

rgw
ceph-rgw-v1-dp.yaml
ceph-rgw-v1-svc.yaml

给存储节点打上标签(必须)

kubectl label node <nodename> node-type=storage
If you want all nodes in your Kubernetes cluster to be a part of your Ceph cluster, label them all.

kubectl label nodes node-type=storage --all
Eventually all pods will be running, including a mon and osd per every labeled node.

kubernetes使用外部持久卷
https://github.com/kubernetes-incubator/external-storage/tree/master/ceph/rbd/deploy/rbac
进入mon pod
kubectl exec -it ceph-mon bash
创建存储池
ceph osd pool create kube 64
创建keyring
ceph auth get-or-create client.kube mon 'allow r' osd \
  'allow class-read object_prefix rbd_children, allow rwx pool=kube' \
  -o ceph.client.kube.keyring
根据生成的keyring创建secret
kubectl --namespace=ceph create secret generic ceph-rbd-kube \
  --from-literal="key=$(grep key ceph.client.kube.keyring  | awk '{ print $3 }')" \
  --type=kubernetes.io/rbd

创建RBD provisioner 使用namespace=ceph
创建rbac授权
kubectl create -f rbd-provisioner/clusterrolebinding.yaml
kubectl create -f rbd-provisioner/clusterrole.yaml
kubectl create -f rbd-provisioner/rolebinding.yaml
kubectl create -f rbd-provisioner/role.yaml
kubectl create -f rbd-provisioner/serviceaccount.yaml
创建RBD provisioner pod
kubectl create -f rbd-provisioner/deployment.yaml
创建RBD存储类
$ kubectl create secret generic ceph-secret-admin --from-file=generator/ceph-client-key --type=kubernetes.io/rbd --namespace=ceph
kubectl create -f rbd-provisioner/storage-class.yaml

POD作为使用者,只需创建pvc获取相应容量的存储，然后挂载即可自动获取存储资源
创建pvc 示例
$ kubectl create -f https://raw.githubusercontent.com/kubernetes/examples/master/staging/persistent-volume-provisioning/claim1.json

POD里面挂载pvc卷,通常是有状态副本集statefulset.示例

rbd-pvc-pod.yaml

生产环境推荐使用有状态副本集statefulset自动创建pvc
  volumeClaimTemplates:
  - metadata:
      name: datadir
    spec:
      accessModes:
      - ReadWriteOnce
      #ceph rbd storageclass
      storageClassName: rbd
      resources:
        requests:
          storage: 10Gi

POD里面挂载CephFS,通常是有状态副本集statefulset.示例
must add the admin client key

kubectl create \
-f ceph-mds-v1-dp.yaml \
-f ceph-mon-v1-svc.yaml \
-f ceph-mon-v1-dp.yaml \
-f ceph-mon-check-v1-dp.yaml \
-f ceph-osd-v1-ds.yaml \
--namespace=ceph

$ kubectl get all --namespace=ceph
NAME                   DESIRED      CURRENT       AGE
ceph-mds               1            1             24s
ceph-mon-check         1            1             24s
NAME                   CLUSTER-IP   EXTERNAL-IP   PORT(S)    AGE
ceph-mon               None         <none>        6789/TCP   24s
NAME                   READY        STATUS        RESTARTS   AGE
ceph-mds-6kz0n         0/1          Pending       0          24s
ceph-mon-check-deek9   1/1          Running       0          24s

$ kubectl get pods --namespace=ceph
NAME                   READY     STATUS    RESTARTS   AGE
ceph-mds-6kz0n         1/1       Running   0          4m
ceph-mon-8wxmd         1/1       Running   2          2m
ceph-mon-c8pd0         1/1       Running   1          2m
ceph-mon-cbno2         1/1       Running   1          2m
ceph-mon-check-deek9   1/1       Running   0          4m
ceph-mon-f9yvj         1/1       Running   1          2m
ceph-osd-3zljh         1/1       Running   2          2m
ceph-osd-d44er         1/1       Running   2          2m
ceph-osd-ieio7         1/1       Running   2          2m
ceph-osd-j1gyd         1/1       Running   2          2m
$ kubectl create -f https://raw.githubusercontent.com/kubernetes/examples/master/staging/persistent-volume-provisioning/claim1.json

Now, try create a claim:

If everything works, expect something like the following:

$ kubectl describe pvc claim1
Name:           claim1
Namespace:      default
StorageClass:
Status:         Bound
Volume:         pvc-a9247186-6e59-11e7-b7b6-00259003b6e8
Labels:         <none>
Capacity:       3Gi
Access Modes:   RWO
Events:
  FirstSeen     LastSeen        Count   From                                                                                    SubObjectPath   Type            Reason                  Message
  ---------     --------        -----   ----                                                                                    -------------   --------        ------                  -------
  6m            6m              2       {persistentvolume-controller }                                                                 Normal           ExternalProvisioning    cannot find provisioner "ceph.com/rbd", expecting that a volume for the claim is provisioned either manually or via external software
  6m            6m              1       {ceph.com/rbd rbd-provisioner-217120805-9dc84 57e293c8-6e59-11e7-a834-ca4351e8550d }           Normal           Provisioning            External provisioner is provisioning volume for claim "default/claim1"
  6m            6m              1       {ceph.com/rbd rbd-provisioner-217120805-9dc84 57e293c8-6e59-11e7-a834-ca4351e8550d }           Normal           ProvisioningSucceeded   Successfully provisioned volume pvc-a9247186-6e59-11e7-b7b6-00259003b6e8
Mounting CephFS in a pod

First you must add the admin client key to your current namespace (or the namespace of your pod).

kubectl create secret generic ceph-client-key --type="kubernetes.io/rbd" --from-file=./generator/ceph-client-key
Now, if skyDNS is set as a resolver for your host nodes then execute the below command as is. Otherwise modify the ceph-mon.ceph host to match the IP address of one of your ceph-mon pods.

kubectl create -f ceph-cephfs-test.yaml --namespace=ceph
You should be able to see the filesystem mounted now

kubectl exec -it --namespace=ceph ceph-cephfs-test df
Mounting a Ceph RBD in a pod

First we have to create an RBD volume.

# This gets a random MON pod.
export PODNAME=`kubectl get pods --selector="app=ceph,daemon=mon" --output=template --template="{{with index .items 0}}{{.metadata.name}}{{end}}" --namespace=ceph`

kubectl exec -it $PODNAME --namespace=ceph -- rbd create ceph-rbd-test --size 20G

kubectl exec -it $PODNAME --namespace=ceph -- rbd info ceph-rbd-test
The same caveats apply for RBDs as Ceph FS volumes. Edit the pod accordingly. Once you're set:

kubectl create -f ceph-rbd-test.yaml --namespace=ceph
And again you should see your mount, but with 20 gigs free

kubectl exec -it --namespace=ceph ceph-rbd-test -- df -h
Common Modifications

Durable Storage

By default emptyDir is used for everything. If you have durable storage on your nodes, replace the emptyDirs with a hostPath to that storage.

参考：
1.https://github.com/ceph/ceph-container/tree/master/examples/kubernetes
2.http://docs.ceph.com/docs/master/start/kube-helm/
3.https://github.com/kubernetes-incubator/external-storage/tree/master/ceph/rbd
4.http://tracker.ceph.com/projects/ceph/wiki/Tuning_for_All_Flash_Deployments
5.http://docs.ceph.com/docs/master/rados/configuration/
