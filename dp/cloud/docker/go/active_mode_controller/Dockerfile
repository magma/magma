FROM golang:1.16.14-alpine3.14 AS build
COPY dp/cloud/go/active_mode_controller /active_mode_controller
WORKDIR /active_mode_controller/cmd
RUN go build

FROM alpine:3.14.3 as final
COPY --from=build /active_mode_controller/cmd/cmd /active_mode_controller/cmd
WORKDIR /active_mode_controller
CMD ["./cmd"]
