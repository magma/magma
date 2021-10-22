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
import os
from typing import Any, Dict, Optional

import sentry_sdk
import snowflake
from magma.configuration.service_configs import get_service_config_value

CONTROL_PROXY = 'control_proxy'
SENTRY_ENABLED = 'sentry_enabled'
SENTRY_URL = 'sentry_url_python'
SENTRY_SAMPLE_RATE = 'sentry_sample_rate'
COMMIT_HASH = 'COMMIT_HASH'
HWID = 'hwid'
SERVICE_NAME = 'service_name'
LOGGING_EXTRA = 'extra'
SEND_TO_SENTRY_KEY = "send_to_sentry"
SEND_TO_SENTRY = {SEND_TO_SENTRY_KEY: True}


def _ignore_if_not_marked(event: Dict[str, Any], hint: Dict[str, Any]) -> Optional[Dict[str, Any]]:
    if event.get(LOGGING_EXTRA) and event.get(LOGGING_EXTRA).get(SEND_TO_SENTRY_KEY):
        return event
    return None


def sentry_init(service_name: str, requires_opt_in=False):
    """ Initialize connection and start piping errors to sentry.io.

    Args:
        service_name: Name of the service that uses Sentry
        requires_opt_in: If set to True, only errors that are explicitly
        marked will be sent to Sentry. Otherwise all log errors will
        be sent. To mark an error, set `extra=SEND_TO_SENTRY` in an error log."""
    sentry_enabled = get_service_config_value(
        service_name,
        SENTRY_ENABLED,
        default=False,
    )
    if not sentry_enabled:
        return

    sentry_url = get_service_config_value(
        CONTROL_PROXY,
        SENTRY_URL,
        default='',
    )
    if not sentry_url:
        return

    sentry_sample_rate = get_service_config_value(
        CONTROL_PROXY,
        SENTRY_SAMPLE_RATE,
        default=1.0,
    )
    sentry_sdk.init(
        dsn=sentry_url,
        release=os.getenv(COMMIT_HASH),
        traces_sample_rate=sentry_sample_rate,
        before_send=_ignore_if_not_marked if requires_opt_in else None,
    )
    sentry_sdk.set_tag(HWID, snowflake.snowflake())
    sentry_sdk.set_tag(SERVICE_NAME, service_name)
