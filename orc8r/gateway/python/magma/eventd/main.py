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
from magma.eventd.event_validator import EventValidator
from magma.eventd.rpc_servicer import EventDRpcServicer
from orc8r.protos.mconfig.mconfigs_pb2 import EventD


def main():
    """ main() for eventd """
    service = MagmaService('eventd', EventD())

    # Optionally pipe errors to Sentry
    sentry_init(service_name=service.name, sentry_mconfig=service.shared_mconfig.sentry_config)

    event_validator = EventValidator(service.config)
    eventd_servicer = EventDRpcServicer(service.config, event_validator)
    eventd_servicer.add_to_server(service.rpc_server)

    # Run the service loop
    service.run()

    # Cleanup the service
    service.close()


if __name__ == "__main__":
    main()
