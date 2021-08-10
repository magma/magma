# Copyright 2021 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

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
        # Used what master was around when D6071833@fb was authored
        # The metrics.proto pulled here should match what we have in orc8r/protos/prometheus/metrics.proto
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

def cpp_redis():
    http_archive(
        name = "tacopie",
        sha256 = "bbdebecdb68d5f9eb64170217000daf844e0aee18b8c4d3dd373d07efd9f7316",
        strip_prefix = "tacopie-master",
        url = "https://github.com/cylix/tacopie/archive/master.zip",
    )

    new_git_repository(
        name = "cpp_redis",
        commit = "bbe38a7f83de943ffcc90271092d689ae02b3489",
        remote = "https://github.com/cpp-redis/cpp_redis.git",
        shallow_since = "1590000158 -0500",
        # TODO(@themarwhal): We do not need a custom BUILD file if we upgrade to a more recent version of cpp_redis - GH8321
        build_file = "//third_party:cpp_redis.BUILD",
    )

def gflags():
    http_archive(
        name = "com_github_gflags_gflags",
        sha256 = "34af2f15cf7367513b352bdcd2493ab14ce43692d2dcd9dfc499492966c64dcf",
        strip_prefix = "gflags-2.2.2",
        urls = ["https://github.com/gflags/gflags/archive/v2.2.2.tar.gz"],
    )

def glog():
    http_archive(
        name = "com_github_google_glog",
        strip_prefix = "glog-0.4.0",
        urls = [
            "https://github.com/google/glog/archive/v0.4.0.tar.gz",
        ],
    )
