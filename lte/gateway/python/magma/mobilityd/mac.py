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


class MacAddress:
    """
    Manage Mac address conversion from various formats.
    """

    def __init__(self, mac: str):
        self.mac_address = mac

    def as_hex(self) -> str:
        """
        Covert Mac address string to binary number format.
        Returns: packed binary number.
        """

        return bytes.fromhex(self.mac_address.replace(':', ''))

    def as_redis_key(self, vlan: str) -> str:
        """
        Convert MAC address string to redis key. Redis does not
        allow ':' in the key, so use '_' instead
        Returns: str to make redis happy
        """
        key = str(self.mac_address).replace(':', '_').lower()
        if vlan:
            return "v{}.{}".format(vlan, key)
        else:
            return key

    def __str__(self):
        return self.mac_address


def create_mac_from_sid(imsi: str) -> MacAddress:
    if imsi.upper().startswith('IMSI'):
        return MacAddress(sid_to_mac(imsi))
    elif imsi.find(':') != -1:
        return MacAddress(imsi)
    elif len(imsi) == 12:
        # Convert a hex number to Mac address string format.
        return MacAddress(hex_to_mac(imsi))
    else:
        raise InvalidMacAddressFormat(imsi)


def sid_to_mac(com_sid: str) -> str:
    """
    Generate MAC address from SID

    https://na.baicells.com/wp-content/uploads/2020/01/
    baicells_configuration__network_admin_guide_srv1.24_2-jan-2020-nxp.pdf
    Bridge - Layer 2 will create a virtual interface for each CPE that
    attaches using a DHCP request to create a 1:1 mapping between the CPE
    IP address (from the EPC) and the LGW IP address. A CPE's MAC address
    is generated from its IMSI: Convert the last 12 digits to hex, and
    then prefix it with "8A". For example, if the IMSI = 311980000002918,
    the MAC would be 8A:E4:2C:8D:53:66.
    """
    sid = com_sid.split('.')[0]
    if not sid.startswith('IMSI'):
        raise InvalidIMSIError(sid)

    sid = sid[4:]  # strip IMSI off of string
    sid = sid[-12:]  # use last 12 digits
    mac_prefix = "8A"
    hex_num = hex(int(sid))[2:].zfill(10)
    mac = "{}:{}{}:{}{}:{}{}:{}{}:{}{}".format(mac_prefix, *hex_num)

    return mac.upper()


def hex_to_mac(sid) -> str:
    """
    Convert a hex number to Mac address string format.
    Args:
        sid: integer
    Returns: Hex Mac address as a string.
    """
    return ':'.join(''.join(x) for x in zip(*[iter(sid)] * 2))


class InvalidIMSIError(Exception):
    """ Exception thrown when a given IP block overlaps with existing ones
    """


class InvalidMacAddressFormat(Exception):
    """ Exception thrown when a given IP block overlaps with existing ones
    """
