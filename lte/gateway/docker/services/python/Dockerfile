FROM ubuntu:20.04 as builder

ENV PYTHON_BUILD=/build
ENV MAGMA_ROOT=/magma
ENV PIP_CACHE_HOME="~/.pipcache"
ENV SWAGGER_CODEGEN_JAR=/var/tmp/codegen/modules/swagger-codegen-cli/target/swagger-codegen-cli.jar

ARG MVN_VERSION=3.5.4
ARG MAGMA_VERSION=master
ARG CODEGEN_VERSION=v2.2.3


RUN apt-get update \
  && DEBIAN_FRONTEND=noninteractive apt-get install -y python3-pip python3-dev python3-eventlet python3-pystemd python3-protobuf git sudo virtualenv lsb-release openjdk-8-jre-headless openjdk-8-jdk wget ruby ruby-dev pkg-config libsystemd-dev libprotobuf-dev\
  && cd /usr/local/bin \
  && ln -s /usr/bin/python3 python

RUN gem install fpm

WORKDIR /

RUN git clone https://github.com/magma/magma.git && \
    cd magma && \
    git checkout $MAGMA_VERSION

RUN git clone https://github.com/swagger-api/swagger-codegen /var/tmp/codegen && \
    cd /var/tmp/codegen && \
    git checkout $CODEGEN_VERSION

WORKDIR /var/tmp/

RUN /magma/third_party/build/bin/aioeventlet_build.sh && \
    dpkg -i python3-aioeventlet*

RUN wget http://mirrors.ibiblio.org/apache/maven/maven-3/$MVN_VERSION/binaries/apache-maven-$MVN_VERSION-bin.tar.gz && \
    tar -xvf apache-maven-$MVN_VERSION-bin.tar.gz

WORKDIR /var/tmp/codegen

RUN ../apache-maven-$MVN_VERSION/bin/mvn clean package -DskipTests

WORKDIR /magma/lte/gateway/python

RUN make buildenv

FROM ubuntu:20.04 AS gateway_python

ENV VIRTUAL_ENV=/build
ARG JSONPOINTER_VERSION=1.13

RUN apt-get update \
  && DEBIAN_FRONTEND=noninteractive apt-get install -y python3-pip python3-venv virtualenv git python3-eventlet python3-pystemd python3-jinja2 python3-pystemd nghttp2-proxy net-tools inetutils-ping redis-server wget ethtool\
  && cd /usr/local/bin \
  && ln -s /usr/bin/python3 python
RUN python3 -m venv $VIRTUAL_ENV
ENV PATH="/magma/orc8r/gateway/python/scripts/:/magma/lte/gateway/python/scripts/:$VIRTUAL_ENV/bin:$PATH"

RUN echo "deb https://packages.fluentbit.io/ubuntu/focal focal main" > /etc/apt/sources.list.d/tda.list
RUN wget -qO - https://packages.fluentbit.io/fluentbit.key | apt-key add -
RUN DEBIAN_FRONTEND=noninteractive apt-get  update && apt-get install -y td-agent-bit

COPY --from=builder /build /build
COPY --from=builder /magma /magma
COPY --from=builder /magma/orc8r/gateway/python/scripts/ /usr/local/bin
COPY --from=builder /magma/lte/gateway/python/scripts/ /usr/local/bin
COPY --from=builder /var/tmp/python3-aioeventlet* /var/tmp/
COPY --from=builder /magma/lte/gateway/configs/templates /etc/magma/templates/
COPY --from=builder /magma/orc8r/gateway/configs/templates/nghttpx.conf.template /etc/magma/templates/nghttpx.conf.template
COPY --from=builder /magma/orc8r/gateway/python/scripts/generate_nghttpx_config.py /usr/local/bin/generate_nghttpx_config.py
COPY --from=builder /magma/orc8r/gateway/python/scripts/generate_service_config.py /usr/local/bin/generate_service_config.py
COPY --from=builder /magma/orc8r/gateway/python/scripts/generate_fluent_bit_config.py /usr/local/bin/generate_fluent_bit_config.py
COPY --from=builder /magma/lte/gateway/deploy/roles/magma/files/set_irq_affinity /usr/local/bin/set_irq_affinity

RUN chmod -R +x /usr/local/bin/generate*

RUN chmod +x /usr/local/bin/set_irq_affinity

RUN dpkg -i /var/tmp/python3-aioeventlet*

RUN mkdir -p /var/opt/magma/

RUN pip install jsonpointer>$JSONPOINTER_VERSION
