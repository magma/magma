from __future__ import annotations

from dp.protos.active_mode_pb2 import (
    ActiveModeConfig,
    Cbsd,
    CbsdState,
    Channel,
    EirpCapabilities,
    FrequencyRange,
    Grant,
    GrantState,
)


class ActiveModeConfigBuilder:
    def __init__(self):
        self.desired_state = None
        self.cbsd_id = None
        self.user_id = None
        self.fcc_id = None
        self.serial_number = None
        self.state = None
        self.grants = None
        self.channels = None
        self.pending_requests = None
        self.last_seen_timestamp = None
        self.eirp_capabilities = None

    def build(self) -> ActiveModeConfig:
        cbsd = Cbsd(
            id=self.cbsd_id,
            user_id=self.user_id,
            fcc_id=self.fcc_id,
            serial_number=self.serial_number,
            state=self.state,
            grants=self.grants,
            channels=self.channels,
            pending_requests=self.pending_requests,
            last_seen_timestamp=self.last_seen_timestamp,
            eirp_capabilities=self.eirp_capabilities,
        )
        return ActiveModeConfig(
            desired_state=self.desired_state,
            cbsd=cbsd,
        )

    def with_desired_state(self, state: CbsdState) -> ActiveModeConfigBuilder:
        self.desired_state = state
        return self

    def with_state(self, state: CbsdState) -> ActiveModeConfigBuilder:
        self.state = state
        return self

    def with_registration(self, prefix: str) -> ActiveModeConfigBuilder:
        self.cbsd_id = f'{prefix}_cbsd_id'
        self.fcc_id = f'{prefix}_fcc_id'
        self.user_id = f'{prefix}_user_id'
        self.serial_number = f'{prefix}_serial_number'
        return self

    def with_eirp_capabilities(
        self,
        min_power: float, max_power: float,
        antenna_gain: float, no_ports: int,
    ) -> ActiveModeConfigBuilder:
        eirp_capabilities = EirpCapabilities(
            min_power=min_power,
            max_power=max_power,
            antenna_gain=antenna_gain,
            number_of_ports=no_ports,
        )
        self.eirp_capabilities = eirp_capabilities
        return self

    def with_grant(
        self,
        grant_id: str, state: GrantState,
        hb_interval_sec: int, last_hb_ts: int,
    ) -> ActiveModeConfigBuilder:
        if not self.grants:
            self.grants = []
        grant = Grant(
            id=grant_id,
            state=state,
            heartbeat_interval_sec=hb_interval_sec,
            last_heartbeat_timestamp=last_hb_ts,
        )
        self.grants.append(grant)
        return self

    def with_channel(
        self,
        low: int, high: int,
        max_eirp: float = None, last_eirp: float = None,
    ) -> ActiveModeConfigBuilder:
        if not self.channels:
            self.channels = []
        channel = Channel(
            frequency_range=FrequencyRange(low=low, high=high),
            max_eirp=max_eirp,
            last_eirp=last_eirp,
        )
        self.channels.append(channel)
        return self

    def with_pending_request(self, payload: str) -> ActiveModeConfigBuilder:
        if not self.pending_requests:
            self.pending_requests = []
        self.pending_requests.append(payload)
        return self

    def with_last_seen(self, last_seen_timestamp: int) -> ActiveModeConfigBuilder:
        self.last_seen_timestamp = last_seen_timestamp
        return self
