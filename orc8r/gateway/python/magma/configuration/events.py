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

import snowflake
from magma.eventd.eventd_client import log_event
from orc8r.protos.eventd_pb2 import Event


def deleted_stored_mconfig():
    log_event(
        Event(
            stream_name="magmad",
            event_type="deleted_stored_mconfig",
            tag=snowflake.snowflake(),
            value="{}",
        ),
    )


def updated_stored_mconfig():
    log_event(
        Event(
            stream_name="magmad",
            event_type="updated_stored_mconfig",
            tag=snowflake.snowflake(),
            value="{}",
        ),
    )
