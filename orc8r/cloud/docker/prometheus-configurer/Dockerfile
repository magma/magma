# Golang image to build prom_alertconfig service
FROM golang:1.13-alpine3.11 as go

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
RUN go build -i -o $GOBIN/prometheus_configurer

FROM alpine:3.11

COPY --from=go /build/bin/prometheus_configurer /bin/prometheus_configurer

# Copy config files
COPY docker/prometheus-configurer/configs /etc/configs

ENTRYPOINT ["prometheus_configurer"]
