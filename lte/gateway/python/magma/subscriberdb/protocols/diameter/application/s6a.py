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

import logging
from enum import IntEnum, unique

from magma.subscriberdb.crypto.utils import CryptoError
from magma.subscriberdb.metrics import (
    S6A_AUTH_FAILURE_TOTAL,
    S6A_AUTH_SUCCESS_TOTAL,
    S6A_LUR_TOTAL,
)
from magma.subscriberdb.protocols.diameter import avp, message
from magma.subscriberdb.store.base import SubscriberNotFoundError

from . import abc


@unique
class S6AApplicationCommands(IntEnum):
    # Command codes defined in this application used in msg header
    AUTHENTICATION_INFORMATION = 318
    UPDATE_LOCATION = 316


class S6AApplication(abc.Application):
    """
    As defined in TS 29.272, the 3GPP S6a/S6d application enables the
    transfer of subscriber-related data between the Mobile Management Entity
    (MME) and the Home Subscriber Server (HSS) on the S6a interface and between
    the Serving GPRS Support Node (SGSN) and the Home Subscriber Server
    (HSS) on the S6d interface.
    """
    # The ID this application uses for messages
    APP_ID = 16777251
    # The Vendor-Specific-Application-Id and VendorId AVPs that
    # the S6a application should advertise
    CAPABILITIES_EXCHANGE_AVPS = [
        avp.AVP('Supported-Vendor-Id', avp.VendorId.TGPP),
        avp.AVP(
            'Vendor-Specific-Application-Id', [
                avp.AVP('Auth-Application-Id', APP_ID),
                avp.AVP('Vendor-Id', avp.VendorId.TGPP),
            ],
        ),
    ]
    # Required fields for requests of each command type
    REQUIRED_FIELDS = {
        S6AApplicationCommands.AUTHENTICATION_INFORMATION:
            [
                'Session-Id',
                'Auth-Session-State',
                'User-Name',
                'Visited-PLMN-Id',
                'Requested-EUTRAN-Authentication-Info',
            ],
        S6AApplicationCommands.UPDATE_LOCATION:
            [
                'Session-Id',
                'Auth-Session-State',
                'User-Name',
                'Visited-PLMN-Id',
                'RAT-Type',
                'ULR-Flags',
            ],
    }

    def __init__(self, lte_processor, realm, host, host_ip, loop=None):
        """Each application has access to a write stream and a collection of
        settings, currently limited to realm and host

        Args:
            lte_processor: A processor instance
            realm: the realm the application should serve
            host: the host name the application should serve
            host_ip: the IP address of the host
        """
        super(S6AApplication, self).__init__(realm, host, host_ip, loop)
        self.lte_processor = lte_processor

    def handle_msg(self, state_id, msg):
        """
        Handle the command of an incoming S6a/S6d request

        Args:
            state_id: the server state identifier
            msg: the message to handle
        Returns:
            None
        """
        if not msg.header.request:
            logging.warning("Received unsolicited answer")
            return

        if msg.header.command_code == \
                S6AApplicationCommands.AUTHENTICATION_INFORMATION:
            self._send_auth(state_id, msg)
        elif msg.header.command_code == S6AApplicationCommands.UPDATE_LOCATION:
            self._send_location_request(state_id, msg)
        else:
            logging.error('Unsupported command: %d', msg.command_code)

    def validate_message(self, state_id, msg):
        """
        Validate a message and send the appropriate error response
        if necessary

        Args:
            state_id: the server state_id
            msg: the message to validate
        Returns:
            True if the message validated
        """
        # Validate we have all required fields
        required_fields = self.REQUIRED_FIELDS[msg.header.command_code]
        if not msg.has_fields(required_fields):
            logging.error(
                "Missing AVP for s6a command %d",
                msg.header.command_code,
            )
            resp = self._gen_response(
                state_id, msg,
                avp.ResultCode.DIAMETER_MISSING_AVP,
            )
            self.writer.send_msg(resp)
            return False
        return True

    def _gen_response(self, state_id, msg, result_code, body_avps=None):
        """
        Generates response message headers to an incoming request and appends
        the response AVPs in the expected order

        Args:
            state_id: the server state id
            msg: the message to respond to
            result_code: the Result-Code of the response
            body_avps: (optional) the AVPs to include in the response body
        Returns:
            a message instance containing the response
        """
        # Generate response message headers
        if body_avps is None:
            body_avps = []
        resp_msg = message.Message.create_response_msg(msg)

        # Session AVPs must come immediately after header RFC3588 8.8
        resp_msg.append_avp(msg.find_avp(*avp.resolve('Session-Id')))

        for body_avp in body_avps:
            resp_msg.append_avp(body_avp)

        # Auth-Session-State is NO_STATE_MAINTAINED (1)
        resp_msg.append_avp(avp.AVP('Auth-Session-State', 1))

        # Host identifiers
        resp_msg.append_avp(avp.AVP('Origin-Host', self.host))
        resp_msg.append_avp(avp.AVP('Origin-Realm', self.realm))
        resp_msg.append_avp(avp.AVP('Origin-State-Id', state_id))

        # Response result
        resp_msg.append_avp(avp.AVP('Result-Code', result_code))
        return resp_msg

    def _send_auth(self, state_id, msg):
        """
        Handles an incoming 3GPP-Authentication-Information-Request
        and writes a 3GPP-Authentication-Information-Answer

        Args:
            state_id: the server state id
            msg: an auth request message
        Returns:
            None
        """
        # Validate the message
        if not self.validate_message(state_id, msg):
            return
        imsi = ""
        try:
            imsi = msg.find_avp(*avp.resolve('User-Name')).value
            plmn = msg.find_avp(*avp.resolve('Visited-PLMN-Id')).value
            request_eutran_info = msg.find_avp(
                *avp.resolve('Requested-EUTRAN-Authentication-Info'),
            )
            re_sync_info = request_eutran_info.find_avp(
                *avp.resolve('Re-Synchronization-Info'),
            )

            if re_sync_info:
                # According to 29.272 7.3.15 this should be concatenation of
                # RAND and AUTS but OAI only sends AUTS so hardcode till fixed
                rand = re_sync_info.value[:16]
                auts = re_sync_info.value[16:]
                self.lte_processor.resync_lte_auth_seq(imsi, rand, auts)

            rand, xres, autn, kasme = \
                self.lte_processor.generate_lte_auth_vector(imsi, plmn)

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
                ],
            )

            S6A_AUTH_SUCCESS_TOTAL.inc()
            resp = self._gen_response(
                state_id, msg,
                avp.ResultCode.DIAMETER_SUCCESS,
                [auth_info],
            )
            logging.info("Auth success: %s", imsi)
        except CryptoError as e:
            S6A_AUTH_FAILURE_TOTAL.labels(
                code=avp.ResultCode.DIAMETER_AUTHENTICATION_REJECTED,
            ).inc()
            resp = self._gen_response(
                state_id, msg, avp.ResultCode.DIAMETER_AUTHENTICATION_REJECTED,
            )
            logging.error("Auth error for %s: %s", imsi, e)
        except SubscriberNotFoundError as e:
            S6A_AUTH_FAILURE_TOTAL.labels(
                code=avp.ResultCode.DIAMETER_ERROR_USER_UNKNOWN,
            ).inc()
            resp = self._gen_response(
                state_id, msg, avp.ResultCode.DIAMETER_ERROR_USER_UNKNOWN,
            )
            logging.warning("Subscriber not found: %s", e)

        self.writer.send_msg(resp)

    def _send_location_request(self, state_id, msg):
        """
        Handles an incoming 3GPP-Update-Location-Request request and writes a
        3GPP-Update-Location-Answer

        Args:
            state_id: the server state id
            msg: an update location request message
        Returns:
            None
        """
        # Validate the message
        if not self.validate_message(state_id, msg):
            return

        ula_flags = avp.AVP('ULA-Flags', 1)

        try:
            imsi = msg.find_avp(*avp.resolve('User-Name')).value
            profile = self.lte_processor.get_sub_profile(imsi)
        except SubscriberNotFoundError as e:
            resp = self._gen_response(
                state_id, msg, avp.ResultCode.DIAMETER_ERROR_USER_UNKNOWN,
            )
            logging.warning('Subscriber not found for ULR: %s', e)
            return

        # Stubbed out Subscription Data from OAI
        subscription_data = avp.AVP(
            'Subscription-Data', [
                avp.AVP('MSISDN', b'333608050011'),
                avp.AVP('Access-Restriction-Data', 47),
                avp.AVP('Subscriber-Status', 0),
                avp.AVP('Network-Access-Mode', 2),
                avp.AVP(
                    'AMBR', [
                        avp.AVP('Max-Requested-Bandwidth-UL', profile.max_ul_bit_rate),
                        avp.AVP('Max-Requested-Bandwidth-DL', profile.max_dl_bit_rate),
                    ],
                ),
                avp.AVP(
                    'APN-Configuration-Profile', [
                        avp.AVP('Context-Identifier', 0),
                        avp.AVP('All-APN-Configurations-Included-Indicator', 0),
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
                                                avp.AVP('Priority-Level', 15),
                                                avp.AVP('Pre-emption-Capability', 1),
                                                avp.AVP('Pre-emption-Vulnerability', 0),
                                            ],
                                        ),
                                    ],
                                ),
                                avp.AVP(
                                    'AMBR', [
                                        avp.AVP(
                                            'Max-Requested-Bandwidth-UL',
                                            profile.max_ul_bit_rate,
                                        ),
                                        avp.AVP(
                                            'Max-Requested-Bandwidth-DL',
                                            profile.max_dl_bit_rate,
                                        ),
                                    ],
                                ),
                            ],
                        ),
                    ],
                ),
            ],
        )

        S6A_LUR_TOTAL.inc()
        resp = self._gen_response(
            state_id, msg,
            avp.ResultCode.DIAMETER_SUCCESS,
            [ula_flags, subscription_data],
        )
        self.writer.send_msg(resp)
