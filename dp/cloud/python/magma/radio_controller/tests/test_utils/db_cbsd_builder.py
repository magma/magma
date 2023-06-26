"""
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""
from __future__ import annotations

import json
from datetime import datetime
from typing import List

from magma.db_service.models import DBCbsd, DBGrant, DBRequest


class DBCbsdBuilder:
    def __init__(self):
        self.cbsd = DBCbsd()

    def build(self) -> DBCbsd:
        return self.cbsd

    def deleted(self):
        self.cbsd.is_deleted = True
        return self

    def updated(self):
        self.cbsd.should_deregister = True
        return self

    def relinquished(self):
        self.cbsd.should_relinquish = True
        return self

    def with_id(self, db_id: int) -> DBCbsdBuilder:
        self.cbsd.id = db_id
        return self

    def with_state(self, state_id: int) -> DBCbsdBuilder:
        self.cbsd.state_id = state_id
        return self

    def with_registration(self, prefix: str) -> DBCbsdBuilder:
        self.cbsd.cbsd_id = f'{prefix}_cbsd_id'
        self.cbsd.user_id = f'{prefix}_user_id'
        self.cbsd.fcc_id = f'{prefix}_fcc_id'
        self.cbsd.cbsd_serial_number = f'{prefix}_serial_number'
        return self

    def with_eirp_capabilities(
        self,
        min_power: float, max_power: float,
        no_ports: int,
    ) -> DBCbsdBuilder:
        self.cbsd.min_power = min_power
        self.cbsd.max_power = max_power
        self.cbsd.number_of_ports = no_ports
        return self

    def with_single_step_enabled(self) -> DBCbsdBuilder:
        self.cbsd.single_step_enabled = True
        return self

    def with_category(self, category: str) -> DBCbsdBuilder:
        self.cbsd.cbsd_category = category
        return self

    def with_antenna_gain(
        self,
        antenna_gain_dbi: float,
    ) -> DBCbsdBuilder:
        self.cbsd.antenna_gain = antenna_gain_dbi
        return self

    def with_installation_params(
        self,
        latitude_deg: float,
        longitude_deg: float,
        height_m: float,
        height_type: str,
        indoor_deployment: bool,
    ) -> DBCbsdBuilder:
        self.cbsd.latitude_deg = latitude_deg
        self.cbsd.longitude_deg = longitude_deg
        self.cbsd.height_m = height_m
        self.cbsd.height_type = height_type
        self.cbsd.indoor_deployment = indoor_deployment
        return self

    def with_last_seen(self, last_seen: int) -> DBCbsdBuilder:
        self.cbsd.last_seen = datetime.fromtimestamp(last_seen)
        return self

    def with_desired_state(self, desired_state_id: int) -> DBCbsdBuilder:
        self.cbsd.desired_state_id = desired_state_id
        return self

    def with_preferences(self, bandwidth_mhz: int, frequencies_mhz: List[int]) -> DBCbsdBuilder:
        self.cbsd.preferred_bandwidth_mhz = bandwidth_mhz
        self.cbsd.preferred_frequencies_mhz = frequencies_mhz
        return self

    def with_available_frequencies(self, frequencies: List[int]):
        self.cbsd.available_frequencies = frequencies
        return self

    def with_carrier_aggregation(self, enabled: bool) -> DBCbsdBuilder:
        self.cbsd.carrier_aggregation_enabled = enabled
        return self

    def with_max_ibw(self, max_ibw_mhz: int) -> DBCbsdBuilder:
        self.cbsd.max_ibw_mhz = max_ibw_mhz
        return self

    def with_grant_redundancy(self, enabled: bool) -> DBCbsdBuilder:
        self.cbsd.grant_redundancy = enabled
        return self

    def with_grant(
        self,
        grant_id: str,
        state_id: int,
        hb_interval_sec: int,
        last_hb_timestamp: int = None,
        low_frequency: int = 3500,
        high_frequency: int = 3700,
    ) -> DBCbsdBuilder:
        last_hb_time = datetime.fromtimestamp(
            last_hb_timestamp,
        ) if last_hb_timestamp else None
        grant = DBGrant(
            grant_id=grant_id,
            state_id=state_id,
            heartbeat_interval=hb_interval_sec,
            last_heartbeat_request_time=last_hb_time,
            low_frequency=low_frequency,
            high_frequency=high_frequency,
            max_eirp=0,
        )
        self.cbsd.grants.append(grant)
        return self

    def with_channel(
        self,
        low: int, high: int,
        max_eirp: float = None,
    ) -> DBCbsdBuilder:
        if not self.cbsd.channels:
            # Default is set on commit, so it might be None at this point.
            self.cbsd.channels = []

        channel = {
            "low_frequency": low,
            "high_frequency": high,
            "max_eirp": max_eirp,
        }
        self.cbsd.channels = self.cbsd.channels + [channel]
        return self

    def with_request(self, type_id: int, payload: str) -> DBCbsdBuilder:
        request = DBRequest(
            type_id=type_id,
            payload=json.loads(payload),
        )
        self.cbsd.requests.append(request)
        return self
