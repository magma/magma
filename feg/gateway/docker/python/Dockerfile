# -----------------------------------------------------------------------------
# Builder image to generate proto files
# -----------------------------------------------------------------------------
FROM ubuntu:focal AS builder

# workaround to avoid interactive tzdata configuration
ARG DEBIAN_FRONTEND=noninteractive

# Install the runtime deps from apt.
RUN apt-get -y update && apt-get -y install curl make virtualenv zip \
    apt-utils software-properties-common apt-transport-https git openjdk-8-jdk ant

# Install protobuf compiler.
RUN curl -Lfs https://github.com/protocolbuffers/protobuf/releases/download/v3.1.0/protoc-3.1.0-linux-x86_64.zip -o protoc3.zip && \
    unzip protoc3.zip -d protoc3 && \
    mv protoc3/bin/protoc /usr/bin/protoc && \
    chmod a+rx /usr/bin/protoc && \
    cp -r protoc3/include/google /usr/include/ && \
    chmod -R a+Xr /usr/include/google && \
    rm -rf protoc3.zip protoc3

RUN apt-get -y update && apt-get -y install python3.8

ENV MAGMA_ROOT /magma
ENV PYTHON_BUILD /build/python
ENV PIP_CACHE_HOME ~/.pipcache
ENV CODEGEN_ROOT /var/tmp/codegen
ENV CODEGEN_VERSION 2.2.3
ENV SWAGGER_CODEGEN_DIR $CODEGEN_ROOT/modules/swagger-codegen-cli/target
ENV SWAGGER_CODEGEN_JAR $SWAGGER_CODEGEN_DIR/swagger-codegen-cli.jar
ENV JAVA_HOME /usr/lib/jvm/java-8-openjdk-amd64/
RUN printenv > /etc/environment

# Download swagger codegen
RUN mkdir -p ${SWAGGER_CODEGEN_DIR}; \
    curl -Lfs https://repo1.maven.org/maven2/io/swagger/swagger-codegen-cli/${CODEGEN_VERSION}/swagger-codegen-cli-${CODEGEN_VERSION}.jar -o ${SWAGGER_CODEGEN_JAR}

# Generate python proto bindings.
COPY cwf/protos $MAGMA_ROOT/cwf/protos
COPY cwf/swagger $MAGMA_ROOT/cwf/swagger
COPY feg/protos $MAGMA_ROOT/feg/protos
COPY feg/swagger $MAGMA_ROOT/feg/swagger
COPY lte/gateway/python/defs.mk $MAGMA_ROOT/lte/gateway/python/defs.mk
COPY lte/gateway/python/Makefile $MAGMA_ROOT/lte/gateway/python/Makefile
COPY lte/protos $MAGMA_ROOT/lte/protos
COPY lte/swagger $MAGMA_ROOT/lte/swagger
COPY orc8r/gateway/python $MAGMA_ROOT/orc8r/gateway/python
COPY orc8r/protos $MAGMA_ROOT/orc8r/protos
COPY orc8r/swagger $MAGMA_ROOT/orc8r/swagger
COPY orc8r/tools/ansible/roles/fluent_bit/files $MAGMA_ROOT/orc8r/tools/ansible/roles/fluent_bit/files
COPY protos $MAGMA_ROOT/protos
ENV PROTO_LIST orc8r_protos lte_protos feg_protos cwf_protos
RUN make -C $MAGMA_ROOT/orc8r/gateway/python protos
ENV SWAGGER_LIST lte_swagger_specs feg_swagger_specs cwf_swagger_specs orc8r_swagger_specs
RUN make -C $MAGMA_ROOT/orc8r/gateway/python swagger

# -----------------------------------------------------------------------------
# Production image
# -----------------------------------------------------------------------------
FROM ubuntu:focal AS gateway_python
ARG MAGMA_BUILD_BRANCH=unknown
ARG MAGMA_BUILD_TAG=unknown
ARG MAGMA_BUILD_COMMIT_HASH=unknonw
ARG MAGMA_BUILD_COMMIT_DATE=unknown

# Add the magma apt repo
RUN apt-get update && \
    apt-get install -y apt-utils software-properties-common apt-transport-https curl
COPY orc8r/tools/ansible/roles/pkgrepo/files/jfrog.pub /tmp/jfrog.pub
RUN apt-key add /tmp/jfrog.pub && \
    apt-add-repository "deb https://artifactory.magmacore.org/artifactory/debian-test focal-ci main"
RUN curl -L http://packages.fluentbit.io/fluentbit.key > /tmp/fluentbit.key
RUN apt-key add /tmp/fluentbit.key && \
    apt-add-repository "deb https://packages.fluentbit.io/ubuntu/focal focal main"

# Install the runtime deps from apt.
RUN apt-get -y update && apt-get -y install \
    iproute2 \
    libc-ares2 \
    libev4 \
    libffi-dev \
    libjansson4 \
    libjemalloc2 \
    libssl-dev \
    libsystemd-dev \
    magma-nghttpx=1.31.1-1 \
    net-tools \
    openssl \
    iputils-ping \
    pkg-config \
    python-cffi \
    python3-pip \
    python3.8 \
    python3.8-dev \
    redis-server \
    git \
    netcat \
    td-agent-bit

# Install docker.
RUN curl -sSL https://get.docker.com/ > /tmp/get_docker.sh && \
    sh /tmp/get_docker.sh && \
    rm /tmp/get_docker.sh

# Install docker-compose for upgrades
RUN curl -L "https://github.com/docker/compose/releases/download/1.25.0-rc1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/bin/docker-compose
RUN chmod 755 /usr/bin/docker-compose


# Install python code.
COPY orc8r/gateway/python /tmp/orc8r
RUN python3.8 -m pip install --no-cache-dir /tmp/orc8r

# update aioh2 since there is no pushed package
RUN pip3 install --force-reinstall git+https://github.com/URenko/aioh2.git

# Copy the build artifacts.
COPY --from=builder /build/python/gen /usr/local/lib/python3.8/dist-packages/

# Copy the configs.
COPY feg/gateway/configs /etc/magma

COPY orc8r/gateway/configs/templates /etc/magma/templates

RUN mkdir -p /var/opt/magma/configs
RUN mkdir -p /var/opt/magma/fluent-bit

# Add commit information
ENV MAGMA_BUILD_BRANCH $MAGMA_BUILD_BRANCH
ENV MAGMA_BUILD_TAG $MAGMA_BUILD_TAG
ENV MAGMA_BUILD_COMMIT_HASH $MAGMA_BUILD_COMMIT_HASH
ENV MAGMA_BUILD_COMMIT_DATE $MAGMA_BUILD_COMMIT_DATE
