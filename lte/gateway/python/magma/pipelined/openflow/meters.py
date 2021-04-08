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

from magma.pipelined.openflow import messages

LOG = logging.getLogger('openflow.meters')


class MeterClass:
    @staticmethod
    def add_meter(datapath, meter_id: int, rate: int, burst_size: int,
                  retries=3):
        LOG.debug("adding meter_id: %d rate: %d burst_size: %d", meter_id, rate,
                 burst_size)
        bands = []
        ofproto, parser = datapath.ofproto, datapath.ofproto_parser
        dropband = parser.OFPMeterBandDrop(rate=rate, burst_size=burst_size)
        bands.append(dropband)
        mod = parser.OFPMeterMod(datapath=datapath,
                                 command=ofproto.OFPMC_ADD,
                                 flags=ofproto.OFPMF_KBPS,
                                 meter_id=meter_id,
                                 bands=bands)
        messages.send_msg(datapath, mod, retries)

    @staticmethod
    def del_meter(datapath, meter_id: int, retries=3):
        LOG.debug("deleting meter_id: %d", meter_id)
        ofproto, parser = datapath.ofproto, datapath.ofproto_parser
        mod = parser.OFPMeterMod(datapath=datapath,
                                 command=ofproto.OFPMC_DELETE,
                                 meter_id=meter_id)
        messages.send_msg(datapath, mod, retries)

    @staticmethod
    def del_all_meters(datapath, retries=3):
        LOG.debug("deleting all meters")
        ofproto, parser = datapath.ofproto, datapath.ofproto_parser
        mod = parser.OFPMeterMod(datapath=datapath,
                                 command=ofproto.OFPMC_DELETE,
                                 meter_id=ofproto.OFPM_ALL)
        messages.send_msg(datapath, mod, retries)

    @staticmethod
    def dump_all_meters(datapath, retries=3):
        LOG.debug("dumping all meters")
        ofproto, parser = datapath.ofproto, datapath.ofproto_parser
        stat = parser.OFPMeterConfigStatsRequest(datapath=datapath,
                                                 meter_id=ofproto.OFPM_ALL)
        messages.send_msg(datapath, stat, retries)

    @staticmethod
    def dump_meter_features(datapath, retries=3):
        LOG.debug("dumping meter features")
        parser = datapath.ofproto_parser
        stat = parser.OFPMeterFeaturesStatsRequest(datapath=datapath)
        messages.send_msg(datapath, stat, retries)
