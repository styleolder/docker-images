#upstream https://github.com/maauso/docker/tree/master/zkui
FROM java:7
MAINTAINER Miguel Auso <m.auso.p@gmail.com> https://github.com/maauso/docker.git
RUN apt-get update && apt-get install -y git maven
RUN git clone https://github.com/DeemOpen/zkui.git
RUN cd /zkui && mvn clean install
WORKDIR /zkui/target/
COPY config.cfg /zkui/target/
COPY entrypoint.sh /zkui/target/
RUN chmod +x entrypoint.sh
EXPOSE 9090
ENTRYPOINT ["./entrypoint.sh"]
CMD ["/usr/bin/java", "-jar", "/zkui/target/zkui-2.0-SNAPSHOT-jar-with-dependencies.jar"]
