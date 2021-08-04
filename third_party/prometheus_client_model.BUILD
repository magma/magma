load("@rules_proto//proto:defs.bzl", "proto_library")
load("@rules_proto_grpc//cpp:defs.bzl", "cpp_proto_library")

package(default_visibility = ["//visibility:public"])

cpp_proto_library(
    name = "prometheus_client_model",
    protos = [":metrics_proto"],
)

proto_library(
    name = "metrics_proto",
    srcs = ["metrics.proto"],
)
