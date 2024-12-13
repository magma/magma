# Copyright 2023 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

################################################################
# Devcontainer Image (for local development and CI)
################################################################
FROM ghcr.io/magma/magma/bazel-base:latest AS devcontainer

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
        clang-12 \
        clang-format-12 \
        clang-tidy-12 \
        clangd-12 \
        g++-9 \
        gcc-9 \
        gdb \
        lcov \
        libclang-12-dev \
        lldb-12 \
        llvm-12-dev \
        make \
        ninja-build \
        openjdk-8-jdk \
        perl \
        pkg-config \
        python3-pip \
        python3-venv \
        redis-server \
        ruby \
        rubygems \
        ruby-dev \
        software-properties-common \
        tzdata \
        virtualenv=20.0.17-1ubuntu0.4
RUN gem install dotenv -v 2.8.1 && \
    gem install fpm

# Install golang
WORKDIR /usr/local
ARG GOLANG_VERSION="1.20.1"
RUN GO_TARBALL="go${GOLANG_VERSION}.linux-amd64.tar.gz" \
    && curl https://linuxfoundation.jfrog.io/artifactory/magma-blob/${GO_TARBALL} --remote-name --location \
    && tar -xzf ${GO_TARBALL} \
    && rm ${GO_TARBALL}
ENV PATH=${PATH}:/usr/local/go/bin

# /home/vscode/go/bin doesn't exist initially, but for example orc8r/cloud/go/Makefile
# populates that folder and expects those binaries to be in PATH
ENV PATH=${PATH}:/home/vscode/go/bin

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
        grpc-dev \
        libprotobuf-dev \
        libprotoc-dev \
        protobuf-compiler \
        prometheus-cpp-dev

##### Useful for logfile modification e.g. pruning all /magma/... prefix from GCC warning logs
RUN GOBIN="/usr/bin/" go install github.com/ezekg/xo@0f7f076932dd && \
    rm --recursive --interactive=never /root/.cache/go-build

# install magma dependencies
RUN apt-get install -y --no-install-recommends \
        libtins-dev \
        magma-cpp-redis \
        python3-aioeventlet

##### libgtpnl
# review https://github.com/OPENAIRINTERFACE/openair-cn/blob/master/build/tools/build_helper.gtpnl
RUN git clone https://gitea.osmocom.org/cellular-infrastructure/libgtpnl && \
    cd libgtpnl && \
    git reset --hard 345d687 && \
    autoreconf -fi && \
    ./configure && \
    make -j"$(nproc)" && \
    make install && \
    ldconfig && \
    cd .. && \
    rm --recursive --interactive=never libgtpnl

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
# Tag 0.15 tracks Clang 12.0 per https://github.com/include-what-you-use/include-what-you-use/tags
RUN git clone https://github.com/include-what-you-use/include-what-you-use && \
    cd include-what-you-use && \
    git checkout 0.15 && \
    cd .. && \
    mkdir build_iwyu && cd build_iwyu && \
    cmake -G "Unix Makefiles" -DCMAKE_PREFIX_PATH=/usr/lib/llvm-12 ../include-what-you-use/ && \
    make && \
    make install && \
    cd .. && \
    rm --recursive --interactive=never include-what-you-use build_iwyu

##### Go language server support for vscode
RUN GOBIN="/usr/bin/" go install -v golang.org/x/tools/gopls@v0.8.3 && \
    rm --recursive --interactive=never /root/.cache/go-build

#### Update shared library configuration
RUN ldconfig -v

WORKDIR $MAGMA_ROOT
