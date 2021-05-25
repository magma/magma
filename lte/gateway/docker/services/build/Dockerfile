FROM ubuntu:20.04 AS builder

ENV DEBIAN_FRONTEND=noninteractive
ENV MAGMA_ROOT=/magma
ENV ARCH=amd64
ENV OS_RELEASE=ubuntu20.04

RUN apt-get update \
  && apt-get install -y python3-pip python3-dev git cmake ninja-build \
  protobuf-c-compiler protobuf-compiler libprotobuf-dev libyaml-dev \
  libyaml-cpp-dev wget ruby ruby-dev autoconf libtool pkg-config \
  libgflags-dev libgoogle-glog-dev libssl-dev libboost-all-dev libghc-double-conversion-dev \
  libsctp-dev libboost-all-dev libevent-dev libdouble-conversion-dev \
  libgoogle-glog-dev libgflags-dev libiberty-dev liblz4-dev liblzma-dev \
  libsctp1 libgcrypt-dev bison libidn11-dev flex libtspi1 libtcmalloc-minimal4 \
  libgoogle-perftools4 sudo libjemalloc2 libsnappy1v5 nlohmann-json3-dev libconfig9 \
  libconfig-dev libxml2-dev libczmq-dev

RUN gem install fpm

WORKDIR /

RUN git clone https://github.com/magma/magma.git && \
    cd magma && \
    git checkout master # TODO, how do we handle versioning for this?

WORKDIR /magma/third_party/build/bin/

RUN ./nettle_build.sh

RUN dpkg -i /magma/third_party/build/bin/oai-nettle_2.5-1ubuntu20.04_amd64.deb

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

RUN dpkg -i oai-gnutls_3.1.23-1ubuntu20.04_amd64.deb

RUN ./prometheus-cpp_build.sh

COPY services/build/freediameter.patch /tmp
RUN patch -p4 < /tmp/freediameter.patch
RUN ./freediameter_build.sh # Patch needed?

RUN mkdir /tmp/packages && \
    mv *.deb /tmp/packages

WORKDIR /tmp/packages

RUN dpkg -i grpc-dev_1.38.0-3ubuntu20.04_amd64.deb libfolly-dev_2018.02.26.00-6ubuntu20.04_amd64.deb \
    liblfds710_7.1.0-0ubuntu20.04_amd64.deb magma-libtacopie_3.2.0.1-1ubuntu20.04_amd64.deb \
    magma-libfluid_0.1.0.6-1ubuntu20.04_amd64.deb oai-asn1c_0~20190423+c0~rf12568d6-0ubuntu20.04_amd64.deb \
    oai-freediameter_0.0.1-1ubuntu20.04_amd64.deb prometheus-cpp-dev_1.0.2-d8326b2ubuntu20.04_amd64.deb

RUN dpkg --force overwrite -i magma-cpp-redis_4.3.1.1-2ubuntu20.04_amd64.deb

WORKDIR /magma/lte/gateway/release/

# See https://github.com/magma/magma/blob/master/lte/gateway/docker/mme/Dockerfile.ubuntu20.04#L91-L92
RUN ln -s /usr/include/nlohmann/json.hpp /usr/include/json.hpp

RUN git clone https://github.com/abseil/abseil-cpp.git && mkdir abseil-cpp/build && cd abseil-cpp/build && \
    cmake -DABSL_USE_GOOGLETEST_HEAD=ON -DCMAKE_CXX_STANDARD=11 .. && \
    make && make install


RUN ./build-magma.sh --os ubuntu
