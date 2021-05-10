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

import grpc
from lte.protos.session_manager_pb2 import (
    UPFNodeState,
    UPFSessionConfigState,
    UPFPagingInfo)
from lte.protos.session_manager_pb2_grpc import SetInterfaceForUserPlaneStub

DEFAULT_GRPC_TIMEOUT = 5

def send_node_state_association_request(node_state_info: UPFNodeState,
                                        setinterface_stub: SetInterfaceForUserPlaneStub):
    """
    Make RPC call to send Node Association Setup/Release request to
    sessionD (SMF)
    """
    try:
        setinterface_stub.SetUPFNodeState(node_state_info, DEFAULT_GRPC_TIMEOUT)
        return True
    except grpc.RpcError as err:
        logging.error(
            "send_node_state_association_request error[%s] %s",
            err.code(),
            err.details())
        return False

def send_periodic_session_update(upf_session_config_state: UPFSessionConfigState,
                                 setinterface_stub: SetInterfaceForUserPlaneStub):
    """
    Make RPC call to send periodic messages to smf about sessions state.
    """
    try:
        setinterface_stub.SetUPFSessionsConfig(upf_session_config_state)
        return True
    except grpc.RpcError as err:
        logging.error(
            "send_periodic_session_update error[%s] %s",
            err.code(),
            err.details())
        return False


def send_paging_intiated_notification(paging_info: UPFPagingInfo,
                                      setinterface_stub: SetInterfaceForUserPlaneStub):
    """
	Make RPC call to send paging initiated notification to sessionD
    """
    try:
        setinterface_stub.SetPagingInitiated(paging_info, DEFAULT_GRPC_TIMEOUT)

    except grpc.RpcError as err:
        logging.error(
            "send_paging_intiated_notification error[%s] %s",
            err.code(),
            err.details())
