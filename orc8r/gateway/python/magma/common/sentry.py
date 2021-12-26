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
from typing import Any, Dict, Optional, Tuple

import sentry_sdk
import snowflake
from magma.configuration.service_configs import get_service_config_value
from orc8r.protos.mconfig import mconfigs_pb2

CONTROL_PROXY = 'control_proxy'
SENTRY_CONFIG = 'sentry'
SENTRY_URL = 'sentry_url_python'
SENTRY_SAMPLE_RATE = 'sentry_sample_rate'
CLOUD_ADDRESS = 'cloud_address'
ORC8R_CLOUD_ADDRESS = 'orc8r_cloud_address'
DEFAULT_SAMPLE_RATE = 1.0
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


# TODO when control_proxy.yml is outdated move to shared mconfig entirely
def get_sentry_dsn_and_sample_rate(sentry_mconfig: mconfigs_pb2.SharedSentryConfig) -> Tuple[str, float]:
    """Get Sentry configs with the following priority

    1) control_proxy.yml (if sentry_python_url is present)
    2) shared mconfig (i.e. first streamed mconfig from orc8r,
    if not present default mconfig in /etc/magma)

    Args:
        sentry_mconfig (SharedSentryConfig): proto message of shared mconfig

    Returns:
        (str, float): sentry url, sentry sample rate
    """
    dsn_python = get_service_config_value(
        CONTROL_PROXY,
        SENTRY_URL,
        default='',
    )

    if not dsn_python:
        dsn_python = sentry_mconfig.dsn_python
        sample_rate = sentry_mconfig.sample_rate
        return dsn_python, sample_rate

    sample_rate = get_service_config_value(
        CONTROL_PROXY,
        SENTRY_SAMPLE_RATE,
        default=DEFAULT_SAMPLE_RATE,
    )
    return dsn_python, sample_rate


def sentry_init(service_name: str, sentry_mconfig: mconfigs_pb2.SharedSentryConfig) -> None:
    """Initialize connection and start piping errors to sentry.io."""

    sentry_status = get_sentry_status(service_name)
    if sentry_status == SentryStatus.DISABLED:
        return

    dsn_python, sample_rate = get_sentry_dsn_and_sample_rate(sentry_mconfig)

    if not dsn_python:
        logging.info(
            'Sentry disabled because of missing dsn_python. '
            'See documentation (Configure > AGW) on how to configure '
            'Sentry dsn.',
        )
        return

    sentry_sdk.init(
        dsn=dsn_python,
        release=os.getenv(COMMIT_HASH),
        traces_sample_rate=sample_rate,
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
