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
import unittest
from unittest.mock import MagicMock, Mock

from magma.pipelined.openflow.messages import MessageHub
from ryu.lib import hub


class MockBarrierRequest(object):
    def __init__(self, datapath):
        self.is_barrier = True
        self.send_msg = MagicMock(return_value=True)

    def set_xid(self, xid):
        self.xid = xid


class MockMessage(object):
    def __init__(self):
        self.is_barrier = False
        self.send_msg = MagicMock(return_value=True)

    def set_xid(self, xid):
        self.xid = xid


class MockDatapath(object):
    def __init__(self, id):
        self._curr_xid = 1
        self.id = id
        self.prev_barrier_xid = None
        self.prev_msg_xid = None
        self.send_msg = MagicMock(return_value=True)

        self.ofproto_parser = Mock()
        self.ofproto_parser.OFPBarrierRequest = MockBarrierRequest

    def set_xid(self, msg):
        xid = self._curr_xid
        msg.set_xid(xid)
        self._curr_xid += 1
        if msg.is_barrier is True:
            self.prev_barrier_xid = xid
        else:
            self.prev_msg_xid = xid
        return xid


class MessageHubTest(unittest.TestCase):
    """
    Tests tracked message sending through Ryu
    """

    def setUp(self):
        self._msg_sender = MessageHub(logging)
        self._mock_datapath = MockDatapath("1234")
        self._mock_datapath

    def test_send_single_msg(self):
        """
        Test sending a single message and receiving a reply
        """
        msg = MockMessage()
        chan = self._msg_sender.send([msg], self._mock_datapath, "1")
        # 1 call for message, 1 for barrier
        self.assertEqual(self._mock_datapath.send_msg.call_count, 2)

        barrier_msg = Mock()
        barrier_msg.xid = self._mock_datapath.prev_barrier_xid
        barrier_msg.datapath = self._mock_datapath
        ev = Mock()
        ev.msg = barrier_msg
        # success = barrier
        self._msg_sender.handle_barrier(ev)
        self._check_reply(chan, "1")

    def test_send_multi_msg(self):
        """
        Test sending multiple messages at once and receiving multiple replies
        """
        msg_list = [MockMessage(), MockMessage(), MockMessage()]
        chan = self._msg_sender.send(msg_list, self._mock_datapath, "1")
        # 3 calls for messages, 1 for barrier
        self.assertEqual(self._mock_datapath.send_msg.call_count, 4)
        ev = self._get_barrier_event()
        # success = barrier
        self._msg_sender.handle_barrier(ev)
        for _ in range(3):
            self._check_reply(chan, "1")

    def test_send_multi_msg_with_err(self):
        """
        Test sending multiple messages, one of which returns an error
        """
        msg_list = [MockMessage(), MockMessage(), MockMessage()]
        chan = self._msg_sender.send(msg_list, self._mock_datapath, "1")
        # 3 calls for messages, 1 for barrier
        self.assertEqual(self._mock_datapath.send_msg.call_count, 4)

        error_ev = self._get_error_event()
        self._msg_sender.handle_error(error_ev)

        ev = self._get_barrier_event()

        self._msg_sender.handle_barrier(ev)
        num_throws = 0
        for _ in range(3):
            reply = chan.get(timeout=1)
            self.assertEqual(reply.txn_id, "1")
            if not reply.ok():
                num_throws += 1
        self.assertEqual(num_throws, 1)

    def test_send_msg_in_parallel(self):
        """
        Test sending messages twice without receiving a barrier
        """
        msg_list1 = [MockMessage(), MockMessage(), MockMessage()]
        msg_list2 = [MockMessage()]
        chan1 = self._msg_sender.send(msg_list1, self._mock_datapath, "1")
        ev1 = self._get_barrier_event()
        chan2 = self._msg_sender.send(msg_list2, self._mock_datapath, "2")
        ev2 = self._get_barrier_event()

        self._msg_sender.handle_barrier(ev1)
        for _ in range(3):
            self._check_reply(chan1, "1")
        self._msg_sender.handle_barrier(ev2)
        self._check_reply(chan2, "2")

    def test_timeout(self):
        """
        Test timeout handling when no response is received
        """
        msg = MockMessage()
        self._msg_sender.send([msg], self._mock_datapath, "1", timeout=0.1)
        hub.sleep(0.2)
        switch = self._msg_sender._switches[self._mock_datapath.id]

        # request and results should be removed
        self.assertEqual(len(switch.results_by_msg), 0)
        self.assertEqual(len(switch.requests_by_barrier), 0)

    def _get_barrier_event(self):
        barrier_msg = Mock()
        barrier_msg.xid = self._mock_datapath.prev_barrier_xid
        barrier_msg.datapath = self._mock_datapath
        ev = Mock()
        ev.msg = barrier_msg
        return ev

    def _get_error_event(self):
        error_msg = Mock()
        error_msg.xid = self._mock_datapath.prev_msg_xid
        error_msg.datapath = self._mock_datapath
        ev = Mock()
        ev.msg = error_msg
        return ev

    def _check_reply(self, chan, txn_id):
        reply = chan.get(timeout=1)
        self.assertEqual(reply.txn_id, txn_id)
        self.assertTrue(reply.ok())


if __name__ == "__main__":
    unittest.main()
