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

import abc


class Application(metaclass=abc.ABCMeta):
    """
    Diameter is built to be extensible for applications to be written on top of.
    This defines the interface for an application that can initiate and handle
    messages.
    """
    # The ID this application uses for messages
    APP_ID = 0
    # Vendor-Specific-Application-Id and VendorId AVPs that
    # this application should advertise in the Capabilities Exchange
    CAPABILITIES_EXCHANGE_AVPS = []

    def __init__(self, realm, host, host_ip, loop=None):
        """Each application has access to a write stream and a collection of
        settings, currently limited to realm, host and host_ip

        Args:
            realm: the realm the application should serve
            host: the FQDN of the host
            host_ip: the IP address of the host
            loop: asyncio loop
        """
        self.writer = None
        self.host = host
        self.realm = realm
        self.host_ip = host_ip
        self._loop = loop

    def set_writer(self, writer):
        """ Set a writer when a connection is made """
        self.writer = writer

    @abc.abstractmethod
    def handle_msg(self, state_id, msg):
        """
        Handles an incoming message bound to the application

        Args:
            state_id: the server state id
            msg: the inbound message to handle
        Returns:
            None
        """
        pass

    @abc.abstractmethod
    def validate_message(self, state_id, msg):
        """
        Validates a message addressed to the application and
        send the error response if necessary

        Args:
            state_id: server state_id
            msg: the message to check
        Returns:
            True if the message validated
        """
        pass
