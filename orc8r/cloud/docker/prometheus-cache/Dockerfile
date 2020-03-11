FROM golang:1.13-alpine3.11 as go

# Use public go modules proxy
ENV GOPROXY https://proxy.golang.org
ENV GOBIN /build/bin

ARG CACHE_FILES=cloud/go/services/metricsd/prometheus/prometheus-cache

COPY ${CACHE_FILES} /go/src/magma/orc8r/cloud/go/services/metricsd/prometheus/prometheus-cache

# Copy go.mod file to cache the downloads
COPY cloud/go/go.mod /gomod/cloud/go/
COPY gateway/go/go.mod /gomod/gateway/go/
COPY lib/go/go.mod /gomod/lib/go/
COPY lib/go/protos/go.mod /gomod/lib/go/protos/
WORKDIR /gomod/cloud/go
RUN go mod download

# Copy magma code
COPY cloud/go/ /go/src/magma/orc8r/cloud/go/
COPY gateway/go/ /go/src/magma/orc8r/gateway/go/
COPY lib/go/ /go/src/magma/orc8r/lib/go/

WORKDIR /go/src/magma/orc8r/cloud/go/services/metricsd/prometheus/prometheus-cache
RUN go build -i -o $GOBIN/prometheus-cache

FROM alpine:3.11

COPY --from=go /build/bin/prometheus-cache /bin/prometheus-cache

EXPOSE 9091

ENTRYPOINT ["prometheus-cache"]
