"""
Copyright (c) 2018-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from collections import defaultdict
from collections import deque
from typing import Dict, List  # noqa
import os
import shlex
import subprocess
import logging


# this code can run in either a docker container(CWAG) or as a native
# python process(AG). When we are running as a root there is no need for
# using sudo. (related to task T63499189 where tc commands failed since 
# sudo wasn't available in the docker container)
def argSplit(cmd: str) -> List[str]:
    args = [] if os.geteuid() == 0 else ["sudo"]
    args.extend(shlex.split(cmd))
    return args


class TrafficClass:
    """
    Creates/Deletes queues in linux. Using Qdiscs for flow based
    rate limiting(traffic shaping) of user traffic.
    """

    @staticmethod
    def delete_class(intf: str, qid: int) -> None:        
        tc_cmd = "tc class del dev {} classid 1:{}".format(intf, qid)
        filter_cmd = "tc filter del dev {} protocol ip parent 1: prio 1 \
                handle {qid} fw flowid 1:{qid}".format(intf, qid=qid)

        args = argSplit(filter_cmd)
        ret = subprocess.call(args)
        logging.debug("add filter ret %d", ret)

        args = argSplit(tc_cmd)
        ret = subprocess.call(args)
        logging.debug("qdisc del q qid %s ret %d", qid, ret)

    @staticmethod
    def create_class(intf: str, qid: int, max_bw: int) -> None:
        tc_cmd = "tc class add dev {} parent 1:fffe classid 1:{} htb \
        rate 12000 ceil {}".format(intf, qid, max_bw)
        qdisc_cmd = "tc qdisc add dev {} parent 1:{} \
                fq_codel".format(intf, qid)
        filter_cmd = "tc filter add dev {} protocol ip parent 1: prio 1 \
                handle {qid} fw flowid 1:{qid}".format(intf, qid=qid)

        args = argSplit(tc_cmd)
        ret = subprocess.call(args)
        logging.debug("create class qid %s ret %d", qid, ret)

        args = argSplit(qdisc_cmd)
        ret = subprocess.call(args)
        logging.debug("create qdisc ret %d", ret)

        args = argSplit(filter_cmd)
        ret = subprocess.call(args)
        logging.debug("add filter ret %d", ret)

    @staticmethod
    def init_qdisc(intf: str) -> None:
        qdisc_cmd = "tc qdisc add dev {} root handle 1: htb".format(intf)
        parent_q_cmd = "tc class add dev {} parent 1: classid 1:fffe htb \
                rate 1Gbit ceil 1Gbit".format(intf)
        tc_cmd = "tc class add dev {} parent 1:fffe classid 1:1 htb \
                rate 12Kbit ceil 1Gbit".format(intf)

        args = argSplit(qdisc_cmd)
        ret = subprocess.call(args)
        logging.debug("qdisc init ret %d", ret)

        args = argSplit(parent_q_cmd)
        ret = subprocess.call(args)
        logging.debug("add class 1: ret %d", ret)

        args = argSplit(tc_cmd)
        ret = subprocess.call(args)
        logging.debug("add class 1:fffe ret %d", ret)

class QosQueueMap:
    """
    Creates/Deletes queues in linux. Using Qdiscs for flow based
    rate limiting(traffic shaping) of user traffic.
    Queues are created on an egress interface and flows
    in OVS are programmed with qid to filter traffic to the queue.
    Traffic matching a specific flow is filtered to a queue and is
    rate limited based on configured value.
    Traffic to flows with no QoS configuration are sent to a
    default queue and are not rate limited.

    (subscriber id, rule_num)->qid mapping is used to
    delete a queue when a flow is removed.
    """
    def __init__(self,
                 nat_iface: str,
                 enodeb_iface: str,
                 enable_qdisc: bool) -> None:
        self._uplink = nat_iface
        self._downlink = enodeb_iface
        self._enable_qdisc = enable_qdisc  # Flag to enable qdisc
        self._qdisc_initialized = False
        self._free_qid = 2
        self._flow_to_queue = defaultdict(dict)  # type: Dict[str, Dict[int, tuple]]
        self._free_qid_list = deque()  # type: deque
        logging.info("Init QoS Module %s", self._uplink)

    def map_flow_to_queue(self,
                          imsi: str,
                          rule_num: int,
                          max_bw: int,
                          is_up_link: bool) -> int:
        """
        Creates a queue in linux and
        adds (imsi, rule_num)->qid mapping

        Args:
            imsi: subscriber id
            rule_num: rule number of the policy
                (imsi, rule_num) uniquely identify a flow
            max_bw: max bandwidth on up link
            is_up_link: Configuration is for uplink or downlink

        Returns:
            Queue ID that should be programmed in the flow
        """
        if not self._enable_qdisc:
            return 0

        if not self._qdisc_initialized:
            TrafficClass.init_qdisc(self._uplink)
            TrafficClass.init_qdisc(self._downlink)
            self._qdisc_initialized = True

        qid = self._create_queue(max_bw, is_up_link)

        self._flow_to_queue[imsi][rule_num] = (qid, is_up_link)
        return qid

    def del_queue_for_flow(self,
                           imsi: str,
                           rule_num: int) -> None:
        """
        Deletes a queue and removes
        (imsi, rule_num)->qid mapping
        Delete the flow using the queue before deleting the flow

        Args:
            imsi: subscriber id
            rule_num: rule number of the policy
                (imsi, rule_num) uniquely identify a flow
        """
        if not self._enable_qdisc:
            return

        if imsi in self._flow_to_queue:
            q = self._flow_to_queue[imsi]
            qid, is_up_link = q.pop(rule_num, (-1, True))
            if qid != -1:
                self._del_queue(qid, is_up_link)

    def del_subscriber_queues(self, imsi: str) -> None:
        """
        Deletes all queues of a subscriber and clears
        (imsi, *)->qid mapping

        Args:
            imsi: subscriber id
        """
        if not self._enable_qdisc:
            return
        if imsi in self._flow_to_queue:
            for rule_num in self._flow_to_queue[imsi]:
                qid, is_up_link = self._flow_to_queue[imsi][rule_num]
                self._del_queue(qid, is_up_link)
            del self._flow_to_queue[imsi]

    def _del_queue(self, qid: int, is_up_link: bool) -> None:
        if is_up_link:
            TrafficClass.delete_class(self._uplink, qid)
        else:
            TrafficClass.delete_class(self._downlink, qid)
        self._add_qid_to_free_list(qid)

    def _create_queue(self, max_bw: int, is_up_link: bool) -> int:
        qid = self._get_qid()
        if is_up_link:
            TrafficClass.create_class(self._uplink, qid, max_bw)
        else:
            TrafficClass.create_class(self._downlink, qid, max_bw)
        return qid

    def _get_qid(self) -> int:
        qid = self._get_free_qid()
        if qid == 0:
            qid = self._free_qid
            self._free_qid += 1

        return qid

    def _get_free_qid(self) -> int:
        if len(self._free_qid_list) != 0:
            qid = self._free_qid_list.popleft()
            return qid
        return 0

    def _add_qid_to_free_list(self, qid: int) -> None:
        self._free_qid_list.append(qid)
