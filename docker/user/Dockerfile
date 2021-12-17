FROM alpine:latest
LABEL maintainer="bmc-user . Authors <qiang.zhuo@zenlayer.com>"

RUN apk add apache2-utils

ADD ./user /opt/user
ADD ./entrypoint.sh /entrypoint.sh

WORKDIR /opt/user

ENTRYPOINT ["/entrypoint.sh"]
