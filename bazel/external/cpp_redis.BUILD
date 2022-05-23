# Copyright 2021 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

load("@rules_cc//cc:defs.bzl", "cc_library")

# We can get rid of this once we upgrade cpp-redis to a more up-to-date version - #8321
cc_library(
    name = "cpp_redis",
    srcs = glob(["sources/**/*.cpp"]),
    hdrs = glob([
        "includes/cpp_redis/**/*.hpp",
    ]) + [
        "includes/cpp_redis/cpp_redis",
        "includes/cpp_redis/impl/client.ipp",
    ],
    strip_include_prefix = "includes",
    visibility = ["//visibility:public"],
    deps = ["@tacopie"],
)
