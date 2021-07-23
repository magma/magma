FROM ubuntu:20.04 as smsd

ENV VIRTUAL_ENV=/build
RUN apt-get update \
  && DEBIAN_FRONTEND=noninteractive apt-get install -y python3-pip python3-venv virtualenv\
  && cd /usr/local/bin \
  && ln -s /usr/bin/python3 python
RUN python3 -m venv $VIRTUAL_ENV
ENV PATH="$VIRTUAL_ENV/bin:$PATH"

COPY --from=pythonbuilder:latest /build /build
COPY --from=pythonbuilder:latest /magma /magma

ENTRYPOINT /usr/bin/env python3 -m magma.smsd.main
