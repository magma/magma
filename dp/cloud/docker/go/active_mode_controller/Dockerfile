FROM golang:1.16 AS build
COPY dp/cloud/go/active_mode_controller /active_mode_controller
WORKDIR /active_mode_controller/cmd
RUN go build

FROM build as final
WORKDIR /active_mode_controller/cmd
CMD ["./cmd"]
