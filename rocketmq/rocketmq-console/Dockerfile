FROM tomcat:9.0

MAINTAINER yechen shan "shanyechen@fold.ren"

COPY . /usr/local/tomcat/webapps/

ENV NAMESRV namesrv1.rocketmq:9876

CMD ["catalina.sh", "run"]

EXPOSE 8080

Run:
docker run -d -e "NAMESRV=10.16.64.28:9876" --name rocketmq-console -p 30081:8080 shanyehcne/rocketmq-console:v0.1
