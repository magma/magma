# AGW test docker private repository location
ARG REPO_LOCATION=agw-test.artifactory.magmacore.org
ARG MAGMA_ROOT=/magma
ARG BASE_VERSION=latest

FROM ${REPO_LOCATION}/agw_gateway_c AS c_builder

ARG GHZ_REPO=https://github.com/bojand/ghz

RUN apt-get update && apt-get install -y \
  git \
  golang \
  build-essential

WORKDIR ${MAGMA_ROOT}

RUN git clone --depth 1 ${GHZ_REPO}

WORKDIR ${MAGMA_ROOT}/ghz

RUN make build

FROM ${REPO_LOCATION}/agw_gateway_python as python_builder

ARG GHZ_REPO=https://github.com/bojand/ghz

RUN apt-get update && apt-get install -y \
  git \
  golang \
  build-essential

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
