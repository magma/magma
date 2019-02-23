"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from typing import Any

import grpc
from magma.enodebd.enodeb_status import get_enodeb_status
from lte.protos.enodebd_pb2 import GetParameterResponse
from lte.protos.enodebd_pb2_grpc import EnodebdServicer, add_EnodebdServicer_to_server
from magma.enodebd.state_machines.enb_acs_pointer import StateMachinePointer
from orc8r.protos.service303_pb2 import ServiceStatus
from magma.common.rpc_utils import return_void


class EnodebdRpcServicer(EnodebdServicer):
    """
    gRPC based server for enodebd
    """
    def __init__(self, state_machine_pointer: StateMachinePointer) -> None:
        self.state_machine_pointer = state_machine_pointer

    def add_to_server(self, server):
        """
        Add the servicer to a gRPC server
        """
        add_EnodebdServicer_to_server(self, server)

    def GetParameter(self, request: Any, context: Any) -> Any:
        """
        Sends a GetParameterValues message. Used for testing only.
        """
        get_parameter_values_response = GetParameterResponse()
        state_machine = self.state_machine_pointer.state_machine
        # Different data models will have different names for the same
        # parameter. Whatever name that the data model uses, we call the
        # 'parameter path', eg. "Device.DeviceInfo.X_BAICELLS_COM_GPS_Status"
        # We denote 'ParameterName' to be a standard string name for
        # equivalent parameters between different data models
        parameter_path = request.parameter_name
        data_model = self.state_machine_pointer.state_machine.data_model
        param_name = data_model.get_parameter_name_from_path(parameter_path)
        param_value = state_machine.get_parameter(param_name)
        get_parameter_values_response.parameters.add(
            name=parameter_path, value=param_value)
        return get_parameter_values_response

    @return_void
    def SetParameter(self, request: Any, context: Any):
        """
        Sends a SetParameterValues message. Used for testing only.
        """
        if request.HasField('value_int'):
            value = (request.value_int, 'int')
        elif request.HasField('value_bool'):
            value = (request.value_bool, 'boolean')
        elif request.HasField('value_string'):
            value = (request.value_string, 'string')
        else:
            context.set_details('SetParameter: Unsupported type %d',
                                request.type)
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            return

        state_machine = self.state_machine_pointer.state_machine
        state_machine.set_parameter_asap(request.parameter_name, value)
        return

    @return_void
    def Reboot(self, _=None, context=None):
        """
        Reboot eNodeB
        """
        state_machine = self.state_machine_pointer.state_machine
        state_machine.reboot_asap()

    def GetStatus(self, _=None, context=None):
        """
        Get eNodeB status
        Note: input variable defaults used so this can be either called locally
        or as an RPC.
        """
        state_machine = self.state_machine_pointer.state_machine
        status = get_enodeb_status(state_machine)
        status_message = ServiceStatus()
        status_message.meta.update(status)
        return status_message
