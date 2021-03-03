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
from typing import Optional  # noqa
import subprocess
import logging
from lte.protos.policydb_pb2 import FlowMatch
from .types import QosInfo
from .utils import IdManager
from .tc_ops_cmd import run_cmd, TcOpsCmd, argSplit
from .tc_ops_pyroute2 import TcOpsPyRoute2

LOG = logging.getLogger('pipelined.qos.qos_tc_impl')
# LOG.setLevel(logging.DEBUG)

# TODO - replace this implementation with pyroute2 tc
ROOT_QID = 65534
DEFAULT_RATE = '12Kbit'
DEFAULT_INTF_SPEED = '1000'


class TrafficClass:
    """
    Creates/Deletes queues in linux. Using Qdiscs for flow based
    rate limiting(traffic shaping) of user traffic.
    """
    tc_ops = None

    @staticmethod
    def delete_class(intf: str, qid: int, skip_filter=False) -> int:
        qid_hex = hex(qid)

        if not skip_filter:
            TrafficClass.tc_ops.del_filter(intf, qid_hex, qid_hex)

        return TrafficClass.tc_ops.del_htb(intf, qid_hex)

    @staticmethod
    def create_class(intf: str, qid: int, max_bw: int, rate=None,
                     parent_qid=None, skip_filter=False) -> int:
        if not rate:
            rate = DEFAULT_RATE

        if not parent_qid:
            parent_qid = ROOT_QID

        if parent_qid == qid:
            # parent qid should only be self for root case, everything else
            # should be the child of root class
            LOG.error('parent and self qid equal, setting parent_qid to root')
            parent_qid = ROOT_QID

        qid_hex = hex(qid)
        parent_qid_hex = '1:' + hex(parent_qid)
        err = TrafficClass.tc_ops.create_htb(intf, qid_hex, max_bw, rate, parent_qid_hex)
        if err < 0 or skip_filter:
            return err

        # add filter
        return TrafficClass.tc_ops.create_filter(intf, qid_hex, qid_hex)

    @staticmethod
    def init_qdisc(intf: str, show_error=False, enable_pyroute2=False) -> int:
        # TODO: Convert this class into an object.
        if TrafficClass.tc_ops is None:
            if enable_pyroute2:
                TrafficClass.tc_ops = TcOpsPyRoute2()
            else:
                TrafficClass.tc_ops = TcOpsCmd()

        cmd_list = []
        speed = DEFAULT_INTF_SPEED
        qid_hex = hex(ROOT_QID)
        fn = "/sys/class/net/{intf}/speed".format(intf=intf)
        try:
            with open(fn) as f:
                speed = f.read().strip()
        except OSError:
            LOG.error('unable to read speed from %s defaulting to %s', fn, speed)

        # qdisc does not support replace, so check it before creating the HTB qdisc.
        qdisc_type = TrafficClass._get_qdisc_type(intf)
        if qdisc_type != "htb":
            qdisc_cmd = "tc qdisc add dev {intf} root handle 1: htb".format(intf=intf)
            cmd_list.append(qdisc_cmd)
            LOG.info("Created root qdisc")

        parent_q_cmd = "tc class replace dev {intf} parent 1: classid 1:{root_qid} htb "
        parent_q_cmd += "rate {speed}Mbit ceil {speed}Mbit"
        parent_q_cmd = parent_q_cmd.format(intf=intf, root_qid=qid_hex, speed=speed)
        cmd_list.append(parent_q_cmd)
        tc_cmd = "tc class replace dev {intf} parent 1:{root_qid} classid 1:1 htb "
        tc_cmd += "rate {rate} ceil {speed}Mbit"
        tc_cmd = tc_cmd.format(intf=intf, root_qid=qid_hex, rate=DEFAULT_RATE,
                               speed=speed)
        cmd_list.append(tc_cmd)
        return run_cmd(cmd_list, show_error)

    @staticmethod
    def read_all_classes(intf: str):
        qid_list = []
        # example output of this command
        # b'class htb 1:1 parent 1:fffe prio 0 rate 12Kbit ceil 1Gbit burst \
        # 1599b cburst 1375b \nclass htb 1:fffe root rate 1Gbit ceil 1Gbit \
        # burst 1375b cburst 1375b \n'
        # we need to parse this output and extract class ids from here
        tc_cmd = "tc class show dev {}".format(intf)
        args = argSplit(tc_cmd)
        try:
            output = subprocess.check_output(args)
            for ln in output.decode('utf-8').split("\n"):
                ln = ln.strip()
                if not ln:
                    continue
                tok = ln.split()
                if len(tok) < 5:
                    continue

                if tok[1] != "htb":
                    continue

                if tok[3] == 'root':
                    continue

                qid_str = tok[2].split(':')[1]
                qid = int(qid_str, 16)

                pqid_str = tok[4].split(':')[1]
                pqid = int(pqid_str, 16)
                qid_list.append((qid, pqid))
                LOG.debug("TC-dump: %s qid %d pqid %d", ln, qid, pqid)

        except subprocess.CalledProcessError as e:
            LOG.error('failed extracting classids from tc %s', e)
        return qid_list

    @staticmethod
    def dump_class_state(intf: str, qid: int):
        qid_hex = hex(qid)
        tc_cmd = "tc -s -d class show dev {} classid 1:{}".format(intf,
                                                                  qid_hex)
        args = argSplit(tc_cmd)
        try:
            output = subprocess.check_output(args)
            print(output.decode())
        except subprocess.CalledProcessError:
            print("Exception dumping Qos State for %s", intf)

    @staticmethod
    def get_class_rate(intf: str, qid: int) -> Optional[str]:
        qid_hex = hex(qid)
        tc_cmd = "tc class show dev {} classid 1:{}".format(intf, qid_hex)
        args = argSplit(tc_cmd)
        try:
            # output: class htb 1:3 parent 1:2 prio 2 rate 250Kbit ceil 500Kbit burst 1600b cburst 1600b
            raw_output = subprocess.check_output(args)
            output = raw_output.decode('utf-8')
            # return all config from 'rate' onwards
            config = output.split("rate")
            try:
                return config[1]
            except IndexError:
                LOG.error("could not find rate: %s", output)
        except subprocess.CalledProcessError:
            LOG.error("Exception dumping Qos State for %s", tc_cmd)

    @staticmethod
    def _get_qdisc_type(intf: str) -> Optional[str]:
        tc_cmd = "tc qdisc show dev {}".format(intf)
        args = argSplit(tc_cmd)
        try:
            # output: qdisc htb 1: root refcnt 2 r2q 10 default 0 direct_packets_stat 314 direct_qlen 1000
            raw_output = subprocess.check_output(args)
            output = raw_output.decode('utf-8')
            config = output.split()
            try:
                return config[1]
            except IndexError:
                LOG.error("could not qdisc type: %s", output)
        except subprocess.CalledProcessError:
            LOG.error("Exception dumping Qos State for %s", tc_cmd)


class TCManager(object):
    """
    Creates/Deletes queues in linux. Using Qdiscs for flow based
    rate limiting(traffic shaping) of user traffic.
    Queues are created on an egress interface and flows
    in OVS are programmed with qid to filter traffic to the queue.
    Traffic matching a specific flow is filtered to a queue and is
    rate limited based on configured value.
    Traffic to flows with no QoS configuration are sent to a
    default queue and are not rate limited.
    """

    def __init__(self,
                 datapath,
                 loop,
                 config) -> None:
        self._datapath = datapath
        self._loop = loop
        self._uplink = config['nat_iface']
        self._downlink = config['enodeb_iface']
        self._max_rate = config["qos"]["max_rate"]
        self._enable_pyroute2 = config["qos"].get('enable_pyroute2', False)
        self._start_idx, self._max_idx = (config['qos']['linux_tc']['min_idx'],
                                          config['qos']['linux_tc']['max_idx'])
        self._id_manager = IdManager(self._start_idx, self._max_idx)
        self._initialized = True
        LOG.info("Init LinuxTC module uplink:%s downlink:%s",
                 config['nat_iface'], config['enodeb_iface'])

    def destroy(self, ):
        LOG.info("destroying existing qos classes")
        # ensure ordering during deletion of classes, children should be deleted
        # prior to the parent class ids
        for intf in [self._uplink, self._downlink]:
            qid_list = TrafficClass.read_all_classes(intf)
            for qid_tuple in qid_list:
                (qid, pqid) = qid_tuple
                if self._start_idx <= qid < (self._max_idx - 1):
                    LOG.info("Attemting to delete class idx %d", qid)
                    TrafficClass.delete_class(intf, qid)
                if self._start_idx <= pqid < (self._max_idx - 1):
                    LOG.info("Attemting to delete parent class idx %d", pqid)
                    TrafficClass.delete_class(intf, pqid)

    def setup(self, ):
        # initialize new qdisc
        TrafficClass.init_qdisc(self._uplink, enable_pyroute2=self._enable_pyroute2)
        TrafficClass.init_qdisc(self._downlink, enable_pyroute2=self._enable_pyroute2)

    def get_action_instruction(self, qid: int):
        # return an action and an instruction corresponding to this qid
        if qid < self._start_idx or qid > (self._max_idx - 1):
            LOG.error("invalid qid %d, no action/inst returned", qid)
            return None, None

        parser = self._datapath.ofproto_parser
        return parser.OFPActionSetField(pkt_mark=qid), None

    def add_qos(self, d: FlowMatch.Direction, qos_info: QosInfo,
                parent=None, skip_filter=False) -> int:
        LOG.debug("add QoS: %s", qos_info)
        qid = self._id_manager.allocate_idx()
        intf = self._uplink if d == FlowMatch.UPLINK else self._downlink
        err = TrafficClass.create_class(intf, qid, qos_info.mbr,
                                        rate=qos_info.gbr,
                                        parent_qid=parent,
                                        skip_filter=skip_filter)
        # typecast to int to avoid MagicMock related error in unit test
        if int(err) < 0:
            return 0
        LOG.debug("assigned qid: %d", qid)
        return qid

    def remove_qos(self, qid: int, d: FlowMatch.Direction,
                   recovery_mode=False, skip_filter=False):
        if not self._initialized and not recovery_mode:
            return

        if qid < self._start_idx or qid > (self._max_idx - 1):
            LOG.error("invalid qid %d, removal failed", qid)
            return

        LOG.debug("deleting qos_handle %s, skip_filter %s", qid, skip_filter)
        intf = self._uplink if d == FlowMatch.UPLINK else self._downlink

        err = TrafficClass.delete_class(intf, qid, skip_filter)
        if err == 0:
            self._id_manager.release_idx(qid)
        else:
            LOG.error('error deleting class %d, not releasing idx', qid)
        return

    def read_all_state(self, ):
        LOG.debug("read_all_state")
        st = {}
        apn_qid_list = set()
        ul_qid_list = TrafficClass.read_all_classes(self._uplink)
        dl_qid_list = TrafficClass.read_all_classes(self._downlink)
        for (d, qid_list) in ((FlowMatch.UPLINK, ul_qid_list),
                              (FlowMatch.DOWNLINK, dl_qid_list)):
            for qid_tuple in qid_list:
                qid, pqid = qid_tuple
                if qid < self._start_idx or qid > (self._max_idx - 1):
                    LOG.debug("qid %d out of range: (%d - %d)", qid, self._start_idx, self._max_idx)
                    continue
                apn_qid = pqid if pqid != self._max_idx else 0
                st[qid] = {
                    'direction': d,
                    'ambr_qid': apn_qid,
                }
                if apn_qid != 0:
                    apn_qid_list.add(apn_qid)

        self._id_manager.restore_state(st)
        return st, apn_qid_list

    def same_qos_config(self, d: FlowMatch.Direction,
                        qid1: int, qid2: int) -> bool:
        intf = self._uplink if d == FlowMatch.UPLINK else self._downlink

        config1 = TrafficClass.get_class_rate(intf, qid1)
        config2 = TrafficClass.get_class_rate(intf, qid2)
        return config1 == config2

