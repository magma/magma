# Copyright 2021 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cc_library(
    name = "sentry",
    srcs = glob([
        "src/**/*.cpp",
    ]) + glob([
        "src/**/*.c",
    ]),
    hdrs = glob([
        "include/**/*.h",
    ]),
    copts= ["-DSENTRY_BUILD_SHARED_LIBS=1"],
    includes = ["include"],
    strip_include_prefix = "sentry-native-0.4.12",
    visibility = ["//visibility:public"],
)
