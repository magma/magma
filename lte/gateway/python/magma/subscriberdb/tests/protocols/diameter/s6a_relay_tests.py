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
from unittest.mock import Mock, patch

import grpc
from feg.protos.s6a_proxy_pb2 import (
    AuthenticationInformationAnswer,
    AuthenticationInformationRequest,
    ErrorCode,
    UpdateLocationAnswer,
    UpdateLocationRequest,
)
from magma.common.service_registry import ServiceRegistry
from magma.subscriberdb.protocols.diameter import avp, message, server
from magma.subscriberdb.protocols.diameter.application import (
    base,
    s6a,
    s6a_relay,
)

from .common import MockTransport


class S6AApplicationTests(unittest.TestCase):
    """
    Tests for the S6a Relay
    """
    REALM = "mai.facebook.com"
    HOST = "hss.mai.facebook.com"
    HOST_ADDR = "127.0.0.1"

    def setUp(self):
        ServiceRegistry.add_service('s6a_proxy', '0.0.0.0', 0)
        proxy_config = {
            'local_port': 1234,
            'cloud_address': 'test',
            'proxy_cloud_connections': True,
        }

        self._base_manager = base.BaseApplication(
            self.REALM, self.HOST, self.HOST_ADDR,
        )
        self._proxy_client = Mock()
        self._s6a_manager = s6a_relay.S6ARelayApplication(
            Mock(),
            self.REALM,
            self.HOST,
            self.HOST_ADDR,
            proxy_client=self._proxy_client,
        )
        self._server = server.S6aServer(
            self._base_manager,
            self._s6a_manager,
            self.REALM,
            self.HOST,
        )

        self.get_proxy_config_patcher = patch.object(
            s6a_relay.ServiceRegistry,
            'get_proxy_config',
            Mock(return_value=proxy_config),
        )
        self.mock_get_proxy_config = self.get_proxy_config_patcher.start()
        self.addCleanup(self.get_proxy_config_patcher.stop)

        # Mock the writes to check responses
        self._writes = Mock()

        def convert_memview_to_bytes(memview):
            """ Deep copy the memoryview for checking later  """
            return self._writes(memview.tobytes())

        self._transport = MockTransport()
        self._transport.write = Mock(side_effect=convert_memview_to_bytes)

        # Here goes nothing..
        self._server.connection_made(self._transport)

    @staticmethod
    def _auth_req(
        user_name, visited_plmn_id, num_request_vectors,
        immediate_response_preferred, resync_info,
    ):
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
        msg.append_avp(avp.AVP('User-Name', user_name))
        msg.append_avp(avp.AVP('Visited-PLMN-Id', visited_plmn_id))
        msg.append_avp(
            avp.AVP(
                'Requested-EUTRAN-Authentication-Info', [
                    avp.AVP(
                        'Number-Of-Requested-Vectors',
                        num_request_vectors,
                    ),
                    avp.AVP(
                        'Immediate-Response-Preferred',
                        1 if immediate_response_preferred else 0,
                    ),
                    avp.AVP('Re-Synchronization-Info', resync_info),
                ],
            ),
        )
        return msg

    @staticmethod
    def _update_location_req(user_name, visited_plmn_id, ulr_flags):
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
        msg.append_avp(avp.AVP('User-Name', user_name))
        msg.append_avp(avp.AVP('Visited-PLMN-Id', visited_plmn_id))
        msg.append_avp(avp.AVP('ULR-Flags', ulr_flags))
        msg.append_avp(avp.AVP('RAT-Type', 1004))
        return msg

    def test_auth_request(self):
        """
        Tests that we convert incoming Diameter AIR to gRPC AuthenticationInformation call
        """
        # Mock out Collect.future
        result = Mock()
        self._proxy_client.AuthenticationInformation.future.side_effect = [
            result,
        ]

        user_name = '1'
        visited_plmn_id = b'(Y'
        num_request_vectors = 1
        immediate_response_preferred = True
        resync_info = b'123456789'

        req = self._auth_req(
            user_name,
            visited_plmn_id,
            num_request_vectors,
            immediate_response_preferred,
            resync_info,
        )
        # Encode request message into buffer
        req_buf = bytearray(req.length)
        req.encode(req_buf, 0)

        request = AuthenticationInformationRequest(
            user_name=user_name,
            visited_plmn=visited_plmn_id,
            num_requested_eutran_vectors=num_request_vectors,
            immediate_response_preferred=immediate_response_preferred,
            resync_info=resync_info,
        )

        self._server.data_received(req_buf)
        self.assertEqual(
            self._proxy_client.AuthenticationInformation.future.call_count,
            1,
        )
        req, _ = self._proxy_client.AuthenticationInformation.future.call_args
        self.assertEqual(repr(request), repr(req[0]))

    def test_auth_answer_success(self):
        """
        Tests that we convert gRPC AuthenticationInformation success response into Diameter AIA
        """
        state_id = 1
        user_name = '1'
        visited_plmn_id = b'(Y'
        num_request_vectors = 2
        immediate_response_preferred = True
        resync_info = b'123456789'
        req = self._auth_req(
            user_name,
            visited_plmn_id,
            num_request_vectors,
            immediate_response_preferred,
            resync_info,
        )

        # response
        rand = b'rand'
        xres = b'xres'
        autn = b'autn'
        kasme = b'kasme'
        auth_info = avp.AVP(
            'Authentication-Info', [
                avp.AVP(
                    'E-UTRAN-Vector', [
                        avp.AVP('RAND', rand),
                        avp.AVP('XRES', xres),
                        avp.AVP('AUTN', autn),
                        avp.AVP('KASME', kasme),
                    ],
                ),
            ] * num_request_vectors,
        )
        resp = self._server._s6a_manager._gen_response(
            state_id, req, avp.ResultCode.DIAMETER_SUCCESS, [auth_info],
        )
        resp_buf = bytearray(resp.length)
        resp.encode(resp_buf, 0)

        result = AuthenticationInformationAnswer(
            error_code=0,
            eutran_vectors=[
                AuthenticationInformationAnswer.EUTRANVector(
                    rand=rand,
                    xres=xres,
                    autn=autn,
                    kasme=kasme,
                ),
            ] * num_request_vectors,
        )
        result_future = unittest.mock.Mock()
        result_future.exception.side_effect = [None]
        result_future.result.side_effect = [result]

        self._server._s6a_manager._relay_auth_answer(
            state_id, req, result_future, 0,
        )
        self._writes.assert_called_once_with(resp_buf)
        self._writes.reset_mock()

    def test_auth_answer_s6a_error(self):
        """
        Tests that we convert gRPC AuthenticationInformation with non-zero ErrorCode to Diameter AIA with error response
        """
        state_id = 1
        user_name = '1'
        visited_plmn_id = b'(Y'
        num_request_vectors = 1
        immediate_response_preferred = True
        resync_info = b'123456789'

        result_info = avp.AVP(
            'Experimental-Result', [
                avp.AVP('Vendor-Id', 10415),
                avp.AVP(
                    'Experimental-Result-Code',
                    avp.ResultCode.DIAMETER_ERROR_USER_UNKNOWN,
                ),
            ],
        )
        req = self._auth_req(
            user_name,
            visited_plmn_id,
            num_request_vectors,
            immediate_response_preferred,
            resync_info,
        )
        resp = self._server._s6a_manager._gen_response(
            state_id, req, avp.ResultCode.DIAMETER_ERROR_USER_UNKNOWN, [
                result_info,
            ],
        )
        resp_buf = bytearray(resp.length)
        resp.encode(resp_buf, 0)

        result = AuthenticationInformationAnswer(
            error_code=ErrorCode.Value('USER_UNKNOWN'),
        )
        result_future = unittest.mock.Mock()
        result_future.exception.side_effect = [None]
        result_future.result.side_effect = [result]

        self._server._s6a_manager._relay_auth_answer(
            state_id, req, result_future, 0,
        )
        self._writes.assert_called_once_with(resp_buf)
        self._writes.reset_mock()

    def test_auth_answer_base_protocol_error(self):
        """
        Tests that we convert gRPC AuthenticationInformation failure to Diameter AIA with base error response
        """
        state_id = 1
        user_name = '1'
        visited_plmn_id = b'(Y'
        num_request_vectors = 1
        immediate_response_preferred = True
        resync_info = b'123456789'

        req = self._auth_req(
            user_name,
            visited_plmn_id,
            num_request_vectors,
            immediate_response_preferred,
            resync_info,
        )
        resp = self._server._s6a_manager._gen_response(
            state_id, req, avp.ResultCode.DIAMETER_UNABLE_TO_COMPLY, [],
        )
        resp_buf = bytearray(resp.length)
        resp.encode(resp_buf, 0)

        grpc_error = unittest.mock.Mock()
        grpc_error.code.side_effect = [grpc.StatusCode.NOT_FOUND]
        grpc_error.details.side_effect = ["Dummy error message."]
        result_future = unittest.mock.Mock()
        result_future.exception.side_effect = [grpc_error]
        result_future.result.side_effect = [None]

        self._server._s6a_manager._relay_auth_answer(
            state_id, req, result_future, 0,
        )
        self._writes.assert_called_once_with(resp_buf)
        self._writes.reset_mock()

    def test_update_location_request(self):
        """
        Tests that we convert incoming Diameter ULR to gRPC UpdateLocation request
        """
        # Mock out Collect.future
        result = Mock()
        self._proxy_client.UpdateLocation.future.side_effect = [result]

        user_name = '1'
        visited_plmn_id = b'(Y'
        ulr_flags = 1 << 2 | 1 << 5
        req = self._update_location_req(user_name, visited_plmn_id, ulr_flags)

        # Encode request message into buffer
        req_buf = bytearray(req.length)
        req.encode(req_buf, 0)

        exp_request = UpdateLocationRequest(
            user_name=user_name,
            visited_plmn=visited_plmn_id,
            skip_subscriber_data=True,
            initial_attach=True,
        )

        self._server.data_received(req_buf)
        self.assertEqual(
            self._proxy_client.UpdateLocation.future.call_count, 1,
        )
        req, _ = self._proxy_client.UpdateLocation.future.call_args
        self.assertEqual(repr(exp_request), repr(req[0]))

    def test_location_update_resp(self):
        """
        Test that gRPC UpdateLocation success response triggers Diameter ULA success
        """
        state_id = 1
        user_name = '1'
        visited_plmn_id = b'(Y'
        ulr_flags = 1 << 2 | 1 << 5
        default_context_id = 0
        total_ambr = {'ul': 10000, 'dl': 50000}
        all_apns_included = True
        req = self._update_location_req(user_name, visited_plmn_id, ulr_flags)

        # Encode request message into buffer
        req_buf = bytearray(req.length)
        req.encode(req_buf, 0)

        apns = [{
            'context_id': i,
            'service_selection': 'apn.%d' % i,
            'qos_profile': {
                'class_id': i,
                'priority_level': i,
                'preemption_capability': True if i % 2 else False,
                'preemption_vulnerability': False if i % 2 else True,
            },
            'ambr': {
                'ul': 1000 * i,
                'dl': 2000 * i,
            },
        } for i in range(2)]

        resp_avps = [
            avp.AVP('ULA-Flags', 1),
            avp.AVP(
                'Subscription-Data', [
                    avp.AVP('MSISDN', b'333608050011'),
                    avp.AVP('Access-Restriction-Data', 47),
                    avp.AVP('Subscriber-Status', 0),
                    avp.AVP('Network-Access-Mode', 2),
                    avp.AVP(
                        'AMBR', [
                            avp.AVP(
                                'Max-Requested-Bandwidth-UL',
                                total_ambr['ul'],
                            ),
                            avp.AVP(
                                'Max-Requested-Bandwidth-DL',
                                total_ambr['dl'],
                            ),
                        ],
                    ),
                    avp.AVP(
                        'APN-Configuration-Profile', [
                            avp.AVP('Context-Identifier', default_context_id),
                            avp.AVP(
                                'All-APN-Configurations-Included-Indicator',
                                1 if all_apns_included else 0,
                            ),
                            *[
                                avp.AVP(
                                    'APN-Configuration', [
                                        avp.AVP(
                                            'Context-Identifier', apn['context_id'],
                                        ),
                                        avp.AVP('PDN-Type', 0),
                                        avp.AVP(
                                            'Service-Selection', apn['service_selection'],
                                        ),
                                        avp.AVP(
                                            'EPS-Subscribed-QoS-Profile', [
                                                avp.AVP(
                                                    'QoS-Class-Identifier',
                                                    apn['qos_profile']['class_id'],
                                                ),
                                                avp.AVP(
                                                    'Allocation-Retention-Priority', [
                                                        avp.AVP(
                                                            'Priority-Level',
                                                            apn['qos_profile']['priority_level'],
                                                        ),
                                                        avp.AVP(
                                                            'Pre-emption-Capability',
                                                            apn['qos_profile']['preemption_capability'],
                                                        ),
                                                        avp.AVP(
                                                            'Pre-emption-Vulnerability',
                                                            apn['qos_profile']['preemption_vulnerability'],
                                                        ),
                                                    ],
                                                ),
                                            ],
                                        ),
                                        avp.AVP(
                                            'AMBR', [
                                                avp.AVP(
                                                    'Max-Requested-Bandwidth-UL',
                                                    apn['ambr']['ul'],
                                                ),
                                                avp.AVP(
                                                    'Max-Requested-Bandwidth-DL',
                                                    apn['ambr']['dl'],
                                                ),
                                            ],
                                        ),
                                    ],
                                ) for apn in apns
                            ],
                        ],
                    ),
                ],
            ),
        ]
        resp = self._server._s6a_manager._gen_response(
            state_id, req, avp.ResultCode.DIAMETER_SUCCESS, resp_avps,
        )
        resp_buf = bytearray(resp.length)
        resp.encode(resp_buf, 0)

        result = UpdateLocationAnswer(
            error_code=0,
            default_context_id=default_context_id,
            total_ambr=UpdateLocationAnswer.AggregatedMaximumBitrate(
                max_bandwidth_ul=total_ambr['ul'],
                max_bandwidth_dl=total_ambr['dl'],
            ),
            msisdn=b'333608050011',
            all_apns_included=all_apns_included,
            apn=[
                UpdateLocationAnswer.APNConfiguration(
                    context_id=apn['context_id'],
                    service_selection=apn['service_selection'],
                    qos_profile=UpdateLocationAnswer.APNConfiguration.QoSProfile(
                        class_id=apn['qos_profile']['class_id'],
                        priority_level=apn['qos_profile']['priority_level'],
                        preemption_capability=apn['qos_profile']['preemption_capability'],
                        preemption_vulnerability=apn['qos_profile']['preemption_vulnerability'],
                    ),
                    ambr=UpdateLocationAnswer.AggregatedMaximumBitrate(
                        max_bandwidth_ul=apn['ambr']['ul'],
                        max_bandwidth_dl=apn['ambr']['dl'],
                    ),
                    pdn=UpdateLocationAnswer.APNConfiguration.IPV4,
                ) for apn in apns
            ],
        )
        result_future = unittest.mock.Mock()
        result_future.exception.side_effect = [None]
        result_future.result.side_effect = [result]

        self._server._s6a_manager._relay_update_location_answer(
            state_id, req, result_future, 0,
        )
        self._writes.assert_called_once_with(resp_buf)
        self._writes.reset_mock()

    def test_ul_answer_s6a_error(self):
        """
        Test that gRPC UpdateLocation response with S6a error triggers Diameter ULA error response
        """
        state_id = 1
        user_name = '1'
        visited_plmn_id = b'(Y'
        ulr_flags = 1 << 2 | 1 << 5
        req = self._update_location_req(user_name, visited_plmn_id, ulr_flags)

        # Encode request message into buffer
        req_buf = bytearray(req.length)
        req.encode(req_buf, 0)

        result_info = avp.AVP(
            'Experimental-Result', [
                avp.AVP('Vendor-Id', 10415),
                avp.AVP(
                    'Experimental-Result-Code',
                    avp.ResultCode.DIAMETER_ERROR_USER_UNKNOWN,
                ),
            ],
        )
        resp = self._server._s6a_manager._gen_response(
            state_id, req, avp.ResultCode.DIAMETER_ERROR_USER_UNKNOWN, [
                result_info,
            ],
        )
        resp_buf = bytearray(resp.length)
        resp.encode(resp_buf, 0)

        result = UpdateLocationAnswer(
            error_code=ErrorCode.Value('USER_UNKNOWN'),
        )
        result_future = unittest.mock.Mock()
        result_future.exception.side_effect = [None]
        result_future.result.side_effect = [result]

        self._server._s6a_manager._relay_update_location_answer(
            state_id, req, result_future, 0,
        )
        self._writes.assert_called_once_with(resp_buf)
        self._writes.reset_mock()

    def test_update_location_answer_base_protocol_error(self):
        """
        Test that gRPC UpdateLocation response failure triggers Diameter base error response
        """
        state_id = 1
        user_name = '1'
        visited_plmn_id = b'(Y'
        ulr_flags = 1 << 2 | 1 << 5
        req = self._update_location_req(user_name, visited_plmn_id, ulr_flags)

        # Encode request message into buffer
        req_buf = bytearray(req.length)
        req.encode(req_buf, 0)

        resp = self._server._s6a_manager._gen_response(
            state_id, req, avp.ResultCode.DIAMETER_UNABLE_TO_COMPLY, [],
        )
        resp_buf = bytearray(resp.length)
        resp.encode(resp_buf, 0)

        grpc_error = unittest.mock.Mock()
        grpc_error.code.side_effect = [grpc.StatusCode.NOT_FOUND]
        grpc_error.details.side_effect = ["Dummy error message."]
        result_future = unittest.mock.Mock()
        result_future.exception.side_effect = [grpc_error]
        result_future.result.side_effect = [None]

        self._server._s6a_manager._relay_update_location_answer(
            state_id, req, result_future, 0,
        )
        self._writes.assert_called_once_with(resp_buf)
        self._writes.reset_mock()


if __name__ == "__main__":
    unittest.main()
