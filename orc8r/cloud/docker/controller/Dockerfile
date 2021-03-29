# ------------------------------------------------------------------------------
# Base: for tests, precommit, codegen, etc.
# ------------------------------------------------------------------------------
FROM ubuntu:xenial as base

ENV GO111MODULE on
ENV GOPATH ${USER}/go
ENV GOBIN /build/bin
ENV PATH ${PATH}:${GOBIN}:${GOPATH}/bin

# Used in Makefiles
ENV MAGMA_ROOT /src/magma

# Apt runtime deps
RUN apt-get update && apt-get install -y \
  bzr \
  curl \
  gcc \
  git \
  make \
  unzip \
  vim

# Golang 1.13
WORKDIR /usr/local
RUN curl https://facebookconnectivity.jfrog.io/artifactory/generic/go1.13.4.linux-amd64.tar.gz -O && \
    tar xf go1.13.4.linux-amd64.tar.gz && \
    cp -r go/bin/* /usr/local/bin/

# Protobuf compiler
# Apt has 2.x but we need 3.x
# See: https://grpc.io/docs/protoc-installation/
RUN curl -Lfs https://github.com/google/protobuf/releases/download/v3.1.0/protoc-3.1.0-linux-x86_64.zip -o protoc3.zip && \
    unzip protoc3.zip -d protoc3 && \
    mv protoc3/bin/protoc /bin/protoc && \
    chmod a+rx /bin/protoc && \
    mv protoc3/include/google /usr/include/ && \
    chmod -R a+Xr /usr/include/google && \
    rm -rf protoc3.zip protoc3

# ------------------------------------------------------------------------------
# Gocache: cache Go modules
# ------------------------------------------------------------------------------
FROM base as gocache

ARG MAGMA_MODULES="orc8r lte feg cwf wifi fbinternal"
RUN echo "export GOCACHE_MODULES=\"$(for m in $MAGMA_MODULES ; do echo -n /gomod/src/magma/$m ; echo -n ' ' ; done)\"" >> /etc/profile.d/env.sh

COPY gomod /gomod
RUN cat /etc/profile.d/env.sh
RUN . /etc/profile.d/env.sh && for m in $GOCACHE_MODULES ; do cd ${m}/cloud/go && echo ${m}/cloud/go && go mod download ; done

# ------------------------------------------------------------------------------
# Src: different src depending on MAGMA_MODULES
# ------------------------------------------------------------------------------
FROM gocache as src

ARG MAGMA_MODULES="orc8r lte feg cwf wifi fbinternal"
RUN echo "export MAGMA_MODULES=\"$(for m in $MAGMA_MODULES ; do echo -n /src/magma/$m ; echo -n ' ' ; done)\"" >> /etc/profile.d/env.sh

# Source code
COPY src /src
WORKDIR /src/magma/orc8r/cloud
RUN . /etc/profile.d/env.sh && make tools

# Configs
COPY configs /etc/magma/configs

# ------------------------------------------------------------------------------
# Builder: compile src
# ------------------------------------------------------------------------------
FROM src as builder

RUN . /etc/profile.d/env.sh && make build

# ------------------------------------------------------------------------------
# Production
# ------------------------------------------------------------------------------
FROM ubuntu:xenial

# Apt runtime deps
RUN apt-get update && apt-get install -y \
  daemontools \
  netcat \
  openssl \
  supervisor \
  unzip \
  wget \
&& rm -rf /var/lib/apt/lists/*

# Swagger UI
# See: https://github.com/swagger-api/swagger-ui
ARG SWAGGER_UI_VERSION="3.1.7"
RUN cd /tmp && \
  wget "https://github.com/swagger-api/swagger-ui/archive/v$SWAGGER_UI_VERSION.zip" && \
  unzip "v$SWAGGER_UI_VERSION.zip" -d swagger-ui && \
  mkdir -p /var/opt/magma/static/swagger-ui && \
  cp -r "swagger-ui/swagger-ui-$SWAGGER_UI_VERSION/dist" /var/opt/magma/static/swagger-ui

# Script to wait for DB to be reachable
COPY src/magma/orc8r/cloud/docker/wait-for-it.sh /usr/local/bin

# Empty envdir for overriding in production
RUN mkdir -p /var/opt/magma/envdir

# Build artifacts
ARG SWAGGER_FILES=src/magma/orc8r/cloud/go/obsidian/swagger
COPY --from=builder /${SWAGGER_FILES}/v1/index.html /var/opt/magma/static/swagger/v1/ui/index.html
COPY --from=builder /${SWAGGER_FILES}/v1/css/sidebar.css /var/opt/magma/static/swagger/v1/static/sidebar.css
COPY --from=builder /${SWAGGER_FILES}/v1/swagger.yml /var/opt/magma/static/swagger/v1/spec/swagger.yml
COPY --from=builder src/magma/orc8r/cloud/swagger /etc/magma/swagger
COPY --from=builder /build/bin /var/opt/magma/bin

# Supervisor configs
ARG CNTLR_FILES=src/magma/orc8r/cloud/docker/controller
COPY ${CNTLR_FILES}/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY ${CNTLR_FILES}/supervisor_logger.py /usr/local/lib/python2.7/dist-packages/supervisor_logger.py

# Scripts for dev mode
COPY ${CNTLR_FILES}/create_test_controller_certs /usr/local/bin/create_test_controller_certs

# Configs
COPY configs /etc/magma/configs
