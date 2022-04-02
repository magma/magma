# Copyright 2021 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""Python toolchain configuration"""

load("@rules_python//python:repositories.bzl", "python_register_toolchains")

def configure_python_toolchain(name = None):
    python_register_toolchains(
        name = "python3_8",
        python_version = "3.8",
    )
