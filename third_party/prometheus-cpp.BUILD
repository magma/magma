package(default_visibility = ["//visibility:public"])

cc_library(
    name = "prometheus-cpp",
    srcs = [
        "lib/counter.cc",
        "lib/counter_builder.cc",
        "lib/exposer.cc",
        "lib/gauge.cc",
        "lib/gauge_builder.cc",
        "lib/handler.cc",
        "lib/histogram.cc",
        "lib/histogram_builder.cc",
        "lib/json_serializer.cc",
        "lib/json_serializer.h",
        "lib/protobuf_delimited_serializer.cc",
        "lib/protobuf_delimited_serializer.h",
        "lib/registry.cc",
        "lib/serializer.h",
        "lib/text_serializer.cc",
        "lib/text_serializer.h",
    ],
    hdrs = glob([
        "include/prometheus/*.h",
    ]),
    strip_include_prefix = "include",
    deps = [
        "@civetweb",
        "@com_google_protobuf//:protobuf",
        "@prometheus_client_model//:metrics_cpp_proto",
    ],
)
