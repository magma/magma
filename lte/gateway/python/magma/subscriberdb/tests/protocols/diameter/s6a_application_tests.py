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

# pylint:disable=protected-access

import unittest
from unittest.mock import Mock

from lte.protos.mconfig.mconfigs_pb2 import SubscriberDB
from magma.subscriberdb.crypto.utils import CryptoError
from magma.subscriberdb.processor import LTEProcessor
from magma.subscriberdb.protocols.diameter import avp, message, server
from magma.subscriberdb.protocols.diameter.application import (
    base,
    s6a,
    s6a_relay,
)
from magma.subscriberdb.store.base import SubscriberNotFoundError

from .common import MockTransport


def _dummy_eutran_vector():
    rand = b'\x00\x01\x02\x03\x04\x05\x06\x07\x08\t\n\x0b\x0c\r\x0e\x0f'
    xres = b'\x2d\xaf\x87\x3d\x73\xf3\x10\xc6'
    autn = b'o\xbf\xa3\x80\x1fW\x80\x00{\xdeY\x88n\x96\xe4\xfe'
    kasme = (
        b'\x87H\xc1\xc0\xa2\x82o\xa4\x05\xb1\xe2~\xa1\x04CJ\xe5V\xc7e'
        b'\xe8\xf0a\xeb\xdb\x8a\xe2\x86\xc4F\x16\xc2'
    )
    return (rand, xres, autn, kasme)


def _dummy_resync_vector():
    mac_s = b'\x01\xcf\xaf\x9e\xc4'
    sqn = 2000
    return (sqn, mac_s)


class MockProcessor(LTEProcessor):
    def generate_lte_auth_vector(self, imsi, plmn):
        if imsi == '1':
            return _dummy_eutran_vector()
        elif imsi == '2':
            raise CryptoError
        else:
            raise SubscriberNotFoundError

    def resync_lte_auth_seq(self, auts, key, rand):
        pass

    def get_next_lte_auth_seq(self, imsi):
        pass

    def set_next_lte_auth_seq(self, imsi, seq):
        pass

    # pylint:disable=unused-argument
    def get_sub_profile(self, imsi):
        return SubscriberDB.SubscriptionProfile(
            max_ul_bit_rate=10000,
            max_dl_bit_rate=50000,
        )


class S6AApplicationTests(unittest.TestCase):
    """
    Tests for the S6a commands implemented. These tests check that
    the server can decode requests, and respond to them with the
    encoded response that is expected
    """
    REALM = "mai.facebook.com"
    HOST = "hss.mai.facebook.com"
    HOST_ADDR = "127.0.0.1"

    def setUp(self):
        base_manager = base.BaseApplication(
            self.REALM, self.HOST, self.HOST_ADDR,
        )
        s6a_manager = s6a_relay.S6AApplication(
            MockProcessor(),
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
            self._writes.assert_called_once_with(resp_bytes)
            self._writes.reset_mock()

    def test_auth_success(self):
        """
        Test that we can respond to auth requests with an auth
        vector
        """
        msg = message.Message()
        msg.header.application_id = s6a.S6AApplication.APP_ID
        msg.header.command_code = s6a.S6AApplicationCommands.AUTHENTICATION_INFORMATION
        msg.header.request = True
        msg.append_avp(
            avp.AVP(
                'Session-Id',
                'enb-Lenovo-Product.openair4G.eur;1475864727;1;apps6a',
            ),
        )
        msg.append_avp(avp.AVP('Auth-Session-State', 1))
        msg.append_avp(avp.AVP('User-Name', '1'))
        msg.append_avp(avp.AVP('Visited-PLMN-Id', b'(Y'))
        msg.append_avp(
            avp.AVP(
                'Requested-EUTRAN-Authentication-Info', [
                    avp.AVP('Number-Of-Requested-Vectors', 1),
                    avp.AVP('Immediate-Response-Preferred', 0),
                ],
            ),
        )
        # Encode request message into buffer
        req_buf = bytearray(msg.length)
        msg.encode(req_buf, 0)

        msg = message.Message()
        msg.header.application_id = s6a.S6AApplication.APP_ID
        msg.header.command_code = \
            s6a.S6AApplicationCommands.AUTHENTICATION_INFORMATION
        msg.header.request = False
        msg.append_avp(
            avp.AVP(
                'Session-Id',
                'enb-Lenovo-Product.openair4G.eur;1475864727;1;apps6a',
            ),
        )
        msg.append_avp(
            avp.AVP(
                'Authentication-Info', [
                    avp.AVP(
                        'E-UTRAN-Vector', [
                            avp.AVP(
                                'RAND', b'\x00\x01\x02\x03\x04\x05'
                                b'\x06\x07\x08\x09\x0a\x0b\x0c\x0d\x0e\x0f',
                            ),
                            avp.AVP('XRES', b'\x2d\xaf\x87\x3d\x73\xf3\x10\xc6'),
                            avp.AVP(
                                'AUTN', b'\x6f\xbf\xa3\x80\x1f\x57\x80'
                                b'\x00\x7b\xde\x59\x88\x6e\x96\xe4\xfe',
                            ),
                            avp.AVP(
                                'KASME', b'\x87\x48\xc1\xc0\xa2\x82'
                                b'\x6f\xa4\x05\xb1\xe2\x7e\xa1\x04\x43\x4a'
                                b'\xe5\x56\xc7\x65\xe8\xf0\x61\xeb\xdb\x8a'
                                b'\xe2\x86\xc4\x46\x16\xc2',
                            ),
                        ],
                    ),
                ],
            ),
        )
        msg.append_avp(avp.AVP('Auth-Session-State', 1))

        # Host identifiers
        msg.append_avp(avp.AVP('Origin-Host', self._server.host))
        msg.append_avp(avp.AVP('Origin-Realm', self._server.realm))
        msg.append_avp(avp.AVP('Origin-State-Id', self._server.state_id))

        # Response result
        msg.append_avp(avp.AVP('Result-Code', avp.ResultCode.DIAMETER_SUCCESS))
        # Encode response into buffer
        resp_buf = bytearray(msg.length)
        msg.encode(resp_buf, 0)

        self._check_reply(req_buf, resp_buf)

    def test_resync(self):
        """
        Test that we can respond to auth requests with an auth
        vector
        """
        msg = message.Message()
        msg.header.application_id = s6a.S6AApplication.APP_ID
        msg.header.command_code = s6a.S6AApplicationCommands.AUTHENTICATION_INFORMATION
        msg.header.request = True
        msg.append_avp(
            avp.AVP(
                'Session-Id',
                'enb-Lenovo-Product.openair4G.eur;1475864727;1;apps6a',
            ),
        )
        msg.append_avp(avp.AVP('Auth-Session-State', 1))
        msg.append_avp(avp.AVP('User-Name', '1'))
        msg.append_avp(avp.AVP('Visited-PLMN-Id', b'(Y'))
        msg.append_avp(
            avp.AVP(
                'Requested-EUTRAN-Authentication-Info', [
                    avp.AVP('Number-Of-Requested-Vectors', 1),
                    avp.AVP('Immediate-Response-Preferred', 0),
                    avp.AVP('Re-Synchronization-Info', 30 * b'\x00'),
                ],
            ),
        )
        # Encode request message into buffer
        req_buf = bytearray(msg.length)
        msg.encode(req_buf, 0)

        processor = self._server._s6a_manager.lte_processor
        with unittest.mock.patch.object(processor, 'resync_lte_auth_seq'):
            self._server.data_received(req_buf)
            processor.resync_lte_auth_seq.assert_called_once_with(
                '1',
                16 * b'\x00', 14 * b'\x00',
            )

    def test_auth_bad_key(self):
        """
        Test that we reject auth requests if the stored key is bad
        """
        msg = message.Message()
        msg.header.application_id = s6a.S6AApplication.APP_ID
        msg.header.command_code = s6a.S6AApplicationCommands.AUTHENTICATION_INFORMATION
        msg.header.request = True
        msg.append_avp(
            avp.AVP(
                'Session-Id',
                'enb-Lenovo-Product.openair4G.eur;1475864727;1;apps6a',
            ),
        )
        msg.append_avp(avp.AVP('Auth-Session-State', 1))
        msg.append_avp(avp.AVP('User-Name', '2'))
        msg.append_avp(avp.AVP('Visited-PLMN-Id', b'(Y'))
        msg.append_avp(
            avp.AVP(
                'Requested-EUTRAN-Authentication-Info', [
                    avp.AVP('Number-Of-Requested-Vectors', 1),
                    avp.AVP('Immediate-Response-Preferred', 0),
                ],
            ),
        )
        # Encode request message into buffer
        req_buf = bytearray(msg.length)
        msg.encode(req_buf, 0)

        msg = message.Message()
        msg.header.application_id = s6a.S6AApplication.APP_ID
        msg.header.command_code = \
            s6a.S6AApplicationCommands.AUTHENTICATION_INFORMATION
        msg.header.request = False
        msg.append_avp(
            avp.AVP(
                'Session-Id',
                'enb-Lenovo-Product.openair4G.eur;1475864727;1;apps6a',
            ),
        )
        msg.append_avp(avp.AVP('Auth-Session-State', 1))

        # Host identifiers
        msg.append_avp(avp.AVP('Origin-Host', self._server.host))
        msg.append_avp(avp.AVP('Origin-Realm', self._server.realm))
        msg.append_avp(avp.AVP('Origin-State-Id', self._server.state_id))

        # Response result
        msg.append_avp(
            avp.AVP(
                'Result-Code',
                avp.ResultCode.DIAMETER_AUTHENTICATION_REJECTED,
            ),
        )
        # Encode response into buffer
        resp_buf = bytearray(msg.length)
        msg.encode(resp_buf, 0)
        self._check_reply(req_buf, resp_buf)

    def test_auth_unknown_subscriber(self):
        """
        Test that we reject auth requests if the subscriber is unknown
        """
        msg = message.Message()
        msg.header.application_id = s6a.S6AApplication.APP_ID
        msg.header.command_code = s6a.S6AApplicationCommands.AUTHENTICATION_INFORMATION
        msg.header.request = True
        msg.append_avp(
            avp.AVP(
                'Session-Id',
                'enb-Lenovo-Product.openair4G.eur;1475864727;1;apps6a',
            ),
        )
        msg.append_avp(avp.AVP('Auth-Session-State', 1))
        msg.append_avp(avp.AVP('User-Name', '3'))
        msg.append_avp(avp.AVP('Visited-PLMN-Id', b'(Y'))
        msg.append_avp(
            avp.AVP(
                'Requested-EUTRAN-Authentication-Info', [
                    avp.AVP('Number-Of-Requested-Vectors', 1),
                    avp.AVP('Immediate-Response-Preferred', 0),
                ],
            ),
        )
        # Encode request message into buffer
        req_buf = bytearray(msg.length)
        msg.encode(req_buf, 0)

        msg = message.Message()
        msg.header.application_id = s6a.S6AApplication.APP_ID
        msg.header.command_code = \
            s6a.S6AApplicationCommands.AUTHENTICATION_INFORMATION
        msg.header.request = False
        msg.append_avp(
            avp.AVP(
                'Session-Id',
                'enb-Lenovo-Product.openair4G.eur;1475864727;1;apps6a',
            ),
        )
        msg.append_avp(avp.AVP('Auth-Session-State', 1))

        # Host identifiers
        msg.append_avp(avp.AVP('Origin-Host', self._server.host))
        msg.append_avp(avp.AVP('Origin-Realm', self._server.realm))
        msg.append_avp(avp.AVP('Origin-State-Id', self._server.state_id))

        # Response result
        msg.append_avp(
            avp.AVP(
                'Result-Code',
                avp.ResultCode.DIAMETER_ERROR_USER_UNKNOWN,
            ),
        )
        # Encode response into buffer
        resp_buf = bytearray(msg.length)
        msg.encode(resp_buf, 0)
        self._check_reply(req_buf, resp_buf)

    def test_location_update(self):
        """
        Test that we can respond to update location request with
        subscriber data
        """
        msg = message.Message()
        msg.header.application_id = s6a.S6AApplication.APP_ID
        msg.header.command_code = s6a.S6AApplicationCommands.UPDATE_LOCATION
        msg.header.request = True
        msg.append_avp(
            avp.AVP(
                'Session-Id',
                'enb-Lenovo-Product.openair4G.eur;1475864727;1;apps6a',
            ),
        )
        msg.append_avp(avp.AVP('Auth-Session-State', 1))
        msg.append_avp(avp.AVP('User-Name', '208950000000001'))
        msg.append_avp(avp.AVP('Visited-PLMN-Id', b'(Y'))
        msg.append_avp(avp.AVP('RAT-Type', 1004))
        msg.append_avp(avp.AVP('ULR-Flags', 34))
        # Encode request message into buffer
        req_buf = bytearray(msg.length)
        msg.encode(req_buf, 0)

        msg = message.Message()
        msg.header.application_id = s6a.S6AApplication.APP_ID
        msg.header.command_code = s6a.S6AApplicationCommands.UPDATE_LOCATION
        msg.header.request = False
        msg.append_avp(
            avp.AVP(
                'Session-Id',
                'enb-Lenovo-Product.openair4G.eur;1475864727;1;apps6a',
            ),
        )
        msg.append_avp(avp.AVP('ULA-Flags', 1))
        msg.append_avp(
            avp.AVP(
                'Subscription-Data', [
                    avp.AVP('MSISDN', b'333608050011'),
                    avp.AVP('Access-Restriction-Data', 47),
                    avp.AVP('Subscriber-Status', 0),
                    avp.AVP('Network-Access-Mode', 2),
                    avp.AVP(
                        'AMBR', [
                            avp.AVP('Max-Requested-Bandwidth-UL', 10000),
                            avp.AVP('Max-Requested-Bandwidth-DL', 50000),
                        ],
                    ),
                    avp.AVP(
                        'APN-Configuration-Profile', [
                            avp.AVP('Context-Identifier', 0),
                            avp.AVP(
                                'All-APN-Configurations-Included-Indicator', 0,
                            ),
                            avp.AVP(
                                'APN-Configuration', [
                                    avp.AVP('Context-Identifier', 0),
                                    avp.AVP('PDN-Type', 0),
                                    avp.AVP('Service-Selection', 'oai.ipv4'),
                                    avp.AVP(
                                        'EPS-Subscribed-QoS-Profile', [
                                            avp.AVP('QoS-Class-Identifier', 9),
                                            avp.AVP(
                                                'Allocation-Retention-Priority', [
                                                    avp.AVP(
                                                        'Priority-Level', 15,
                                                    ),
                                                    avp.AVP(
                                                        'Pre-emption-Capability', 1,
                                                    ),
                                                    avp.AVP(
                                                        'Pre-emption-Vulnerability', 0,
                                                    ),
                                                ],
                                            ),
                                        ],
                                    ),
                                    avp.AVP(
                                        'AMBR', [
                                            avp.AVP(
                                                'Max-Requested-Bandwidth-UL', 10000,
                                            ),
                                            avp.AVP(
                                                'Max-Requested-Bandwidth-DL', 50000,
                                            ),
                                        ],
                                    ),
                                ],
                            ),
                        ],
                    ),
                ],
            ),
        )
        msg.append_avp(avp.AVP('Auth-Session-State', 1))

        # Host identifiers
        msg.append_avp(avp.AVP('Origin-Host', self._server.host))
        msg.append_avp(avp.AVP('Origin-Realm', self._server.realm))
        msg.append_avp(avp.AVP('Origin-State-Id', self._server.state_id))

        # Response result
        msg.append_avp(avp.AVP('Result-Code', avp.ResultCode.DIAMETER_SUCCESS))
        # Encode response into buffer
        resp_buf = bytearray(msg.length)
        msg.encode(resp_buf, 0)

        self._check_reply(req_buf, resp_buf)


if __name__ == "__main__":
    unittest.main()
