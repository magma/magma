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

from feg.protos import s6a_proxy_pb2, s6a_proxy_pb2_grpc
from google.protobuf.json_format import MessageToJson
from magma.subscriberdb import metrics
from magma.subscriberdb.crypto.utils import CryptoError
from magma.subscriberdb.store.base import SubscriberNotFoundError


class S6aProxyRpcServicer(s6a_proxy_pb2_grpc.S6aProxyServicer):
    """
    gRPC based server for the S6aProxy.
    """

    def __init__(self, lte_processor, print_grpc_payload: bool = False):
        self.lte_processor = lte_processor
        logging.info("starting s6a_proxy servicer")
        self._print_grpc_payload = print_grpc_payload

    def add_to_server(self, server):
        """
        Add the servicer to a gRPC server
        """
        s6a_proxy_pb2_grpc.add_S6aProxyServicer_to_server(self, server)

    def AuthenticationInformation(self, request, context):
        self._print_grpc(request)
        imsi = request.user_name
        aia = s6a_proxy_pb2.AuthenticationInformationAnswer()
        try:
            plmn = request.visited_plmn

            re_sync_info = request.resync_info
            # resync_info =
            #  rand + auts, rand is of 16 bytes + auts is of 14 bytes
            sizeof_resync_info = 30
            if re_sync_info and (re_sync_info != b'\x00' * sizeof_resync_info):
                rand = re_sync_info[:16]
                auts = re_sync_info[16:]
                self.lte_processor.resync_lte_auth_seq(imsi, rand, auts)

            rand, xres, autn, kasme = \
                self.lte_processor.generate_lte_auth_vector(imsi, plmn)

            metrics.S6A_AUTH_SUCCESS_TOTAL.inc()

            # Generate and return response message
            aia.error_code = s6a_proxy_pb2.SUCCESS
            eutran_vector = aia.eutran_vectors.add()
            eutran_vector.rand = bytes(rand)
            eutran_vector.xres = xres
            eutran_vector.autn = autn
            eutran_vector.kasme = kasme
            logging.info("Auth success: %s", imsi)
            self._print_grpc(aia)
            return aia

        except CryptoError as e:
            logging.error("Auth error for %s: %s", imsi, e)
            metrics.S6A_AUTH_FAILURE_TOTAL.labels(
                code=metrics.DIAMETER_AUTHENTICATION_REJECTED,
            ).inc()
            aia.error_code = metrics.DIAMETER_AUTHENTICATION_REJECTED
            self._print_grpc(aia)
            return aia

        except SubscriberNotFoundError as e:
            logging.warning("Subscriber not found: %s", e)
            metrics.S6A_AUTH_FAILURE_TOTAL.labels(
                code=metrics.DIAMETER_ERROR_USER_UNKNOWN,
            ).inc()
            aia.error_code = metrics.DIAMETER_ERROR_USER_UNKNOWN
            self._print_grpc(aia)
            return aia

    def UpdateLocation(self, request, context):
        self._print_grpc(request)
        imsi = request.user_name
        ula = s6a_proxy_pb2.UpdateLocationAnswer()
        try:
            profile = self.lte_processor.get_sub_profile(imsi)
        except SubscriberNotFoundError as e:
            ula.error_code = s6a_proxy_pb2.USER_UNKNOWN
            logging.warning('Subscriber not found for ULR: %s', e)
            self._print_grpc(ula)
            return ula

        try:
            sub_data = self.lte_processor.get_sub_data(imsi)
        except SubscriberNotFoundError as e:
            ula.error_code = s6a_proxy_pb2.USER_UNKNOWN
            logging.warning("Subscriber not found for ULR: %s", e)
            return ula
        ula.error_code = s6a_proxy_pb2.SUCCESS
        ula.default_context_id = 0
        ula.total_ambr.max_bandwidth_ul = profile.max_ul_bit_rate
        ula.total_ambr.max_bandwidth_dl = profile.max_dl_bit_rate
        ula.all_apns_included = 0
        ula.msisdn = self.encode_msisdn(sub_data.non_3gpp.msisdn)

        context_id = 0
        for apn in sub_data.non_3gpp.apn_config:
            sec_apn = ula.apn.add()
            sec_apn.context_id = context_id
            context_id += 1
            sec_apn.service_selection = apn.service_selection
            sec_apn.qos_profile.class_id = apn.qos_profile.class_id
            sec_apn.qos_profile.priority_level = apn.qos_profile.priority_level
            sec_apn.qos_profile.preemption_capability = (
                apn.qos_profile.preemption_capability
            )
            sec_apn.qos_profile.preemption_vulnerability = (
                apn.qos_profile.preemption_vulnerability
            )

            sec_apn.ambr.max_bandwidth_ul = apn.ambr.max_bandwidth_ul
            sec_apn.ambr.max_bandwidth_dl = apn.ambr.max_bandwidth_dl
            sec_apn.ambr.unit = (
                s6a_proxy_pb2.UpdateLocationAnswer
                .AggregatedMaximumBitrate.BitrateUnitsAMBR.BPS
            )
            sec_apn.pdn = (
                apn.pdn
                if apn.pdn
                else s6a_proxy_pb2.UpdateLocationAnswer.APNConfiguration.IPV4
            )

        self._print_grpc(ula)
        return ula

    def PurgeUE(self, request, context):
        logging.warning(
            "Purge request not implemented: %s %s",
            request.DESCRIPTOR.full_name, MessageToJson(request),
        )
        res = s6a_proxy_pb2.PurgeUEAnswer()
        self._print_grpc(res)
        return res

    def _print_grpc(self, message):
        if self._print_grpc_payload:
            try:
                log_msg = "{} {}".format(
                    message.DESCRIPTOR.full_name,
                    MessageToJson(message),
                )
                # add indentation
                padding = 2 * ' '
                log_msg = ''.join(
                    "{}{}".format(padding, line)
                    for line in log_msg.splitlines(True)
                )

                log_msg = "GRPC message:\n{}".format(log_msg)
                logging.info(log_msg)
            except Exception as e:  # pylint: disable=broad-except
                logging.debug("Exception while trying to log GRPC: %s", e)

    @staticmethod
    def encode_msisdn(msisdn: str) -> bytes:
        # Mimic how the MSISDN is encoded in ULA : 3GPP TS 29.329-f10
        # For odd length MSISDN pad it with an extra 'F'/'1111'
        if len(msisdn) % 2 != 0:
            msisdn = msisdn + "F"
        result = []
        # Treat each 2 characters as a byte and flip the order
        for i in range(len(msisdn) // 2):
            first = int(msisdn[2 * i])
            second = int(msisdn[2 * i + 1], 16)
            flipped = first + (second << 4)
            result.append(flipped)
        return bytes(result)
