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

import unittest
from unittest.mock import Mock

from magma.subscriberdb.protocols.diameter import avp, message, server
from magma.subscriberdb.protocols.diameter.application import (
    base,
    s6a,
    s6a_relay,
)

from .common import MockTransport


class BaseApplicationTests(unittest.TestCase):
    """
    Tests for the Base Protocol commands implemented
    """
    REALM = "mai.facebook.com"
    HOST = "hss.mai.facebook.com"
    HOST_ADDR = "127.0.0.1"

    def setUp(self):
        base_manager = base.BaseApplication(
            self.REALM, self.HOST, self.HOST_ADDR,
        )
        s6a_manager = s6a_relay.S6AApplication(
            Mock(),
            self.REALM,
            self.HOST,
            self.HOST_ADDR,
        )
        base_manager.register(s6a_manager)
        self._server = server.S6aServer(
            base_manager,
            s6a_manager,
            self.REALM,
            self.HOST,
        )

        # Mock the writes to check responses
        self._writes = Mock()

        def convert_memview_to_bytes(memview):
            """ Deep copy the memoryview for checking later  """
            return self._writes(memview.tobytes())

        self._transport = MockTransport()
        self._transport.write = Mock(side_effect=convert_memview_to_bytes)

        # Here goes nothing..
        self._server.connection_made(self._transport)

    def _check_reply(self, req_bytes, resp_bytes):
        """
        Send data to the protocol in different step lengths to
        verify that we assemble all segments and parse correctly.

        Args:
            req_bytes (bytes): request which would be sent
                multiple times with different step sizes
            resp_bytes (bytes): response which needs to be
                received each time
        Returns:
            None
        """
        for step in range(1, len(req_bytes) + 1):
            offset = 0
            while offset < len(req_bytes):
                self._server.data_received(req_bytes[offset:offset + step])
                offset += step
            if len(resp_bytes):
                self._writes.assert_called_once_with(resp_bytes)
            self._writes.reset_mock()

    def test_watchdog(self):
        """Test that we can respond to watchdog requests"""
        msg = message.Message()
        msg.header.command_code = base.BaseApplicationCommands.DEVICE_WATCHDOG
        msg.header.request = True
        # Encode request message into buffer
        req_buf = bytearray(msg.length)
        msg.encode(req_buf, 0)

        msg = message.Message()
        msg.header.command_code = base.BaseApplicationCommands.DEVICE_WATCHDOG
        msg.header.request = False
        msg.append_avp(avp.AVP('Result-Code', avp.ResultCode.DIAMETER_SUCCESS))
        msg.append_avp(avp.AVP('Origin-Host', self._server.host))
        msg.append_avp(avp.AVP('Origin-Realm', self._server.realm))
        msg.append_avp(avp.AVP('Origin-State-Id', self._server.state_id))
        # Encode response into buffer
        resp_buf = bytearray(msg.length)
        msg.encode(resp_buf, 0)

        self._check_reply(req_buf, resp_buf)

    def test_capability_exchange(self):
        """Test that we can respond to capability exchange requests"""
        msg = message.Message()
        msg.header.command_code = base.BaseApplicationCommands.CAPABILITIES_EXCHANGE
        msg.header.request = True
        msg.append_avp(avp.AVP('Host-IP-Address', '127.0.0.1'))
        msg.append_avp(avp.AVP('Inband-Security-Id', 0))
        msg.append_avp(avp.AVP('Supported-Vendor-Id', 0))
        msg.append_avp(avp.AVP('Vendor-Id', 0))
        msg.append_avp(avp.AVP('Vendor-Specific-Application-Id', []))
        # Encode request message into buffer
        req_buf = bytearray(msg.length)
        msg.encode(req_buf, 0)

        msg = message.Message()
        msg.header.command_code = base.BaseApplicationCommands.CAPABILITIES_EXCHANGE
        msg.header.request = False
        msg.append_avp(avp.AVP('Result-Code', avp.ResultCode.DIAMETER_SUCCESS))
        msg.append_avp(avp.AVP('Origin-Host', self._server.host))
        msg.append_avp(avp.AVP('Origin-Realm', self._server.realm))
        msg.append_avp(avp.AVP('Origin-State-Id', self._server.state_id))
        msg.append_avp(avp.AVP('Host-IP-Address', self.HOST_ADDR))
        msg.append_avp(avp.AVP('Vendor-Id', 0))
        msg.append_avp(avp.AVP('Supported-Vendor-Id', avp.VendorId.TGPP))
        msg.append_avp(
            avp.AVP(
                'Vendor-Specific-Application-Id', [
                    avp.AVP('Auth-Application-Id', s6a.S6AApplication.APP_ID),
                    avp.AVP('Vendor-Id', avp.VendorId.TGPP),
                ],
            ),
        )
        msg.append_avp(avp.AVP('Product-Name', 'magma'))
        # Encode response message into buffer
        resp_buf = bytearray(msg.length)
        msg.encode(resp_buf, 0)

        self._check_reply(req_buf, resp_buf)

    def test_disconnect(self):
        """Test that we can respond to disconnect requests"""
        msg = message.Message()
        msg.header.command_code = base.BaseApplicationCommands.DISCONNECT_PEER
        msg.header.request = True
        # Encode request message into buffer
        req_buf = bytearray(msg.length)
        msg.encode(req_buf, 0)

        msg = message.Message()
        msg.header.command_code = base.BaseApplicationCommands.DISCONNECT_PEER
        msg.header.request = False
        msg.append_avp(avp.AVP('Result-Code', avp.ResultCode.DIAMETER_SUCCESS))
        msg.append_avp(avp.AVP('Origin-Host', self._server.host))
        msg.append_avp(avp.AVP('Origin-Realm', self._server.realm))
        msg.append_avp(avp.AVP('Origin-State-Id', self._server.state_id))
        # Encode response into buffer
        resp_buf = bytearray(msg.length)
        msg.encode(resp_buf, 0)

        self._check_reply(req_buf, resp_buf)

    def test_unsolicited_response(self):
        """Test that we ignore unsolicited responses"""
        msg = message.Message()
        msg.header.command_code = base.BaseApplicationCommands.CAPABILITIES_EXCHANGE
        msg.header.request = False
        # Encode request message into buffer
        req_buf = bytearray(msg.length)
        msg.encode(req_buf, 0)

        self._check_reply(req_buf, b'')

    def test_invalid_command(self):
        """Test that we ignore invalid commands"""
        msg = message.Message()
        msg.header.command_code = 0xfa4e
        msg.header.request = True
        # Encode request message into buffer
        req_buf = bytearray(msg.length)
        msg.encode(req_buf, 0)

        self._check_reply(req_buf, b'')


if __name__ == "__main__":
    unittest.main()
