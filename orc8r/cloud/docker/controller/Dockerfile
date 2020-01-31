# ------------------------------------------------------------------------------
# Base image (for tests, precommit, codegen, etc.)
# ------------------------------------------------------------------------------
ARG baseImage="ubuntu:xenial"
FROM ${baseImage} as base

# Install the runtime deps from apt
RUN apt-get update && \
    apt-get install -y curl make git gcc bzr unzip vim

# Install Golang 1.12
WORKDIR /usr/local
RUN curl https://dl.google.com/go/go1.12.linux-amd64.tar.gz -O && \
    tar xf go1.12.linux-amd64.tar.gz && \
    cp -r go/bin/* /usr/local/bin/
ENV GO111MODULE on

# Env variables for specifying the output directories
ENV GOBIN /build/bin
ENV PLUGIN_DIR /build/plugins
ENV PATH ${PATH}:${GOBIN}
# Use public go modules proxy
ENV GOPROXY https://proxy.golang.org

# Install protobuf compiler
RUN curl -Lfs https://github.com/google/protobuf/releases/download/v3.1.0/protoc-3.1.0-linux-x86_64.zip -o protoc3.zip && \
    unzip protoc3.zip -d protoc3 && \
    mv protoc3/bin/protoc /bin/protoc && \
    chmod a+rx /bin/protoc && \
    mv protoc3/include/google /usr/include/ && \
    chmod -R a+Xr /usr/include/google && \
    rm -rf protoc3.zip protoc3

# Copy just the go.mod files to download the golang deps from internet.
# This step would allow us to cache the downloads in a layer, and prevent
# reaching out to the internet unless any of the go.mod file is changed.
COPY gomod /gomod
RUN mkdir -p /gomod/src/modules
RUN for module in `ls -d /gomod/src/magma/**`; do cd ${module}/cloud/go && go mod download; done
RUN find /gomod/src/modules -maxdepth 2 -mindepth 2 -type d -exec bash -c "cd {}/cloud/go && go mod download" \;

# Copy the configs
COPY configs /etc/magma/configs

# Copy the code
ENV MAGMA_ROOT /src/magma
COPY src /src
RUN mkdir -p /src/modules
WORKDIR /src/magma/orc8r/cloud
RUN echo export INTERNAL_MAGMA_MODULES="\"`ls -d /src/magma/** | tr '\n' ' '`\"" >> /etc/profile.d/env.sh
RUN echo export EXTERNAL_MAGMA_MODULES="\"`find /src/modules -maxdepth 2 -mindepth 2 -type d | tr '\n' ' '`\"" >> /etc/profile.d/env.sh
RUN . /etc/profile.d/env.sh && echo export MAGMA_MODULES="\"$INTERNAL_MAGMA_MODULES $EXTERNAL_MAGMA_MODULES\"" >> /etc/profile.d/env.sh
RUN . /etc/profile.d/env.sh && make tools

# ------------------------------------------------------------------------------
# Builder image with binaries and plugins
# ------------------------------------------------------------------------------
FROM base as builder

RUN . /etc/profile.d/env.sh && make build

# -----------------------------------------------------------------------------
# Go-cache base image
# -----------------------------------------------------------------------------
FROM ubuntu:xenial as gocache
COPY --from=builder /root/.cache /root/.cache

# ------------------------------------------------------------------------------
# Production image
# ------------------------------------------------------------------------------
FROM ubuntu:xenial

# Install the runtime deps from apt
RUN apt-get update && \
    apt-get install -y daemontools netcat openssl supervisor wget unzip

ARG CNTLR_FILES=src/magma/orc8r/cloud/docker/controller

# Setup Swagger UI
COPY ${CNTLR_FILES}/setup_swagger_ui /usr/local/bin/setup_swagger_ui
RUN /usr/local/bin/setup_swagger_ui

# Create an empty envdir for overriding in production
RUN mkdir -p /var/opt/magma/envdir
# TODO: automatically read the services from the plugins
ENV CONTROLLER_SERVICES STREAMER,SUBSCRIBERDB,EPS_AUTHENTICATION,METRICSD,CERTIFIER,BOOTSTRAPPER,METERINGD_RECORDS,ACCESSD,OBSIDIAN,LOGGER,DISPATCHER,DIRECTORYD,FEG_RELAY,HEALTH,VPNSERVICE,STATE,DEVICE,ANALYTICS,TENANTS
# TODO: fix dispatcher to work on containers
ENV SERVICE_HOST_NAME localhost
# TODO: fix metricsd for fetching host level metrics
ENV HOST_NAME dev

COPY src/magma/orc8r/cloud/docker/wait-for-it.sh /usr/local/bin

# Copy the configs
COPY configs /etc/magma/configs

# Copy the build artifacts
COPY --from=builder /build/bin /var/opt/magma/bin
COPY --from=builder /build/plugins /var/opt/magma/plugins
COPY --from=builder /${CNTLR_FILES}/apidocs /var/opt/magma/static/apidocs

# Copy the scripts from the context
COPY ${CNTLR_FILES}/create_test_controller_certs /usr/local/bin/create_test_controller_certs

# Copy the supervisor configs
COPY ${CNTLR_FILES}/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY ${CNTLR_FILES}/supervisor_logger.py /usr/local/lib/python2.7/dist-packages/supervisor_logger.py
CMD ["/usr/bin/supervisord"]
