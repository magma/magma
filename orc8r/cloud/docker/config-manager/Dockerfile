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
WORKDIR /go/src/magma/orc8r/cloud/go/services/metricsd/prometheus/configmanager/prometheus/
RUN go build -i -o $GOBIN/prometheus_configmanager
WORKDIR /go/src/magma/orc8r/cloud/go/services/metricsd/prometheus/configmanager/alertmanager/
RUN go build -i -o $GOBIN/alertmanager_configmanager

FROM ubuntu:xenial

RUN apt-get update && \
    apt-get install -y supervisor

COPY --from=go /build/bin/prometheus_configmanager /bin/prometheus_configmanager
COPY --from=go /build/bin/alertmanager_configmanager /bin/alertmanager_configmanager

# Copy config files
COPY docker/config-manager/configs /etc/configs

COPY docker/config-manager/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
CMD ["/usr/bin/supervisord"]
