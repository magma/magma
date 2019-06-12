FROM golang:1.11

ENV GO111MODULE on
# Use public go modules proxy
ENV GOPROXY https://proxy.golang.org

ARG CACHE_FILES=go/services/metricsd/prometheus/prometheus-cache

COPY ${CACHE_FILES} /go/src/magma/orc8r/cloud/go/services/metricsd/prometheus/prometheus-cache

WORKDIR /go/src/magma/orc8r/cloud/go/services/metricsd/prometheus/prometheus-cache

RUN go get -d -v ./...

RUN go install -v ./...

EXPOSE 9091

CMD ["prometheus-cache"]
