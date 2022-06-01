FROM golang:alpine as builder
RUN apk add git gcc musl-dev bash protobuf

ENV MAGMA_ROOT /magma
ENV GOPROXY https://proxy.golang.org

COPY feg/radius/lib/go/ $MAGMA_ROOT/feg/radius/lib/go
COPY feg/radius/src/go.* $MAGMA_ROOT/feg/radius/src/
COPY orc8r/lib/go/ $MAGMA_ROOT/orc8r/lib/go/
WORKDIR $MAGMA_ROOT/feg/radius/src
RUN go mod download

COPY feg/radius/src/ $MAGMA_ROOT/feg/radius/src/
RUN true # workaround for moby issue #37965
COPY feg/radius/lib/go/ $MAGMA_ROOT/feg/radius/lib/go/
RUN ./run.sh build

FROM alpine
RUN apk add gettext musl

COPY feg/radius/src/config/samples/*template /app/
COPY feg/radius/src/docker-entrypoint.sh /app/docker-entrypoint.sh
RUN chmod 0755 /app/docker-entrypoint.sh

COPY --from=builder /magma/feg/radius/src/radius /app/
WORKDIR /app
# Add version file with default BUILD_NUM unless set otherwise in build command
ARG BUILD_NUM=1.0.0
RUN echo "${BUILD_NUM}" > /app/VERSION
# ENTRYPOINT [ "./docker-entrypoint.sh" ]
