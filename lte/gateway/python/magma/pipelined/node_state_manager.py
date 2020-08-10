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
import datetime

from lte.protos.session_manager_pb2 import UPFNodeState
from lte.protos.session_manager_pb2 import NodeID 
from google.protobuf.timestamp_pb2 import Timestamp
from magma.common.service import MagmaService
from magma.pipelined.set_interface_client import set_interface_node_state_association 

ASSOC_DELETE = -1

class NodeStateManager:
    """
    This controller manages node state information
    and reports to SMF.
    """

    TEID_RANGE_INDICATION = 1
    TEID_RANGE_VALUE = 0
    NETWORK_INSTANCE = "internet"
    NODE_IDENTIFIER = "192.168.200.1"

    def __init__(self, magma_service: MagmaService):
        self._magma_service = magma_service
        self._downlink = magma_service.config.get('enodeb_iface')
        logging.info (" Node State Manager class launched ")
        #self._downlink = self._magma_service.config.get('enode_iface', eth1) 
        self.assoc_message_version = 0
        self.assoc_message_count = 0
        self.node_report_version = 0
        self.node_report_count = 0
        self._recovery_timestamp = Timestamp()
        self._recovery_timestamp.FromDatetime(datetime.datetime.now())

        # to be filled while sending even based node reports
        #self.node_report_thread = hub.spawn(self._monitor, poll_interval)

    def _send_messsage_wrapper (self, node_message):
        set_interface_node_state_association (node_message)

    def _send_association_setup_message (self, version):
        #Build the message
        node_message = UPFNodeState()
        assoc_message = node_message.associaton_state
        assoc_message.node.idtype = NodeID.IPv4

        assoc_message.node.identifier = self.NODE_IDENTIFIER
        assoc_message.state_version = version
        assoc_message.recovery_time_stamp.CopyFrom(self._recovery_timestamp)

        resource_schema = assoc_message.ip_resource_schema.add()
        def get_enodeb_if_ip (iface):
            enode_if_ip = netifaces.ifaddresses(iface)
            return enode_if_ip[netifaces.AF_INET][0]['addr']
  
        resource_schema.ipv4_address = get_enodeb_if_ip(self._downlink)
        resource_schema.teid_range_indication = self.TEID_RANGE_INDICATION
        resource_schema.teid_range = self.TEID_RANGE_VALUE
        resource_schema.assoc_network_instance = self.NETWORK_INSTANCE

        if self._send_messsage_wrapper (node_message) == True:
            if version == ASSOC_DELETE:
                self.assoc_message_version = 0
            else:
                self.assoc_message_version += 1

            self.assoc_message_count += 1

    def send_association_setup_message (self):
        self._send_association_setup_message (self.assoc_message_version + 1)

    def send_association_release_message (self):
        self._send_association_setup_message (ASSOC_DELETE)
