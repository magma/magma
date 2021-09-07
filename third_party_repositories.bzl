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
        sha256 = "43e6a9fcb146ad871515f0d0873947e5d497a1c9c60c58cb102a97b47208b7c3",
        urls = ["https://github.com/jbeder/yaml-cpp/archive/refs/tags/yaml-cpp-0.7.0.tar.gz"],
    )

def nlohmann_json():
    http_archive(
        name = "github_nlohmann_json",
        build_file = "//third_party:nlohmann_json.BUILD",
        sha256 = "69cc88207ce91347ea530b227ff0776db82dcb8de6704e1a3d74f4841bc651cf",
        urls = ["https://github.com/nlohmann/json/releases/download/v3.6.1/include.zip"],
    )

def boost():
    git_repository(
        name = "com_github_nelhage_rules_boost",
        commit = "1e3a69bf2d5cd10c34b74f066054cd335d033d71",
        remote = "https://github.com/nelhage/rules_boost",
        shallow_since = "1591047380 -0700",
    )

def protobuf():
    http_archive(
        # The name is protobuf here as that is what prometheus-cpp expects
        # See https://github.com/jupp0r/prometheus-cpp.git @ d8326b2bba945a435f299e7526c403d7a1f68c1f
        name = "protobuf",
        strip_prefix = "protobuf-3.15.0",
        sha256 = "6aff9834fd7c540875e1836967c8d14c6897e3785a2efac629f69860fb7834ff",
        # TODO(@themarwhal): Upgrade to latest release once we resolve GH8457
        urls = ["https://github.com/protocolbuffers/protobuf/archive/refs/tags/v3.15.0.tar.gz"],
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
        urls = ["https://github.com/google/googletest/archive/ba96d0b1161f540656efdaed035b3c062b60e006.zip"],
    )

def prometheus_cpp_deps():
    new_git_repository(
        name = "prometheus_client_model",
        build_file = "//:third_party/prometheus_client_model.BUILD",
        # Used what master was around when D6071833@fb was authored (Oct 18, 2017)
        # The metrics.proto pulled here should match what we have in orc8r/protos/prometheus/metrics.proto
        commit = "fa8ad6fec33561be4280a8f0514318c79d7f6cb6",
        remote = "https://github.com/prometheus/client_model.git",
    )

    new_git_repository(
        name = "civetweb",
        build_file = "//:third_party/civetweb.BUILD",
        commit = "fbdee7440be24f904208c15a1fc9e2582b866049",
        remote = "https://github.com/civetweb/civetweb.git",
        shallow_since = "1474835570 +0200",
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

    http_archive(
        name = "cpp_redis",
        sha256 = "3859289d8254685fc775bda73de03dad27df923423b8ceb375b02d036c03b02f",
        strip_prefix = "cpp_redis-4.3.1",
        # TODO(@themarwhal): We do not need a custom BUILD file if we upgrade to a more recent version of cpp_redis - GH8321
        build_file = "//third_party:cpp_redis.BUILD",
        url = "https://github.com/cpp-redis/cpp_redis/archive/refs/tags/4.3.1.tar.gz",
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
        sha256 = "f28359aeba12f30d73d9e4711ef356dc842886968112162bc73002645139c39c",
        urls = ["https://github.com/google/glog/archive/v0.4.0.tar.gz"],
    )
