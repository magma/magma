FROM ubuntu:20.04 as redirectd

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update \
  && apt-get install -y python3-pip python3-dev git \
  && cd /usr/local/bin \
  && ln -s /usr/bin/python3 python \

COPY python/magma/redirectd /usr/lib/python3.8/

ENTRYPOINT python3.8 -m redirectd.main
