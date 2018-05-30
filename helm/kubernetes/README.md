site https://github.com/kubernetes/helm
wget https://kubernetes-helm.storage.googleapis.com/helm-v2.8.2-linux-amd64.tar.gz

install steps
自Kubernetes 1.6版本开始，API Server启用了RBAC授权。而目前的Tiller部署没有定义授权的ServiceAccount，这会导致访问API Server时被拒绝。我们可以采用如下方法，明确为Tiller部署添加授权。
kubectl create serviceaccount --namespace kube-system tiller
kubectl create clusterrolebinding tiller-cluster-rule --clusterrole=cluster-admin --serviceaccount=kube-system:tiller
kubectl patch deploy --namespace kube-system tiller-deploy -p '{"spec":{"template":{"spec":{"serviceAccount":"tiller"}}}}'

请执行如下命令利用阿里云的镜像来配置 Helm
helm init --upgrade -i slpcat/tiller:v2.8.2 --stable-repo-url https://kubernetes.oss-cn-hangzhou.aliyuncs.com/charts
helm search
若要更新charts列表以获取最新版本
helm repo update 
若要查看在群集上安装的Charts列表
helm list 
helm ls

示例
helm repo add stable https://kubernetes-charts.storage.googleapis.com
helm install --name wordpress-test --set "persistence.enabled=false,mariadb.persistence.enabled=false" stable/wordpress
https://kubeapps.com/ 你可以寻找和发现已有的Charts
