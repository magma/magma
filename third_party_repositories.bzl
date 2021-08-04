load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository", "new_git_repository")

def yaml_cpp():
    http_archive(
        name = "yaml-cpp",
        strip_prefix = "yaml-cpp-yaml-cpp-0.7.0",
        urls = ["https://github.com/jbeder/yaml-cpp/archive/refs/tags/yaml-cpp-0.7.0.tar.gz"],
    )

def zlib():
    http_archive(
        name = "zlib",
        build_file = "//:third_party/zlib.BUILD",
        sha256 = "c3e5e9fdd5004dcb542feda5ee4f0ff0744628baf8ed2dd5d66f8ca1197cb1a1",
        strip_prefix = "zlib-1.2.11",
        urls = [
            "https://mirror.bazel.build/zlib.net/zlib-1.2.11.tar.gz",
            "https://zlib.net/zlib-1.2.11.tar.gz",
        ],
    )

def nlohmann_json():
    http_archive(
        name = "github_nlohmann_json",
        build_file = "//third_party:nlohmann_json.BUILD",
        sha256 = "69cc88207ce91347ea530b227ff0776db82dcb8de6704e1a3d74f4841bc651cf",
        urls = [
            "https://github.com/nlohmann/json/releases/download/v3.6.1/include.zip",
        ],
    )

def boost():
    git_repository(
        name = "com_github_nelhage_rules_boost",
        commit = "1e3a69bf2d5cd10c34b74f066054cd335d033d71",
        remote = "https://github.com/nelhage/rules_boost",
        shallow_since = "1591047380 -0700",
    )

def protobuf():
    git_repository(
        # The name is protobuf here as that is what prometheus-cpp expects
        # See https://github.com/jupp0r/prometheus-cpp.git @ d8326b2bba945a435f299e7526c403d7a1f68c1f
        name = "protobuf",
        # TODO(@themarwhal): Upgrade to latest release once we resolve GH8457
        commit = "ae50d9b9902526efd6c7a1907d09739f959c6297",
        remote = "https://github.com/protocolbuffers/protobuf",
        shallow_since = "1624681439 -0700",
    )

def grpc():
    # see https://rules-proto-grpc.aliddell.com/en/latest/index.html
    http_archive(
        name = "rules_proto_grpc",
        sha256 = "7954abbb6898830cd10ac9714fbcacf092299fda00ed2baf781172f545120419",
        strip_prefix = "rules_proto_grpc-3.1.1",
        urls = ["https://github.com/rules-proto-grpc/rules_proto_grpc/archive/3.1.1.tar.gz"],
    )

def cpp_testing_deps():
    http_archive(
        name = "com_google_googletest",
        sha256 = "5cf189eb6847b4f8fc603a3ffff3b0771c08eec7dd4bd961bfd45477dd13eb73",
        strip_prefix = "googletest-609281088cfefc76f9d0ce82e1ff6c30cc3591e5",
        urls = ["https://github.com/google/googletest/archive/609281088cfefc76f9d0ce82e1ff6c30cc3591e5.zip"],
    )

def prometheus_cpp_deps():
    new_git_repository(
        name = "prometheus_client_model",
        build_file = "//:third_party/prometheus_client_model.BUILD",
        # Used what master probably was when D6071833@fb was authored
        commit = "fa8ad6fec33561be4280a8f0514318c79d7f6cb6",
        remote = "https://github.com/prometheus/client_model.git",
    )

    new_git_repository(
        name = "civetweb",
        build_file = "//:third_party/civetweb.BUILD",
        commit = "fbdee74",
        remote = "https://github.com/civetweb/civetweb.git",
    )

    git_repository(
        name = "prometheus_cpp",
        commit = "d8326b2bba945a435f299e7526c403d7a1f68c1f",
        remote = "https://github.com/jupp0r/prometheus-cpp.git",
    )
