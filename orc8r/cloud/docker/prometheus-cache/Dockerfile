FROM golang:1.11 as go

# Use public go modules proxy
ENV GOPROXY https://proxy.golang.org
ENV GOBIN /build/bin

ARG CACHE_FILES=go/services/metricsd/prometheus/prometheus-cache

COPY ${CACHE_FILES} /go/src/magma/orc8r/cloud/go/services/metricsd/prometheus/prometheus-cache

WORKDIR /go/src/magma/orc8r/cloud/go/services/metricsd/prometheus/prometheus-cache

RUN go get -d -v ./...
RUN CGO_ENABLED=0 go install .

FROM alpine:latest

COPY --from=go /build/bin/prometheus-cache /bin/prometheus-cache

EXPOSE 9091

ENTRYPOINT ["prometheus-cache"]
