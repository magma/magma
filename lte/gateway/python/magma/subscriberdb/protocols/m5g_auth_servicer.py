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

from grpc import StatusCode
from lte.protos import (  # type: ignore[attr-defined]
    diam_errors_pb2,
    subscriberauth_pb2,
    subscriberauth_pb2_grpc,
    subscriberdb_pb2,
    subscriberdb_pb2_grpc,
)
from magma.common.rpc_utils import print_grpc, set_grpc_err
from magma.subscriberdb import metrics
from magma.subscriberdb.crypto.ECIES import ECIES_HN
from magma.subscriberdb.crypto.utils import CryptoError
from magma.subscriberdb.store.base import (
    SubscriberNotFoundError,
    SubscriberServerTooBusy,
    SuciProfileNotFoundError,
)
from magma.subscriberdb.subscription.utils import ServiceNotActive


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
        subscriberauth_pb2_grpc.add_M5GSubscriberAuthenticationServicer_to_server(
            self, server,
        )

    def M5GAuthenticationInformation(self, request, context):
        """
        Process M5GAuthenticationInformation Request
        """
        print_grpc(
            request, self._print_grpc_payload,
            "M5GAuthenticationInformation Request:",
        )
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
                self.lte_processor.generate_m5g_auth_vector(
                    imsi,
                    request.serving_network_name.encode(
                        'utf-8',
                    ),
                )

            metrics.M5G_AUTH_SUCCESS_TOTAL.inc()

            # Generate and return response message
            aia.error_code = diam_errors_pb2.SUCCESS
            m5gauth_vector = aia.m5gauth_vectors.add()
            m5gauth_vector.rand = bytes(m5g_ran_auth_vectors.rand)
            m5gauth_vector.xres_star = m5g_ran_auth_vectors.xres_star[16:]
            m5gauth_vector.autn = m5g_ran_auth_vectors.autn
            m5gauth_vector.kseaf = m5g_ran_auth_vectors.kseaf
            return aia

        except CryptoError as e:
            logging.error("Auth error for %s: %s", imsi, e)
            metrics.M5G_AUTH_FAILURE_TOTAL.labels(
                code=metrics.DIAMETER_AUTHENTICATION_REJECTED,
            ).inc()
            aia.error_code = metrics.DIAMETER_AUTHENTICATION_REJECTED
            return aia

        except SubscriberNotFoundError as e:
            logging.warning("Subscriber not found: %s", e)
            metrics.M5G_AUTH_FAILURE_TOTAL.labels(
                code=metrics.DIAMETER_ERROR_USER_UNKNOWN,
            ).inc()
            aia.error_code = metrics.DIAMETER_ERROR_USER_UNKNOWN
            return aia
        except ServiceNotActive as e:
            logging.error("Service not active for %s: %s", imsi, e)
            metrics.M5G_AUTH_FAILURE_TOTAL.labels(
                code=metrics.DIAMETER_ERROR_UNAUTHORIZED_SERVICE,
            ).inc()
            aia.error_code = metrics.DIAMETER_ERROR_UNAUTHORIZED_SERVICE
            return aia
        except SubscriberServerTooBusy as e:
            logging.error("Sqlite3 DB is locked for %s: %s", imsi, e)
            metrics.M5G_AUTH_FAILURE_TOTAL.labels(
                code=metrics.DIAMETER_TOO_BUSY,
            ).inc()
            aia.error_code = metrics.DIAMETER_TOO_BUSY
            return aia
        finally:
            print_grpc(
                aia, self._print_grpc_payload,
                "M5GAuthenticationInformation Response:",
            )


class M5GSUCIRegRpcServicer(subscriberdb_pb2_grpc.M5GSUCIRegistrationServicer):
    """
    gRPC based server for the AMF
    """

    def __init__(
        self, lte_processor, suciprofile_db: dict,
        print_grpc_payload: bool = False,
    ):
        self.lte_processor = lte_processor
        self.suciprofile_db = suciprofile_db
        logging.info("starting M5GSUCIRegRpcServicer servicer")
        self._print_grpc_payload = print_grpc_payload

    def add_to_server(self, server):
        """
        Add the servicer to a gRPC server
        """
        subscriberdb_pb2_grpc.add_M5GSUCIRegistrationServicer_to_server(
            self, server,
        )

    def M5GDecryptMsinSUCIRegistration(self, request, context):
        """
        M5GDecryptMsinSUCIRegistration
        """
        print_grpc(
            request, self._print_grpc_payload,
            "M5GDecryptMsinSUCIRegistration Request:",
        )
        aia = subscriberdb_pb2.M5GSUCIRegistrationAnswer()

        try:
            suciprofile = self.suciprofile_db.get(request.ue_pubkey_identifier)
            if suciprofile is None:
                set_grpc_err(
                    context,
                    StatusCode.NOT_FOUND,
                    f"identifier {request.ue_pubkey_identifier} not found",
                )
                return aia

            if suciprofile.protection_scheme == 0 and len(request.ue_pubkey) == 32:
                profile = 'A'
            elif suciprofile.protection_scheme == 1 and len(request.ue_pubkey) == 33:
                profile = 'B'
            else:
                set_grpc_err(
                    context,
                    StatusCode.INVALID_ARGUMENT,
                    "Public key length or protection scheme is invalid",
                )
                return aia

            home_network_info = ECIES_HN(
                suciprofile.home_network_private_key,
                profile,
            )

            msin_recv = home_network_info.unprotect(
                request.ue_pubkey, request.ue_ciphertext,
                request.ue_encrypted_mac,
            )

            if msin_recv is not None:
                aia.ue_msin_recv = msin_recv[:10]
                logging.info("Deconcealed MSIN: %s", aia.ue_msin_recv)
            else:
                set_grpc_err(
                    context,
                    StatusCode.INVALID_ARGUMENT,
                    "Deconcealing MSIN failed.",
                )
            return aia

        except SuciProfileNotFoundError as e:
            logging.warning("Suciprofile not found: %s", e)
            return aia

        finally:
            print_grpc(
                aia, self._print_grpc_payload,
                "M5GDecryptMsinSUCIRegistration Response:",
            )
