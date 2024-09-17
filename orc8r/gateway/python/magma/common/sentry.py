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
from typing import Any, Callable, Dict, List, Optional

import sentry_sdk
import snowflake
from magma.configuration.service_configs import get_service_config_value
from orc8r.protos.mconfig import mconfigs_pb2
from sentry_sdk.integrations.redis import RedisIntegration

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
EXCLUDE_FROM_ERROR_MONITORING_KEY = 'exclude_from_error_monitoring'
# Dictionary constant for convenience, must not be mutated
EXCLUDE_FROM_ERROR_MONITORING = {EXCLUDE_FROM_ERROR_MONITORING_KEY: True}  # noqa: WPS407


@dataclass
class SharedSentryConfig(object):
    """Sentry configuration shared by all Python services,
    taken from shared mconfig or control_proxy.yml"""
    dsn: str
    sample_rate: float
    exclusion_patterns: List[str]


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


def _ignore_if_marked(event: Event) -> Optional[Event]:
    if event.get(LOGGING_EXTRA) and event.get(LOGGING_EXTRA).get(EXCLUDE_FROM_ERROR_MONITORING_KEY):
        return None
    return event


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


def _get_before_send_hook(patterns_to_exclude: List[str]) -> SentryHook:

    def filter_excluded_and_marked_messages(
            event: Event, hint: Hint,
    ) -> Optional[Event]:
        event = _ignore_if_marked(event)
        if event:
            return _filter_excluded_messages(event, hint, patterns_to_exclude)
        return None

    def filter_marked_messages(
            event: Event, _: Hint,
    ) -> Optional[Event]:
        return _ignore_if_marked(event)

    if patterns_to_exclude:
        return filter_excluded_and_marked_messages
    return filter_marked_messages


def sentry_init(service_name: str, sentry_mconfig: mconfigs_pb2.SharedSentryConfig) -> None:
    """Initialize connection and start piping errors to sentry.io."""

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
        before_send=_get_before_send_hook(sentry_config.exclusion_patterns),
        integrations=[
            RedisIntegration(),
        ],
    )

    cloud_address = get_service_config_value(
        CONTROL_PROXY,
        CLOUD_ADDRESS,
        default=None,
    )
    sentry_sdk.set_tag(ORC8R_CLOUD_ADDRESS, cloud_address)
    sentry_sdk.set_tag(HWID, snowflake.snowflake())
    sentry_sdk.set_tag(SERVICE_NAME, service_name)
