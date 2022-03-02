load("@rules_cc//cc:defs.bzl", "cc_binary", "cc_library", "cc_test")

ASAN_COPTS = ["-fsanitize=address", "-fsanitize=undefined", "-O0", "-fno-omit-frame-pointer"]
ASAN_LINKOPTS = ["-fsanitize=address", "-fsanitize=undefined"]
ASAN_DEFINES = ["ASAN_OPTIONS=detect_leaks=1:color=always"]

LSAN_COPTS = ["-fsanitize=leak", "-fno-omit-frame-pointer"]
LSAN_LINKOPTS = ["-fsanitize=leak"]

MAGMA_TEST_DEFINES = ["MME_UNIT_TEST"]

def _magma_copts(copts):
    return copts + select({
        "//bazel:enable_asan": ASAN_COPTS,
        "//conditions:default": [],
    }) + select({
        "//bazel:enable_lsan": LSAN_COPTS,
        "//conditions:default": [],
    }) 

def _magma_linkopts(linkopts):
    return linkopts + select({
        "//bazel:enable_asan": ASAN_LINKOPTS,
        "//conditions:default": [],
    }) + select({
        "//bazel:enable_lsan": LSAN_LINKOPTS,
        "//conditions:default": [],
    })

def _magma_defines(defines):
    return defines + select({
        "//bazel:enable_asan": ASAN_DEFINES,
        "//conditions:default": [],
    })+ select ({
         "//bazel:set_test_flag": MAGMA_TEST_DEFINES,
         "//conditions:default": [],
    })

def magma_cc_library(
        name,
        srcs = [],
        hdrs = [],
        copts = [],
        linkopts = [],
        visibility = None,
        tags = [],
        deps = [],
        strip_include_prefix = None,
        include_prefix = None,
        defines = []):
    """TODO: add doc"""
    cc_library(
        name = name,
        srcs = srcs,
        hdrs = hdrs,
        copts = _magma_copts(copts),
        linkopts = _magma_linkopts(linkopts),
        visibility = visibility,
        tags = tags,
        deps = deps,
        strip_include_prefix = strip_include_prefix,
        include_prefix = include_prefix,
        defines = _magma_defines(defines),
    )

def magma_cc_binary(
        name,
        srcs = [],
        copts = [],
        linkopts = [],
        visibility = None,
        tags = [],
        deps = [],
        strip_include_prefix = None,
        include_prefix = None,
        defines = [],
        linkstatic = True):
    """TODO: add doc"""
    cc_binary(
        name = name,
        srcs = srcs,
        copts = _magma_copts(copts),
        linkopts = _magma_linkopts(linkopts),
        visibility = visibility,
        tags = tags,
        deps = deps,
        strip_include_prefix = strip_include_prefix,
        include_prefix = include_prefix,
        defines = _magma_defines(defines),
        linkstatic = linkstatic,
    )

def magma_cc_test(
        name,
        srcs = [],
        copts = [],
        linkopts = [],
        visibility = None,
        tags = [],
        deps = [],
        strip_include_prefix = None,
        include_prefix = None,
        defines = [],
        size = None,
        flaky = None):
    """TODO: add doc"""
    cc_test(
        name = name,
        srcs = srcs,
        copts = _magma_copts(copts),
        linkopts = _magma_linkopts(linkopts),
        visibility = visibility,
        tags = tags,
        deps = deps,
        size = size,
        flaky = None,
        strip_include_prefix = strip_include_prefix,
        include_prefix = include_prefix,
        defines = _magma_defines(defines),
    )
