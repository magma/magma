# Copyright 2021 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""Base images to be loaded from Dockerfiles"""

load(
    "@io_bazel_rules_docker//container:container.bzl",
    "container_load",
)
load("@io_bazel_rules_docker//contrib:dockerfile_build.bzl", "dockerfile_image")

def load_docker_images():
    dockerfile_image(
        name = "dockerfile_docker_agw_c_base",
        dockerfile = "//lte/gateway/docker/services/c:Dockerfile.c_base",
    )

    # Load the image tarball into Docker
    container_load(
        name = "loaded_dockerfile_docker_agw_c_base",
        file = "@dockerfile_docker_agw_c_base//image:dockerfile_image.tar",
    )
