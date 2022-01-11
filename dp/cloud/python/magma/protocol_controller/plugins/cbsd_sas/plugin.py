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

from magma.protocol_controller.plugin import ProtocolPlugin
from magma.protocol_controller.plugins.cbsd_sas.wsgi import application


class CBSDSASProtocolPlugin(ProtocolPlugin):
    """
    Protocol Controller plugin for CBSD-SAS protocol
    """

    def initialize(self):
        """
        Initialize CBSD-SAS protocol plugin. Start HTTP service
        """
        application.run(host='0.0.0.0', port=8080)  # noqa: S104
