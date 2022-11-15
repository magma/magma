# AGW test docker private repository location
ARG REPO_LOCATION=agw-test.artifactory.magmacore.org
ARG MAGMA_ROOT=/magma
ARG BASE_VERSION=latest

FROM ${REPO_LOCATION}/agw_gateway_c AS c_builder

ARG GHZ_REPO=https://github.com/bojand/ghz

RUN apt-get update && apt-get install -y \
  curl \
  git \
  build-essential

# Install golang
WORKDIR /usr/local
ARG GOLANG_VERSION="1.19.3"
RUN GO_TARBALL="go${GOLANG_VERSION}.linux-amd64.tar.gz" \
 && curl https://artifactory.magmacore.org/artifactory/generic/${GO_TARBALL} --remote-name --location \
 && tar -xzf ${GO_TARBALL} \
 && ln -s /usr/local/go/bin/go /usr/local/bin/go \
 && rm ${GO_TARBALL}

WORKDIR ${MAGMA_ROOT}

RUN git clone --depth 1 ${GHZ_REPO}

WORKDIR ${MAGMA_ROOT}/ghz

RUN make build

FROM ${REPO_LOCATION}/agw_gateway_python as python_builder

ARG GHZ_REPO=https://github.com/bojand/ghz

RUN apt-get update && apt-get install -y \
  curl \
  git \
  build-essential

# Install golang
WORKDIR /usr/local
ARG GOLANG_VERSION="1.19.3"
RUN GO_TARBALL="go${GOLANG_VERSION}.linux-amd64.tar.gz" \
 && curl https://artifactory.magmacore.org/artifactory/generic/${GO_TARBALL} --remote-name --location \
 && tar -xzf ${GO_TARBALL} \
 && ln -s /usr/local/go/bin/go /usr/local/bin/go \
 && rm ${GO_TARBALL}

ENV MAGMA_ROOT /magma

WORKDIR ${MAGMA_ROOT}

RUN git clone --depth 1 ${GHZ_REPO}

WORKDIR ${MAGMA_ROOT}/ghz

RUN make build

FROM ${REPO_LOCATION}/agw_gateway_c as agw_c_ghz

COPY --from=c_builder /magma/ghz/dist/ghz /usr/local/bin/

WORKDIR ${MAGMA_ROOT}/lte/gateway/python/load_tests

FROM ${REPO_LOCATION}/agw_gateway_python as agw_python_ghz

COPY --from=python_builder /magma/ghz/dist/ghz /usr/local/bin/

WORKDIR ${MAGMA_ROOT}/lte/gateway/python/load_tests
