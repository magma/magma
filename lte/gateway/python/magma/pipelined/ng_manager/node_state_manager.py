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
import netifaces
from typing import NamedTuple, Dict

from lte.protos.session_manager_pb2 import (
    UPFNodeState,
    UPFAssociationState,
    UPFFeatureSet,
    UserPlaneIPResourceSchema)

from google.protobuf.timestamp_pb2 import Timestamp
from magma.pipelined.set_interface_client import send_node_state_association_request
from ryu.lib import hub

EXP_BASE = 3

class NodeStateManager:
    """
    This controller manages node state information
    and reports to SMF.
    """

    TEID_RANGE_INDICATION = 0
    TEID_RANGE_VALUE = 0
    ASSOC_MAX_RETRIES = 40

    LocalNodeConfig = NamedTuple(
        'LocalNodeConfig',
        [('downlink_ip', str), ('node_identifier', str)])

    def __init__(self, loop, sessiond_setinterface, config):
        self.config = self._get_config(config)
        self._loop = loop

        #setinterface for sending node association message
        self._sessiond_setinterface = sessiond_setinterface

        #Local counters and state information
        self._smf_assoc_version = 1
        self._assoc_message_count = 0
        self._smf_assoc_state = UPFAssociationState.STARTED

        #Fill the node ID for sending association message
        self._node_id=self.config.node_identifier
        self._assoc_mon_thread = None
        self._recovery_timestamp = Timestamp()

        logging.info(" NGServicer : Node state manager launched ")

    def _get_teid_pool_range(self):
        #TEID_RANGE_INDICATION = 0 shows wild card as per 15.8, section 8.2.82
        return self.TEID_RANGE_INDICATION, self.TEID_RANGE_VALUE

    def _get_config(self, config_dict: Dict) -> NamedTuple:
        def get_enodeb_if_ip(ng_params):
            if ng_params and ng_params.get('downlink_ip_address', None):
                return ng_params['downlink_ip_address']
            enode_if_ip = netifaces.ifaddresses(config_dict['enodeb_iface'])
            return enode_if_ip[netifaces.AF_INET][0]['addr']

        def get_node_identifier(ng_params):
            if ng_params and ng_params.get('node_identifier', None):
                return ng_params['node_identifier']
            return get_enodeb_if_ip(ng_params)

        return self.LocalNodeConfig(
            downlink_ip=get_enodeb_if_ip(config_dict['5G_feature_set']),
            node_identifier=get_node_identifier(config_dict['5G_feature_set'])
        )

    def _send_messsage_wrapper(self, node_message):
        return send_node_state_association_request(node_message,\
                                                   self._sessiond_setinterface)

    def _send_association_request_message(self, assoc_message):
        #Build the message
        node_message=UPFNodeState(upf_id=self._node_id)
        node_message.associaton_state.CopyFrom(assoc_message)

        if self._send_messsage_wrapper(node_message) == True:
            self._smf_assoc_state = assoc_message.assoc_state

            if self._smf_assoc_state != UPFAssociationState.RELEASE:
                self._smf_assoc_version += 1

            self._assoc_message_count += 1
            return True

        return False

    def send_association_setup_message(self):

        teid_range_indicate, teid_range_value =\
                  self._get_teid_pool_range()

        #Create Node association setup message
        assoc_message= \
           UPFAssociationState(
                           state_version=self._smf_assoc_version,
                           assoc_state=UPFAssociationState.ESTABLISHED,
                           feature_set=UPFFeatureSet(f_teid=True),
                           recovery_time_stamp=self._recovery_timestamp.GetCurrentTime(),
                           ip_resource_schema=\
                             [UserPlaneIPResourceSchema(ipv4_address=self.config.downlink_ip,
                                                        teid_range_indication=teid_range_indicate,
                                                        teid_range=teid_range_value)]
                           )

        self._assoc_mon_thread = hub.spawn(self._monitor_association, assoc_message)

    def _monitor_association(self, assoc_message: UPFAssociationState, poll_interval: int = 3):
        """
        Polling to establish smf association
        """
        retry_count = 0
        assoc_established = False

        while assoc_established == False:
            assoc_established =\
                    self._send_association_request_message(assoc_message)

            if assoc_established == False:
                retry_count += 1

                if retry_count == self.ASSOC_MAX_RETRIES:
                    logging.info(" Max Attempt to SMF connection failed. Reattempting..")
                    retry_count = 0

                poll_interval = pow(EXP_BASE, retry_count)
                hub.sleep(poll_interval)

    def send_association_release_message(self):
        # If setup is not established no need to release
        if self._smf_assoc_state != UPFAssociationState.ESTABLISHED:
            return

        #Create Node association release message
        assoc_message=UPFAssociationState(
                            state_version=self._smf_assoc_version+1,
                            assoc_state=UPFAssociationState.RELEASE)

        self._send_association_request_message(assoc_message)
        self._smf_assoc_state = UPFAssociationState.RELEASE
        self._smf_assoc_version = 0

    #In case of restarts
    def get_node_assoc_message(self):
        node_message=UPFNodeState(upf_id=self._node_id)

        teid_range_indicate, teid_range_value =\
                         self._get_teid_pool_range()

        #Create Node association setup message
        assoc_message = \
           UPFAssociationState(
                           state_version=self._smf_assoc_version,
                           assoc_state=UPFAssociationState.ESTABLISHED,
                           feature_set=UPFFeatureSet(f_teid=True),
                           recovery_time_stamp=self._recovery_timestamp.GetCurrentTime(),
                           ip_resource_schema=\
                             [UserPlaneIPResourceSchema(ipv4_address=self.config.downlink_ip,
                                                        teid_range_indication=teid_range_indicate,
                                                        teid_range=teid_range_value)]
                           )

        node_message.associaton_state.CopyFrom(assoc_message)
        return node_message
