# Copyright 2021 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

load("@rules_cc//cc:defs.bzl", "cc_library")

cc_library(
    name = "lfds710",
    srcs = glob(["liblfds/liblfds7.1.0/liblfds710/src/**"]),
    hdrs = glob(["liblfds/liblfds7.1.0/liblfds710/inc/**"]),
    includes = ["iblfds/liblfds7.1.0/liblfds710/inc"],
    visibility = ["//visibility:public"],
)
