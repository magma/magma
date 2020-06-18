# -----------------------------------------------------------------------------
# Builder image to generate proto files
# -----------------------------------------------------------------------------
FROM ubuntu:bionic AS builder

# Install the runtime deps from apt.
RUN apt-get -y update && apt-get -y install curl make virtualenv zip \
 apt-utils software-properties-common apt-transport-https git openjdk-8-jdk ant

# Install protobuf compiler.
RUN curl -Lfs https://github.com/google/protobuf/releases/download/v3.1.0/protoc-3.1.0-linux-x86_64.zip -o protoc3.zip && \
  unzip protoc3.zip -d protoc3 && \
  mv protoc3/bin/protoc /usr/bin/protoc && \
  chmod a+rx /usr/bin/protoc && \
  cp -r protoc3/include/google /usr/include/ && \
  chmod -R a+Xr /usr/include/google && \
  rm -rf protoc3.zip protoc3

# Install python3.5 since it's not native on bionic
RUN apt-add-repository "ppa:deadsnakes/ppa"
RUN apt-get -y update && apt-get -y install python3.5

ENV MAGMA_ROOT /magma
ENV PYTHON_BUILD /build/python
ENV PIP_CACHE_HOME ~/.pipcache
ENV CODEGEN_ROOT /var/tmp/codegen
ENV SWAGGER_CODEGEN_JAR $CODEGEN_ROOT/modules/swagger-codegen-cli/target/swagger-codegen-cli.jar
ENV MVN_VERSION 3.5.4
ENV MVN_DIR /var/tmp/apache-maven-$MVN_VERSION
ENV M2_HOME $MVN_DIR
ENV JAVA_HOME /usr/lib/jvm/java-8-openjdk-amd64/
RUN printenv > /etc/environment

# Install swagger codegen
RUN git clone -b v2.2.3 'https://github.com/swagger-api/swagger-codegen' $CODEGEN_ROOT
RUN curl -L http://mirrors.ibiblio.org/apache/maven/maven-3/$MVN_VERSION/binaries/apache-maven-$MVN_VERSION-bin.tar.gz > $MVN_DIR-bin.tar.gz
RUN tar -xvf $MVN_DIR-bin.tar.gz -C /var/tmp
COPY lte/gateway/deploy/roles/dev_common/files/build_swagger_codegen.sh /var/tmp/build_swagger_codegen.sh
RUN ./var/tmp/build_swagger_codegen.sh

# Generate python proto bindings.
COPY feg/protos $MAGMA_ROOT/feg/protos
COPY lte/gateway/python/defs.mk $MAGMA_ROOT/lte/gateway/python/defs.mk
COPY lte/gateway/python/Makefile $MAGMA_ROOT/lte/gateway/python/Makefile
COPY lte/protos $MAGMA_ROOT/lte/protos
COPY lte/swagger $MAGMA_ROOT/lte/swagger
COPY orc8r/gateway/python/python.mk $MAGMA_ROOT/orc8r/gateway/python/python.mk
COPY orc8r/protos $MAGMA_ROOT/orc8r/protos
COPY orc8r/swagger $MAGMA_ROOT/orc8r/swagger
COPY orc8r/tools/ansible/roles/fluent_bit/files $MAGMA_ROOT/orc8r/tools/ansible/roles/fluent_bit/files
COPY protos $MAGMA_ROOT/protos
RUN make -C $MAGMA_ROOT/lte/gateway/python protos
RUN make -C $MAGMA_ROOT/lte/gateway/python swagger

# -----------------------------------------------------------------------------
# Production image
# -----------------------------------------------------------------------------
FROM ubuntu:bionic AS gateway_python

# Add the magma apt repo
RUN apt-get update && \
    apt-get install -y apt-utils software-properties-common apt-transport-https curl
COPY orc8r/tools/ansible/roles/pkgrepo/files/jfrog.pub /tmp/jfrog.pub
RUN apt-key add /tmp/jfrog.pub && \
    apt-add-repository "deb https://facebookconnectivity.jfrog.io/artifactory/list/dev/ xenial main"
RUN curl -L http://packages.fluentbit.io/fluentbit.key > /tmp/fluentbit.key
RUN apt-key add /tmp/fluentbit.key && \
    apt-add-repository "deb https://packages.fluentbit.io/ubuntu/bionic bionic main"

# Install python3.5 since it's not native on bionic
RUN apt-add-repository "ppa:deadsnakes/ppa"
RUN apt-get -y update && apt-get -y install python3.5

# Install the runtime deps from apt.
RUN apt-get -y update && apt-get -y install \
  iproute2 \
  libc-ares2 \
  libev4 \
  libffi-dev \
  libjansson4 \
  libjemalloc1 \
  libssl-dev \
  libsystemd-dev \
  magma-nghttpx=1.31.1-1 \
  net-tools \
  openssl \
  iputils-ping \
  pkg-config \
  python-cffi \
  python3-pip \
  python3.5-dev \
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
RUN python3.5 -m pip install /tmp/orc8r

# Copy the build artifacts.
COPY --from=builder /build/python/gen /usr/local/lib/python3.5/dist-packages/

# Copy the configs.
COPY feg/gateway/configs /etc/magma

COPY orc8r/gateway/configs/templates /etc/magma/templates
COPY orc8r/cloud/docker/proxy/magma_headers.rb /etc/nghttpx/magma_headers.rb

RUN mkdir -p /var/opt/magma/configs
RUN mkdir -p /var/opt/magma/fluent-bit
