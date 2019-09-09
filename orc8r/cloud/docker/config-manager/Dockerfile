# Golang image to build prom_alertconfig service
FROM golang:1.12 as go

ENV GO111MODULE on
ENV GOBIN /build/bin
ENV PATH ${PATH}:${GOBIN}
# Use public go modules proxy
ENV GOPROXY https://proxy.golang.org

# Copy go.mod file to cache the downloads
COPY go/go.mod /gomod/
WORKDIR /gomod
RUN go mod download

# Copy magma code
COPY go/ /go/src/magma/orc8r/cloud/go/

# Build prom_alertconfig service
WORKDIR /go/src/magma/orc8r/cloud/go/services/metricsd/prometheus/alerting/prom_alertconfig
RUN go install .
WORKDIR /go/src/magma/orc8r/cloud/go/services/metricsd/prometheus/alerting/alertmanager_config
RUN go install .

FROM ubuntu:xenial

RUN apt-get update && \
    apt-get install -y supervisor

COPY --from=go /build/bin/prom_alertconfig /bin/prom_alertconfig
COPY --from=go /build/bin/alertmanager_config /bin/alertmanager_config

# Copy config files
COPY docker/config-manager/configs /etc/configs

COPY docker/config-manager/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
CMD ["/usr/bin/supervisord"]
