load("@rules_cc//cc:defs.bzl", "cc_library")

# Dependency of prometheus-cpp
# modified from https://github.com/jupp0r/prometheus-cpp WORKSPACE @ d8326b2bba945a435f299e7526c403d7a1f68c1f
package(default_visibility = ["//visibility:public"])

cc_library(
    name = "civetweb",
    srcs = [
        "src/CivetServer.cpp",
        "src/civetweb.c",
    ],
    hdrs = [
        "include/CivetServer.h",
        "include/civetweb.h",
        "src/handle_form.inl",
        "src/md5.inl",
    ],
    copts = [
        "-DUSE_IPV6",
        "-DNDEBUG",
        "-DNO_CGI",
        "-DNO_CACHING",
        "-DNO_SSL",
        "-DNO_FILES",
    ],
    includes = [
        "include",
    ],
)
