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
from magma.pipelined.openflow.meters import MeterClass
from .utils import IdManager
from .types import QosInfo
import subprocess

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
        self._qos_impl_broken = False
        self._fut = self._loop.create_future()

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

    def add_qos(self, _, qos_info: QosInfo, parent=None, __=False) -> int:
        if self._qos_impl_broken:
            raise RuntimeError(BROKEN_KERN_ERROR_MSG)

        if parent:
            #TODO add ovs meter logic to handle APN AMBR
            pass

        meter_id = self._id_manager.allocate_idx()
        rate_in_kbps = int(qos_info.mbr / 1000)
        MeterClass.add_meter(self._datapath, meter_id, rate=rate_in_kbps,
                             burst_size=0)
        LOG.debug("Adding meter_id %d", meter_id)
        return meter_id

    def remove_qos(self, meter_id: int, d, recovery_mode=False, _=False):
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
        MeterClass.dump_all_meters(self._datapath)
        return {}, []

    def handle_meter_config_stats(self, ev_body):
        LOG.debug("handle_meter_config_stats %s", ev_body)
        meter_id_map = {}
        for stat in ev_body:
            meter_id_map[stat.meter_id] = {
                'direction': 0,
                'ambr_qid': 0,
            }
        self._id_manager.restore_state(meter_id_map)
        self._fut.set_result(meter_id_map)

    def handle_meter_feature_stats(self, ev_body):
        LOG.debug("handle_meter_feature_stats %s", ev_body)
        for stat in ev_body:
            if stat.max_meter == 0:
                self._qos_impl_broken = True
                LOG.error("kernel module has a broken meter implementation")

    @staticmethod
    def dump_meter_state(meter_id):
        try:
            output = subprocess.check_output(["ovs-ofctl", "-O", "OpenFlow15",
                                              "meter-stats", "cwag_br0",
                                              "meter=%s" % meter_id])
            print(output.decode())
        except subprocess.CalledProcessError:
            print("Exception dumping meter state for %s", meter_id)

    # pylint: disable=unused-argument
    def same_qos_config(self, d,
                        qid1: int, qid2: int) -> bool:
        # Once APN AMBR support is added, implement this methode
        return False
