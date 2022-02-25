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

import copy
import logging
import re

from ovstest import util  # pylint: disable=import-error
from ovstest import vswitch  # pylint: disable=import-error

logger = logging.getLogger(__name__)  # pylint: disable=invalid-name


class UseAfterFreeException(Exception):
    """
    Exception thrown when the underlying system resource is being accessed
    after it has been freed.
    e.g. trying to up an interface that has been destroyed.
    """
    pass


def is_valid(func):
    """
    Decorator to check if the object is valid for bound functions.
    """

    def nested(self, *args, **kwargs):
        """
        Wrapper decorator calls wrapped func with passed in args and kwargs
        """
        if not self.valid:
            raise UseAfterFreeException(
                "Underlying system resource has "
                "been freed invalid access to "
                "function %s on class %s" %
                (
                    func.__name__,
                    self.__class__.__name__,
                ),
            )

        return func(self, *args, **kwargs)

    return nested


class OvsException(Exception):
    """
    Class encapsulating all exceptions while interacting with OVS.
    Refer to ovs-vswitchd or ovsdb logs to root cause issue
    """
    pass


class Port(object):
    """
    Test representation of a virtual bridge port.
    """
    UNINITIALIZED_PORT_NO = -1

    # ovs interface table schema constants
    INT_TABLE = "Interface"  # ovsdb interface table
    INT_TYPE = "type"  # type of interface
    INT_OFPORT = "ofport"  # bridge port num associated with the interface
    INT_LINK_STATE = "link_state"  # link state attribute of interface

    def __init__(self, iface_name, bridge_name, port_type):
        """
        Create a bridge port of the same name as the iface and sets the port
        type in the interface table.
        Initialize the list of port attributes.
        Internal ports are special in the sense that they are opened as tap
        devices by the bridge implementation.
        Bridge needs to already exist for the creation to succeed.
        Interface needs to exist for the port status to be up.
        Args:
            iface_name: Interface to connect to the bridge
            bridge_name: Name of the bridge.
            port_type: Type of port defaults to internal

        Raises:
            OvsException in case of error
        """
        self._iface_name = iface_name
        self._bridge_name = bridge_name
        self._port_type = port_type
        self._of_port_no = self.UNINITIALIZED_PORT_NO  # Lazy initialization.

        # Create the interface
        ret_val, out, err = util.start_process([
            "ovs-vsctl", "add-port",
            bridge_name, iface_name,
        ])

        if ret_val:
            raise OvsException(
                "Error creating port on bridge %s output %s, "
                "error %s" % (bridge_name, out, err),
            )

        ret_val = vswitch.ovs_vsctl_set(
            self.INT_TABLE, iface_name,
            self.INT_TYPE, None, port_type,
        )
        if ret_val:
            raise OvsException(
                "Error setting interface type for interface "
                "%s" % iface_name,
            )
        self._valid = True

    @staticmethod
    def _extract_of_port_no(iface_name):
        """
        Extract the ofport number for the interface.

        Returns:
            return of ofport number of the specified interface.
        """

        # Query the db for the ofport number.
        ret, out, err = util.start_process([
            "ovs-vsctl", "get",
            Port.INT_TABLE, iface_name,
            Port.INT_OFPORT,
        ])
        if ret:
            raise OvsException(
                "Failed to read port number for interface %s, "
                "message %s, error no %s" % (
                    iface_name,
                    out, err,
                ),
            )

        return int(out)

    @property
    def port_no(self):
        """
        Get the ofport number associated with the vswitch port. Explicitly
        calls out to ovs to obtain port number if not already fetched.

        Returns:
            of port number of the interface.
        """
        if self._of_port_no == self.UNINITIALIZED_PORT_NO:
            self._of_port_no = self._extract_of_port_no(self._iface_name)

        return self._of_port_no

    @property
    def valid(self):
        """
        Return if the object is valid
        Returns:
            True if the object is still valid
        """
        return self._valid

    @property
    def iface_name(self):
        """ Accessor for the bound iface name """
        return self._iface_name

    @property
    def bridge_name(self):
        """ Accessor for the bound bridge name """
        return self._bridge_name

    @is_valid
    def destroy(self, free_resource):
        """
        Marks the port as invalid, the underlying resource is freed by the
        bridge when it is deleted.
        Args:
            free_resource: Free the underlying system resource along with
            marking the device as invalid
        """
        if free_resource:
            ret_val = vswitch.ovs_vsctl_del_port_from_bridge(self._iface_name)
            if ret_val:
                raise OvsException(
                    "Error deleting port %s on bridge %s" % (
                        self._iface_name, self._bridge_name,
                    ),
                )

        self._valid = False

    @is_valid
    def sanity_check(self):
        """
        Check that a port is in linked up state.
        Returns:
            True if port is in linked up state False otherwise
        Raises:
            OvsException if port cannot be accessed.
        """
        up_str = b'up\n'
        ret, out, err = util.start_process([
            "ovs-vsctl", "get", Port.INT_TABLE,
            self._iface_name,
            Port.INT_LINK_STATE,
        ])
        if ret:
            raise OvsException(
                "Failed to determine interface link state %s, "
                "output %s, erro %s" % (
                    self._iface_name, out,
                    err,
                ),
            )

        return out == up_str


class Bridge(object):
    """
    Test representation of a linux bridge
    """

    def __init__(self, br_name):
        """
        Create a ovs bridge with the given name
        Initialize the list of ports. Virtual and physical ports are maintained
        independently for convenience.
        Args:
            br_name: Bridge name, string
        Raises:
            OvsException on port creation failure.
        """
        self._br_name = br_name
        ret_val = vswitch.ovs_vsctl_add_bridge(br_name)
        if ret_val:
            raise OvsException("Error creating ovs bridge %s" % self._br_name)

        # List of physical adapters connected to the switch e.g. eth0
        self._phy_port = set()
        # Dictionary of virtual port name to virtual port objects connected to
        # the switch.
        self._virt_port = {}
        self._valid = True

    @property
    def name(self):
        """ Accessor for bridge name """
        return self._br_name

    @property
    def valid(self):
        """ Accessor for object validity """
        return self._valid

    @is_valid
    def add_virtual_port(self, iface_name, port_type):
        """
        Add a virtual port to the ovs bridge and "bind" it to the iface.
        Args:
            iface_name: name of the port to add
            port_type: Type of port, refer to ovs documentation for port types.
        Returns:
            Port object corresponding to the created port.
        Raises:
            OvsException on port creation failure.
        """
        port = Port(iface_name, self._br_name, port_type)
        self._virt_port[iface_name] = port
        return port

    @is_valid
    def add_physical_port(self, p_nic):
        """
        Add a physical interface to the bridge
        TODO: Support nic bonding in case we need multiple pNics to be linked
        to the same bridge.
        TODO: Add support for migrating ip config associated with a pnic to
        a virtual interface.
        Args:
            p_nic: Physical network card name.
        Raises:
            ValueError: If the bridge already has a pnic attached or the pnic
            has a configured IP address.
            OvsException on device addition failure.
        """
        # No support for nic bonding yet, so bail early
        if self._phy_port:
            assert len(self._phy_port) == 1
            raise ValueError(
                "Bridge already has a pnic %s attached to it" %
                next(iter(self._phy_port)),
            )

        # No ip migration support yet, so bail instead of blackholing box.
        (ip_addr, _) = Interface.interface_get_ip(p_nic)
        if ip_addr != "0.0.0.0":
            raise ValueError(
                "pnic %s has an ip address assigned to it "
                "use a migrate ip workflow to move ip address "
                "to a virtual interface" % p_nic,
            )

        vswitch.ovs_vsctl_add_port_to_bridge(self._br_name, p_nic)
        self._phy_port.add(p_nic)

    @is_valid
    def destroy(self):
        """
        Iterate through port objects on the bridge and delete them
        Delete the bridge.
        For performance optimization the delete of the port doesn't actually
        free up the underlying system bridge port resource as the subsequent
        delete will clear it out anyway.
        """
        ret_val = vswitch.ovs_vsctl_del_bridge(self._br_name)
        free_resource = False
        for port in self._virt_port.values():
            if port.valid:
                port.destroy(free_resource)
        if ret_val:
            raise OvsException("Failed to delete bridge %s" % self._br_name)
        self._valid = False

    @is_valid
    def sanity_check(self):
        """
        Iterate through the list of ports and determine their link state. Iff
        all the ports are linked up then the port is considered up
        Returns:
            return True if all the ports are sanity checked.
            return False otherwise.
        """
        for port in self._virt_port.values():
            if not port.sanity_check():
                logging.warning(
                    "Sanity check for port %s failed", port.iface_name,
                )
                return False
        return True


class Interface(object):
    """
    Test class wrapping interface configuration.
    """

    def __init__(self, iface, ip_address, netmask):
        """
        Create a network interface with the provided configuration
        Args:
            iface: iface name
            ip_address: string address
            netmask: Netmask
        TODO: Support v6.
        """
        self._iface = iface
        self._ip_address = ip_address
        self._netmask = netmask

        ip_cidr = "%s/%s" % (
            self._ip_address,
            Interface.dotdec_to_cidr(self._netmask),
        )
        args = ["ip", "addr", "add", ip_cidr, "dev", self._iface]
        ret, _, _ = util.start_process(args)
        if ret:
            raise OvsException("Failed to create interface %s" % self._iface)

        self._valid = True

    @is_valid
    def up(self):  # pylint: disable=invalid-name
        """
        Bring up the interface
        Raises:
            OvsException if the interface bring up failed.
        """
        ret, _, _ = util.start_process([
            "ip", "link", "set",
            self._iface, "up",
        ])

        if ret:
            raise OvsException("Failed to bring up interface %s" % self._iface)

    @is_valid
    def destroy(self):
        """
        Bring the interface down.
        Raises
            OvsException if interface destroy fails.
        """
        ret, _, _ = util.start_process(["ip", "link", "show", self._iface])
        if ret == 1:  # iface doesn't exist
            self._valid = False
            return

        ret, _out, _err = util.start_process([
            "ip", "link", "set",
            self._iface, "down",
        ])
        if ret:
            raise OvsException(
                "Failed to bring down interface %s"
                % self._iface,
            )
        self._valid = False

    @property
    def name(self):
        """ Accessor for the iface name """
        return self._iface

    @property
    def ip_address(self):
        """ Accessor for the ip_address """
        return self._ip_address

    @property
    def valid(self):
        """ Accessor for the validity of the object """
        return self._valid

    @property
    def netmask(self):
        """ Configured netmask """
        return self._netmask

    @is_valid
    def sanity_check(self):
        """
        Check that the ip address and netmask configured matches the
        Return:
            True if configured matches assigned.
        Raises:
            OvsException if the interface is not accessible
        """

        # Poor interface, returns error code on failure but tuple on success.
        ip_addr, netmask = Interface.interface_get_ip(self._iface)
        if ip_addr == self._ip_address and netmask == self._netmask:
            return True

        logging.warning(
            "Configured ipaddress/netmask %s/%s differs from "
            "system ipaddress/netmask %s/%s for interface %s",
            self._ip_address, self._netmask, ip_addr, netmask,
            self._iface,
        )
        return False

    @staticmethod
    def dotdec_to_cidr(dot_dec_netmask):
        """
        Convert a dot-decimal netmask into a CIDR-style one.

        e.g.: 255.255.255.0 -> 24

        Args:
            dot_dec_netmask: Dotted decimal netmask (e.g., 255.255.255.0)
        Returns:
            Integer representing number of bits netmask represents.
        Raises:
            ValueError if input is invalid.
        """
        octets = dot_dec_netmask.split('.')
        if len(octets) != 4:
            raise ValueError("Invalid netmask: %s" % dot_dec_netmask)
        return sum([bin(int(x)).count('1') for x in octets])

    @staticmethod
    def cidr_to_dotdec(cidr_netmask):
        """
        Convert the netmask portion of a CIDR-style ip/netmask to a dot-decimal
        format.

        e.g., 24 -> 255.255.255.0

        Args:
            cidr_netmask: An integer representing the netmask portion of a
                          CIDR-style IP address/netmask combination.
        Returns:
            Dotted-decimal format netmask (e.g., 255.255.255.0)
        Raise:
            ValueError if cidr_netmask is invalid
        """
        MAXLEN = 32
        if cidr_netmask > MAXLEN:
            raise ValueError("Invalid netmask: %s" % cidr_netmask)

        mask_bits = 0
        for i in range(0, cidr_netmask):
            mask_bits |= (1 << i)
        mask_bits = mask_bits << (MAXLEN - cidr_netmask)

        # build return
        octets = [((mask_bits & (0xff << (MAXLEN - 8 - i))) >> (MAXLEN - 8 - i))
                  for i in range(0, MAXLEN, 8)]
        return ".".join(["%d" % b for b in octets])

    @staticmethod
    def interface_get_ip(iface):
        """
        Returns ip address and netmask (in dot decimal format) for the given
        iface.

        Args:
            iface: Interface name
        Return:
            (ip_address, netmask) tuple
        Raises:
            OvsException if iface doesn't exist
        """
        args = ["ip", "-f", "inet", "-o", "addr", "show", iface]
        ret, out, _err = util.start_process(args)
        if ret:
            raise OvsException("Can't get ip address for %s" % iface)

        try:
            res = re.search(br'inet (\S+)/(\S+)', out)
            if res:
                ip = res.group(1).decode('utf-8')
                netmask = Interface.cidr_to_dotdec(int(res.group(2)))
                return (ip, netmask)
        except ValueError as e:
            raise OvsException("Can't get ip address for %s" % iface) from e


class TopologyBuilder(object):
    """
    Class to help build a topology with a test bridge.
    """

    def __init__(self):
        """
        Initialize topology state.
        """
        self._bridges = {}
        self._interfaces = {}

    def create_bridge(self, br_name):
        """
        Creates an OVS bridge.
        Args:
            br_name: Bridge name (str)
        Returns:
            created Bridge object
        """
        bridge = Bridge(br_name)
        self._bridges[br_name] = bridge
        return bridge

    def create_interface(self, iface_name, ip_address, netmask):
        """
        Create an iface given the iface_name, ip_address and netmask
        Args:
            iface_name: name of the iface to be created.
            ip_address: iface ip_address
            netmask: netmask
        Returns:
            created iface object
        """
        iface = Interface(iface_name, ip_address, netmask)
        self._interfaces[iface_name] = iface
        iface.up()
        return iface

    def bind(self, iface_name, bridge, port_type="internal"):
        """
        Bind an iface to a bridge
        Args:
            iface_name: iface to bind to the bridge
            bridge: bridge to bind iface to.
            port_type: internal
        Returns:
            The created port that the interface was bound to.
        """
        assert bridge.name in self._bridges
        # Create a virtual port and associate it with the iface
        port = bridge.add_virtual_port(iface_name, port_type)
        return port

    def destroy(self):
        """
        Destroy bridges and interfaces created as part of this topology.
        Aborts destroy on first exception.
        """
        bridges = copy.deepcopy(self._bridges)
        for bridge in bridges.values():
            if bridge.valid:
                bridge.destroy()
            del self._bridges[bridge.name]

        interfaces = copy.deepcopy(self._interfaces)
        for iface in interfaces.values():
            if iface.valid:
                iface.destroy()
            del self._interfaces[iface.name]

    def invalid_devices(self):
        """
        Sanity checks that the various network objects are in the appropriate
        state
        Returns:
            The list of devices that are not in the expected state.
        """
        device_list = []
        for bridge in self._bridges.values():
            if not bridge.sanity_check():
                device_list.append("Bridge %s" % bridge.name)
        for iface in self._interfaces.values():
            if not iface.sanity_check():
                device_list.append("iface %s" % iface.name)

        return device_list
