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

from google.protobuf.json_format import MessageToJson
from magma.subscriberdb import metrics
from magma.subscriberdb.crypto.utils import CryptoError
from magma.subscriberdb.store.base import SubscriberNotFoundError
from lte.protos import subscriberauth_pb2_grpc, subscriberauth_pb2

class M5GAuthRpcServicer(subscriberauth_pb2_grpc.M5GSubscriberAuthenticationServicer):
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
        subscriberauth_pb2_grpc.add_M5GSubscriberAuthenticationServicer_to_server(self, server)

    def M5GAuthenticationInformation(self, request, context):
        imsi = request.user_name
        aia = subscriberauth_pb2.M5GAuthenticationInformationAnswer()
        try:

            re_sync_info = request.resync_info
            # resync_info =
            #  rand + auts, rand is of 16 bytes + auts is of 14 bytes
            sizeof_resync_info = 30
            if re_sync_info and (re_sync_info != b'\x00' * sizeof_resync_info):
                rand = re_sync_info[:16]
                auts = re_sync_info[16:]
                self.lte_processor.resync_lte_auth_seq(imsi, rand, auts)

            m5g_ran_auth_vectors = \
                self.lte_processor.generate_m5g_auth_vector(imsi, request.serving_network_name.encode('utf-8'))

            metrics.M5G_AUTH_SUCCESS_TOTAL.inc()

            # Generate and return response message
            aia.error_code = subscriberauth_pb2.SUCCESS
            m5gauth_vector = aia.m5gauth_vectors.add()
            m5gauth_vector.rand = bytes(m5g_ran_auth_vectors.rand)
            m5gauth_vector.xres_star = m5g_ran_auth_vectors.xres_star[16:]
            m5gauth_vector.autn = m5g_ran_auth_vectors.autn
            m5gauth_vector.kseaf = m5g_ran_auth_vectors.kseaf
            return aia

        except CryptoError as e:
            logging.error("Auth error for %s: %s", imsi, e)
            metrics.M5G_AUTH_FAILURE_TOTAL.labels(
                code=metrics.DIAMETER_AUTHENTICATION_REJECTED).inc()
            aia.error_code = metrics.DIAMETER_AUTHENTICATION_REJECTED
            return aia

        except SubscriberNotFoundError as e:
            logging.warning("Subscriber not found: %s", e)
            metrics.M5G_AUTH_FAILURE_TOTAL.labels(
                code=metrics.DIAMETER_ERROR_USER_UNKNOWN).inc()
            aia.error_code = metrics.DIAMETER_ERROR_USER_UNKNOWN
            return aia

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
                logging.warning("Exception while trying to log GRPC: %s", e)
