load("@rules_cc//cc:defs.bzl", "cc_library")
load("@rules_proto//proto:defs.bzl", "proto_library")

cc_proto_library(
    name = "common_cc_proto",
    deps = [":common_proto"],
)
proto_library(
    name = "common_proto",
    srcs = ["common.proto"],
	visibility = ["//visibility:public"],
)