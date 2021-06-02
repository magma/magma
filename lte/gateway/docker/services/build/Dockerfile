FROM ubuntu:20.04 AS builder

ENV DEBIAN_FRONTEND=noninteractive
ENV MAGMA_ROOT=/magma
ENV ARCH=amd64

RUN apt-get update \
  && apt-get install -y python3-pip python3-dev git cmake ninja-build \
  protobuf-c-compiler protobuf-compiler libprotobuf-dev libyaml-dev \
  libyaml-cpp-dev wget ruby ruby-dev autoconf libtool pkg-config \
  libgflags-dev libgoogle-glog-dev libssl-dev libboost-all-dev libghc-double-conversion-dev nettle-dev

RUN gem install fpm

WORKDIR /

RUN git clone https://github.com/magma/magma.git && \
    cd magma && \
    git checkout master # TODO, how do we handle versioning for this?

WORKDIR /magma/third_party/build/bin/

RUN ./nettle_build.sh

RUN dpkg -i /magma/third_party/build/bin/oai-nettle_2.5-1_amd64.deb

RUN ./asn1c_build.sh

RUN ./folly_build.sh

RUN ./magma-libtacopie_build.sh

RUN ./magma-cpp-redis_build.sh

RUN ./liblfds_build.sh

RUN ./libfluid_build.sh

# temp patch grpc version
WORKDIR /magma
COPY services/build/grpc.patch /tmp
RUN patch -p1 < /tmp/grpc.patch

WORKDIR /magma/third_party/build/bin/

RUN ./grpc_build.sh

#temp patch gnutls
COPY services/build/gnutls.patch /tmp
WORKDIR /tmp/build-oai-gnutls/gnutls-3.1.23/gl/

RUN ./gnutls_build.sh

RUN ./prometheus-cpp_build.sh

RUN ./freediameter_build.sh # Patch needed?

RUN ./build-magma.sh
