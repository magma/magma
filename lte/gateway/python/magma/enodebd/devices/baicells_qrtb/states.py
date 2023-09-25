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
from distutils.util import strtobool
from typing import Any

from dp.protos.cbsd_pb2 import CBSDStateResult
from magma.enodebd.data_models.data_model_parameters import ParameterName
from magma.enodebd.device_config.cbrs_consts import (
    BAND,
    SAS_MAX_POWER_SPECTRAL_DENSITY,
    SAS_MIN_POWER_SPECTRAL_DENSITY,
)
from magma.enodebd.device_config.configuration_util import (
    calc_bandwidth_mhz,
    calc_bandwidth_rbs,
    calc_earfcn,
)
from magma.enodebd.device_config.enodeb_configuration import EnodebConfiguration
from magma.enodebd.devices.baicells_qrtb.params import (
    CarrierAggregationParameters,
)
from magma.enodebd.dp_client import (
    build_enodebd_update_cbsd_request,
    enodebd_update_cbsd,
)
from magma.enodebd.exceptions import ConfigurationError
from magma.enodebd.logger import EnodebdLogger
from magma.enodebd.state_machines.acs_state_utils import process_inform_message
from magma.enodebd.state_machines.enb_acs import EnodebAcsStateMachine
from magma.enodebd.state_machines.enb_acs_states import (
    AcsMsgAndTransition,
    AcsReadMsgResult,
    EnodebAcsState,
    NotifyDPState,
    WaitInformMRebootState,
)
from magma.enodebd.state_machines.timer import StateMachineTimer
from magma.enodebd.tr069 import models

logger = EnodebdLogger


class BaicellsQRTBEndSessionState(EnodebAcsState):
    """ To end a TR-069 session, send an empty HTTP response

    For Baicells QRTB we can expect an inform message on
    End Session state, either a queued one or a periodic one
    """

    def __init__(
            self,
            acs: EnodebAcsStateMachine,
            when_done: str,
    ):
        super().__init__()
        self.acs = acs
        self.done_transition = when_done

    def get_msg(self, message: Any) -> AcsMsgAndTransition:
        """
        Send back a message to enb

        Args:
            message (Any): TR069 message

        Returns:
            AcsMsgAndTransition
        """
        request = models.DummyInput()
        return AcsMsgAndTransition(msg=request, next_state=self.done_transition)

    def state_description(self) -> str:
        """
        Describe the state

        Returns:
            str
        """
        return 'Completed provisioning eNB. Notifying DP.'


class BaicellsQRTBQueuedEventsWaitState(EnodebAcsState):
    """
    We've already received an Inform message. This state is to handle a
    Baicells eNodeB issue.

    After eNodeB is rebooted, hold off configuring it for some time.

    In this state, just hang at responding to Inform, and then ending the
    TR-069 session.
    """

    CONFIG_DELAY_AFTER_BOOT = 60

    def __init__(self, acs: EnodebAcsStateMachine, when_done: str):
        super().__init__()
        self.acs = acs
        self.done_transition = when_done
        self.wait_timer = None

    def enter(self):
        """
        Perform additional actions on state enter
        """
        self.wait_timer = StateMachineTimer(self.CONFIG_DELAY_AFTER_BOOT)
        logger.info(
            'Holding off of eNB configuration for %s seconds. ',
            self.CONFIG_DELAY_AFTER_BOOT,
        )

    def exit(self):
        """
        Perform additional actions on state exit
        """
        self.wait_timer = None

    def read_msg(self, message: Any) -> AcsReadMsgResult:
        """
        Read incoming message

        Args:
            message (Any): TR069 message

        Returns:
            AcsReadMsgResult
        """
        if not isinstance(message, models.Inform):
            return AcsReadMsgResult(msg_handled=False, next_state=None)
        process_inform_message(
            message, self.acs.data_model,
            self.acs.device_cfg,
        )
        return AcsReadMsgResult(msg_handled=True, next_state=None)

    def get_msg(self, message: Any) -> AcsMsgAndTransition:
        """
        Send back a message to enb

        Args:
            message (Any): TR069 message

        Returns:
            AcsMsgAndTransition
        """
        if not self.wait_timer:
            logger.error('wait_timer is None.')
            raise ValueError('wait_timer is None.')

        if self.wait_timer.is_done():
            return AcsMsgAndTransition(
                msg=models.DummyInput(),
                next_state=self.done_transition,
            )
        remaining = self.wait_timer.seconds_remaining()
        logger.info(
            'Waiting with eNB configuration for %s more seconds. ',
            remaining,
        )
        return AcsMsgAndTransition(msg=models.DummyInput(), next_state=None)

    def state_description(self) -> str:
        """
        Describe the state

        Returns:
            str
        """
        if not self.wait_timer:
            logger.error('wait_timer is None.')
            raise ValueError('wait_timer is None.')

        remaining = self.wait_timer.seconds_remaining()
        return 'Waiting for eNB REM to run for %d more seconds before ' \
               'resuming with configuration.' % remaining


class BaicellsQRTBWaitInformRebootState(WaitInformMRebootState):
    """
    BaicellsQRTB WaitInformRebootState implementation
    """
    INFORM_EVENT_CODE = '1 BOOT'


class BaicellsQRTBNotifyDPState(NotifyDPState):
    """
        BaicellsQRTB NotifyDPState implementation
    """

    def enter(self):
        """
        Enter the state
        """
        serial_number = self.acs.device_cfg.get_parameter(ParameterName.SERIAL_NUMBER)

        # NOTE: In case GPS scan is not completed, eNB reports LAT and LONG values as 0.
        #       Only update CBSD in Domain Proxy when all params are available.
        gps_status = strtobool(self.acs.device_cfg.get_parameter(ParameterName.GPS_STATUS))
        if gps_status:
            enodebd_update_cbsd_request = build_enodebd_update_cbsd_request(
                serial_number=serial_number,
                latitude_deg=self.acs.device_cfg.get_parameter(ParameterName.GPS_LAT),
                longitude_deg=self.acs.device_cfg.get_parameter(ParameterName.GPS_LONG),
                indoor_deployment=self.acs.device_cfg.get_parameter(ParameterName.INDOOR_DEPLOYMENT),
                antenna_height=self.acs.device_cfg.get_parameter(ParameterName.ANTENNA_HEIGHT),
                antenna_height_type=self.acs.device_cfg.get_parameter(ParameterName.ANTENNA_HEIGHT_TYPE),
                cbsd_category=self.acs.device_cfg.get_parameter(ParameterName.CBSD_CATEGORY),
            )
            state = enodebd_update_cbsd(enodebd_update_cbsd_request)
            qrtb_update_desired_config_from_cbsd_state(state, self.acs.desired_cfg)
        else:
            EnodebdLogger.debug("Waiting for GPS to sync, before updating CBSD params in Domain Proxy.")


def _qrtb_check_state_compatibility_with_ca(state: CBSDStateResult) -> bool:
    """
    Check if state returned by Domain Proxy contains data that can be applied
    to BaiCells QRTB BS in Carrier Aggregation Mode.
    BaiCells QRTB can apply carrier aggregation if:
    * 2 channels are available:
      * with symmetric bandwidths: 5+5, 10+10 or 20+20
      * with asymmetric bandwidths: 20+10 (10+20 theoretically supported)
    * Max IBW of the channels is 100MHz
        (IBW == max(high_frequency_hz) - min(low_frequency_hz), max and min taken from parameters of the 2 channels)

    Additionally, such channels may be available but Domain Proxy may explicitly disable CA

    Only check the first 2 channels (Domain Proxy may return more)
    """
    _MAX_IBW_HZ = 100_000_000
    _CA_SUPPORTED_BANDWIDTHS_MHZ = (
        (5, 5),
        (10, 10),
        (20, 20),
        (20, 10),
    )
    num_of_channels = len(state.channels)
    # Check if CA explicitly disabled, or not enough channels
    if num_of_channels < 2 or not state.carrier_aggregation_enabled:
        logger.debug(f"Domain Proxy state {num_of_channels=}, {state.carrier_aggregation_enabled=}.")
        return False

    ch1 = state.channels[0]
    ch2 = state.channels[1]

    # Check Max IBW
    high_frequency_hz = max(ch1.high_frequency_hz, ch2.high_frequency_hz)
    low_frequency_hz = min(ch1.low_frequency_hz, ch2.low_frequency_hz)
    if high_frequency_hz - low_frequency_hz > _MAX_IBW_HZ:
        logger.debug(f"Domain Proxy channel1 {ch1}, channel2 {ch2} exceed max IBW {_MAX_IBW_HZ}.")
        return False

    # Check supported bandwidths of the channels
    bw1 = calc_bandwidth_mhz(low_freq_hz=ch1.low_frequency_hz, high_freq_hz=ch1.high_frequency_hz)
    bw2 = calc_bandwidth_mhz(low_freq_hz=ch2.low_frequency_hz, high_freq_hz=ch2.high_frequency_hz)
    if not (bw1, bw2) in _CA_SUPPORTED_BANDWIDTHS_MHZ:
        logger.debug(f"Domain Proxy channel1 {ch1}, channel2 {ch2}, bandwidth configuration not in {_CA_SUPPORTED_BANDWIDTHS_MHZ}.")
        return False

    return True


def qrtb_update_desired_config_from_cbsd_state(state: CBSDStateResult, config: EnodebConfiguration) -> None:
    """
    Call grpc endpoint on the Domain Proxy to update the desired config based on sas grant

    Args:
        state (CBSDStateResult): state result as received from DP
        config (EnodebConfiguration): configuration to update
    """
    logger.debug("Updating desired config based on Domain Proxy state.")
    num_of_channels = len(state.channels)
    radio_enabled = num_of_channels > 0 and state.radio_enabled
    config.set_parameter(ParameterName.SAS_RADIO_ENABLE, radio_enabled)

    if not radio_enabled:
        return

    # FAPService.1
    channel = state.channels[0]
    earfcn = calc_earfcn(channel.low_frequency_hz, channel.high_frequency_hz)
    bandwidth_mhz = calc_bandwidth_mhz(channel.low_frequency_hz, channel.high_frequency_hz)
    bandwidth_rbs = calc_bandwidth_rbs(bandwidth_mhz)
    psd = _calc_psd(channel.max_eirp_dbm_mhz)
    logger.debug(f"Channel1: {earfcn=}, {bandwidth_rbs=}, {psd=}")

    can_enable_carrier_aggregation = _qrtb_check_state_compatibility_with_ca(state)
    logger.debug(f"Should Carrier Aggregation be enabled on eNB: {can_enable_carrier_aggregation=}")
    # Enabling Carrier Aggregation on QRTB eNB means:
    # 1. Set CA_ENABLE to 1
    # 2. Set CA_NUM_OF_CELLS to 2
    # 3. Configure appropriate TR nodes for FAPSerivce.2 like EARFCNDL/UL etc
    # Otherwise we need to disable Carrier Aggregation on eNB and switch to Single Carrier configuration
    # 1. Set CA_ENABLE to 0
    # 2. Set CA_NUM_OF_CELLS to 1
    # Those two nodes should handle everything else accordingly.
    # (NOTE: carrier aggregation may still be enabled on Domain Proxy, but Domain Proxy may not have 2 channels granted by SAS.
    #        In such case, we still have to switch eNB to Single Carrier)
    num_of_cells = 2 if can_enable_carrier_aggregation else 1
    ca_enable = 1 if can_enable_carrier_aggregation else 0

    params_to_set = {
        ParameterName.SAS_RADIO_ENABLE: True,
        ParameterName.BAND: BAND,
        ParameterName.DL_BANDWIDTH: bandwidth_rbs,
        ParameterName.UL_BANDWIDTH: bandwidth_rbs,
        ParameterName.EARFCNDL: earfcn,
        ParameterName.EARFCNUL: earfcn,
        ParameterName.POWER_SPECTRAL_DENSITY: psd,
        CarrierAggregationParameters.CA_ENABLE: ca_enable,
        CarrierAggregationParameters.CA_NUM_OF_CELLS: num_of_cells,
    }
    if can_enable_carrier_aggregation:
        # Configure FAPService.2
        # NOTE: We set PCI and CELL_ID to the values of FAP1 "+1"
        #       This was suggested by BaiCells
        channel = state.channels[1]
        earfcn = calc_earfcn(channel.low_frequency_hz, channel.high_frequency_hz)
        bandwidth_mhz = calc_bandwidth_mhz(channel.low_frequency_hz, channel.high_frequency_hz)
        bandwidth_rbs = calc_bandwidth_rbs(bandwidth_mhz)
        psd = _calc_psd(channel.max_eirp_dbm_mhz)
        logger.debug(f"Channel2: {earfcn=}, {bandwidth_rbs=}, {psd=}")
        params_to_set.update({
            CarrierAggregationParameters.CA_DL_BANDWIDTH: bandwidth_rbs,
            CarrierAggregationParameters.CA_UL_BANDWIDTH: bandwidth_rbs,
            CarrierAggregationParameters.CA_BAND: BAND,
            CarrierAggregationParameters.CA_EARFCNDL: earfcn,
            CarrierAggregationParameters.CA_EARFCNUL: earfcn,
            CarrierAggregationParameters.CA_PCI: config.get_parameter(ParameterName.PCI) + 1,
            CarrierAggregationParameters.CA_CELL_ID: config.get_parameter(ParameterName.CELL_ID) + 1,
            CarrierAggregationParameters.CA_RADIO_ENABLE: True,
        })

    for param, value in params_to_set.items():
        config.set_parameter(param, value)


def _calc_psd(eirp: float) -> int:
    psd = int(eirp)
    if not SAS_MIN_POWER_SPECTRAL_DENSITY <= psd <= SAS_MAX_POWER_SPECTRAL_DENSITY:  # noqa: WPS508
        raise ConfigurationError(
            'Power Spectral Density %d exceeds allowed range [%d, %d]' %
            (psd, SAS_MIN_POWER_SPECTRAL_DENSITY, SAS_MAX_POWER_SPECTRAL_DENSITY),
        )
    return psd
