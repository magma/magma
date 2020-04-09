FROM golang:alpine as builder
RUN apk add git gcc musl-dev bash protobuf netcat-openbsd

RUN go get -d -u github.com/golang/protobuf/protoc-gen-go
RUN git -C "$(go env GOPATH)"/src/github.com/golang/protobuf checkout v1.2.0
RUN go install github.com/golang/protobuf/protoc-gen-go

# Use public go modules proxy
ENV GOPROXY https://proxy.golang.org

RUN printenv > /etc/environment

# Copy just the go.mod and go.sum files to download the golang deps.
# This step allows us to cache the downloads, and prevents reaching out to
# the internet unless any of the go.mod or go.sum files are changed.
COPY feg/radius/lib/go/ /radius/lib/go
COPY feg/radius/src/go.* /radius/src/
WORKDIR /radius/src
RUN go mod download

COPY feg/radius/src /radius/src
COPY feg/radius/src/config/samples/radius.cwf.config.json.template /radius/src/radius.config.json.template
COPY feg/gateway/services/aaa /gateway/services/aaa
RUN ./run.sh build

FROM alpine
RUN apk add gettext musl
COPY --from=builder /radius/src/radius /app/
COPY --from=builder /radius/src/radius.config.json.template /app/

WORKDIR /app
CMD ["sh", "-c", "envsubst < radius.config.json.template > radius.config.json && ./radius"]

