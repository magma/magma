"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import logging

from magma.subscriberdb import metrics
from magma.subscriberdb.crypto.utils import CryptoError
from magma.subscriberdb.store.base import SubscriberNotFoundError

from feg.protos import s6a_proxy_pb2, s6a_proxy_pb2_grpc

class S6aProxyRpcServicer(s6a_proxy_pb2_grpc.S6aProxyServicer):
    """
    gRPC based server for the S6aProxy.
    """

    def __init__(self, lte_processor):
        """
        Store should be thread-safe since we use a thread pool for requests.
        """
        self.lte_processor = lte_processor
        logging.info("starting s6a_proxy servicer")

    def add_to_server(self, server):
        """
        Add the servicer to a gRPC server
        """
        s6a_proxy_pb2_grpc.add_S6aProxyServicer_to_server(self, server)

    def AuthenticationInformation(self, request, context):
        """
        Adds a subscriber to the store
        """
        imsi = request.user_name
        aia = s6a_proxy_pb2.AuthenticationInformationAnswer()
        try:
            plmn = request.visited_plmn

            re_sync_info = request.resync_info
            #resync_info =
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
            return aia

        except CryptoError as e:
            logging.error("Auth error for %s: %s", imsi, e)
            metrics.S6A_AUTH_FAILURE_TOTAL.labels(
                code=metrics.DIAMETER_AUTHENTICATION_REJECTED).inc()
            aia.error_code = metrics.DIAMETER_AUTHENTICATION_REJECTED
            return aia

        except SubscriberNotFoundError as e:
            logging.warning("Subscriber not found: %s", e)
            metrics.S6A_AUTH_FAILURE_TOTAL.labels(
                code=metrics.DIAMETER_AUTHORIZATION_REJECTED).inc()
            aia.error_code = metrics.DIAMETER_AUTHORIZATION_REJECTED
            return aia

    def UpdateLocation(self, request, context):
        imsi = request.user_name
        ula = s6a_proxy_pb2.UpdateLocationAnswer()
        try:
            profile = self.lte_processor.get_sub_profile(imsi)
        except SubscriberNotFoundError as e:
            ula.error_code = s6a_proxy_pb2.USER_UNKNOWN
            logging.warning('Subscriber not found for ULR: %s', e)
            return ula

        try:
            sub_data = self.lte_processor.get_sub_data(imsi)
        except SubscriberNotFoundError as e:
            ula.error_code = s6a_proxy_pb2.USER_UNKNOWN
            logging.warning('Subscriber not found for ULR: %s', e)
            return ula
        ula.error_code = s6a_proxy_pb2.SUCCESS
        ula.default_context_id = 0
        ula.total_ambr.max_bandwidth_ul = profile.max_ul_bit_rate
        ula.total_ambr.max_bandwidth_dl = profile.max_dl_bit_rate
        ula.all_apns_included = 0

        apn = ula.apn.add()
        apn.context_id = 0
        apn.service_selection = 'oai.ipv4'
        apn.qos_profile.class_id = 9
        apn.qos_profile.priority_level = 15
        apn.qos_profile.preemption_capability = 1
        apn.qos_profile.preemption_vulnerability = 0

        apn.ambr.max_bandwidth_ul = profile.max_ul_bit_rate
        apn.ambr.max_bandwidth_dl = profile.max_dl_bit_rate
        apn.pdn = s6a_proxy_pb2.UpdateLocationAnswer.APNConfiguration.IPV4

        num_apn = len(sub_data.non_3gpp.apn_config)
        for i in range(num_apn):
            apn_ims = ula.apn.add()
            # Context id 0 is assigned to oai.ipv4 apn. So start from 1
            apn_ims.context_id = i+1
            apn_ims.service_selection = sub_data.non_3gpp.apn_config[i].service_selection
            apn_ims.qos_profile.class_id = sub_data.non_3gpp.apn_config[i].qos_profile.class_id
            apn_ims.qos_profile.priority_level = 15
            apn_ims.qos_profile.preemption_capability = 1
            apn_ims.qos_profile.preemption_vulnerability = 0

            apn_ims.ambr.max_bandwidth_ul = profile.max_ul_bit_rate
            apn_ims.ambr.max_bandwidth_dl = profile.max_dl_bit_rate
            apn_ims.pdn = s6a_proxy_pb2.UpdateLocationAnswer.APNConfiguration.IPV4
        return ula
