gohadoop
========

gohadoop is native go clients for Apache Hadoop YARN.

It includes an early version of Hadoop IPC client and requisite YARN client libraries to implement YARN applications completely in go (both YARN application-client and application-master).

###Notes: 
* Set HADOOP_CONF_DIR environment variable, and ensure the conf directory contains both *-default.xml and *-site.xml files.
* hadoop_yarn/examples/dist_shell is an example go YARN application: client.go is the submission client and applicationmaster.go is the application-master. 

To run:
  

### Run regular dist-shell with AM in-cluster

$ HADOOP_CONF_DIR=conf go run hadoop_yarn/examples/dist_shell/client.go


### Run dist-shell with unmanaged AM

$ HADOOP_CONF_DIR=conf go run hadoop_yarn/examples/dist_shell/unmanaged.go
