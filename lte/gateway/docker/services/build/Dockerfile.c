FROM mme_builder:latest

ARG MAGMA_VERSION=master

RUN git clone https://github.com/magma/magma.git && \
    cd magma && \
    git checkout $MAGMA_VERSION

WORKDIR /magma/lte/gateway

RUN make build_common build_oai build_sctpd build_session_manager
