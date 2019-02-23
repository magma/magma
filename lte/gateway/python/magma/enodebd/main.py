"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from threading import Thread
from magma.enodebd.devices.baicells import BaicellsHandler
from magma.enodebd.enodeb_status import get_enodeb_status
from magma.enodebd.state_machines.enb_acs_pointer import StateMachinePointer
from .rpc_servicer import EnodebdRpcServicer
from .stats_manager import StatsManager
from .tr069.server import tr069_server
from .enodebd_iptables_rules import set_enodebd_iptables_rule
from magma.common.service import MagmaService


def main():
    """
    Top-level function for enodebd
    """
    service = MagmaService('enodebd')

    # Statistics manager
    stats_mgr = StatsManager()
    stats_mgr.run()

    # We incorrectly assume that we are dealing with a Baicells device here.
    # When this assumption is invalidated (after the device reports its
    # make and model, then we recreate a new handler to use
    acs_state_machine = BaicellsHandler(service, stats_mgr)
    state_machine_pointer = StateMachinePointer(acs_state_machine)

    # Start TR-069 thread
    server_thread = Thread(target=tr069_server,
                           args=(state_machine_pointer, ),
                           daemon=True)
    server_thread.start()

    # Add all servicers to the server
    enodebd_servicer = EnodebdRpcServicer(state_machine_pointer)
    enodebd_servicer.add_to_server(service.rpc_server)

    # Register function to get service status
    def get_status():
        return get_enodeb_status(state_machine_pointer.state_machine)
    service.register_get_status_callback(get_status)

    # Set eNodeBD iptables rules due to exposing public IP to eNodeB
    service.loop.create_task(set_enodebd_iptables_rule())

    # Run the service loop
    service.run()

    # Cleanup the service
    service.close()


def call_repeatedly(loop, interval, function, *args, **kwargs):
    """
    Wrapper function to schedule function periodically
    """
    # Schedule next call
    loop.call_later(interval, call_repeatedly, loop, interval, function,
                    *args, **kwargs)
    # Call function
    function(*args, **kwargs)

if __name__ == "__main__":
    main()
