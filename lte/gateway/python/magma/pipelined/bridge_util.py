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
import binascii
from collections import defaultdict
import re
import logging
import subprocess
from typing import Optional, Dict, List, TYPE_CHECKING

# Prevent circular import
if TYPE_CHECKING:
    from magma.pipelined.service_manager import Tables


class DatapathLookupError(Exception):
    pass


class BridgeTools:
    """
    BridgeTools

    Use ovs-vsctl commands to get bridge info and setup bridges for testing.
    """
    TABLE_NUM_REGEX = r'table=(\d+)'

    @staticmethod
    def get_datapath_id(bridge_name):
        """
        Gets the datapath_id by bridge_name

        Hacky, call vsctl, decode output to str, strip '\n', remove '' around
        the output, convert to int.

        This gives the integer datapath_id that we want to run apps on, this is
        needed when 2 bridges are setup, gtp_br0(main bridge) and testing_br)
        """
        try:
            output = subprocess.check_output(["ovs-vsctl", "get", "bridge",
                                              bridge_name, "datapath_id"])
            output_str = str(output, 'utf-8').strip()[1:-1]
            output_hex = int(output_str, 16)
        except subprocess.CalledProcessError as e:
            raise DatapathLookupError(
                'Error: ovs-vsctl bridge({}) datapath id lookup: {}'.format(
                    bridge_name, e
                )
            )
        return output_hex

    @staticmethod
    def get_ofport(interface_name):
        """
        Gets the ofport name ofport number of a interface
        """
        try:
            port_num = subprocess.check_output(["ovs-vsctl", "get", "interface",
                                                interface_name, "ofport"])
        except subprocess.CalledProcessError as e:
            raise DatapathLookupError(
                'Error: ovs-vsctl interface({}) of port lookup: {}'.format(
                    interface_name, e
                )
            )
        return int(port_num)

    @staticmethod
    def create_internal_iface(bridge_name, iface_name, ip):
        """
        Creates a simple bridge, sets up an interface.
        Used when running unit tests
        """
        subprocess.Popen(["ovs-vsctl", "add-port", bridge_name, iface_name,
                          "--", "set", "Interface", iface_name,
                          "type=internal"]).wait()
        if ip is not None:
            subprocess.Popen(["ifconfig", iface_name, ip]).wait()

        subprocess.Popen(["ifconfig", iface_name, "up"]).wait()

    @staticmethod
    def add_ovs_port(bridge_name: str, iface_name: str, ofp_port: str):
        """
            Add interface to ovs bridge
        """
        try:
            add_port_cmd = ["ovs-vsctl", "--may-exist",
                            "add-port",
                            bridge_name,
                            iface_name,
                            "--", "set", "interface",
                            iface_name,
                            "ofport_request=" + ofp_port]
            subprocess.check_call(add_port_cmd)
            logging.debug("add_port_cmd %s", add_port_cmd)
        except subprocess.CalledProcessError as e:
            logging.warning("Error while adding ports: %s", e)

        try:
            if_up_cmd = ["ip", "link", "set", "dev",
                         iface_name, "up"]
            subprocess.check_call(if_up_cmd)
            logging.debug("if_up_cmd %s", if_up_cmd)
        except subprocess.CalledProcessError as e:
            logging.warning("Error while if up interface: %s", e)

    @staticmethod
    def create_veth_pair(port1: str, port2: str):
        try:
            create_veth = ["ip", "link", "add",
                           port1,
                           "type", "veth",
                           "peer", "name", port2]
            subprocess.check_call(create_veth)
            logging.debug("if_up_cmd %s", create_veth)
        except subprocess.CalledProcessError as e:
            logging.debug("Error while creating veth pair: %s", e)

    @staticmethod
    def create_bridge(bridge_name, iface_name):
        """
        Creates a simple bridge, sets up an interface.
        Used when running unit tests
        """
        subprocess.Popen(["ovs-vsctl", "--if-exists", "del-br",
                          bridge_name]).wait()
        subprocess.Popen(["ovs-vsctl", "add-br", bridge_name]).wait()
        subprocess.Popen(["ovs-vsctl", "set", "bridge", bridge_name,
                          "protocols=OpenFlow10,OpenFlow13,OpenFlow14",
                          "other-config:disable-in-band=true"]).wait()
        subprocess.Popen(["ovs-vsctl", "set-controller", bridge_name,
                          "tcp:127.0.0.1:6633", "tcp:127.0.0.1:6654"]).wait()
        subprocess.Popen(["ovs-vsctl", "set-manager", "ptcp:6640"]).wait()
        subprocess.Popen(["ifconfig", iface_name, "192.168.1.1/24"]).wait()

    @staticmethod
    def flush_conntrack():
        """
        Cleanup the conntrack state
        """
        subprocess.Popen(["ovs-dpctl", "flush-conntrack"]).wait()

    @staticmethod
    def destroy_bridge(bridge_name):
        """
        Removes the bridge.
        Used when unit test finishes
        """
        subprocess.Popen(["ovs-vsctl", "del-br", bridge_name]).wait()

    @staticmethod
    def get_controllers_for_bridge(bridge_name):
        curr_controllers = subprocess.check_output(
            ["ovs-vsctl", "get-controller", bridge_name],
        ).decode("utf-8").replace(' ', '').split('\n')
        return list(filter(None, curr_controllers))

    @staticmethod
    def add_controller_to_bridge(bridge_name, port_num):
        curr_controllers = BridgeTools.get_controllers_for_bridge(bridge_name)
        ctlr_ip = "tcp:127.0.0.1:{}".format(port_num)
        if ctlr_ip in curr_controllers:
            return
        curr_controllers.append(ctlr_ip)
        BridgeTools.set_controllers_for_bridge(bridge_name, curr_controllers)

    @staticmethod
    def remove_controller_from_bridge(bridge_name, port_num):
        curr_controllers = BridgeTools.get_controllers_for_bridge(bridge_name)
        ctlr_ip = 'tcp:127.0.0.1:{}'.format(port_num)
        curr_controllers.remove(ctlr_ip)
        BridgeTools.set_controllers_for_bridge(bridge_name, curr_controllers)

    @staticmethod
    def set_controllers_for_bridge(bridge_name, ctlr_list):
        set_cmd = ["ovs-vsctl", "set-controller", bridge_name]
        set_cmd.extend(ctlr_list)
        subprocess.Popen(set_cmd).wait()

    @staticmethod
    def get_flows_for_bridge(bridge_name, table_num=None, include_stats=True):
        """
        Returns a flow dump of the given bridge from ovs-ofctl. If table_num is
        specified, then only the flows for the table will be returned.
        """
        if include_stats:
            set_cmd = ["ovs-ofctl", "dump-flows", bridge_name]
        else:
            set_cmd = ["ovs-ofctl", "dump-flows", bridge_name, "--no-stats"]
        if table_num:
            set_cmd.append("table=%s" % table_num)

        flows = \
            subprocess.check_output(set_cmd).decode('utf-8').split('\n')
        flows = list(filter(lambda x: (x is not None and
                                       x != '' and
                                       x.find("NXST_FLOW") == -1),
                            flows))
        return flows

    @staticmethod
    def _get_annotated_name_by_table_num(
            table_assignments: 'Dict[str, Tables]') -> Dict[int, str]:
        annotated_tables = {}
        # A main table may be used by multiple apps
        apps_by_main_table_num = defaultdict(list)
        for name in table_assignments:
            apps_by_main_table_num[table_assignments[name].main_table].append(
                name)
            # Scratch tables are used for only one app
            for ind, scratch_num in enumerate(
                    table_assignments[name].scratch_tables):
                annotated_tables[scratch_num] = '{}(scratch_table_{})'.format(
                    name,
                    ind)
        for table, apps in apps_by_main_table_num.items():
            annotated_tables[table] = '{}(main_table)'.format(
                '/'.join(sorted(apps)))
        return annotated_tables

    @classmethod
    def get_annotated_flows_for_bridge(cls, bridge_name: str,
                                       table_assignments: 'Dict[str, Tables]',
                                       apps: Optional[List[str]] = None,
                                       include_stats: bool = True
                                       ) -> List[str]:
        """
        Returns an annotated flow dump of the given bridge from ovs-ofctl.
        table_assignments is used to annotate table number with its
        corresponding app. If a note exists, the note will be decoded.
        If apps is not None, then only the flows for the given apps will be
        returned.
        """
        annotated_tables = cls._get_annotated_name_by_table_num(
            table_assignments)

        def annotated_table_num(num):
            if int(num) in annotated_tables:
                return annotated_tables[int(num)]
            return num

        def parse_resubmit_action(match):
            """
            resubmit(port,1) => resubmit(port,app_name(main_table))
            """
            ret = ''
            # We can have more than one resubmit per flow
            actions = [a for a in match.group().split('resubmit') if a]
            for action in actions:
                resubmit_tokens = re.search(r'\((.*?)\)', action)\
                                    .group(1).split(',')
                in_port, table = resubmit_tokens[0], resubmit_tokens[1]
                if ret:
                    ret += ','
                ret += 'resubmit({},{})'.format(in_port,
                                                annotated_table_num(table))
            return ret

        def parse_flow(flow):
            sub_rules = [
                # Annotate table number with app name
                (cls.TABLE_NUM_REGEX,
                 lambda match: 'table={}'.format(annotated_table_num(
                     match.group(1)))),
                (r'resubmit\((.*)\)', parse_resubmit_action),
                # Decode the note
                (r'note:([\d\.a-fA-F]*)',
                 lambda match: 'note:{}'.format(
                               str(binascii.unhexlify(match.group(1)
                                                      .replace('00', '')
                                                      .replace('.', ''))))),
            ]
            for rule in sub_rules:
                flow = re.sub(rule[0], rule[1], flow)
            return flow

        def filter_apps(flows):
            if apps is None:
                yield from flows
                return

            selected_tables = []
            for app in apps:
                selected_tables.append(table_assignments[app].main_table)
                selected_tables.extend(table_assignments[app].scratch_tables)

            for flow in flows:
                table_num = int(re.search(cls.TABLE_NUM_REGEX, flow).group(1))
                if table_num in selected_tables or not selected_tables:
                    yield flow

        return [parse_flow(flow) for flow in
                filter_apps(cls.get_flows_for_bridge(bridge_name,
                    include_stats=include_stats))]
