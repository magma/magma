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
# ENTRYPOINT [ "./docker-entrypoint.sh" ]
