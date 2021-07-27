# Dependency of prometheus-cpp, modified from https://github.com/jupp0r/prometheus-cpp/blob/master/bazel/civetweb.BUILD
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
