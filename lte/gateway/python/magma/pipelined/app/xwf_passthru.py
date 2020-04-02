"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
from magma.pipelined.app.inout import INGRESS
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch

from .base import MagmaController, ControllerType


class XWFPassthruController(MagmaController):

    APP_NAME = "xwf_passthru"
    APP_TYPE = ControllerType.SPECIAL

    def __init__(self, *args, **kwargs):
        super(XWFPassthruController, self).__init__(*args, **kwargs)
        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)

    def initialize_on_connect(self, datapath):
        self.delete_all_flows(datapath)
        self._install_passthrough_flow(datapath)

    def cleanup_on_disconnect(self, datapath):
        self.delete_all_flows(datapath)

    def delete_all_flows(self, datapath):
        flows.delete_all_flows_from_table(datapath, self.tbl_num)

    def _install_passthrough_flow(self, datapath):
        flows.add_resubmit_next_service_flow(
            datapath,
            self.tbl_num,
            MagmaMatch(),
            actions=[],
            priority=flows.DEFAULT_PRIORITY,
            resubmit_table=self._service_manager.get_table_num(INGRESS),
        )
