"""
Copyright (c) 2018-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import logging
from magma.pipelined.openflow.meters import MeterClass
from .utils import IdManager
from .types import QosInfo

LOG = logging.getLogger('pipelined.qos.qos_meter_impl')
BROKEN_KERN_ERROR_MSG = "kernel module has a broken meter implementation"


class MeterManager(object):
    def __init__(self, datapath, loop, config):
        LOG.info("Init OVS Meter module")
        self._datapath = datapath
        self._loop = loop
        meter_config = config['qos']['ovs_meter']
        self._start_idx, self._max_idx = (meter_config['min_idx'],
                                          meter_config['max_idx'])
        self._max_rate = config["qos"]["max_rate"]
        self._id_manager = IdManager(self._start_idx, self._max_idx)
        self.initial_meter_id_map = {}
        self._qos_impl_broken = False
        self.fut = self._loop.create_future()

        # dump meter features and check if max_meters = 0
        self.check_broken_kernel_impl()

    def setup(self,):
        # create default meter
        pass

    def destroy(self,):
        LOG.info("Destroying all meters")
        MeterClass.del_all_meters(self._datapath)

    def get_action_instruction(self, meter_id: int):
        if self._qos_impl_broken:
            raise RuntimeError(BROKEN_KERN_ERROR_MSG)

        LOG.debug("get_action_instruction for %d", meter_id)
        if meter_id < self._start_idx or meter_id > (self._max_idx - 1):
            LOG.error("invalid meter_id %d, no action/inst returned", meter_id)
            return

        parser = self._datapath.ofproto_parser
        ofproto = self._datapath.ofproto
        return None, parser.OFPInstructionMeter(meter_id, ofproto.OFPIT_METER)

    def add_qos(self, _, qos_info: QosInfo) -> int:
        if self._qos_impl_broken:
            raise RuntimeError(BROKEN_KERN_ERROR_MSG)

        meter_id = self._id_manager.allocate_idx()
        rate_in_kbps = int(qos_info.mbr / 1000)
        MeterClass.add_meter(self._datapath, meter_id, rate=rate_in_kbps,
                             burst_size=0)
        LOG.debug("Adding meter_id %d", meter_id)
        return meter_id

    def remove_qos(self, meter_id: int, d, recovery_mode=False):
        LOG.debug("Removing meter %d d %d recovery_mode %s", meter_id,
                  d, recovery_mode)
        if meter_id < self._start_idx or meter_id > (self._max_idx - 1):
            LOG.error("invalid meter_id %d, removal failed", meter_id)
            return
        MeterClass.del_meter(self._datapath, meter_id)
        self._id_manager.release_idx(meter_id)

    def check_broken_kernel_impl(self, ):
        LOG.info("check_broken_kernel_impl")
        MeterClass.dump_meter_features(self._datapath)

    def read_all_state(self, ):
        LOG.debug("read_all_state")
        # TODO update ID manager
        MeterClass.dump_all_meters(self._datapath)
        return self.fut

    def handle_meter_config_stats(self, ev_body):
        LOG.debug("handle_meter_config_stats %s", ev_body)
        meter_id_map = {stat.meter_id: 0 for stat in ev_body}
        self._id_manager.restore_state(meter_id_map)
        self.fut.set_result(meter_id_map)

    def handle_meter_feature_stats(self, ev_body):
        LOG.debug("handle_meter_feature_stats %s", ev_body)
        for stat in ev_body:
            if stat.max_meter == 0:
                self._qos_impl_broken = True
                LOG.error("kernel module has a broken meter implementation")
