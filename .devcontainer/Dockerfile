FROM ubuntu:focal as magma-mme-builder

# [Option] Install zsh
ARG INSTALL_ZSH="true"
# [Option] Upgrade OS packages to their latest versions
ARG UPGRADE_PACKAGES="false"
# [Option] Enable non-root Docker access in container
ARG ENABLE_NONROOT_DOCKER="true"
# [Option] Use the OSS Moby CLI instead of the licensed Docker CLI
ARG USE_MOBY="true"

# Install needed packages and setup non-root user. Use a separate RUN statement to add your
# own dependencies. A user of "automatic" attempts to reuse an user ID if one already exists.
ARG USERNAME=automatic
ARG USER_UID=1000
ARG USER_GID=$USER_UID
COPY .devcontainer/library-scripts/*.sh /tmp/library-scripts/
RUN apt-get update \
    && /bin/bash /tmp/library-scripts/common-debian.sh "${INSTALL_ZSH}" "${USERNAME}" "${USER_UID}" "${USER_GID}" "${UPGRADE_PACKAGES}" \
    # Use Docker script from script library to set things up
    && /bin/bash /tmp/library-scripts/docker-debian.sh "${ENABLE_NONROOT_DOCKER}" "/var/run/docker-host.sock" "/var/run/docker.sock" "${USERNAME}" "${USE_MOBY}" \
    # Clean up
    && apt-get autoremove -y && apt-get clean -y && rm -rf /var/lib/apt/lists/* /tmp/library-scripts/

# Setting the ENTRYPOINT to docker-init.sh will configure non-root access to
# the Docker socket if "overrideCommand": false is set in devcontainer.json.
# The script will also execute CMD if you need to alter startup behaviors.
ENTRYPOINT [ "/usr/local/share/docker-init.sh" ]
CMD [ "sleep", "infinity" ]

ARG GIT_PROXY
ARG FEATURES=mme_oai
ARG PYTHON_VERSION=3.8
# PYTHON_VERSION must be in sync with "python.defaultInterpreterPath" and "python.analysis.extraPaths" in .devcontainer/devcontainer.json
ENV MAGMA_ROOT=/workspaces/magma
ENV BUILD_TYPE=Debug
ENV C_BUILD=/workspaces/magma/build/c
ENV DEBIAN_FRONTEND=noninteractive


RUN echo "Install general purpose packages" && \
    apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get upgrade -y && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y \
        apt-transport-https \
        apt-utils \
        autoconf \
        autoconf automake \
        automake \
        build-essential \
        ca-certificates \
        clang-11 \
        clang-format-11 \
        clang-tidy-11 \
        curl \
        g++ \
        g++-9 \
        gcc-9 \
        # iputils-ping is a dependency of python services (e.g. magmad)
        iputils-ping \
        lcov \
        gdb \
        git \
        gnupg \
        golang \
        libclang-11-dev \
        libgmp3-dev \
        libpcap-dev \
        libssl-dev \
        libtool \
        lld \
        lldb \
        make \
        netbase \
        ninja-build \
        openjdk-8-jdk \
        perl \
        pkg-config \
        python${PYTHON_VERSION} \
        python3-pip \
        python-is-python3 \
        redis-server \
        ruby \
        rubygems \
        ruby-dev \
        software-properties-common \
        tzdata \
        unzip \
        vim \
        virtualenv=20.0.17-1ubuntu0.4 \
        wget && \
        gem install fpm && \
        update-alternatives --install /usr/bin/clang clang /usr/lib/llvm-11/bin/clang 10 && \
        update-alternatives --install /usr/bin/clang++ clang++ /usr/lib/llvm-11/bin/clang++ 10 && \
        update-alternatives --install /usr/bin/clang-format clang-format /usr/lib/llvm-11/bin/clang-format 10 && \
        update-alternatives --install /usr/bin/clang-tidy clang-tidy /usr/lib/llvm-11/bin/clang-tidy 10 && \
        update-alternatives --install /usr/bin/clang-apply-replacements clang-apply-replacements /usr/lib/llvm-11/bin/clang-apply-replacements 10

RUN DEBIAN_FRONTEND=noninteractive apt-get install -y clangd-12

# Install bazel
WORKDIR /usr/sbin
RUN wget --progress=dot:giga https://github.com/bazelbuild/bazelisk/releases/download/v1.10.0/bazelisk-linux-amd64 && \
    chmod +x bazelisk-linux-amd64 && \
    ln -s /usr/sbin/bazelisk-linux-amd64 /usr/sbin/bazel

RUN echo "Install 3rd party dependencies" && \
    apt-get update && \
    echo "Install CMake" && \
    apt-get -y install cmake && \
    echo "Install FMT lib requirements" && \
    apt-get -y install libunwind8-dev libelf-dev libdwarf-dev bzip2 && \
    echo "Install Folly requirements" && \
    apt-get -y install libboost-all-dev libevent-dev libdouble-conversion-dev \
    libgoogle-glog-dev libgflags-dev libiberty-dev liblz4-dev liblzma-dev \
    libsnappy-dev binutils-dev libjemalloc-dev libssl-dev pkg-config libunwind-dev && \
    echo "Install check for test support" && \
    apt-get -y install check && \
    echo "Install gtest for test support" && \
    apt-get -y install libgtest-dev && \
    echo "Install FreeDiameter requirements" && \
    apt-get -y install libsctp1 libsctp-dev libgcrypt-dev \
    bison flex libidn11-dev && \
    echo "Install libgtpnl requirements" && \
    apt-get -y install libmnl-dev && \
    echo "Install Nettle requirements" && \
    apt-get install -y libgoogle-glog-dev libconfig-dev libxml2-dev \
    libyaml-cpp-dev nlohmann-json3-dev && \
    echo "Install ZeroMQ" && \
    apt-get install -y libczmq-dev && \
    echo "Install dependencies for Sentry Native" && \
    apt-get install -y libcurl4-openssl-dev

RUN ["/bin/bash", "-c", "if [[ -v GIT_PROXY ]]; then git config --global http.proxy $GIT_PROXY; fi"]

##### NETTLE and GNUTLS
# TODO Upgrade these - requires us to update our use of libnettle due to API migration.
#  see https://gist.github.com/electronjoe/a899e4bfbc2904cb353444386296c38e
# Note the CFLAGS define below due to glibc deprecation of critical flag,
#  see https://github.com/rdslw/openwrt/blob/e5d47f32131849a69a9267de51a30d6be1f0d0ac/tools/bison/patches/110-glibc-change-work-around.patch
RUN wget --progress=dot:giga https://ftp.gnu.org/gnu/nettle/nettle-2.5.tar.gz && \
    tar -xf nettle-2.5.tar.gz && \
    cd nettle-2.5 && \
    mkdir build && \
    cd build/ && \
    ../configure --disable-openssl --enable-shared --libdir=/usr/lib && \
    make -j"$(nproc)" && \
    make install && \
    ldconfig -v && \
    cd / && \
    wget --progress=dot:giga https://www.gnupg.org/ftp/gcrypt/gnutls/v3.1/gnutls-3.1.23.tar.xz && \
    tar xf gnutls-3.1.23.tar.xz && \
    cd gnutls-3.1.23 && \
    mkdir build && cd build && \
    CFLAGS=-D_IO_ftrylockfile ../configure --libdir=/usr/lib --with-libnettle-prefix=/usr && \
    make -j"$(nproc)" && \
    make install && \
    ldconfig -v && \
    cd / && \
    rm -rf nettle* && \
    rm -rf gnutls*

##### Useful for logfile modification e.g. pruning all /magma/... prefix from GCC warning logs
RUN GOBIN="/usr/bin/" go get github.com/ezekg/xo

##### GRPC and it's dependencies
RUN git clone --recurse-submodules -b v1.35.0 https://github.com/grpc/grpc && \
    cd grpc && \
    mkdir -p cmake/build && \
    cd cmake/build && \
    cmake -DgRPC_INSTALL=ON -DgRPC_BUILD_TESTS=OFF -DBUILD_SHARED_LIBS=ON ../.. && \
    make -j"$(nproc)" && \
    make install && \
    cd / && \
    rm -rf grpc

##### libprotobuf-mutator is used for randomized proto unit tests / property tests
RUN git clone -b v1.0 https://github.com/google/libprotobuf-mutator && \
    mkdir -p libprotobuf-mutator/build && \
    cd libprotobuf-mutator/build && \
    cmake .. -GNinja -DCMAKE_C_COMPILER=gcc -DCMAKE_CXX_COMPILER=g++ -DCMAKE_BUILD_TYPE=Debug && \
    ninja && \
    ninja install && \
    cd / && \
    rm -rf libprotobuf-mutator

##### Prometheus CPP
RUN git clone https://github.com/jupp0r/prometheus-cpp.git && \
    cd prometheus-cpp && \
    git checkout d8326b2bba945a435f299e7526c403d7a1f68c1f && \
    git submodule init && git submodule update && \
    mkdir _build && \
    cd _build/ && \
    cmake .. && \
    make -j"$(nproc)" && \
    make install && \
    rm -rf /prometheus-cpp

##### Redis CPP
RUN git clone https://github.com/cpp-redis/cpp_redis.git && \
    cd cpp_redis && \
    git checkout bbe38a7f83de943ffcc90271092d689ae02b3489 && \
    git submodule init && git submodule update && \
    mkdir build && cd build && \
    cmake .. -DCMAKE_BUILD_TYPE=Release && \
    make -j"$(nproc)" && \
    make install && \
    rm -rf /cpp_redis

##### liblfds
# https://www.liblfds.org/mediawiki/index.php?title=r7.1.0:Building_Guide_(liblfds)
RUN git clone https://github.com/liblfds/liblfds.git && \
    cd liblfds/liblfds/liblfds7.1.0/liblfds710/build/gcc_gnumake/ && \
    make -j"$(nproc)" && \
    make ar_install && \
    cd / && \
    rm -rf liblfds

##### libgtpnl
# review https://github.com/OPENAIRINTERFACE/openair-cn/blob/master/build/tools/build_helper.gtpnl
RUN git clone https://git.osmocom.org/libgtpnl && \
    cd libgtpnl && \
    git reset --hard 345d687 && \
    autoreconf -fi && \
    ./configure && \
    make -j"$(nproc)" && \
    make install && \
    ldconfig && \
    cd / && \
    rm -rf libgtpnl

#####  asn1c
RUN git clone https://gitlab.eurecom.fr/oai/asn1c.git && \
    cd asn1c && \
    git checkout f12568d617dbf48497588f8e227d70388fa217c9 && \
    autoreconf -iv && \
    ./configure && \
    make -j"$(nproc)" && \
    make install && \
    ldconfig

## Install Fmt (Folly Dep)
RUN git clone https://github.com/fmtlib/fmt.git && cd fmt && \
    mkdir _build && cd _build && \
    cmake -DBUILD_SHARED_LIBS=ON .. && \
    make -j"$(nproc)" && \
    make install && \
    cd / && \
    rm -rf fmt

##### Facebook Folly C++ lib
# Note: "Because folly does not provide any ABI compatibility guarantees from
#        commit to commit, we generally recommend building folly as a static library."
# Here we checkout the hash for v2021.02.22.00 (arbitrary recent version)
RUN git clone https://github.com/facebook/folly && cd folly && \
    git checkout tags/v2021.02.15.00 && \
    mkdir _build && cd _build && \
    cmake -DBUILD_SHARED_LIBS=ON .. && \
    make -j"$(nproc)" && \
    make install && \
    cd / && \
    rm -rf folly

##### Build and install libgtest and gmock
RUN cd /usr/src/googletest && \
    mkdir build && \
    cd build && \
    cmake -DBUILD_SHARED_LIBS=ON .. && \
    echo "Build gtest and gmock" && \
    make && \
    echo "Install gtest and gmock" && \
    make install && \
    ldconfig -v

#### libtins is required to build the connection_tracker
# TODO(@themarwhal): this might not be hard to bazelify if we can generate the config
RUN git clone --branch v4.3 https://github.com/mfontanini/libtins.git && \
    cd libtins && \
    mkdir build && \
    cd build && \
    cmake ../ -DLIBTINS_ENABLE_CXX11=1 && \
    make -j"$(nproc)" && \
    make install && \
    cd / && \
    rm -rf libtins && \
    # symlink to /usr/lib to match Magma VM
    ln -s /usr/local/lib/libtins.so /usr/lib/libtins.so

###### Install Include What You Use for c/cpp header include fixup tooling
# Tag 0.15 tracks Clang 11.0 per https://github.com/include-what-you-use/include-what-you-use/tags
RUN git clone https://github.com/include-what-you-use/include-what-you-use && \
    cd include-what-you-use && \
    git checkout 0.15 && \
    cd .. && \
    mkdir build_iwyu && cd build_iwyu && \
    cmake -G "Unix Makefiles" -DIWHYU_LLVM_ROOT_PATH=/usr/lib/llvm-11 ../include-what-you-use/ && \
    make && \
    make install && \
    cd / && \
    rm -rf include-what-you-use && \
    rm -rf build_liwyu

#### libfluid is requird to build MME with OVS support
COPY third_party /tmp/third_party/
RUN /tmp/third_party/build/bin/libfluid_build.sh && \
    find . -name magma-libfluid_0.1* -exec dpkg -i {} \; && \
    rm -rf /tmp/*

##### FreeDiameter
COPY lte/gateway/c/core/oai/patches/ /tmp/
RUN git clone https://github.com/OPENAIRINTERFACE/opencord.org.freeDiameter.git freediameter && \
    cd freediameter && \
    patch -p1 < /tmp/0001-opencoord.org.freeDiameter.patch && \
    mkdir build && \
    cd build && \
    cmake -DDISABLE_SCTP:BOOL=ON ../ && \
    make -j"$(nproc)" && \
    make install && \
    rm -rf /freediameter /tmp/*

##### Go language server support for vscode
RUN go get -v golang.org/x/tools/gopls

#### Update shared library configuration
RUN ldconfig -v


##### Install Python requirements

### create virtualenv
ARG PYTHON_VENV=/home/vscode/build/python
ENV PYTHON_VENV_EXECUTABLE=${PYTHON_VENV}/bin/python${PYTHON_VERSION}
# PYTHON_VENV must by in sync with "python.defaultInterpreterPath", "python.analysis.extraPaths" and magtivate path in "postCreateCommand" in .devcontainer/devcontainer.json

RUN virtualenv --system-site-packages --python=/usr/bin/python${PYTHON_VERSION} ${PYTHON_VENV}
RUN ${PYTHON_VENV_EXECUTABLE} -m pip install --quiet --upgrade --no-cache-dir "setuptools==49.6.0"

### install python3-aioeventlet and bcc-tools from the magma apt repo
COPY orc8r/tools/ansible/roles/pkgrepo/files/jfrog.pub /tmp/
RUN apt-key add /tmp/jfrog.pub && \
    echo "deb https://facebookconnectivity.jfrog.io/artifactory/list/dev-focal/ focal main" > /etc/apt/sources.list.d/fbc.list && \
    apt-get update -y && \
    apt-get install -y python3-aioeventlet bcc-tools && \
    rm -f /tmp/* && \
    rm -rf /var/lib/apt/lists/*  # delete lists that are downloaded by apt-get update

# add patch that is missing in jfrog version of aioeventlet (it only comes with 1 of 2 patches)
COPY lte/gateway/deploy/roles/magma/files/patches/aioeventlet_fd_exception.patch /tmp/
RUN patch -N -s -f /usr/local/lib/python${PYTHON_VERSION}/dist-packages/aioeventlet.py < /tmp/aioeventlet_fd_exception.patch && \
    rm -rf /tmp/*

### install eggs (lte, orc8r)
COPY /lte/gateway/python/ ${MAGMA_ROOT}/lte/gateway/python/
WORKDIR ${MAGMA_ROOT}/lte/gateway/python/
RUN ${PYTHON_VENV_EXECUTABLE} -m pip install --quiet --upgrade --no-build-isolation --no-cache-dir --verbose --editable .[dev] && \
    rm -rf lte.egg-info

COPY /orc8r/gateway/python/ ${MAGMA_ROOT}/orc8r/gateway/python/
WORKDIR ${MAGMA_ROOT}/orc8r/gateway/python/
RUN ${PYTHON_VENV_EXECUTABLE} -m pip install --quiet --upgrade --no-build-isolation --no-cache-dir --verbose --editable .[dev] && \
    rm -rf orc8r.egg-info

### install formatter autopep8
RUN ${PYTHON_VENV_EXECUTABLE} -m pip install --no-cache-dir autopep8

#### protos
ARG GEN_DIR=lib/python${PYTHON_VERSION}/gen

COPY /protos/ ${MAGMA_ROOT}/protos/
COPY /lte/protos/ ${MAGMA_ROOT}/lte/protos/
COPY /orc8r/protos/ ${MAGMA_ROOT}/orc8r/protos/
COPY /feg/protos/ ${MAGMA_ROOT}/feg/protos/
COPY /dp/protos/ ${MAGMA_ROOT}/dp/protos/
WORKDIR ${MAGMA_ROOT}
RUN ${PYTHON_VENV_EXECUTABLE} -m pip install --no-cache-dir "mypy-protobuf==2.4" && \
    mkdir ${PYTHON_VENV}/${GEN_DIR} && \
    for PROTO_SRC in orc8r lte feg dp; \
    do \
    ${PYTHON_VENV_EXECUTABLE} protos/gen_protos.py ${PROTO_SRC}/protos ${MAGMA_ROOT},orc8r/protos/prometheus ${MAGMA_ROOT} ${PYTHON_VENV}/${GEN_DIR} && \
    ${PYTHON_VENV_EXECUTABLE} protos/gen_prometheus_proto.py ${MAGMA_ROOT} ${PYTHON_VENV}/${GEN_DIR}; \
    done && \
    echo "${PYTHON_VENV}/${GEN_DIR}" > ${PYTHON_VENV}/lib/python${PYTHON_VERSION}/site-packages/magma_gen.pth

### swagger
ENV SWAGGER_CODEGEN_JAR=/var/tmp/codegen/modules/swagger-codegen-cli/target/swagger-codegen-cli.jar
ARG MVN_VERSION=3.5.4
ARG CODEGEN_VERSION=v2.2.3

WORKDIR /var/tmp/
RUN wget --progress=dot:giga http://mirrors.ibiblio.org/apache/maven/maven-3/${MVN_VERSION}/binaries/apache-maven-${MVN_VERSION}-bin.tar.gz && \
    tar --extract --verbose --file apache-maven-${MVN_VERSION}-bin.tar.gz

WORKDIR /var/tmp/codegen
RUN git clone --branch ${CODEGEN_VERSION} https://github.com/swagger-api/swagger-codegen /var/tmp/codegen
RUN /var/tmp/apache-maven-${MVN_VERSION}/bin/mvn clean package -D skipTests

# Copy swagger specs over to the build directory,
# so that eventd can access them at runtime
COPY lte/swagger/*.yml ${PYTHON_VENV}/${GEN_DIR}/lte/swagger/specs/
COPY orc8r/swagger/*.yml ${PYTHON_VENV}/${GEN_DIR}/orc8r/swagger/specs/
RUN for SWAGGER_SRC in lte orc8r; \
    do \
    # Generate the files
    ls ${PYTHON_VENV}/${GEN_DIR}/${SWAGGER_SRC}/swagger/specs/*.yml \
    | xargs -t -I% /usr/bin/java -jar ${SWAGGER_CODEGEN_JAR} generate \
    -i % \
    -o ${PYTHON_VENV}/${GEN_DIR}/${SWAGGER_SRC}/swagger \
    -l python \
    -D models && \
    # Flatten and clean up directory
    mv ${PYTHON_VENV}/${GEN_DIR}/${SWAGGER_SRC}/swagger/swagger_client/* ${PYTHON_VENV}/${GEN_DIR}/${SWAGGER_SRC}/swagger/ && \
    rmdir ${PYTHON_VENV}/${GEN_DIR}/${SWAGGER_SRC}/swagger/swagger_client && \
    rm -r ${PYTHON_VENV}/${GEN_DIR}/${SWAGGER_SRC}/swagger/test; \
    done

RUN ln -s $MAGMA_ROOT/bazel/bazelrcs/devcontainer.bazelrc /etc/bazelrc

WORKDIR $MAGMA_ROOT
