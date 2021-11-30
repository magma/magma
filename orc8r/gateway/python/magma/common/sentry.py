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
import os
from enum import Enum
from functools import wraps
from typing import Any, Dict, Optional

import sentry_sdk
import snowflake
from magma.configuration.service_configs import get_service_config_value

CONTROL_PROXY = 'control_proxy'
SENTRY_CONFIG = 'sentry'
SENTRY_URL = 'sentry_url_python'
SENTRY_SAMPLE_RATE = 'sentry_sample_rate'
CLOUD_ADDRESS = 'cloud_address'
ORC8R_CLOUD_ADDRESS = 'orc8r_cloud_address'
COMMIT_HASH = 'COMMIT_HASH'
HWID = 'hwid'
SERVICE_NAME = 'service_name'
LOGGING_EXTRA = 'extra'
SEND_TO_ERROR_MONITORING_KEY = "send_to_error_monitoring"
SEND_TO_ERROR_MONITORING = {SEND_TO_ERROR_MONITORING_KEY: True}


class SentryStatus(Enum):
    SEND_ALL_ERRORS = 'send_all_errors'
    SEND_SELECTED_ERRORS = 'send_selected_errors'
    DISABLED = 'disabled'


def send_uncaught_errors_to_monitoring(enabled: bool):
    """Reusable decorator for logging unexpected exceptions."""
    def error_logging_wrapper(func):

        @wraps(func)
        def wrapper(*args, **kwargs):
            try:
                return func(*args, **kwargs)
            except Exception as err:
                logging.error("Uncaught error", exc_info=True, extra=SEND_TO_ERROR_MONITORING)
                raise err

        if enabled:
            return wrapper
        return func

    return error_logging_wrapper


def _ignore_if_not_marked(
        event: Dict[str, Any],
        hint: Dict[str, Any],  # pylint: disable=unused-argument
) -> Optional[Dict[str, Any]]:
    if event.get(LOGGING_EXTRA) and event.get(LOGGING_EXTRA).get(SEND_TO_ERROR_MONITORING_KEY):
        del event[LOGGING_EXTRA][SEND_TO_ERROR_MONITORING_KEY]
        return event
    return None


def get_sentry_status(service_name: str) -> SentryStatus:
    """Get Sentry status from service config value"""
    try:
        return SentryStatus(
            get_service_config_value(
                service_name,
                SENTRY_CONFIG,
                default=SentryStatus.DISABLED.value,
            ),
        )
    except ValueError:
        return SentryStatus.DISABLED


def sentry_init(service_name: str):
    """Initialize connection and start piping errors to sentry.io."""

    sentry_status = get_sentry_status(service_name)
    if sentry_status == SentryStatus.DISABLED:
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
        before_send=_ignore_if_not_marked if sentry_status == SentryStatus.SEND_SELECTED_ERRORS else None,
    )

    cloud_address = get_service_config_value(
        CONTROL_PROXY,
        CLOUD_ADDRESS,
        default=None,
    )
    sentry_sdk.set_tag(ORC8R_CLOUD_ADDRESS, cloud_address)
    sentry_sdk.set_tag(HWID, snowflake.snowflake())
    sentry_sdk.set_tag(SERVICE_NAME, service_name)
