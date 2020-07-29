"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import os


def is_dev_mode() -> bool:
    """
    Returns whether the environment is set for dev mode
    """
    return os.environ.get('MAGMA_DEV_MODE') == '1'


def is_docker_network_mode() -> bool:
    """
    Returns whether the environment is set for dev mode
    """
    return os.environ.get('DOCKER_NETWORK_MODE') == '1'
