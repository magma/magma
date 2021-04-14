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

import abc
from collections import namedtuple
from concurrent.futures import Future

from integ_tests.s1aptests.ovs import LOCALHOST
from integ_tests.s1aptests.ovs.rest_api import get_datapath, get_flows
from ryu.lib import hub

FlowStats = namedtuple('FlowData', ['packets', 'bytes', 'duration_sec',
                                    'cookie'])


def _generate_ryu_req(table_id, match, cookie):
    query = {"table_id": table_id}
    if match is not None:
        query["match"] = match
    else:
        query["match"] = None
    if cookie is not None:
        query["cookie"] = cookie
    return query


class FlowQuery(metaclass=abc.ABCMeta):
    """Flow Lookup interface"""

    @abc.abstractmethod
    def lookup(self):
        """
        Lookup flow rules based on some match criteria
        Returns:
            [FlowStats]: Flow rules information
        """
        raise NotImplementedError()


# REST API is deprecated transition to RyuDirectFlowQuery
class RyuRestFlowQuery(FlowQuery):
    """
    RyuRestFlowQuery uses ryu REST api requests to get ovs flow stats.

    Flows are matched on cookie and match fields. When the FlowQuery is created
    the lookup method can be used to check the stats of the mathed flow.
    """

    def __init__(self, table_id, ovs_ip=LOCALHOST, match=None, cookie=None):
        self._table_id = table_id
        self._datapath = get_datapath(ovs_ip)
        self._ovs_ip = ovs_ip
        self._match = match
        self._cookie = cookie

    def lookup(self, match=None, cookie=None):
        return [
            FlowStats(
                flow["packet_count"], flow["byte_count"], flow["duration_sec"],
                flow["cookie"]
            ) for flow in get_flows(
                self._datapath,
                _generate_ryu_req(self._table_id, self._match, self._cookie),
                self._ovs_ip
            )
        ]


class RyuDirectFlowQuery(FlowQuery):
    """
    RyuDirectFlowQuery
        uses ryu.hub and the test_controller app to get ovs flow stats.

    Flows are matched on cookie and match fields. When the FlowQuery is created
    the lookup method can be used to check the stats of the mathed flow.
    """

    def __init__(self, table_id, test_controller, match=None, cookie=None):
        self._table_id = table_id
        self._tc = test_controller
        self._match = match
        self._cookie = cookie

    def lookup(self):
        queue = hub.Queue()

        def get_stats():
            self._tc.ryu_query_lookup(
                _generate_ryu_req(self._table_id, self._match, self._cookie),
                queue
            )

        hub.joinall([hub.spawn(get_stats)])
        flows = queue.get(block=True)
        return [FlowStats(flow.packet_count, flow.byte_count,
                flow.duration_sec, flow.cookie)
                for flow in flows]

    @staticmethod
    def get_table_stats(test_controller):
        """
        Send an ovs request to retrieve all table stats, wait on queue
        """
        queue = hub.Queue()

        def request_table_stats():
            test_controller.table_stats_lookup(
                queue
            )
        hub.joinall([hub.spawn(request_table_stats)])
        return queue.get(block=True)
