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

from feg.protos.s6a_proxy_pb2 import (
    AuthenticationInformationRequest,
    UpdateLocationRequest,
)
from feg.protos.s6a_proxy_pb2_grpc import S6aProxyStub
from magma.common.service_registry import ServiceRegistry
from magma.subscriberdb.metrics import (
    S6A_AUTH_FAILURE_TOTAL,
    S6A_AUTH_SUCCESS_TOTAL,
    S6A_LUR_TOTAL,
)
from magma.subscriberdb.protocols.diameter import avp

from .s6a import S6AApplication


class S6ARelayApplication(S6AApplication):
    """
    As defined in TS 29.272, the 3GPP S6a/S6d application enables the
    transfer of subscriber-related data between the Mobile Management Entity
    (MME) and the Home Subscriber Server (HSS) on the S6a interface and between
    the Serving GPRS Support Node (SGSN) and the Home Subscriber Server
    (HSS) on the S6d interface.
    """
    grpc_timeout = 60

    def __init__(
        self, lte_processor, realm, host, host_ip,
        loop=None, proxy_client=None, retry_count=0,
    ):
        super(S6ARelayApplication, self).__init__(
            lte_processor, realm, host, host_ip, loop,
        )
        self.retry_count = retry_count
        if proxy_client is None:
            chan = ServiceRegistry.get_rpc_channel(
                's6a_proxy',
                ServiceRegistry.CLOUD,
            )
            self._client = S6aProxyStub(chan)
        else:
            self._client = proxy_client

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
        return self._send_auth_with_retries(state_id, msg, self.retry_count)

    def _send_auth_with_retries(self, state_id, msg, retries_left):
        user_name = msg.find_avp(*avp.resolve('User-Name')).value
        visited_plmn = msg.find_avp(*avp.resolve('Visited-PLMN-Id')).value
        request_eutran_info = msg.find_avp(
            *avp.resolve('Requested-EUTRAN-Authentication-Info'),
        )

        num_requested_eutran_vectors = request_eutran_info.find_avp(
            *avp.resolve('Number-Of-Requested-Vectors'),
        ).value
        immediate_response_preferred = request_eutran_info.find_avp(
            *avp.resolve('Immediate-Response-Preferred'),
        ).value
        resync_info = request_eutran_info.find_avp(
            *avp.resolve('Re-Synchronization-Info'),
        )

        request = AuthenticationInformationRequest(
            user_name=user_name,
            visited_plmn=visited_plmn,
            num_requested_eutran_vectors=num_requested_eutran_vectors,
            immediate_response_preferred=immediate_response_preferred,
            resync_info=resync_info.value if resync_info else None,
        )
        future = self._client.AuthenticationInformation.future(
            request,
            self.grpc_timeout,
        )
        future.add_done_callback(
            lambda answer:
            self._loop.call_soon_threadsafe(
                self._relay_auth_answer,
                state_id,
                msg,
                answer,
                retries_left,
            ),
        )

    def _relay_auth_answer(self, state_id, msg, answer_future, retries_left):
        user_name = msg.find_avp(*avp.resolve('User-Name')).value
        err = answer_future.exception()
        if err and retries_left > 0:
            # TODO: retry only on network failure and not application failures
            logging.info(
                "Auth %s Error! [%s] %s, retrying...",
                user_name, err.code(), err.details(),
            )
            self._send_auth_with_retries(
                state_id,
                msg,
                retries_left - 1,
            )
            return
        elif err:
            logging.warning(
                "Auth %s Error! [%s] %s",
                user_name, err.code(), err.details(),
            )
            resp = self._gen_response(
                state_id, msg, avp.ResultCode.DIAMETER_UNABLE_TO_COMPLY,
            )
            S6A_AUTH_FAILURE_TOTAL.labels(
                code=avp.ResultCode.DIAMETER_UNABLE_TO_COMPLY,
            ).inc()
        else:
            answer = answer_future.result()
            error_code = answer.error_code
            if answer.error_code:
                result_info = avp.AVP(
                    'Experimental-Result', [
                        avp.AVP('Vendor-Id', 10415),
                        avp.AVP('Experimental-Result-Code', error_code),
                    ],
                )
                resp = self._gen_response(
                    state_id, msg, error_code, [result_info],
                )
                logging.warning(
                    "Auth S6a %s Error! [%s]",
                    user_name, error_code,
                )
                S6A_AUTH_FAILURE_TOTAL.labels(
                    code=error_code,
                ).inc()
            else:
                auth_info = avp.AVP(
                    'Authentication-Info', [
                        avp.AVP(
                            'E-UTRAN-Vector', [
                                avp.AVP('RAND', vector.rand),
                                avp.AVP('XRES', vector.xres),
                                avp.AVP('AUTN', vector.autn),
                                avp.AVP('KASME', vector.kasme),
                            ],
                        ) for vector in answer.eutran_vectors
                    ],
                )

                resp = self._gen_response(
                    state_id, msg,
                    avp.ResultCode.DIAMETER_SUCCESS,
                    [auth_info],
                )
                S6A_AUTH_SUCCESS_TOTAL.inc()
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
        return self._send_location_request_with_retries(
            state_id, msg, self.retry_count,
        )

    def _send_location_request_with_retries(self, state_id, msg, retries_left):
        user_name = msg.find_avp(*avp.resolve('User-Name')).value
        visited_plmn = msg.find_avp(*avp.resolve('Visited-PLMN-Id')).value
        ulr_flags = msg.find_avp(*avp.resolve('ULR-Flags')).value

        request = UpdateLocationRequest(
            user_name=user_name,
            visited_plmn=visited_plmn,
            skip_subscriber_data=ulr_flags & 1 << 2,
            initial_attach=ulr_flags & 1 << 5,
        )
        future = self._client.UpdateLocation.future(request, self.grpc_timeout)
        future.add_done_callback(
            lambda answer:
            self._loop.call_soon_threadsafe(
                self._relay_update_location_answer,
                state_id,
                msg,
                answer,
                retries_left,
            ),
        )

    def _relay_update_location_answer(
            self, state_id, msg, answer_future, retries_left,
    ):
        err = answer_future.exception()
        if err and retries_left > 0:
            # TODO: retry only on network failure and not application failures
            user_name = msg.find_avp(*avp.resolve('User-Name')).value
            logging.info(
                "Location Update %s Error! [%s] %s, retrying...",
                user_name, err.code(), err.details(),
            )
            self._send_location_request_with_retries(
                state_id,
                msg,
                retries_left - 1,
            )
            return
        elif err:
            user_name = msg.find_avp(*avp.resolve('User-Name')).value
            logging.warning(
                "Location Update %s Error! [%s] %s",
                user_name, err.code(), err.details(),
            )
            resp = self._gen_response(
                state_id, msg, avp.ResultCode.DIAMETER_UNABLE_TO_COMPLY,
            )
        else:
            answer = answer_future.result()
            error_code = answer.error_code
            if error_code:
                result_info = avp.AVP(
                    'Experimental-Result', [
                        avp.AVP('Vendor-Id', 10415),
                        avp.AVP('Experimental-Result-Code', error_code),
                    ],
                )
                resp = self._gen_response(
                    state_id, msg, error_code, [result_info],
                )
            else:

                # Stubbed out Subscription Data from OAI
                subscription_data = avp.AVP(
                    'Subscription-Data', [
                        avp.AVP('MSISDN', answer.msisdn),
                        avp.AVP('Access-Restriction-Data', 47),
                        avp.AVP('Subscriber-Status', 0),
                        avp.AVP('Network-Access-Mode', 2),
                        avp.AVP(
                            'AMBR', [
                                avp.AVP(
                                    'Max-Requested-Bandwidth-UL',
                                    answer.total_ambr.max_bandwidth_ul,
                                ),
                                avp.AVP(
                                    'Max-Requested-Bandwidth-DL',
                                    answer.total_ambr.max_bandwidth_dl,
                                ),
                            ],
                        ),
                        avp.AVP(
                            'APN-Configuration-Profile', [
                                avp.AVP(
                                    'Context-Identifier',
                                    answer.default_context_id,
                                ),
                                avp.AVP(
                                    'All-APN-Configurations-Included-Indicator',
                                    1 if answer.all_apns_included else 0,
                                ),
                                *[
                                    avp.AVP(
                                        'APN-Configuration', [
                                            avp.AVP(
                                                'Context-Identifier',
                                                apn.context_id,
                                            ),
                                            avp.AVP('PDN-Type', apn.pdn),
                                            avp.AVP(
                                                'Service-Selection',
                                                apn.service_selection,
                                            ),
                                            avp.AVP(
                                                'EPS-Subscribed-QoS-Profile', [
                                                    avp.AVP(
                                                        'QoS-Class-Identifier',
                                                        apn.qos_profile.class_id,
                                                    ),
                                                    avp.AVP(
                                                        'Allocation-Retention-Priority', [
                                                            avp.AVP(
                                                                'Priority-Level', apn.qos_profile.priority_level,
                                                            ),
                                                            avp.AVP(
                                                                'Pre-emption-Capability',
                                                                1 if apn.qos_profile.preemption_capability else 0,
                                                            ),
                                                            avp.AVP(
                                                                'Pre-emption-Vulnerability',
                                                                1 if apn.qos_profile.preemption_vulnerability else 0,
                                                            ),
                                                        ],
                                                    ),
                                                ],
                                            ),
                                            avp.AVP(
                                                'AMBR', [
                                                    avp.AVP(
                                                        'Max-Requested-Bandwidth-UL',
                                                        apn.ambr.max_bandwidth_ul,
                                                    ),
                                                    avp.AVP(
                                                        'Max-Requested-Bandwidth-DL',
                                                        apn.ambr.max_bandwidth_dl,
                                                    ),
                                                ],
                                            ),
                                        ],
                                    ) for apn in answer.apn
                                ],
                            ],
                        ),
                    ],
                )

                ula_flags = avp.AVP('ULA-Flags', 1)
                resp = self._gen_response(
                    state_id, msg,
                    avp.ResultCode.DIAMETER_SUCCESS,
                    [ula_flags, subscription_data],
                )
        self.writer.send_msg(resp)
        S6A_LUR_TOTAL.inc()
