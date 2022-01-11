"""
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

from abc import abstractmethod


class ProtocolPlugin(object):
    """
    A very simple plugin class.

    Since protocol controllers do not really have any generic shareable API or design,
    we can only assume that each plugin will have its own 'initialize' method, which
    will be responsible for running the plugin.

    For instance, in the case of the `cbsd_sas` plugin, the plugin will start a flask server.

    The interface with other services (like Radio Controller) is based on gRPC so any other plugin
    is free to implement their own gRPC client and interface with RC.
    """

    @abstractmethod
    def initialize(self):
        """
        Initialize plugin. Abstract method, redefine in real plugin.
        """
        pass
