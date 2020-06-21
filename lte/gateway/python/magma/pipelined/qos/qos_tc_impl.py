"""
Copyright (c) 2018-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from typing import List  # noqa
import os
import shlex
import subprocess
import logging
from lte.protos.policydb_pb2 import FlowMatch
from .types import QosInfo
from .utils import IdManager

LOG = logging.getLogger('pipelined.qos.qos_tc_impl')


# this code can run in either a docker container(CWAG) or as a native
# python process(AG). When we are running as a root there is no need for
# using sudo. (related to task T63499189 where tc commands failed since
# sudo wasn't available in the docker container)
def argSplit(cmd: str) -> List[str]:
    args = [] if os.geteuid() == 0 else ["sudo"]
    args.extend(shlex.split(cmd))
    return args


def run_cmd(cmd_list, throw_except=False, show_error=True):
    for cmd in cmd_list:
        LOG.debug("running %s", cmd)
        try:
            LOG.debug(cmd)
            args = argSplit(cmd)
            subprocess.check_call(args)
        except subprocess.CalledProcessError as e:
            if show_error:
                LOG.error("error running %s", cmd)
            if throw_except:
                raise e


# TODO - replace this implementation with pyroute2 tc
class TrafficClass:
    """
    Creates/Deletes queues in linux. Using Qdiscs for flow based
    rate limiting(traffic shaping) of user traffic.
    """
    @staticmethod
    def delete_class(intf: str, qid: int, throw_except=True,
                     show_error=True) -> None:
        qid_hex = hex(qid)
        tc_cmd = "tc class del dev {} classid 1:{}".format(intf, qid_hex)
        filter_cmd = "tc filter del dev {} protocol ip parent 1: prio 1 \
                handle {qid} fw flowid 1:{qid}".format(intf, qid=qid_hex)
        run_cmd((filter_cmd, tc_cmd), throw_except, show_error)

    @staticmethod
    def create_class(intf: str, qid: int, max_bw: int, throw_except=True,
                     show_error=True) -> None:
        qid_hex = hex(qid)
        tc_cmd = "tc class add dev {} parent 1:fffe classid 1:{} htb \
        rate 12000 ceil {}".format(intf, qid_hex, max_bw)
        qdisc_cmd = "tc qdisc add dev {} parent 1:{} \
                fq_codel".format(intf, qid_hex)
        filter_cmd = "tc filter add dev {} protocol ip parent 1: prio 1 \
                handle {qid} fw flowid 1:{qid}".format(intf, qid=qid_hex)

        # try to delete if exists
        TrafficClass.delete_class(intf, qid, throw_except=False,
                                  show_error=False)

        # create class
        run_cmd((tc_cmd, qdisc_cmd, filter_cmd), throw_except, show_error)

    @staticmethod
    def init_qdisc(intf: str, throw_except=False, show_error=False) -> None:
        qdisc_cmd = "tc qdisc add dev {} root handle 1: htb".format(intf)
        parent_q_cmd = "tc class add dev {} parent 1: classid 1:fffe htb \
                rate 1Gbit ceil 1Gbit".format(intf)
        tc_cmd = "tc class add dev {} parent 1:fffe classid 1:1 htb \
                rate 12Kbit ceil 1Gbit".format(intf)

        run_cmd((qdisc_cmd, parent_q_cmd, tc_cmd), throw_except, show_error)

    @staticmethod
    def read_all_classes(intf: str):
        tc_cmd = "tc class show dev {}".format(intf)
        args = argSplit(tc_cmd)

        # example output of this command
        # b'class htb 1:1 parent 1:fffe prio 0 rate 12Kbit ceil 1Gbit burst \
        # 1599b cburst 1375b \nclass htb 1:fffe root rate 1Gbit ceil 1Gbit \
        # burst 1375b cburst 1375b \n'
        # we need to parse this output and extract class ids from here
        output = subprocess.check_output(args)
        qid_list = []
        for ln in output.decode('utf-8').split("\n"):
            ln = ln.strip()
            if not ln:
                continue
            tok = ln.split()
            if tok[1] != "htb":
                continue
            qid_str = tok[2].split(':')[1]
            qid = int(qid_str, 16)
            qid_list.append(qid)
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
        self._start_idx, self._max_idx = (config['qos']['linux_tc']['min_idx'],
                                          config['qos']['linux_tc']['max_idx'])
        self._id_manager = IdManager(self._start_idx, self._max_idx)
        self._initialized = True
        LOG.info("Init LinuxTC module uplink:%s downlink:%s",
                 config['nat_iface'], config['enodeb_iface'])

    def destroy(self,):
        LOG.info("destroying existing qos classes")
        ul_qid_list = TrafficClass.read_all_classes(self._uplink)
        dl_qid_list = TrafficClass.read_all_classes(self._downlink)
        for qid in ul_qid_list:
            if qid >= self._start_idx and qid < (self._max_idx - 1):
                LOG.info("deleting class idx %d", qid)
                TrafficClass.delete_class(self._uplink, qid,
                                          throw_except=False, show_error=False)
        for qid in dl_qid_list:
            if qid >= self._start_idx and qid < (self._max_idx - 1):
                LOG.info("deleting class idx %d", qid)
                TrafficClass.delete_class(self._downlink, qid,
                                          throw_except=False, show_error=False)

    def setup(self,):
        # initialize new qdisc
        TrafficClass.init_qdisc(self._uplink)
        TrafficClass.init_qdisc(self._downlink)

    def get_action_instruction(self, qid: int):
        # return an action and an instruction corresponding to this qid
        if qid < self._start_idx or qid > (self._max_idx - 1):
            LOG.error("invalid qid %d, no action/inst returned", qid)
            return

        parser = self._datapath.ofproto_parser
        return (parser.OFPActionSetField(pkt_mark=qid), None)

    def add_qos(self, d: FlowMatch.Direction, qos_info: QosInfo) -> int:
        qid = self._id_manager.allocate_idx()
        intf = self._uplink if d == FlowMatch.UPLINK else self._downlink

        # currently all Qos policies inherit from top level
        # in case we have hierarchy then we should create a child
        # class of the appropriate parent class
        TrafficClass.create_class(intf, qid, qos_info.mbr)

        return qid

    def remove_qos(self, qid: int, d: FlowMatch.Direction,
                   recovery_mode=False):
        if not self._initialized and not recovery_mode:
            return

        if qid < self._start_idx or qid > (self._max_idx - 1):
            LOG.error("invalid qid %d, removal failed", qid)
            return

        LOG.debug("deleting qos_handle %s", qid)
        intf = self._uplink if d == FlowMatch.UPLINK else self._downlink
        TrafficClass.delete_class(intf, qid)
        self._id_manager.release_idx(qid)

    def read_all_state(self, ):
        LOG.debug("read_all_state")
        st = {}
        ul_qid_list = TrafficClass.read_all_classes(self._uplink)
        dl_qid_list = TrafficClass.read_all_classes(self._downlink)

        for(d, qid_list) in ((FlowMatch.UPLINK, ul_qid_list),
                             (FlowMatch.DOWNLINK, dl_qid_list)):
            for qid in qid_list:
                if qid < self._start_idx or qid > (self._max_idx - 1):
                    continue
                st[qid] = d
        self._id_manager.restore_state(st)
        fut = self._loop.create_future()
        LOG.debug("map -> %s", st)
        fut.set_result(st)
        return fut
