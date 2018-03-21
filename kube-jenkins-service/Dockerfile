# upsteam https://github.com/mesosphere/dcos-jenkins-service
FROM jenkins/jenkins:lts

MAINTAINER 若虚 <slpcat@qq.com>

USER root

# Container variables
ENV \ 
    TERM="xterm" \
    LANG="en_US.UTF-8" \ 
    LANGUAGE="en_US:en" \
    LC_ALL="en_US.UTF-8" \
    TIMEZONE="Asia/Shanghai"

# Set timezone and locales
RUN \
    echo "${TIMEZONE}" > /etc/timezone \
    && echo "$LANG UTF-8" > /etc/locale.gen \
    && apt-get update -q \
    && DEBIAN_FRONTEND=noninteractive apt-get install -yq apt-utils dialog vim-tiny locales \ 
    && ln -sf /usr/share/zoneinfo/${TIMEZONE} /etc/localtime \
    && update-locale LANG=$LANG \
    && locale-gen $LANG \
    && DEBIAN_FRONTEND=noninteractive dpkg-reconfigure locales

# Install required packages
RUN \
    apt-get dist-upgrade -y

WORKDIR /tmp

# Environment variables used throughout this Dockerfile
#
# $JENKINS_HOME     will be the final destination that Jenkins will use as its
#                   data directory. This cannot be populated before Marathon
#                   has a chance to create the host-container volume mapping.
#
ENV JENKINS_FOLDER /usr/share/jenkins

# Build Args
ARG BLUEOCEAN_VERSION=1.4.2
ARG JENKINS_STAGING=/usr/share/jenkins/ref/
ARG CURL_RETRY=20
ARG CURL_RETRY_MAX_TIME=900

# install dependencies
RUN apt-get update && apt-get install -y python zip jq git

# Override the default property for DNS lookup caching
RUN echo 'networkaddress.cache.ttl=60' >> ${JAVA_HOME}/jre/lib/security/java.security

RUN mkdir -p "$JENKINS_HOME" "${JENKINS_FOLDER}/war"

# jenkins setup
COPY conf/jenkins/config.xml "${JENKINS_STAGING}/config.xml"
#COPY conf/jenkins/jenkins.model.JenkinsLocationConfiguration.xml "${JENKINS_STAGING}/jenkins.model.JenkinsLocationConfiguration.xml"
#COPY conf/jenkins/nodeMonitors.xml "${JENKINS_STAGING}/nodeMonitors.xml"
#COPY conf/jenkins/hudson.model.UpdateCenter.xml "${JENKINS_HOME}/hudson.model.UpdateCenter.xml"

# add plugins
RUN /usr/local/bin/install-plugins.sh       \
  blueocean-commons:${BLUEOCEAN_VERSION}    \
  blueocean-config:${BLUEOCEAN_VERSION}     \
  blueocean-dashboard:${BLUEOCEAN_VERSION}  \
  blueocean-display-url:2.2.0               \
  blueocean-events:${BLUEOCEAN_VERSION}     \
  blueocean-git-pipeline:${BLUEOCEAN_VERSION}          \
  blueocean-github-pipeline:${BLUEOCEAN_VERSION}       \
  blueocean-i18n:${BLUEOCEAN_VERSION}       \
  blueocean-jwt:${BLUEOCEAN_VERSION}        \
  blueocean-personalization:${BLUEOCEAN_VERSION}    \
  blueocean-pipeline-api-impl:${BLUEOCEAN_VERSION}  \
  blueocean-pipeline-editor:${BLUEOCEAN_VERSION}           \
  blueocean-pipeline-scm-api:${BLUEOCEAN_VERSION}   \
  blueocean-rest-impl:${BLUEOCEAN_VERSION}  \
  blueocean-rest:${BLUEOCEAN_VERSION}       \
  blueocean-web:${BLUEOCEAN_VERSION}        \
  blueocean:${BLUEOCEAN_VERSION}            \
  android-emulator:3.0           \
  android-lint:2.6               \
  ant:1.8                        \
  ace-editor:1.1                 \
  analysis-core:1.93             \
  ansible:0.8                    \
  ansicolor:0.5.2                \
  antisamy-markup-formatter:1.5  \
  artifactory:2.15.0             \
  audit-trail:2.2                \
  authentication-tokens:1.3      \
  azure-credentials:1.5.0        \
  azure-vm-agents:0.6.2          \
  branch-api:2.0.18              \
  build-failure-analyzer:1.19.2  \
  build-name-setter:1.6.9        \
  build-pipeline-plugin:1.5.8    \
  build-timeout:1.19             \
  build-token-root:1.4           \
  cloudbees-folder:6.4           \
  cloverphp:0.5                  \
  conditional-buildstep:1.3.6    \
  config-file-provider:2.18      \
  copy-to-slave:1.4.4            \
  copyartifact:1.39              \
  cvs:2.14                       \
  dashboard-view:2.9.11          \
  delivery-pipeline-plugin:1.1.0 \
  description-setter:1.10        \
  dingding-notifications:1.4     \
  docker-build-publish:1.3.2     \
  docker-workflow:1.15.1         \
  durable-task:1.22              \
  ec2:1.39                       \
  email-ext:2.61                 \
  embeddable-build-status:1.9    \
  envinject:2.1.5                \
  extended-read-permission:2.0   \
  external-monitor-job:1.7       \
  favorite:2.3.1                 \
  ghprb:1.39.0                   \
  git:3.8.0                      \
  git-client:2.7.1               \
  git-server:1.7                 \
  github:1.29.0                  \
  github-api:1.90                \
  github-organization-folder:1.6 \
  gitlab:1.5.2                   \
  gitlab-hook:1.4.2              \
  gitlab-oauth:1.3               \
  gradle:1.28                    \
  greenballs:1.15                \
  handlebars:1.1.1               \
  ivy:1.28                       \
  jackson2-api:2.8.10.1          \
  jenkins-multijob-plugin:1.29   \
  job-dsl:1.68                   \
  jobConfigHistory:2.18          \
  jquery:1.12.4-0                \
  kpp-management-plugin:1.0.0    \
  kubernetes:1.3.3               \
  kubernetes-ci:1.3              \
  kubernetes-cd:0.1.5            \
  kubernetes-pipeline-steps:1.5  \
  ldap:1.20                      \
  mapdb-api:1.0.9.0              \
  marathon:1.6.0                 \
  matrix-auth:2.2                \
  matrix-project:1.12            \
  maven-plugin:3.1               \
  mesos:0.15.1                   \
  metrics:3.1.2.10               \
  momentjs:1.1.1                 \
  monitoring:1.71.0              \
  msbuild:1.29                   \
  nant:1.4.3                     \
  node-iterator-api:1.5.0        \
  pam-auth:1.3                   \
  parameterized-trigger:2.35.2   \
  pipeline-build-step:2.6        \
  pipeline-github-lib:1.0        \
  pipeline-input-step:2.8        \
  pipeline-milestone-step:1.3.1  \
  pipeline-model-api:1.2.6       \
  pipeline-model-definition:1.2.6 \
  pipeline-model-extensions:1.2.6 \
  pipeline-rest-api:2.9          \
  pipeline-stage-step:2.3        \
  pipeline-stage-view:2.9        \
  plain-credentials:1.4          \
  postbuildscript:2.2.1          \
  publish-over-cifs:0.6          \
  publish-over-ftp:1.13          \
  publish-over-ssh:1.19.1        \
  puppet:1.1                     \
  rebuild:1.27                   \
  role-strategy:2.7.0            \
  run-condition:1.0              \
  s3:0.11.0                      \
  saferestart:0.3                \
  saml:1.0.4                     \
  saltstack:3.1.6                \
  scm-api:2.2.6                  \
  ssh-agent:1.15                 \
  ssh-slaves:1.25                \
  slave-setup:1.10               \
  subversion:2.10.3              \
  timestamper:1.8.9              \
  translation:1.16               \
  uno-choice:2.1                 \
  variant:1.1                    \
  view-job-filters:1.27          \
  windows-slaves:1.3.1           \
  workflow-aggregator:2.5        \
  workflow-api:2.24              \
  workflow-basic-steps:2.6       \
  workflow-cps:2.45              \
  workflow-cps-global-lib:2.9    \
  workflow-durable-task-step:2.17 \
  workflow-job:2.17              \
  workflow-multibranch:2.17      \
  workflow-scm-step:2.6          \
  workflow-step-api:2.14         \
  workflow-support:2.17          \
  ws-cleanup:0.34                \
  xcode-plugin:2.0.0

# disable first-run wizard
RUN echo 2.0 > /usr/share/jenkins/ref/jenkins.install.UpgradeWizard.state

# remove executors in master
COPY src/main/docker/master-executors.groovy /usr/share/jenkins/ref/init.groovy.d/

# ENV JAVA_OPTS="-Djava.util.logging.config.file=/var/jenkins_home/log.properties"
ENV JAVA_OPTS="-XX:+UnlockExperimentalVMOptions -XX:+UseCGroupMemoryLimitForHeap -XX:MaxRAMFraction=1 -XshowSettings:vm \
	-Djenkins.install.runSetupWizard=false -Dhudson.udp=-1 -Djava.awt.headless=true -Dhudson.DNSMultiCast.disabled=true"
ENV JENKINS_OPTS="--webroot=${JENKINS_FOLDER}/war --httpListenAddress=0.0.0.0"
