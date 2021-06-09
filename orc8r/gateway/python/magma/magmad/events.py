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
import json

import snowflake
from google.protobuf.json_format import MessageToDict
from magma.eventd.eventd_client import log_event
from orc8r.protos.eventd_pb2 import Event
from orc8r.swagger.models.restarted_services import RestartedServices


def processed_updates(configs_by_service):
    # Convert to dicts for JSON serializability
    configs = {}
    for srv, config in configs_by_service.items():
        configs[srv] = MessageToDict(config)
    log_event(
        Event(
            stream_name="magmad",
            event_type="processed_updates",
            tag=snowflake.snowflake(),
            value=json.dumps(configs),
        ),
    )


def restarted_services(services):
    # Convert to a list for JSON serializability
    services = list(services)
    log_event(
        Event(
            stream_name="magmad",
            event_type="restarted_services",
            tag=snowflake.snowflake(),
            value=json.dumps(RestartedServices(services=services).to_dict()),
        ),
    )


def established_sync_rpc_stream():
    log_event(
        Event(
            stream_name="magmad",
            event_type="established_sync_rpc_stream",
            tag=snowflake.snowflake(),
            value="{}",
        ),
    )


def disconnected_sync_rpc_stream():
    log_event(
        Event(
            stream_name="magmad",
            event_type="disconnected_sync_rpc_stream",
            tag=snowflake.snowflake(),
            value="{}",
        ),
    )
