load("@rules_cc//cc:defs.bzl", "cc_library")
load("@rules_proto//proto:defs.bzl", "proto_library")


# Doing this now to save time, but is this OK in practice?
package(default_visibility = ["//visibility:public"])

# Convention:
# A cc_proto_library that wraps a proto_library named foo_proto
# should be called foo_cc_proto.
cc_proto_library(
    name = "abort_session_cc_proto",
    deps = [":abort_session_proto"],
)

# Conventions:
# 1. One proto_library rule per .proto file.
# 2. A file named foo.proto will be in a rule named foo_proto.
proto_library(
    name = "abort_session_proto",
    srcs = ["abort_session.proto"],
    deps = [],
)

cc_proto_library(
    name = "apn_cc_proto",
    deps = [":apn_proto"],
)
proto_library(
    name = "apn_proto",
    srcs = ["apn.proto"],
    deps = [],
)

cc_proto_library(
    name = "subscriberdb_cc_proto",
    deps = [":subscriberdb_proto"],
)
proto_library(
    name = "subscriberdb_proto",
    srcs = ["subscriberdb.proto"],
    deps = [
        ":apn_proto",
        "//orc8r/protos:common_proto",
        "@com_google_protobuf//:field_mask_proto",
    ],
)

cc_proto_library(
    name = "mobilityd_cc_proto",
    deps = [":mobilityd_proto"],
)
proto_library(
    name = "mobilityd_proto",
    srcs = ["mobilityd.proto"],
    deps = [
        ":subscriberdb_proto",
        "//orc8r/protos:common_proto",
    ],
)

cc_proto_library(
    name = "policydb_cc_proto",
    deps = [":policydb_proto"],
    visibility = ["//visibility:public"],
)
proto_library(
    name = "policydb_proto",
    srcs = ["policydb.proto"],
    deps = [
        ":mobilityd_proto",
        "//orc8r/protos:common_proto",
    ],
)

cc_proto_library(
    name = "session_manager_cc_proto",
    deps = [":session_manager_proto"],
    visibility = ["//visibility:public"],
)
proto_library(
    name = "session_manager_proto",
    srcs = ["session_manager.proto"],
    deps = [
        ":apn_proto",
        ":policydb_proto",
        ":subscriberdb_proto",
        "//orc8r/protos:common_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)

cc_proto_library(
    name = "pipelined_cc_proto",
    deps = [":pipelined_proto"],
)
proto_library(
    name = "pipelined_proto",
    srcs = ["pipelined.proto"],
    deps = [
        ":apn_proto",
        ":mobilityd_proto",
        ":policydb_proto",
        ":subscriberdb_proto",
        ":session_manager_proto",
        "//orc8r/protos:common_proto",
    ],
)

cc_proto_library(
    name = "spgw_service_cc_proto",
    deps = [":spgw_service_proto"],
)
proto_library(
    name = "spgw_service_proto",
    srcs = ["spgw_service.proto"],
    deps = [
        ":policydb_proto",
        ":subscriberdb_proto",
    ],
)