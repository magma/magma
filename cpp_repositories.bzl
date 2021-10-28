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

def cpp_repositories():
    """All external repositories used for C++/C dependencies"""
    http_archive(
        name = "com_github_gflags_gflags",
        sha256 = "34af2f15cf7367513b352bdcd2493ab14ce43692d2dcd9dfc499492966c64dcf",
        strip_prefix = "gflags-2.2.2",
        urls = ["https://github.com/gflags/gflags/archive/v2.2.2.tar.gz"],
    )

    http_archive(
        name = "com_github_google_glog",
        strip_prefix = "glog-0.4.0",
        sha256 = "f28359aeba12f30d73d9e4711ef356dc842886968112162bc73002645139c39c",
        urls = ["https://github.com/google/glog/archive/v0.4.0.tar.gz"],
    )

    rules_boost_commit = "fb9f3c9a6011f966200027843d894923ebc9cd0b"
    http_archive(
        name = "com_github_nelhage_rules_boost",
        sha256 = "046f774b185436d506efeef8be6979f2c22f1971bfebd0979bafa28088bf28d0",
        strip_prefix = "rules_boost-{}".format(rules_boost_commit),
        urls = [
            "https://github.com/nelhage/rules_boost/archive/{}.tar.gz".format(rules_boost_commit),
        ],
    )

    http_archive(
        name = "yaml-cpp",
        strip_prefix = "yaml-cpp-yaml-cpp-0.7.0",
        sha256 = "43e6a9fcb146ad871515f0d0873947e5d497a1c9c60c58cb102a97b47208b7c3",
        urls = ["https://github.com/jbeder/yaml-cpp/archive/refs/tags/yaml-cpp-0.7.0.tar.gz"],
    )

    http_archive(
        name = "github_nlohmann_json",
        build_file = "//third_party:nlohmann_json.BUILD",
        sha256 = "69cc88207ce91347ea530b227ff0776db82dcb8de6704e1a3d74f4841bc651cf",
        urls = ["https://github.com/nlohmann/json/releases/download/v3.6.1/include.zip"],
    )

    # prometheus_cpp dependency
    new_git_repository(
        name = "prometheus_client_model",
        build_file = "//:third_party/prometheus_client_model.BUILD",
        # Used what master was around when D6071833@fb was authored (Oct 18, 2017)
        # The metrics.proto pulled here should match what we have in orc8r/protos/prometheus/metrics.proto
        commit = "fa8ad6fec33561be4280a8f0514318c79d7f6cb6",
        shallow_since = "1423736264 +0000",
        remote = "https://github.com/prometheus/client_model.git",
    )

    # prometheus_cpp dependency
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
        shallow_since = "1485901529 +0100",
    )

    # cpp_redis dependency
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

    http_archive(
        name = "com_google_googletest",
        sha256 = "5cf189eb6847b4f8fc603a3ffff3b0771c08eec7dd4bd961bfd45477dd13eb73",
        strip_prefix = "googletest-609281088cfefc76f9d0ce82e1ff6c30cc3591e5",
        urls = ["https://github.com/google/googletest/archive/609281088cfefc76f9d0ce82e1ff6c30cc3591e5.zip"],
    )

    new_git_repository(
        name = "sentry_native",
        build_file = "//third_party:sentry_native.BUILD",
        # 0.4.12 tag
        commit = "3436a29d839aa7437548be940ab62a85ca699635",
        # This is important, we pull in get_sentry/breakpad this way
        init_submodules = True,
        remote = "https://github.com/getsentry/sentry-native",
        shallow_since = "1627998929 +0000",
    )
