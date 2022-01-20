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
import re
from dataclasses import dataclass
from enum import Enum
from functools import wraps
from typing import Any, Callable, Dict, List, Optional

import sentry_sdk
import snowflake
from magma.configuration.service_configs import get_service_config_value
from orc8r.protos.mconfig import mconfigs_pb2

Event = Dict[str, Any]
Hint = Dict[str, Any]
SentryHook = Callable[[Event, Hint], Optional[Event]]

CONTROL_PROXY = 'control_proxy'
SENTRY_CONFIG = 'sentry'
SENTRY_URL = 'sentry_url_python'
SENTRY_EXCLUDED = 'sentry_excluded_errors'
SENTRY_SAMPLE_RATE = 'sentry_sample_rate'
CLOUD_ADDRESS = 'cloud_address'
ORC8R_CLOUD_ADDRESS = 'orc8r_cloud_address'
DEFAULT_SAMPLE_RATE = 1.0
COMMIT_HASH = 'COMMIT_HASH'
HWID = 'hwid'
SERVICE_NAME = 'service_name'
LOGGING_EXTRA = 'extra'
SEND_TO_ERROR_MONITORING_KEY = "send_to_error_monitoring"
# Dictionary constant for convenience, must not be mutated
SEND_TO_ERROR_MONITORING = {SEND_TO_ERROR_MONITORING_KEY: True}  # noqa: WPS407


class SentryStatus(Enum):
    """Describes which kind of Sentry monitoring is configured"""
    SEND_ALL_ERRORS = 'send_all_errors'
    SEND_SELECTED_ERRORS = 'send_selected_errors'
    DISABLED = 'disabled'


@dataclass
class SharedSentryConfig(object):
    """Sentry configuration shared by all Python services,
    taken from shared mconfig or control_proxy.yml"""
    dsn: str
    sample_rate: float
    exclusion_patterns: List[str]


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
def _get_shared_sentry_config(sentry_mconfig: mconfigs_pb2.SharedSentryConfig) -> SharedSentryConfig:
    """Get Sentry configs with the following priority

    1) control_proxy.yml (if sentry_python_url is present)
    2) shared mconfig (i.e. first try streamed mconfig from orc8r,
    if empty: default mconfig in /etc/magma)

    Args:
        sentry_mconfig (SharedSentryConfig): proto message of shared mconfig

    Returns:
        (str, float): sentry url, sentry sample rate
    """
    dsn = get_service_config_value(
        CONTROL_PROXY,
        SENTRY_URL,
        default='',
    )

    if not dsn:
        # Here, we assume that `dsn` and `sample_rate` should be pulled
        # from the same source, that is the source where the user has
        # entered the `dsn`.
        # Without this coupling `dsn` and `sample_rate` could possibly
        # be pulled from different sources.
        dsn = sentry_mconfig.dsn_python
        sample_rate = sentry_mconfig.sample_rate
    else:
        logging.info("Sentry config: dsn_python and sample_rate are pulled from control_proxy.yml.")
        sample_rate = get_service_config_value(
            CONTROL_PROXY,
            SENTRY_SAMPLE_RATE,
            default=DEFAULT_SAMPLE_RATE,
        )

    # Exclusion patterns only exist in mconfig, not in control_proxy.yml
    exclusion_patterns = sentry_mconfig.exclusion_patterns
    return SharedSentryConfig(dsn, sample_rate, exclusion_patterns)


def _ignore_if_not_marked(event: Event) -> Optional[Event]:
    if event.get(LOGGING_EXTRA) and event.get(LOGGING_EXTRA).get(SEND_TO_ERROR_MONITORING_KEY):
        logging.info("Sending because of present tag")
        del event[LOGGING_EXTRA][SEND_TO_ERROR_MONITORING_KEY]
        return event
    logging.info("Ignoring because of missing tag")
    return None


def _filter_excluded_messages(event: Event, hint: Hint, patterns_to_exclude: List[str]) -> Optional[Event]:
    explicit_message = event.get('message')

    log_entry = event.get("logentry")
    log_message = log_entry.get('message') if log_entry else None

    exc_info = hint.get("exc_info")
    exception_message = str(exc_info[1]) if exc_info else None

    messages = [msg for msg in (explicit_message, log_message, exception_message) if msg]
    if not messages:
        return event

    for pattern in patterns_to_exclude:
        for message in messages:
            if re.search(pattern, message):
                return None

    return event


def _get_before_send_hook(
        sentry_status: SentryStatus,
        patterns_to_exclude: List[str],
) -> Optional[SentryHook]:

    def filter_excluded_and_unmarked_messages(
            event: Event, hint: Hint,
    ) -> Optional[Event]:
        event = _ignore_if_not_marked(event)
        if event and patterns_to_exclude:
            return _filter_excluded_messages(event, hint, patterns_to_exclude)
        return None

    def filter_excluded_messages(
            event: Event, hint: Hint,
    ) -> Optional[Event]:
        return _filter_excluded_messages(event, hint, patterns_to_exclude)

    if sentry_status == SentryStatus.SEND_SELECTED_ERRORS:
        return filter_excluded_and_unmarked_messages
    if patterns_to_exclude:
        return filter_excluded_messages
    return None


def sentry_init(service_name: str, sentry_mconfig: mconfigs_pb2.SharedSentryConfig) -> None:
    """Initialize connection and start piping errors to sentry.io."""

    sentry_status = get_sentry_status(service_name)
    if sentry_status == SentryStatus.DISABLED:
        return

    sentry_config = _get_shared_sentry_config(sentry_mconfig)

    if not sentry_config.dsn:
        logging.info(
            'Sentry disabled because of missing dsn_python. '
            'See documentation (Configure > AGW) on how to configure '
            'Sentry dsn.',
        )
        return

    sentry_sdk.init(
        dsn=sentry_config.dsn,
        release=os.getenv(COMMIT_HASH),
        traces_sample_rate=sentry_config.sample_rate,
        before_send=_get_before_send_hook(sentry_status, sentry_config.exclusion_patterns),
    )

    cloud_address = get_service_config_value(
        CONTROL_PROXY,
        CLOUD_ADDRESS,
        default=None,
    )
    sentry_sdk.set_tag(ORC8R_CLOUD_ADDRESS, cloud_address)
    sentry_sdk.set_tag(HWID, snowflake.snowflake())
    sentry_sdk.set_tag(SERVICE_NAME, service_name)
