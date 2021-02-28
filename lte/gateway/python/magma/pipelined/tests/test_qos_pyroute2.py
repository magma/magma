from pyroute2 import IPRoute
from pyroute2 import NetlinkError
from pyroute2 import protocols

import unittest
import socket
import logging
import traceback
import time
import pprint
import subprocess
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.qos.qos_tc_impl import TrafficClass
from magma.pipelined.qos.tc_ops_pyroute2 import TcOpsPyRoute2
from magma.pipelined.qos.tc_ops_cmd import TcOpsCmd

LOG = logging.getLogger('pipelined.qos.tc_rtnl')

QUEUE_PREFIX = '1:'
PROTOCOL = 3


class TcSetypTest(unittest.TestCase):
    BRIDGE = 'testing_qos'
    IFACE = 'dev_qos'

    @classmethod
    def setUpClass(cls):
        BridgeTools.create_bridge(cls.BRIDGE, cls.BRIDGE)
        BridgeTools.create_internal_iface(cls.BRIDGE, cls.IFACE, None)
        TrafficClass.init_qdisc(cls.IFACE, True)

    @classmethod
    def tearDownClass(cls):
        BridgeTools.destroy_bridge(cls.BRIDGE)
        pass

    def check_qid_in_tc(self, qid):
        cmd = "tc filter show dev dev_qos"
        exe_cmd = cmd.split(" ")
        output = subprocess.check_output(exe_cmd)
        found = False
        for ln in output.decode('utf-8').split("\n"):
            ln = ln.strip()
            if not ln:
                continue
            #print(ln)
            tokens = ln.split(" ")

            if len(tokens) > 10 and tokens[9] == qid:
                found = True

        return found

    def test_basic(self):
        cls = self.__class__
        t1 = TcOpsPyRoute2()
        iface = cls.IFACE
        qid = "0xae"
        max_bw = 10000
        rate = 1000
        parent_qid = '1:fffe'

        err1 = t1.create(iface, qid, max_bw, rate, parent_qid)
        self.assertTrue(self.check_qid_in_tc(qid))
        err = t1.delete(iface, qid)
        self.assertFalse(self.check_qid_in_tc(qid))
        self.assertEqual(err, 0)
        self.assertEqual(err1, 0)

    def test_basic6(self):
        cls = self.__class__
        t1 = TcOpsPyRoute2()
        iface = cls.IFACE
        qid = "0xae"
        max_bw = 10000
        rate = 1000
        parent_qid = '1:fffe'

        err1 = t1.create(iface, qid, max_bw, rate, parent_qid, proto=0x86DD)
        self.assertTrue(self.check_qid_in_tc(qid))
        err = t1.delete(iface, qid, proto=0x86DD)
        self.assertFalse(self.check_qid_in_tc(qid))
        self.assertEqual(err, 0)
        self.assertEqual(err1, 0)

    def test_hierarchy(self):
        cls = self.__class__
        t1 = TcOpsPyRoute2()
        # First queue

        iface1 = cls.IFACE
        qid1 = "0xae"
        max_bw = 10000
        rate = 1000
        parent_qid1 = '1:fffe'

        err1 = t1.create(iface1, qid1, max_bw, rate, parent_qid1)
        self.assertTrue(self.check_qid_in_tc(qid1))

        # Second queue

        qid2 = "0x1ae"
        max_bw = 10000
        rate = 1000
        parent_qid2 = '1:' + qid1

        err1 = t1.create(iface1, qid2, max_bw, rate, parent_qid2)
        self.assertTrue(self.check_qid_in_tc(qid2))
        # t1._print_classes(iface1)
        # t1._print_filters(iface1)

        err = t1.delete(iface1, qid2)
        self.assertEqual(err, 0)
        self.assertFalse(self.check_qid_in_tc(qid2))

        err = t1.delete(iface1, qid1)
        self.assertFalse(self.check_qid_in_tc(qid1))

        self.assertEqual(err, 0)
        self.assertEqual(err1, 0)

    def test_mix1(self):
        cls = self.__class__
        t1 = TcOpsPyRoute2()
        t2 = TcOpsCmd()
        iface = cls.IFACE
        qid = "0xae"
        max_bw = 10000
        rate = 1000
        parent_qid = '1:fffe'

        err1 = t1.create(iface, qid, max_bw, rate, parent_qid)
        self.assertTrue(self.check_qid_in_tc(qid))

        err = t2.del_filter(iface, qid, qid)
        self.assertEqual(err, 0)
        err = t2.del_htb(iface, qid)
        self.assertEqual(err, 0)

        self.assertFalse(self.check_qid_in_tc(qid))
        self.assertEqual(err, 0)
        self.assertEqual(err1, 0)

    def test_mix2(self):
        cls = self.__class__
        t1 = TcOpsPyRoute2()
        t2 = TcOpsCmd()
        iface = cls.IFACE
        qid = "0xae"
        max_bw = 10000
        rate = 1000
        parent_qid = '1:fffe'

        err1 = t2.create_htb(iface, qid, max_bw, rate, parent_qid)
        self.assertEqual(err1, 0)
        err1 = t2.create_filter(iface, qid, qid)
        self.assertEqual(err1, 0)

        self.assertTrue(self.check_qid_in_tc(qid))

        err = t1.del_filter(iface, qid, qid)
        self.assertEqual(err, 0)
        err = t1.del_htb(iface, qid)
        self.assertEqual(err, 0)

        self.assertFalse(self.check_qid_in_tc(qid))
        self.assertEqual(err, 0)
        self.assertEqual(err1, 0)


if __name__ == "__main__":
    unittest.main()
