# Copyright 2021 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""All external repositories used for C++/C dependencies"""

load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository", "new_git_repository")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def cpp_repositories():
    """Entry point for all external repositories used for C++/C dependencies"""
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

    http_archive(
        name = "yaml-cpp",
        strip_prefix = "yaml-cpp-yaml-cpp-0.7.0",
        sha256 = "43e6a9fcb146ad871515f0d0873947e5d497a1c9c60c58cb102a97b47208b7c3",
        urls = ["https://github.com/jbeder/yaml-cpp/archive/refs/tags/yaml-cpp-0.7.0.tar.gz"],
    )

    http_archive(
        name = "github_nlohmann_json",
        build_file = "//bazel/external:nlohmann_json.BUILD",
        sha256 = "69cc88207ce91347ea530b227ff0776db82dcb8de6704e1a3d74f4841bc651cf",
        urls = ["https://github.com/nlohmann/json/releases/download/v3.6.1/include.zip"],
    )

    # prometheus_cpp dependency
    new_git_repository(
        name = "prometheus_client_model",
        build_file = "//bazel/external:prometheus_client_model.BUILD",
        # Used what master was around when D6071833@fb was authored (Oct 18, 2017)
        # The metrics.proto pulled here should match what we have in orc8r/protos/prometheus/metrics.proto
        commit = "fa8ad6fec33561be4280a8f0514318c79d7f6cb6",
        shallow_since = "1423736264 +0000",
        remote = "https://github.com/prometheus/client_model.git",
    )

    # prometheus_cpp dependency
    new_git_repository(
        name = "civetweb",
        build_file = "//bazel/external:civetweb.BUILD",
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
        build_file = "//bazel/external:cpp_redis.BUILD",
        url = "https://github.com/cpp-redis/cpp_redis/archive/refs/tags/4.3.1.tar.gz",
    )

    http_archive(
        name = "com_google_googletest",
        sha256 = "5cf189eb6847b4f8fc603a3ffff3b0771c08eec7dd4bd961bfd45477dd13eb73",
        strip_prefix = "googletest-609281088cfefc76f9d0ce82e1ff6c30cc3591e5",
        urls = ["https://github.com/google/googletest/archive/609281088cfefc76f9d0ce82e1ff6c30cc3591e5.zip"],
    )

    http_archive(
        name = "sentry_native",
        sha256 = "d7fa804995124c914a3abe077a73307960bbcadfbba9021e8ccbd05c7ba45f88",
        build_file = "//bazel/external:sentry_native.BUILD",
        url = "https://github.com/getsentry/sentry-native/releases/download/0.4.12/sentry-native.zip",
    )

    http_archive(
        name = "libtins",
        build_file = "//bazel/external:libtins.BUILD",
        url = "https://github.com/mfontanini/libtins/archive/refs/tags/v4.2.tar.gz",
        strip_prefix = "libtins-4.2",
        sha256 = "a9fed73e13f06b06a4857d342bb30815fa8c359d00bd69547e567eecbbb4c3a1",
    )

    new_git_repository(
        name = "liblfds",
        build_file = "//bazel/external:liblfds.BUILD",
        commit = "b813a0e546ed54e54b3873bdf180cf885c39bbca",
        remote = "https://github.com/liblfds/liblfds.git",
        shallow_since = "1464682027 +0300",
        patches = ["//third_party/build/patches/liblfds:0001-arm64-support.patch"],
        patch_args = ["--strip=1"],
    )

    new_git_repository(
        name = "libfluid_base",
        build_file = "//bazel/external:libfluid_base.BUILD",
        commit = "56df5e20c49387ab8e6b5cd363c6c10d309f263e",
        remote = "https://github.com/OpenNetworkingFoundation/libfluid_base",
        shallow_since = "1448037833 -0200",
        patches = [
            "//third_party/build/patches/libfluid/libfluid_base_patches:EVLOOP_NO_EXIT_ON_EMPTY_compat.patch",
            "//third_party/build/patches/libfluid/libfluid_base_patches:ExternalEventPatch.patch",
        ],
        patch_args = ["--strip=1"],
    )

    new_git_repository(
        name = "libfluid_msg",
        build_file = "//bazel/external:libfluid_msg.BUILD",
        commit = "71a4fccdedfabece730082fbe87ef8ae5f92059f",
        remote = "https://github.com/OpenNetworkingFoundation/libfluid_msg.git",
        shallow_since = "1487696730 +0000",
        patches = [
            "//third_party/build/patches/libfluid/libfluid_msg_patches:0001-Add-TunnelIPv4Dst-support.patch",
            "//third_party/build/patches/libfluid/libfluid_msg_patches:0002-Add-support-for-setting-OVS-reg8.patch",
            "//third_party/build/patches/libfluid/libfluid_msg_patches:0003-Add-Reg-field-match-support.patch",
            "//third_party/build/patches/libfluid/libfluid_msg_patches:0004-Add-TunnelIPv6Dst-support.patch",
        ],
        patch_args = ["--strip=1"],
    )
