# Copyright 2022 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

load("@rules_python//python:defs.bzl", "py_library")

package(default_visibility = ["//visibility:public"])

# The s1ap framework is only present on the magma_test environment
config_setting(
    name = "expect_s1ap_framework",
    values = {"define": "on_magma_test=1"},
)

py_library(
    name = "s1ap_types",
    srcs = select({
        ":expect_s1ap_framework": ["home/vagrant/s1ap-tester/bin/s1ap_types.py"],
        "//conditions:default": [],
    }),
    imports = select({
        ":expect_s1ap_framework": ["home/vagrant/s1ap-tester/bin/"],
        "//conditions:default": [],
    }),
    visibility = ["//visibility:public"],
)
