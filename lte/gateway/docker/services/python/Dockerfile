FROM ubuntu:20.04 as control_proxy

ENV VIRTUAL_ENV=/build
ARG JSONPOINTER_VERSION=1.13

RUN apt-get update \
  && DEBIAN_FRONTEND=noninteractive apt-get install -y python3-pip python3-venv virtualenv git python3-eventlet python3-pystemd net-tools inetutils-ping\\
  && cd /usr/local/bin \
  && ln -s /usr/bin/python3 python
RUN python3 -m venv $VIRTUAL_ENV
ENV PATH="$VIRTUAL_ENV/bin:$PATH"

COPY --from=pythonbuilder:latest /build /build
COPY --from=pythonbuilder:latest /magma /magma

COPY --from=pythonbuilder:latest /var/tmp/python3-aioeventlet* /var/tmp/
RUN dpkg -i /var/tmp/python3-aioeventlet*

RUN mkdir -p /var/opt/magma/

RUN pip install jsonpointer>$JSONPOINTER_VERSION
