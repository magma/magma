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

import ipaddress
import shlex
import subprocess
import unittest.mock
from concurrent import futures

import grpc
from lte.protos.mobilityd_pb2 import (
    AllocateIPRequest,
    GWInfo,
    IPAddress,
    IPBlock,
    IPLookupRequest,
    ListAddedIPBlocksResponse,
    ListAllocatedIPsResponse,
    ReleaseIPRequest,
    RemoveIPBlockRequest,
    RemoveIPBlockResponse,
    SubscriberIPTableEntry,
)
from lte.protos.mobilityd_pb2_grpc import MobilityServiceStub
from magma.common.redis.client import get_default_client
from magma.mobilityd.ip_address_man import IPAddressManager
from magma.mobilityd.ip_allocator_pool import IpAllocatorPool
from magma.mobilityd.ipv6_allocator_pool import IPv6AllocatorPool
from magma.mobilityd.mobility_store import MobilityStore
from magma.mobilityd.rpc_servicer import MobilityServiceRpcServicer
from magma.subscriberdb.sid import SIDUtils
from orc8r.protos.common_pb2 import Void


class RpcTests(unittest.TestCase):
    """
    Tests for the IPAllocator rpc servicer and stub
    """

    def setUp(self):
        # Bind the rpc server to a free port
        thread_pool = futures.ThreadPoolExecutor(max_workers=10)
        self._rpc_server = grpc.server(thread_pool)
        port = self._rpc_server.add_insecure_port('0.0.0.0:0')

        store = MobilityStore(get_default_client(), False, 3980)
        store.dhcp_gw_info.read_default_gw()
        ip_allocator = IpAllocatorPool(store)
        ipv6_allocator = IPv6AllocatorPool(
            store,
            session_prefix_alloc_mode='RANDOM',
        )
        self._allocator = IPAddressManager(
            ip_allocator,
            ipv6_allocator,
            store,
        )

        # Add the servicer
        self._servicer = MobilityServiceRpcServicer(self._allocator, False)
        self._servicer.add_to_server(self._rpc_server)
        self._rpc_server.start()

        # Create a rpc stub
        channel = grpc.insecure_channel('0.0.0.0:{}'.format(port))
        self._stub = MobilityServiceStub(channel)

        # variables shared across tests
        self._netaddr = '192.168.0.0'
        self._prefix_len = 28
        ip_bytes = bytes(map(int, self._netaddr.split('.')))
        self._block_msg = IPBlock(
            version=IPBlock.IPV4,
            net_address=ip_bytes,
            prefix_len=self._prefix_len,
        )
        self._ipv6_block = ipaddress.ip_network('fdee:5:6c::/48')
        self._ipv6_netaddr = self._ipv6_block.network_address.packed
        self._ipv6_block_msg = IPBlock(
            version=IPBlock.IPV6,
            net_address=self._ipv6_netaddr,
            prefix_len=self._ipv6_block.prefixlen,
        )
        self._block = ipaddress.ip_network(
            "%s/%s" % (self._netaddr, self._prefix_len),
        )
        self._sid0 = SIDUtils.to_pb('IMSI0')
        self._sid1 = SIDUtils.to_pb('IMSI1')
        self._sid2 = SIDUtils.to_pb('IMSI2')
        self._apn0 = 'Internet'
        self._apn1 = 'IMS'

    def tearDown(self):
        self._rpc_server.stop(0)

    def test_add_invalid_ip_block(self):
        """ adding invalid ipblock should raise INVALID_ARGUMENT """
        block_msg = IPBlock()
        with self.assertRaises(grpc.RpcError) as err:
            self._stub.AddIPBlock(block_msg)
        self.assertEqual(
            err.exception.code(),
            grpc.StatusCode.INVALID_ARGUMENT,
        )

    def test_list_invalid_ip_block(self):
        """ listing invalid ipblock should raise INVALID_ARGUMENT """
        block_msg = IPBlock()
        with self.assertRaises(grpc.RpcError) as err:
            self._stub.ListAllocatedIPs(block_msg)
        self.assertEqual(
            err.exception.code(),
            grpc.StatusCode.INVALID_ARGUMENT,
        )

    def test_add_overlapped_ip_block(self):
        """ overlaping IPBlocks should raise FAILED_PRECONDITION """
        self._stub.AddIPBlock(self._block_msg)

        # overlaped block
        ip_bytes = bytes(map(int, self._netaddr.split('.')))
        block_msg = IPBlock(
            version=IPBlock.IPV4,
            net_address=ip_bytes,
            prefix_len=30,
        )
        with self.assertRaises(grpc.RpcError) as err:
            self._stub.AddIPBlock(block_msg)
        self.assertEqual(
            err.exception.code(),
            grpc.StatusCode.FAILED_PRECONDITION,
        )

    def test_list_added_ip_blocks(self):
        """ List IP blocks added to the allocator """
        # return empty list before adding an IP block
        resp = self._stub.ListAddedIPv4Blocks(Void())
        self.assertEqual(len(resp.ip_block_list), 0)

        # list one assigned IP blocks
        self._stub.AddIPBlock(self._block_msg)
        resp = self._stub.ListAddedIPv4Blocks(Void())
        self.assertEqual(len(resp.ip_block_list), 1)
        self.assertEqual(resp.ip_block_list[0], self._block_msg)

    def test_list_allocated_ips(self):
        """ test list allocated IPs from a IP block """
        self._stub.AddIPBlock(self._block_msg)

        # list empty allocated IPs
        resp = self._stub.ListAllocatedIPs(self._block_msg)
        self.assertEqual(len(resp.ip_list), 0)

        # list after allocating one IP
        request = AllocateIPRequest(
            sid=self._sid0,
            version=AllocateIPRequest.IPV4,
            apn=self._apn0,
        )
        ip_msg0 = self._stub.AllocateIPAddress(request)
        resp = self._stub.ListAllocatedIPs(self._block_msg)
        self.assertNotEqual(resp, None)
        tmp = ListAllocatedIPsResponse()
        tmp.ip_list.extend([ip_msg0.ip_list[0]])
        self.assertEqual(resp, tmp)

    def test_list_allocated_ips_from_unknown_ipblock(self):
        """ test list allocated IPs from an unknown IP block """
        self._stub.AddIPBlock(self._block_msg)

        ip_bytes = bytes(map(int, '10.0.0.0'.split('.')))
        block_msg = IPBlock(
            version=IPBlock.IPV4,
            net_address=ip_bytes,
            prefix_len=30,
        )
        with self.assertRaises(grpc.RpcError) as err:
            self._stub.ListAllocatedIPs(block_msg)
        self.assertEqual(
            err.exception.code(),
            grpc.StatusCode.FAILED_PRECONDITION,
        )

    def test_allocate_ip_address(self):
        """ test AllocateIPAddress and ListAllocatedIPs """
        self._stub.AddIPBlock(self._block_msg)

        # allocate 1st IP
        request = AllocateIPRequest(
            sid=self._sid0,
            version=AllocateIPRequest.IPV4,
            apn=self._apn0,
        )
        ip_msg0 = self._stub.AllocateIPAddress(request)
        self.assertEqual(ip_msg0.ip_list[0].version, AllocateIPRequest.IPV4)
        ip0 = ipaddress.ip_address(ip_msg0.ip_list[0].address)
        self.assertTrue(ip0 in self._block)

        # TODO: uncomment the code below when ip_allocator
        # actually rejects with DuplicatedIPAllocationError

        # with self.assertRaises(grpc.RpcError) as err:
        #     self._stub.AllocateIPAddress(request)
        # self.assertEqual(err.exception.code(),
        #                  grpc.StatusCode.ALREADY_EXISTS)

        # allocate 2nd IP
        request.sid.CopyFrom(self._sid1)
        ip_msg2 = self._stub.AllocateIPAddress(request)
        self.assertEqual(ip_msg2.ip_list[0].version, AllocateIPRequest.IPV4)
        ip2 = ipaddress.ip_address(ip_msg2.ip_list[0].address)
        self.assertTrue(ip2 in self._block)

    def test_multiple_apn_ipallocation(self):
        """ test AllocateIPAddress for multiple APNs """
        self._stub.AddIPBlock(self._block_msg)

        # allocate 1st IP
        request = AllocateIPRequest(
            sid=self._sid0,
            version=AllocateIPRequest.IPV4,
            apn=self._apn0,
        )
        ip_msg0 = self._stub.AllocateIPAddress(request)
        self.assertEqual(ip_msg0.ip_list[0].version, AllocateIPRequest.IPV4)
        ip0 = ipaddress.ip_address(ip_msg0.ip_list[0].address)
        self.assertTrue(ip0 in self._block)

        # allocate 2nd IP from another APN to the same user
        request.apn = self._apn1
        ip_msg1 = self._stub.AllocateIPAddress(request)
        self.assertEqual(ip_msg1.ip_list[0].version, AllocateIPRequest.IPV4)
        ip1 = ipaddress.ip_address(ip_msg1.ip_list[0].address)
        self.assertTrue(ip1 in self._block)

    def test_run_out_of_ip(self):
        """ should raise RESOURCE_EXHAUSTED when running out of IP """
        #  The subnet is provisioned with 16 addresses
        #  Inside ip_address_man.py 11 addresses are reserved,
        #  2 addresses are not usable (all zeros and all ones)
        #  Thus, we have a usable pool of 3 IP addresses;
        #  first three allocations should succeed, while the fourth
        #  request should raise RESOURCE_EXHAUSTED error
        self._stub.AddIPBlock(self._block_msg)

        request = AllocateIPRequest(
            sid=self._sid0,
            version=AllocateIPRequest.IPV4,
            apn=self._apn0,
        )
        self._stub.AllocateIPAddress(request)
        request.apn = self._apn1
        self._stub.AllocateIPAddress(request)

        request.sid.CopyFrom(self._sid1)
        self._stub.AllocateIPAddress(request)

        request.sid.CopyFrom(self._sid2)

        with self.assertRaises(grpc.RpcError) as err:
            self._stub.AllocateIPAddress(request)
        self.assertEqual(
            err.exception.code(),
            grpc.StatusCode.RESOURCE_EXHAUSTED,
        )

    def test_release_ip_address(self):
        """ test ReleaseIPAddress """
        self._stub.AddIPBlock(self._block_msg)

        alloc_request0 = AllocateIPRequest(
            sid=self._sid0,
            version=AllocateIPRequest.IPV4,
            apn=self._apn0,
        )
        ip_msg0 = self._stub.AllocateIPAddress(alloc_request0)
        alloc_request1 = AllocateIPRequest(
            sid=self._sid1,
            version=AllocateIPRequest.IPV4,
            apn=self._apn0,
        )
        ip_msg1 = self._stub.AllocateIPAddress(alloc_request1)

        # release ip_msg0
        release_request0 = ReleaseIPRequest(
            sid=self._sid0,
            ip=ip_msg0.ip_list[0],
            apn=self._apn0,
        )
        resp = self._stub.ReleaseIPAddress(release_request0)
        self.assertEqual(resp, Void())
        resp = self._stub.ListAllocatedIPs(self._block_msg)
        tmp = ListAllocatedIPsResponse()
        tmp.ip_list.extend([ip_msg1.ip_list[0]])
        self.assertEqual(resp, tmp)

        # release ip_msg1
        release_request1 = ReleaseIPRequest(
            sid=self._sid1,
            ip=ip_msg1.ip_list[0],
            apn=self._apn0,
        )
        resp = self._stub.ReleaseIPAddress(release_request1)
        resp = self._stub.ListAllocatedIPs(self._block_msg)
        self.assertEqual(len(resp.ip_list), 0)

    def test_release_unknown_sid_apn_ip_tuple(self):
        """ releasing unknown sid-apn-ip tuple should raise NOT_FOUND """
        self._stub.AddIPBlock(self._block_msg)

        alloc_request0 = AllocateIPRequest(
            sid=self._sid0,
            version=AllocateIPRequest.IPV4,
            apn=self._apn0,
        )
        ip_msg0 = self._stub.AllocateIPAddress(alloc_request0)

        request = ReleaseIPRequest(
            sid=SIDUtils.to_pb("IMSI12345"),
            ip=ip_msg0.ip_list[0],
            apn=self._apn0,
        )
        with self.assertRaises(grpc.RpcError) as err:
            self._stub.ReleaseIPAddress(request)
        self.assertEqual(
            err.exception.code(),
            grpc.StatusCode.NOT_FOUND,
        )

        ip_bytes = bytes(map(int, '10.0.0.0'.split('.')))
        request.ip.CopyFrom(
            IPAddress(
                version=IPAddress.IPV4,
                address=ip_bytes,
            ),
        )
        with self.assertRaises(grpc.RpcError) as err:
            self._stub.ReleaseIPAddress(request)
        self.assertEqual(
            err.exception.code(),
            grpc.StatusCode.NOT_FOUND,
        )

        request = ReleaseIPRequest(
            sid=self._sid0,
            ip=ip_msg0.ip_list[0],
            apn=self._apn1,
        )
        with self.assertRaises(grpc.RpcError) as err:
            self._stub.ReleaseIPAddress(request)
        self.assertEqual(
            err.exception.code(),
            grpc.StatusCode.NOT_FOUND,
        )

    def test_get_ip_for_subscriber(self):
        """ test GetIPForSubscriber """
        self._stub.AddIPBlock(self._block_msg)

        alloc_request0 = AllocateIPRequest(
            sid=self._sid0,
            version=AllocateIPRequest.IPV4,
            apn=self._apn0,
        )
        ip_msg0 = self._stub.AllocateIPAddress(alloc_request0)
        ip0 = ipaddress.ip_address(ip_msg0.ip_list[0].address)

        lookup_request0 = IPLookupRequest(
            sid=self._sid0,
            apn=self._apn0,
            version=IPAddress.IPV4,
        )
        ip_msg0_returned = self._stub.GetIPForSubscriber(lookup_request0)
        ip0_returned = ipaddress.ip_address(ip_msg0_returned.address)
        self.assertEqual(ip0, ip0_returned)

    def test_get_ip_for_unknown_subscriber(self):
        """ Getting ip for non existent subscriber should return NOT_FOUND
        status code """
        lookup_request0 = IPLookupRequest(
            sid=self._sid0,
            apn=self._apn0,
            version=IPAddress.IPV4,
        )
        with self.assertRaises(grpc.RpcError) as err:
            self._stub.GetIPForSubscriber(lookup_request0)
        self.assertEqual(
            err.exception.code(),
            grpc.StatusCode.NOT_FOUND,
        )

    def test_get_gw_info(self):
        def_gw_cmd = shlex.split("ip route show")
        p = subprocess.Popen(def_gw_cmd, stdout=subprocess.PIPE)
        output = p.stdout.read().decode("utf-8")
        def_ip = None
        for line in output.splitlines():
            if 'default ' in line:
                tokens = line.split()
                def_ip = tokens[2]
                break

        gw_info_list = self._stub.ListGatewayInfo(Void())
        for gw_info in gw_info_list.gw_list:
            gw_ip_get = str(ipaddress.ip_address(gw_info.ip.address))
            if gw_ip_get == def_ip:
                return

        assert 0

    def test_set_gw_info(self):
        mac1 = "22:22:c6:d0:02:3c"
        ipaddr1 = ipaddress.ip_address("10.1.1.11")
        gwinfo_msg = GWInfo()
        gwinfo_msg.ip.version = IPBlock.IPV4
        gwinfo_msg.ip.address = ipaddr1.packed
        gwinfo_msg.mac = mac1
        gwinfo_msg.vlan = ""

        self._stub.SetGatewayInfo(gwinfo_msg)
        gw_info_list = self._stub.ListGatewayInfo(Void())
        for gw_info in gw_info_list.gw_list:
            gw_ip_get = ipaddress.ip_address(gw_info.ip.address)
            self.assertEqual(ipaddr1, gw_ip_get)
            self.assertEqual(mac1, gw_info.mac)

    def test_set_gw_info_vlan(self):
        mac1 = "22:22:c6:d0:02:3c"
        ipaddr1 = ipaddress.ip_address("10.1.1.11")
        gwinfo_msg = GWInfo(
            ip=IPAddress(
                version=IPAddress.IPV4,
                address=ipaddr1.packed,
            ),
            mac=mac1, vlan="1",
        )

        self._stub.SetGatewayInfo(gwinfo_msg)
        gw_info_list = self._stub.ListGatewayInfo(Void())
        for gw_info in gw_info_list.gw_list:
            if gw_info.vlan == "1":
                gw_ip_get = ipaddress.ip_address(gw_info.ip.address)
                self.assertEqual(ipaddr1, gw_ip_get)
                self.assertEqual(mac1, gw_info.mac)
                return

        assert 0

    def test_set_gw_info_vlan2(self):
        mac1 = "22:22:c6:d0:02:3c"
        ipaddr1 = ipaddress.ip_address("10.1.1.11")
        gwinfo_msg1 = GWInfo(
            ip=IPAddress(
                version=IPAddress.IPV4,
                address=ipaddr1.packed,
            ),
            mac=mac1, vlan="1",
        )

        self._stub.SetGatewayInfo(gwinfo_msg1)

        mac2 = "33:22:c6:d0:02:3c"
        ipaddr2 = ipaddress.ip_address("20.1.1.11")
        gwinfo_msg2 = GWInfo(
            ip=IPAddress(
                version=IPAddress.IPV4,
                address=ipaddr2.packed,
            ),
            mac=mac2, vlan="2",
        )

        self._stub.SetGatewayInfo(gwinfo_msg2)

        gw_info_list = self._stub.ListGatewayInfo(Void())
        count1 = 0
        count2 = 0
        for gw_info in gw_info_list.gw_list:
            if gw_info == gwinfo_msg1:
                count1 = count1 + 1
            if gw_info == gwinfo_msg2:
                count2 = count2 + 1

        self.assertEqual(count1, 1)
        self.assertEqual(count2, 1)

    def test_get_subscriber_id_from_ip(self):
        """ test GetSubscriberIDFromIP """
        self._stub.AddIPBlock(self._block_msg)
        alloc_request0 = AllocateIPRequest(
            sid=self._sid0,
            version=AllocateIPRequest.IPV4,
            apn=self._apn0,
        )
        ip_msg0 = self._stub.AllocateIPAddress(alloc_request0)
        sid_pb_returned = self._stub.GetSubscriberIDFromIP(ip_msg0.ip_list[0])
        self.assertEqual(
            SIDUtils.to_str(self._sid0),
            SIDUtils.to_str(sid_pb_returned),
        )

    def test_get_subscriber_id_from_unknown_ip(self):
        """
        Getting subscriber id for non-allocated ip address should return
        NOT_FOUND error code
        """
        ip_pb = IPAddress(
            version=IPAddress.IPV4,
            address=ipaddress.ip_address('1.1.1.1').packed,
        )
        with self.assertRaises(grpc.RpcError) as err:
            self._stub.GetSubscriberIDFromIP(ip_pb)
        self.assertEqual(
            err.exception.code(),
            grpc.StatusCode.NOT_FOUND,
        )

    def test_get_subscriber_ip_table(self):
        """ test GetSubscriberIPTable """
        self._stub.AddIPBlock(self._block_msg)

        resp = self._stub.GetSubscriberIPTable(Void())
        self.assertEqual(len(resp.entries), 0)

        alloc_request0 = AllocateIPRequest(
            sid=self._sid0,
            version=AllocateIPRequest.IPV4,
            apn=self._apn0,
        )
        ip_msg0 = self._stub.AllocateIPAddress(alloc_request0)
        entry0 = SubscriberIPTableEntry(
            sid=self._sid0,
            ip=ip_msg0.ip_list[0],
            apn=self._apn0,
        )
        resp = self._stub.GetSubscriberIPTable(Void())
        self.assertTrue(entry0 in resp.entries)

        alloc_request1 = AllocateIPRequest(
            sid=self._sid1,
            version=AllocateIPRequest.IPV4,
            apn=self._apn1,
        )
        ip_msg1 = self._stub.AllocateIPAddress(alloc_request1)
        entry1 = SubscriberIPTableEntry(
            sid=self._sid1,
            ip=ip_msg1.ip_list[0],
            apn=self._apn1,
        )
        resp = self._stub.GetSubscriberIPTable(Void())
        self.assertTrue(entry0 in resp.entries)
        self.assertTrue(entry1 in resp.entries)

        # keep in table after in release
        release_request0 = ReleaseIPRequest(
            sid=self._sid0,
            ip=ip_msg0.ip_list[0],
            apn=self._apn0,
        )
        resp = self._stub.ReleaseIPAddress(release_request0)
        resp = self._stub.GetSubscriberIPTable(Void())
        self.assertTrue(entry0 in resp.entries)
        self.assertTrue(entry1 in resp.entries)

    def test_remove_no_assigned_blocks(self):
        """ remove should return nothing """
        remove_request0 = RemoveIPBlockRequest(
            ip_blocks=[self._block_msg],
            force=False,
        )
        resp = self._stub.RemoveIPBlock(remove_request0)

        expect = RemoveIPBlockResponse()
        self.assertEqual(expect, resp)

        resp = self._stub.ListAddedIPv4Blocks(Void())
        expect = ListAddedIPBlocksResponse()
        self.assertEqual(expect, resp)

    def test_remove_unallocated_assigned_block(self):
        """ remove should return nothing """
        self._stub.AddIPBlock(self._block_msg)

        remove_request0 = RemoveIPBlockRequest(
            ip_blocks=[self._block_msg],
            force=False,
        )
        resp = self._stub.RemoveIPBlock(remove_request0)

        expect = RemoveIPBlockResponse()
        expect.ip_blocks.extend([self._block_msg])
        self.assertEqual(expect, resp)

        resp = self._stub.ListAddedIPv4Blocks(Void())
        expect = ListAddedIPBlocksResponse()
        self.assertEqual(expect, resp)

    def test_remove_allocated_assigned_block_without_force(self):
        """ remove should return nothing """
        self._stub.AddIPBlock(self._block_msg)

        alloc_request0 = AllocateIPRequest(
            sid=self._sid0,
            version=AllocateIPRequest.IPV4,
            apn=self._apn0,
        )
        self._stub.AllocateIPAddress(alloc_request0)

        remove_request0 = RemoveIPBlockRequest(
            ip_blocks=[self._block_msg],
            force=False,
        )
        resp = self._stub.RemoveIPBlock(remove_request0)

        expect = RemoveIPBlockResponse()
        self.assertEqual(expect, resp)

        resp = self._stub.ListAddedIPv4Blocks(Void())
        expect = ListAddedIPBlocksResponse()
        expect.ip_block_list.extend([self._block_msg])
        self.assertEqual(expect, resp)

    def test_remove_allocated_assigned_block_with_force(self):
        """ remove should return nothing """
        self._stub.AddIPBlock(self._block_msg)

        alloc_request0 = AllocateIPRequest(
            sid=self._sid0,
            version=AllocateIPRequest.IPV4,
            apn=self._apn0,
        )
        self._stub.AllocateIPAddress(alloc_request0)

        remove_request0 = RemoveIPBlockRequest(
            ip_blocks=[self._block_msg],
            force=True,
        )
        resp = self._stub.RemoveIPBlock(remove_request0)

        expect = RemoveIPBlockResponse()
        expect.ip_blocks.extend([self._block_msg])
        self.assertEqual(expect, resp)

        resp = self._stub.ListAddedIPv4Blocks(Void())
        expect = ListAddedIPBlocksResponse()
        self.assertEqual(expect, resp)

    def test_remove_after_releasing_all_addresses(self):
        """ remove after releasing all addresses should remove block """
        # Assign IP block
        self._stub.AddIPBlock(self._block_msg)

        # Allocate 2 IPs
        alloc_request0 = AllocateIPRequest(
            sid=self._sid0,
            version=AllocateIPRequest.IPV4,
            apn=self._apn0,
        )
        ip_msg0 = self._stub.AllocateIPAddress(alloc_request0)

        alloc_request1 = AllocateIPRequest(
            sid=self._sid1,
            version=AllocateIPRequest.IPV4,
            apn=self._apn0,
        )
        ip_msg1 = self._stub.AllocateIPAddress(alloc_request1)

        # Test remove without force -- should not remove block
        remove_request0 = RemoveIPBlockRequest(
            ip_blocks=[self._block_msg],
            force=False,
        )
        resp = self._stub.RemoveIPBlock(remove_request0)

        expect = RemoveIPBlockResponse()
        self.assertEqual(expect, resp)

        # Ensure that block has not been removed
        resp = self._stub.ListAddedIPv4Blocks(Void())
        expect = ListAddedIPBlocksResponse()
        expect.ip_block_list.extend([self._block_msg])
        self.assertEqual(expect, resp)

        # Release the allocated IPs
        release_request0 = ReleaseIPRequest(
            sid=self._sid0,
            ip=ip_msg0.ip_list[0],
            apn=self._apn0,
        )
        resp = self._stub.ReleaseIPAddress(release_request0)
        self.assertEqual(resp, Void())

        release_request1 = ReleaseIPRequest(
            sid=self._sid1,
            ip=ip_msg1.ip_list[0],
            apn=self._apn0,
        )
        resp = self._stub.ReleaseIPAddress(release_request1)
        self.assertEqual(resp, Void())

        # Test remove without force -- should remove block
        remove_request1 = RemoveIPBlockRequest(
            ip_blocks=[self._block_msg],
            force=False,
        )
        resp = self._stub.RemoveIPBlock(remove_request1)

        expect = RemoveIPBlockResponse()
        expect.ip_blocks.extend([self._block_msg])
        self.assertEqual(expect, resp)

        # Ensure that block has been removed
        resp = self._stub.ListAddedIPv4Blocks(Void())
        expect = ListAddedIPBlocksResponse()
        self.assertEqual(expect, resp)

    def test_remove_after_releasing_some_addresses(self):
        """ remove after releasing some addresses shouldn't remove block """
        # Assign IP block
        self._stub.AddIPBlock(self._block_msg)

        # Allocate 2 IPs
        alloc_request0 = AllocateIPRequest(
            sid=self._sid0,
            version=AllocateIPRequest.IPV4,
            apn=self._apn0,
        )
        ip_msg0 = self._stub.AllocateIPAddress(alloc_request0)

        alloc_request1 = AllocateIPRequest(
            sid=self._sid1,
            version=AllocateIPRequest.IPV4,
            apn=self._apn0,
        )
        self._stub.AllocateIPAddress(alloc_request1)

        # Test remove without force -- should not remove block
        remove_request0 = RemoveIPBlockRequest(
            ip_blocks=[self._block_msg],
            force=False,
        )
        resp = self._stub.RemoveIPBlock(remove_request0)

        expect = RemoveIPBlockResponse()
        self.assertEqual(expect, resp)

        # Ensure that block has not been removed
        resp = self._stub.ListAddedIPv4Blocks(Void())
        expect = ListAddedIPBlocksResponse()
        expect.ip_block_list.extend([self._block_msg])
        self.assertEqual(expect, resp)

        # Release the allocated IPs
        release_request0 = ReleaseIPRequest(
            sid=self._sid0,
            ip=ip_msg0.ip_list[0],
            apn=self._apn0,
        )
        resp = self._stub.ReleaseIPAddress(release_request0)
        self.assertEqual(resp, Void())

        # Test remove without force -- should not remove block
        remove_request1 = RemoveIPBlockRequest(
            ip_blocks=[self._block_msg],
            force=False,
        )
        resp = self._stub.RemoveIPBlock(remove_request1)

        expect = RemoveIPBlockResponse()
        self.assertEqual(expect, resp)

        # Ensure that block has not been removed
        resp = self._stub.ListAddedIPv4Blocks(Void())
        expect = ListAddedIPBlocksResponse()
        expect.ip_block_list.extend([self._block_msg])
        self.assertEqual(expect, resp)

    def test_ipv6(self):
        """ ipv6 requests should work for allocate / release IP requests """
        # Assign IP block
        self._stub.AddIPBlock(self._ipv6_block_msg)
        # AllocateIPAddress
        request = AllocateIPRequest(
            sid=self._sid1,
            version=AllocateIPRequest.IPV6,
            apn=self._apn0,
        )

        ip_msg = self._stub.AllocateIPAddress(request)
        self.assertTrue(
            ipaddress.ip_address(ip_msg.ip_list[0].address) in
            self._ipv6_block,
        )

        # ReleaseIPAddress
        release_request = ReleaseIPRequest(
            sid=self._sid1,
            ip=ip_msg.ip_list[0],
            apn=self._apn0,
        )
        self._stub.ReleaseIPAddress(release_request)
