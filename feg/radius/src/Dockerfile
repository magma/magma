FROM golang:alpine as builder
RUN apk add git gcc musl-dev bash protobuf
COPY radius/src /src/radius
COPY radius/src/config/samples/radius.cwf.config.json.template /src/radius/radius.config.json.template
COPY radius/lib/go/ /src/lib/go
WORKDIR /src/radius
ENV GOPROXY https://proxy.golang.org
RUN go mod download
RUN ./run.sh build 

FROM alpine
RUN apk add gettext musl
COPY --from=builder /src/radius/radius /app/
COPY --from=builder /src/radius/radius.config.json.template /app/
COPY radius/src/docker-entrypoint.sh /app/docker-entrypoint.sh
RUN chmod 0755 /app/docker-entrypoint.sh
WORKDIR /app
# ENTRYPOINT [ "./docker-entrypoint.sh" ]
