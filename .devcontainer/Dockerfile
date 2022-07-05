################################################################
# Devcontainer Image (for local development and CI)
################################################################
FROM ghcr.io/magma/magma/bazel-base:latest as devcontainer

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

ARG PYTHON_VERSION=3.8
# PYTHON_VERSION must be in sync with "python.defaultInterpreterPath" and "python.analysis.extraPaths" in .devcontainer/devcontainer.json
ENV MAGMA_ROOT=/workspaces/magma
ENV BUILD_TYPE=Debug
ENV C_BUILD=/workspaces/magma/build/c


RUN echo "Install general purpose packages" && \
    apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y --no-install-recommends \
        autoconf \
        automake \
        clang-11 \
        clang-format-11 \
        clang-tidy-11 \
        clangd-12 \
        g++-9 \
        gcc-9 \
        gdb \
        lcov \
        libclang-11-dev \
        lldb \
        llvm-11-dev \
        make \
        ninja-build \
        openjdk-8-jdk \
        perl \
        pkg-config \
        python3-pip \
        redis-server \
        ruby \
        rubygems \
        ruby-dev \
        software-properties-common \
        tzdata \
        virtualenv=20.0.17-1ubuntu0.4 && \
    gem install fpm && \
    update-alternatives --install /usr/bin/clang clang /usr/lib/llvm-11/bin/clang 10 && \
    update-alternatives --install /usr/bin/clang++ clang++ /usr/lib/llvm-11/bin/clang++ 10 && \
    update-alternatives --install /usr/bin/clang-format clang-format /usr/lib/llvm-11/bin/clang-format 10 && \
    update-alternatives --install /usr/bin/clang-tidy clang-tidy /usr/lib/llvm-11/bin/clang-tidy 10 && \
    update-alternatives --install /usr/bin/clang-apply-replacements clang-apply-replacements /usr/lib/llvm-11/bin/clang-apply-replacements 10

# Install golang 1.18
WORKDIR /usr/local
ARG GOLANG_VERSION="1.18"
RUN GO_TARBALL="go${GOLANG_VERSION}.linux-amd64.tar.gz" \
 && curl https://artifactory.magmacore.org/artifactory/generic/${GO_TARBALL} --remote-name --location \
 && tar -xzf ${GO_TARBALL} \
 && rm ${GO_TARBALL}
ENV PATH=$PATH:/usr/local/go/bin

RUN echo "Install 3rd party dependencies" && \
    apt-get update && \
    echo "Install FMT lib requirements" && \
    apt-get -y install --no-install-recommends \
        libunwind8-dev \
        libelf-dev \
        libdwarf-dev \
        bzip2 && \
    echo "Install Folly requirements" && \
    apt-get -y install --no-install-recommends \
        liblz4-dev \
        liblzma-dev \
        libsnappy-dev \
        binutils-dev \
        libjemalloc-dev \
        pkg-config \
        libunwind-dev && \
    echo "Install check for test support" && \
    apt-get -y install --no-install-recommends \
        check && \
    echo "Install gtest for test support" && \
    apt-get -y install --no-install-recommends \
        libgtest-dev && \
    echo "Install Nettle requirements" && \
    apt-get -y install --no-install-recommends \
        libxml2-dev \
        libyaml-cpp-dev \
        nlohmann-json3-dev && \
    echo "GRPC and it's dependencies" && \
    apt-get -y install --no-install-recommends \
        grpc-dev

##### Useful for logfile modification e.g. pruning all /magma/... prefix from GCC warning logs
RUN GOBIN="/usr/bin/" go install github.com/ezekg/xo@0f7f076932dd

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

# install magma dependencies
RUN apt-get install -y --no-install-recommends \
        libtins-dev \
        magma-cpp-redis \
        magma-libfluid \
        python3-aioeventlet

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

###### Install Include What You Use for c/cpp header include fixup tooling
# Tag 0.15 tracks Clang 11.0 per https://github.com/include-what-you-use/include-what-you-use/tags
RUN git clone https://github.com/include-what-you-use/include-what-you-use && \
    cd include-what-you-use && \
    git checkout 0.15 && \
    cd .. && \
    mkdir build_iwyu && cd build_iwyu && \
    cmake -G "Unix Makefiles" -DCMAKE_PREFIX_PATH=/usr/lib/llvm-11 ../include-what-you-use/ && \
    make && \
    make install && \
    cd / && \
    rm -rf include-what-you-use && \
    rm -rf build_liwyu

##### Go language server support for vscode
RUN GOBIN="/usr/bin/" go install -v golang.org/x/tools/gopls@v0.8.3

#### Update shared library configuration
RUN ldconfig -v


##### Install Python requirements

### create virtualenv
ARG PYTHON_VENV=/home/vscode/build/python
ENV PYTHON_VENV_EXECUTABLE=${PYTHON_VENV}/bin/python${PYTHON_VERSION}
# PYTHON_VENV must by in sync with "python.defaultInterpreterPath", "python.analysis.extraPaths" and magtivate path in "postCreateCommand" in .devcontainer/devcontainer.json

RUN virtualenv --system-site-packages --python=/usr/bin/python${PYTHON_VERSION} ${PYTHON_VENV}
RUN ${PYTHON_VENV_EXECUTABLE} -m pip install --quiet --upgrade --no-cache-dir "setuptools==49.6.0"

# add patch that is missing in jfrog version of aioeventlet (it only comes with 1 of 2 patches)
COPY lte/gateway/deploy/roles/magma/files/patches/aioeventlet_fd_exception.patch /tmp/
RUN (patch -N -s -f /usr/local/lib/python${PYTHON_VERSION}/dist-packages/aioeventlet.py < /tmp/aioeventlet_fd_exception.patch) || ( true && echo "skipping aioeventlet patch since it was already applied") && \
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
ENV SWAGGER_CODEGEN_DIR=/var/tmp/codegen
ENV SWAGGER_CODEGEN_JAR=${SWAGGER_CODEGEN_DIR}/swagger-codegen-cli.jar
ARG CODEGEN_VERSION=2.2.3

RUN mkdir -p ${SWAGGER_CODEGEN_DIR}; \
    wget --no-verbose https://repo1.maven.org/maven2/io/swagger/swagger-codegen-cli/${CODEGEN_VERSION}/swagger-codegen-cli-${CODEGEN_VERSION}.jar -O ${SWAGGER_CODEGEN_JAR}

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

WORKDIR $MAGMA_ROOT
