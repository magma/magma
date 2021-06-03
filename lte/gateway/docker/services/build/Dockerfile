FROM ubuntu:20.04 AS builder

ENV DEBIAN_FRONTEND=noninteractive
ENV MAGMA_ROOT=/magma
ENV ARCH=amd64

RUN apt-get update \
  && apt-get install -y python3-pip python3-dev git cmake ninja-build \
  protobuf-c-compiler protobuf-compiler libprotobuf-dev libyaml-dev \
  libyaml-cpp-dev wget ruby ruby-dev autoconf libtool pkg-config \
  libgflags-dev libgoogle-glog-dev libssl-dev libboost-all-dev libghc-double-conversion-dev \
  libsctp-dev libboost-all-dev libevent-dev libdouble-conversion-dev \
  libgoogle-glog-dev libgflags-dev libiberty-dev liblz4-dev liblzma-dev \
  libsctp1 libgcrypt-dev bison libidn11-dev flex libtspi1 libtcmalloc-minimal4 libgoogle-perftools4

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
COPY services/build/gnutls-build.patch /tmp
RUN patch -p4 < /tmp/gnutls-build.patch
COPY services/build/gnutls-glibc.patch /tmp

RUN ./gnutls_build.sh

RUN dpkg -i oai-gnutls_3.1.23-1_amd64.deb

RUN ./prometheus-cpp_build.sh

RUN ./freediameter_build.sh # Patch needed?

RUN mkdir /tmp/packages && \
    mv *.deb /tmp/packages

WORKDIR /tmp/packages

RUN ./build-magma.sh
