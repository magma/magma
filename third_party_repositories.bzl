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
