"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from magma.pipelined.app.base import MagmaController, ControllerType
from magma.pipelined.apn import encode_apn



class APNController(MagmaController):
    """
    APN Controller.


    The APN controller is responsible for marking a flow with an encoded APN name. The APN name should be stored in
    register.

    """

    APP_NAME = "apn"
    APP_TYPE = ControllerType.LOGICAL

    def __init__(self, *args, **kwargs):
        super(APNController, self).__init__(*args, **kwargs)
        self._clean_start = True # get from config file
        # set Table_num


    def initialize_on_connect(self, datapath):
        """
        Install the default flows on datapath connect event.

        Args:
            datapath: ryu datapath struct
        """
        # In case wee need to clean all existing  buggy / orphaned flows before start the controller
        # if self._clean_start
        #   self.delete_existing_flows()
        pass

    def add_apn_flow_for_ue(self, ue_ip_addr, apn):
        """ Add flow which match all IN traffic with specified UE_IP and set APN hash in to register.

        Args:
            ue_ip_addr: ip addr allocated for the UE in scope of connection to specific APN
            apn: APN UE is connected to with specified IP addr
        """
        # TODO(119vik): same IP is reused for several bearers connected to the same APN - take care about duplications
        # For all packets which have:
        #  - direction == Direction.IN
        #  - ipv4_src == ue_ip_addr
        #  Set
        #  - registers.APN_TAG_REG to encode_apn(apn)

        pass

    def delete_apn_flow_for_ue(self, ue_ip_addr, apn):
        """ Delete flow been created in scope of add_apn_flow_for_ue.

        Args:
            ue_ip_addr: ip addr allocated for the UE in scope of connection to specific APN
            apn: APN UE is connected to with specified IP addr
        """
        # TODO(119vik): same IP is reused for several bearers connected to the same APN - take care about duplications
        # flow delete
        pass

    def delete_existing_flows(self):
        """ Delete all flows which set APN register"""

        # for flow in flows:
        # delete flow
        pass