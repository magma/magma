FROM ubuntu:20.04 as magmad

ENV VIRTUAL_ENV=/build
ARG JSONPOINTER_VERSION=1.13

RUN apt-get update \
  && DEBIAN_FRONTEND=noninteractive apt-get install -y python3-pip python3-venv python3-eventlet python3-pystemd virtualenv net-tools inetutils-ping\
  && cd /usr/local/bin \
  && ln -s /usr/bin/python3 python
RUN python3 -m venv $VIRTUAL_ENV
ENV PATH="$VIRTUAL_ENV/bin:$PATH"

COPY --from=pythonbuilder:latest /build /build
COPY --from=pythonbuilder:latest /magma /magma
COPY --from=pythonbuilder:latest /var/tmp/python3-aioeventlet* /var/tmp/

RUN dpkg -i /var/tmp/python3-aioeventlet*

RUN pip install jsonpointer>$JSONPOINTER_VERSION
ENTRYPOINT /usr/bin/env python3 -m magma.magmad.main
