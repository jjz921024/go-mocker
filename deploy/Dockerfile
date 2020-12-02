FROM alpine:latest

MAINTAINER JJZ <470623352@qq.com>

RUN echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.4/main/" > /etc/apk/repositories \
    && apk update \
    && apk upgrade \
    && apk add --no-cache bash \
    bash-doc \
    bash-completion \
    && rm -rf /var/cache/apk/* \
    && /bin/bash \
    && echo 'hosts: files mdns4_minimal [NOTFOUND=return] dns mdns4' >> /etc/nsswitch.conf

ADD mocker-go /bin

EXPOSE 8080
CMD ["/bin/mocker-go"]
