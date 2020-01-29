FROM golang:1.13-alpine3.11 as go

# Use public go modules proxy
ENV GOPROXY https://proxy.golang.org
ENV GOBIN /build/bin

ARG CACHE_FILES=go/services/metricsd/prometheus/prometheus-cache

COPY ${CACHE_FILES} /go/src/magma/orc8r/cloud/go/services/metricsd/prometheus/prometheus-cache

# Copy go.mod file to cache the downloads
COPY go/go.mod /gomod/
WORKDIR /gomod
RUN go mod download

# Copy magma code
COPY go/ /go/src/magma/orc8r/cloud/go/

WORKDIR /go/src/magma/orc8r/cloud/go/services/metricsd/prometheus/prometheus-cache
RUN go build -i -o $GOBIN/prometheus-cache

FROM alpine:3.11

COPY --from=go /build/bin/prometheus-cache /bin/prometheus-cache

EXPOSE 9091

ENTRYPOINT ["prometheus-cache"]
