# Kubernetes-YARN

A version of [Kubernetes](https://github.com/GoogleCloudPlatform/kubernetes) using [Apache Hadoop YARN](http://hadoop.apache.org/docs/current/hadoop-yarn/hadoop-yarn-site/YARN.html) as the scheduler. Integrating Kubernetes with YARN lets users run [Docker](https://www.docker.com/whatisdocker/) containers packaged as pods (using Kubernetes) and YARN applications (using YARN), while ensuring common resource management across these (PaaS and data) workloads. 

## Kubernetes-YARN is currently in the protoype/alpha phase
This integration is under development. Please expect bugs and significant changes as we work towards making things more stable and adding additional features.


## Getting started
### Dev Environment
Kubernetes and Kubernetes-YARN are written in [Go](http://golang.org). Currently, [vagrant](http://www.vagrantup.com/) and [ansible](http://docs.ansible.com/) based setup mechanims are supported. The instructions below are for creating a vagrant based cluster. For ansible instructions, see [here](https://github.com/hortonworks/kubernetes-yarn/blob/master/cluster/ansible/README.md). 

Please ensure you have [boot2docker](http://boot2docker.io/), Go (at least 1.3), Vagrant (at least 1.6), VirtualBox (at least 4.3.x) and git installed. Run boot2docker to bring up a VM with a running docker daemon (this is used for building release binaries for Kubernetes). 

```
$ $(boot2docker shellinit) #sets up docker env vars
$ echo $GOPATH
/home/user/goproj
$ mkdir -p $GOPATH/src/github.com/hortonworks/
$ cd $GOPATH/src/github.com/hortonworks/
$ git clone git@github.com:hortonworks/kubernetes-yarn.git
$ cd kubernetes-yarn
$ build/release.sh #builds kubernetes release binaries 
$ hack/build-go.sh #builds kubernetes client binaries
$ cluster/kube-up.sh #brings up kubernetes cluster
```
Following these steps will bring up a multi-VM cluster (1 master and 3 minions, by default) running Kubernetes and YARN. Please note that, depending on your local hardware and available bandwidth, bringing the cluster up could take a while to complete.
### YARN Dashboard
By default, the kubernetes master is assigned the IP 10.245.1.2. The YARN resource manager runs on the name host. Once the vagrant cluster is running, the YARN dashboard accessible at http://10.245.1.2:8088/

### HDFS Dashboard
The HDFS dashboard is accessible at http://10.245.1.2:50070/

## Interacting with the Kubernetes-YARN cluster
### Creating pods/running Docker containters
For instructions on creating pods, running containers and other interactions with the cluster, please see Kubernetes' vagrant instructions [here](https://github.com/GoogleCloudPlatform/kubernetes/blob/master/docs/getting-started-guides/vagrant.md#running-containers)

### Running a test map-reduce job
In order to run a test map-reduce job, log into the cluster (ensure that you are in the `kubernetes-yarn` directory) and run the included test script.

```
$ vagrant ssh master
[vagrant@kubernetes-master ~]$ cd hadoop/install/hadoop-2.6.0-SNAPSHOT/
[vagrant@kubernetes-master hadoop-2.6.0-SNAPSHOT]$ sudo su #you need to be root to run this test script
[root@kubernetes-master hadoop-2.6.0-SNAPSHOT]# ./test-pi-yarn.sh
```
