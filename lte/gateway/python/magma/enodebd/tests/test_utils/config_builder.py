"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from lte.protos.mconfig import mconfigs_pb2
from magma.enodebd.devices.device_utils import EnodebDeviceName


class EnodebConfigBuilder:
    @classmethod
    def get_mconfig(
        cls,
        device: EnodebDeviceName = EnodebDeviceName.BAICELLS,
    ) -> mconfigs_pb2.EnodebD:
        mconfig = mconfigs_pb2.EnodebD()
        mconfig.bandwidth_mhz = 20
        mconfig.special_subframe_pattern = 7
        # This earfcndl is actually unused, remove later
        mconfig.earfcndl = 44490
        mconfig.log_level = 1
        mconfig.plmnid_list = "00101"
        mconfig.pci = 260
        mconfig.allow_enodeb_transmit = False
        mconfig.subframe_assignment = 2
        mconfig.tac = 1
        if device is EnodebDeviceName.BAICELLS_QAFB:
            # fdd config
            mconfig.fdd_config.earfcndl = 9211
        else:
            # tdd config
            mconfig.tdd_config.earfcndl = 39150
            mconfig.tdd_config.subframe_assignment = 2
            mconfig.tdd_config.special_subframe_pattern = 7

        return mconfig

    @classmethod
    def get_service_config(cls):
        return {
            "tr069": {
                "interface": "eth1",
                "port": 48080,
                "perf_mgmt_port": 8081,
                "public_ip": "192.88.99.142",
            },
            "reboot_enodeb_on_mme_disconnected": True,
            "s1_interface": "eth1",
        }
