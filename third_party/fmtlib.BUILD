load("@rules_cc//cc:defs.bzl", "cc_library")

cc_library(
    name = "libfmt",
    srcs = glob([
        "src/*.cc",
    ]),
    hdrs = glob([
        "include/fmt/*.h",
    ]),
    copts = ["-D__cpp_modules"],
    includes = ["include"],
    visibility = ["//visibility:public"],
)
