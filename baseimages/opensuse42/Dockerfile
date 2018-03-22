# upstream https://github.com/openSUSE/docker-containers-build/blob/6252e7bc2a411b3c5cfd5ed3b74e5839f7f7615b/x86_64/Dockerfile
FROM opensuse:leap

MAINTAINER 若虚 <slpcat@qq.com>

# Container variables
ENV \ 
    TERM="xterm" \
    LANG="en_US.UTF-8" \ 
    LANGUAGE="en_US.UTF-8" \
    #LC_ALL="en_US.UTF-8" \
    TIMEZONE="Asia/Shanghai"

RUN  \
    zypper refresh -f && \
    zypper -n update && \
    zypper -n install timezone && \
    zypper clean

# set timezone
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

CMD ["/bin/bash"]
