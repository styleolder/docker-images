source env.sh
export YARN_EXAMPLES_JAR=./share/hadoop/mapreduce/hadoop-mapreduce-examples-2.6.0-SNAPSHOT.jar
bin/yarn jar $YARN_EXAMPLES_JAR pi 4 10000
