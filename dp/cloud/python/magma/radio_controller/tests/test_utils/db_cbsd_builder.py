from __future__ import annotations

import json
from datetime import datetime
from typing import List

from magma.db_service.models import (
    DBActiveModeConfig,
    DBCbsd,
    DBChannel,
    DBGrant,
    DBRequest,
)


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
        antenna_gain: float, no_ports: int,
    ) -> DBCbsdBuilder:
        self.cbsd.min_power = min_power
        self.cbsd.max_power = max_power
        self.cbsd.antenna_gain = antenna_gain
        self.cbsd.number_of_ports = no_ports
        return self

    def with_last_seen(self, last_seen: int) -> DBCbsdBuilder:
        self.cbsd.last_seen = datetime.fromtimestamp(last_seen)
        return self

    def with_grant_attempts(self, grant_attempts: int) -> DBCbsdBuilder:
        self.cbsd.grant_attempts = grant_attempts
        return self

    def with_active_mode_config(self, desired_state_id: int) -> DBCbsdBuilder:
        config = DBActiveModeConfig(
            desired_state_id=desired_state_id,
        )
        self.cbsd.active_mode_config.append(config)
        return self

    def with_preferences(self, bandwidth_mhz: int, frequencies_mhz: List[int]) -> DBCbsdBuilder:
        self.cbsd.preferred_bandwidth_mhz = bandwidth_mhz
        self.cbsd.preferred_frequencies_mhz = frequencies_mhz
        return self

    def with_grant(
        self,
        grant_id: str, state_id: int,
        hb_interval_sec: int, last_hb_timestamp: int = None,
    ) -> DBCbsdBuilder:
        last_hb_time = datetime.fromtimestamp(
            last_hb_timestamp,
        ) if last_hb_timestamp else None
        grant = DBGrant(
            grant_id=grant_id,
            state_id=state_id,
            heartbeat_interval=hb_interval_sec,
            last_heartbeat_request_time=last_hb_time,
            low_frequency=0,
            high_frequency=0,
            max_eirp=0,
        )
        self.cbsd.grants.append(grant)
        return self

    def with_channel(
        self,
        low: int, high: int,
        max_eirp: float = None,
    ) -> DBCbsdBuilder:
        channel = DBChannel(
            low_frequency=low,
            high_frequency=high,
            max_eirp=max_eirp,
            channel_type='channel_type',
            rule_applied='rule',
        )
        self.cbsd.channels.append(channel)
        return self

    def with_request(self, type_id: int, payload: str) -> DBCbsdBuilder:
        request = DBRequest(
            type_id=type_id,
            payload=json.loads(payload),
        )
        self.cbsd.requests.append(request)
        return self
