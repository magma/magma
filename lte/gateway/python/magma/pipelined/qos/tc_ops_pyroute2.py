"""
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""


import logging
import pprint
from typing import Optional, Union

from pyroute2 import IPRoute, NetlinkError  # pylint: disable=no-name-in-module

from .tc_ops import TcOpsBase

LOG = logging.getLogger('pipelined.qos.tc_pyroute2')

QUEUE_PREFIX = '1:'
PROTOCOL = 0x0800
PARENT_ID = 0x10000


class TcOpsPyRoute2(TcOpsBase):
    """
    Create TC scheduler and corresponding filter
    """

    def __init__(self):
        self._ipr = IPRoute()
        self._iface_if_index = {}
        LOG.info("initialized")

    def create_htb(
        self, iface: str, qid: str, max_bw: int, rate: str,
        units: str, parent_qid: Optional[str] = None,
    ) -> int:
        """
        Create HTB class for a UE session.

        Args:
            iface: Egress interface name.
            qid: qid number.
            max_bw: ceiling in bits per sec.
            rate: rate limiting.
            units: bit/kbit
            parent_qid: HTB parent queue.

        Returns:
            zero on success.
        """

        LOG.debug("Create HTB iface %s qid %s max_bw %s%s rate %s", iface, qid, max_bw, units, rate)
        try:
            # API needs ceiling in bytes per sec.
            if_index = self._get_if_index(iface)
            htb_queue = QUEUE_PREFIX + qid
            ret = self._ipr.tc(
                "add-class", "htb", if_index,
                htb_queue, parent=parent_qid,
                rate=str(rate).lower(), ceil=str(max_bw) + units, prio=1,
            )
            LOG.debug("Return: %s", ret)
        except (ValueError, NetlinkError) as ex:
            return log_error_and_get_code(ex, "create-htb")
        return 0

    def del_htb(self, iface: str, qid: str) -> int:
        """
        Delete given queue from HTB classed

        Args:
            iface: interface name
            qid: queue-id of the HTB class

        Returns:
        """
        LOG.debug("Delete HTB iface %s qid %s", iface, qid)

        try:
            if_index = self._get_if_index(iface)
            htb_queue = QUEUE_PREFIX + qid

            ret = self._ipr.tc("del-class", "htb", if_index, htb_queue)
            LOG.debug("Return: %s", ret)
        except (ValueError, NetlinkError) as ex:
            return log_error_and_get_code(ex, "del-htb")
        return 0

    def create_filter(self, iface: str, mark: str, qid: str, proto: int = PROTOCOL) -> int:
        """
        Create TC Filter for given HTB class.
        """

        LOG.debug("Create Filter iface %s qid %s", iface, qid)
        try:
            if_index = self._get_if_index(iface)

            class_id = int(PARENT_ID) | int(qid, 16)
            ret = self._ipr.tc(
                "add-filter", "fw", if_index, int(mark, 16),
                parent=PARENT_ID,
                prio=1,
                protocol=proto,
                classid=class_id,
            )
            LOG.debug("Return: %s", ret)

        except (ValueError, NetlinkError) as ex:
            return log_error_and_get_code(ex, "create-filter")
        return 0

    def del_filter(self, iface: str, mark: str, qid: str, proto: int = PROTOCOL) -> int:
        """
        Delete TC filter.
        """

        LOG.debug("Del Filter iface %s qid %s", iface, qid)
        try:
            if_index = self._get_if_index(iface)

            class_id = int(PARENT_ID) | int(qid, 16)

            ret = self._ipr.tc(
                "del-filter", "fw", if_index, int(mark, 16),
                parent=PARENT_ID,
                prio=1,
                protocol=proto,
                classid=class_id,
            )
            LOG.debug("Return: %s", ret)
        except (ValueError, NetlinkError) as ex:
            return log_error_and_get_code(ex, "del-filter")
        return 0

    def create(
        self, iface: str, qid: str, max_bw: int, units: str, rate=None,
        parent_qid: Optional[str] = None, proto=PROTOCOL,
    ) -> int:
        err = self.create_htb(iface, qid, max_bw, rate, units, parent_qid)
        if err:
            return err
        err = self.create_filter(iface, qid, qid, proto)
        if err:
            return err
        return 0

    def delete(self, iface: str, qid: str, proto=PROTOCOL) -> int:
        err = self.del_filter(iface, qid, qid, proto)
        if err:
            return err

        err = self.del_htb(iface, qid)
        if err:
            return err

        return 0

    def _get_if_index(self, iface: str):
        if_index = self._iface_if_index.get(iface, -1)
        if if_index == -1:
            if_index = self._ipr.link_lookup(ifname=iface)
            self._iface_if_index[iface] = if_index

        return if_index

    def _print_classes(self, iface):
        if_index = self._get_if_index(iface)

        pprint.pprint(self._ipr.get_classes(if_index))

    def _print_filters(self, iface):
        if_index = self._get_if_index(iface)

        pprint.pprint(self._ipr.get_filters(if_index))


def log_error_and_get_code(
        ex: Union[ValueError, NetlinkError],
        error_type: str,
) -> int:
    code = getattr(ex, 'code', -1)
    LOG.error("%s error : %s", error_type, code)
    LOG.debug(ex, exc_info=True)
    return code
