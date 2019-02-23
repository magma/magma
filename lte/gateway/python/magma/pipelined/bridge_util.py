"""
Copyright (c) 2018-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import subprocess


class DatapathLookupError(Exception):
    pass


class BridgeTools:
    """
    BridgeTools

    Use ovs-vsctl commands to get bridge info and setup bridges for testing.
    """

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
    def create_bridge(bridge_name, iface_name):
        """
        Creates a simple bridge, sets up an interface.
        Used when running unit tests
        """
        subprocess.Popen(["ovs-vsctl", "add-br", bridge_name]).wait()
        subprocess.Popen(["ovs-vsctl", "set", "bridge", bridge_name,
                          "protocols=OpenFlow10,OpenFlow13,OpenFlow14",
                          "other-config:disable-in-band=true"]).wait()
        subprocess.Popen(["ovs-vsctl", "set-controller", bridge_name,
                          "tcp:127.0.0.1:6633", "tcp:127.0.0.1:6654"]).wait()
        subprocess.Popen(["ifconfig", iface_name, "192.168.1.1/24"]).wait()

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
    def get_flows_for_bridge(bridge_name, table_num=None):
        set_cmd = ["ovs-ofctl", "dump-flows", bridge_name]
        if table_num:
            set_cmd.append("table=%s" % table_num)
        flows = subprocess.check_output(set_cmd).decode('utf-8').split('\n')
        return flows
