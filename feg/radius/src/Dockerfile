FROM golang:alpine as builder
RUN apk add git gcc musl-dev bash protobuf
COPY radius/src /src/radius
COPY radius/src/config/samples/*template /src/radius/config/
COPY radius/lib/go/ /src/lib/go
WORKDIR /src/radius
ENV GOPROXY https://proxy.golang.org
RUN go mod download
RUN ./run.sh build

FROM alpine
RUN apk add gettext musl
COPY --from=builder /src/radius/radius /app/
COPY --from=builder /src/radius/config/ /app/
COPY radius/src/docker-entrypoint.sh /app/docker-entrypoint.sh
RUN chmod 0755 /app/docker-entrypoint.sh
WORKDIR /app
# Add version file with default BUILD_NUM unless set otherwise in build command
ARG BUILD_NUM=1.0.0
RUN echo "${BUILD_NUM}" > /app/VERSION
# ENTRYPOINT [ "./docker-entrypoint.sh" ]
