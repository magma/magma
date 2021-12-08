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

from magma.common.sentry import sentry_init
from magma.common.service import MagmaService
from magma.common.service_registry import ServiceRegistry
from magma.ctraced.rpc_servicer import CtraceDRpcServicer
from magma.ctraced.trace_manager import TraceManager
from orc8r.protos.ctraced_pb2_grpc import CallTraceControllerStub
from orc8r.protos.mconfig.mconfigs_pb2 import CtraceD


def main():
    """ main() for ctraced """
    service = MagmaService('ctraced', CtraceD())

    # Optionally pipe errors to Sentry
    sentry_init(service_name=service.name, sentry_mconfig=service.shared_mconfig.sentry_config)

    orc8r_chan = ServiceRegistry.get_rpc_channel(
        'ctraced',
        ServiceRegistry.CLOUD,
    )
    ctraced_stub = CallTraceControllerStub(orc8r_chan)

    trace_manager = TraceManager(service.config, ctraced_stub)

    ctraced_servicer = CtraceDRpcServicer(trace_manager)
    ctraced_servicer.add_to_server(service.rpc_server)

    # Run the service loop
    service.run()

    # Cleanup the service
    service.close()


if __name__ == "__main__":
    main()
