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

from ryu.lib import hub

FlowStats = namedtuple(
    'FlowStats', [
        'packets', 'bytes', 'duration_sec',
        'cookie',
    ],
)


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
                queue,
            )

        hub.joinall([hub.spawn(get_stats)])
        flows = queue.get(block=True)
        return [
            FlowStats(
                flow.packet_count, flow.byte_count,
                flow.duration_sec, flow.cookie,
            )
            for flow in flows
        ]

    @staticmethod
    def get_table_stats(test_controller):
        """
        Send an ovs request to retrieve all table stats, wait on queue
        """
        queue = hub.Queue()

        def request_table_stats():
            test_controller.table_stats_lookup(
                queue,
            )
        hub.joinall([hub.spawn(request_table_stats)])
        return queue.get(block=True)
