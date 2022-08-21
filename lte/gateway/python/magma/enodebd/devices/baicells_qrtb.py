"""
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""
from distutils.util import strtobool
from typing import Any, Callable, Dict, List, Optional

from dp.protos.cbsd_pb2 import CBSDStateResult
from magma.common.service import MagmaService
from magma.enodebd.data_models import transform_for_magma
from magma.enodebd.data_models.data_model import DataModel, TrParam
from magma.enodebd.data_models.data_model_parameters import (
    ParameterName,
    TrParameterType,
)
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
from magma.enodebd.device_config.enodeb_config_postprocessor import (
    EnodebConfigurationPostProcessor,
)
from magma.enodebd.device_config.enodeb_configuration import EnodebConfiguration
from magma.enodebd.devices.device_utils import EnodebDeviceName
from magma.enodebd.dp_client import (
    build_enodebd_update_cbsd_request,
    enodebd_update_cbsd,
)
from magma.enodebd.exceptions import ConfigurationError
from magma.enodebd.logger import EnodebdLogger
from magma.enodebd.state_machines.acs_state_utils import process_inform_message
from magma.enodebd.state_machines.enb_acs import EnodebAcsStateMachine
from magma.enodebd.state_machines.enb_acs_impl import BasicEnodebAcsStateMachine
from magma.enodebd.state_machines.enb_acs_states import (
    AcsMsgAndTransition,
    AcsReadMsgResult,
    AddObjectsState,
    CheckFirmwareUpgradeDownloadState,
    DeleteObjectsState,
    EnbSendRebootState,
    EnodebAcsState,
    ErrorState,
    FirmwareUpgradeDownloadState,
    GetObjectParametersState,
    GetParametersState,
    NotifyDPState,
    SendGetTransientParametersState,
    SetParameterValuesState,
    WaitEmptyMessageState,
    WaitForFirmwareUpgradeDownloadResponse,
    WaitGetObjectParametersState,
    WaitGetParametersState,
    WaitGetTransientParametersState,
    WaitInformMRebootState,
    WaitInformState,
    WaitRebootResponseState,
    WaitSetParameterValuesState,
)
from magma.enodebd.state_machines.timer import StateMachineTimer
from magma.enodebd.tr069 import models

logger = EnodebdLogger


class BaicellsQRTBHandler(BasicEnodebAcsStateMachine):
    """
    BaicellsQRTB State Machine
    """

    def __init__(
            self,
            service: MagmaService,
    ) -> None:
        self._state_map: Dict[str, Any] = {}
        super().__init__(service, use_param_key=False)

    def reboot_asap(self) -> None:
        """
        Transition to 'reboot' state
        """
        self.transition('reboot')

    def is_enodeb_connected(self) -> bool:
        """
        Check if enodebd has received an Inform from the enodeb

        Returns:
            bool
        """
        return not isinstance(self.state, WaitInformState)

    def _init_state_map(self) -> None:
        self._state_map = {
            # RemWait state seems not needed for QRTB
            'wait_inform': WaitInformState(self, when_done='wait_empty', when_boot=None),
            'wait_empty': WaitEmptyMessageState(self, when_done='check_fw_upgrade_download'),

            # Download flow
            'check_fw_upgrade_download': CheckFirmwareUpgradeDownloadState(
                self,
                when_download='fw_upgrade_download',
                when_skip='get_transient_params',
            ),
            'fw_upgrade_download': FirmwareUpgradeDownloadState(
                self,
                when_done='wait_fw_upgrade_download_response',
            ),
            'wait_fw_upgrade_download_response': WaitForFirmwareUpgradeDownloadResponse(
                self,
                when_done='get_transient_params',
                when_skip='get_transient_params',
            ),
            # Download flow ends

            'get_transient_params': SendGetTransientParametersState(self, when_done='wait_get_transient_params'),
            'wait_get_transient_params': WaitGetTransientParametersState(
                self,
                when_get='get_params',
                when_get_obj_params='get_obj_params',
                when_delete='delete_objs',
                when_add='add_objs',
                when_set='set_params',
                when_skip='end_session',
                request_all_params=True,
            ),
            'get_params': GetParametersState(self, when_done='wait_get_params', request_all_params=True),
            'wait_get_params': WaitGetParametersState(self, when_done='get_obj_params'),
            'get_obj_params': GetObjectParametersState(self, when_done='wait_get_obj_params', request_all_params=True),
            'wait_get_obj_params': WaitGetObjectParametersState(
                self, when_delete='delete_objs', when_add='add_objs',
                when_set='set_params', when_skip='end_session',
            ),
            'delete_objs': DeleteObjectsState(self, when_add='add_objs', when_skip='set_params'),
            'add_objs': AddObjectsState(self, when_done='set_params'),
            'set_params': SetParameterValuesState(self, when_done='wait_set_params'),
            'wait_set_params': WaitSetParameterValuesState(
                self, when_done='check_get_params',
                when_apply_invasive='reboot',
            ),
            'check_get_params': GetParametersState(
                self,
                when_done='check_wait_get_params',
                request_all_params=True,
            ),
            'check_wait_get_params': WaitGetParametersState(self, when_done='end_session'),
            'end_session': BaicellsQRTBEndSessionState(self, when_done='notify_dp'),
            'notify_dp': BaicellsQRTBNotifyDPState(self, when_inform='wait_inform'),
            'reboot': EnbSendRebootState(self, when_done='wait_reboot'),
            'wait_reboot': WaitRebootResponseState(self, when_done='wait_post_reboot_inform'),
            'wait_post_reboot_inform': BaicellsQRTBWaitInformRebootState(
                self,
                when_done='wait_queued_events_post_reboot',
                when_timeout='wait_inform_post_reboot',
            ),
            "wait_queued_events_post_reboot": BaicellsQRTBQueuedEventsWaitState(
                self,
                when_done='wait_inform_post_reboot',
            ),
            'wait_inform_post_reboot': WaitInformState(self, when_done='wait_empty_post_reboot', when_boot=None),
            'wait_empty_post_reboot': WaitEmptyMessageState(
                self, when_done='get_transient_params',
                when_missing='check_optional_params',
            ),
            # The states below are entered when an unexpected message type is
            # received
            'unexpected_fault': ErrorState(self, inform_transition_target='wait_inform'),
        }

    @property
    def device_name(self) -> str:
        """
        Return the device name

        Returns:
            device name
        """
        return EnodebDeviceName.BAICELLS_QRTB

    @property
    def data_model_class(self) -> DataModel:
        """
        Return the class of the data model

        Returns:
            DataModel
        """
        return BaicellsQRTBTrDataModel

    @property
    def config_postprocessor(self) -> EnodebConfigurationPostProcessor:
        """
        Return the instance of config postprocessor

        Returns:
            EnodebConfigurationPostProcessor
        """
        return BaicellsQRTBTrConfigurationInitializer()

    @property
    def state_map(self) -> Dict[str, EnodebAcsState]:
        """
        Return the state map for the State Machine

        Returns:
            Dict[str, EnodebAcsState]
        """
        return self._state_map

    @property
    def disconnected_state_name(self) -> str:
        """
        Return the string representation of a disconnected state

        Returns:
            str
        """
        return 'wait_inform'

    @property
    def unexpected_fault_state_name(self) -> str:
        """
        Return the string representation of an unexpected fault state

        Returns:
            str
        """
        return 'unexpected_fault'


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


class CarrierAggregationParameters:
    """
    Class defines additional TR parameters used to configure Carrier Aggregation

    Currently there is no good way of achieving parameter extensions in data models.
    Idea taken from FreedomFi one model, where data model PARAMETERS
    is updated on the class definition level (bad).
    """
    FAPSERVICE2_PATH = "Device.Services.FAPService.2."

    CA_ENABLE = 'Carrier Aggregation Enabled'
    CA_NUM_OF_CELLS = 'CA Number of Cells'
    CA_CELL_ID = 'CA Cell ID'
    CA_BAND = 'CA Band'
    CA_DL_BANDWIDTH = 'CA DL bandwidth'
    CA_UL_BANDWIDTH = 'CA UL bandwidth'
    CA_PCI = 'CA PCI'
    CA_EARFCNDL = 'CA EARFCNDL'
    CA_EARFCNUL = 'CA EARFCNUL'
    CA_ADMIN_STATE = 'CA Admin State'
    CA_OP_STATE = 'CA Op State'
    CA_RF_TX_STATUS = 'CA RF TX status'
    CA_RADIO_ENABLE = 'CA Radio Enable'

    CA_PLMN_CELL_RESERVED = 'CA PLMN 1 cell reserved'
    CA_PLMN_ENABLE = 'CA PLMN 1 enable'
    CA_PLMN_PRIMARY = 'CA PLMN 1 primary'
    CA_PLMN_PLMNID = 'CA PLMN 1 PLMNID'

    CA_PARAMETERS = {
        CA_ENABLE: TrParam(
            path='Device.Services.FAPService.1.CellConfig.LTE.RAN.CA.CaEnable',
            is_invasive=False,
            type=TrParameterType.INT,
            is_optional=False,
        ),
        CA_NUM_OF_CELLS: TrParam(
            path='FAPService.1.CellConfig.LTE.RAN.CA.PARAMS.NumOfCells',
            is_invasive=False,
            type=TrParameterType.UNSIGNED_INT,
            is_optional=False,
        ),
        CA_CELL_ID: TrParam(
            path=FAPSERVICE2_PATH + 'CellConfig.LTE.RAN.Common.CellIdentity',
            is_invasive=True,
            type=TrParameterType.UNSIGNED_INT,
            is_optional=False,
        ),
        CA_BAND: TrParam(
            path=FAPSERVICE2_PATH + 'CellConfig.LTE.RAN.RF.FreqBandIndicator',
            is_invasive=True,
            type=TrParameterType.UNSIGNED_INT,
            is_optional=False,
        ),
        CA_DL_BANDWIDTH: TrParam(
            path=FAPSERVICE2_PATH + 'CellConfig.LTE.RAN.RF.DLBandwidth',
            is_invasive=True,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        CA_UL_BANDWIDTH: TrParam(
            path=FAPSERVICE2_PATH + 'CellConfig.LTE.RAN.RF.ULBandwidth',
            is_invasive=True,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        CA_PCI: TrParam(
            path=FAPSERVICE2_PATH + 'CellConfig.LTE.RAN.RF.PhyCellID',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        CA_EARFCNDL: TrParam(
            path=FAPSERVICE2_PATH + 'CellConfig.LTE.RAN.RF.EARFCNDL',
            is_invasive=True,
            type=TrParameterType.UNSIGNED_INT,
            is_optional=False,
        ),
        CA_EARFCNUL: TrParam(
            path=FAPSERVICE2_PATH + 'CellConfig.LTE.RAN.RF.EARFCNUL',
            is_invasive=True,
            type=TrParameterType.UNSIGNED_INT,
            is_optional=False,
        ),
        CA_ADMIN_STATE: TrParam(
            path=FAPSERVICE2_PATH + 'FAPControl.LTE.AdminState',
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        CA_OP_STATE: TrParam(
            path=FAPSERVICE2_PATH + 'FAPControl.LTE.OpState',
            is_invasive=True,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        CA_RF_TX_STATUS: TrParam(
            path=FAPSERVICE2_PATH + 'FAPControl.LTE.RFTxStatus',
            is_invasive=True,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        # X_COM_RadioEnable is invasive in Single Carrier for FAPService.1
        # But for Carrier Aggregation in FAPService.2 it appears to take effect
        # immediately - and so we set it as non-invasive
        CA_RADIO_ENABLE: TrParam(
            path=FAPSERVICE2_PATH + 'CellConfig.LTE.RAN.RF.X_COM_RadioEnable',
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        CA_PLMN_CELL_RESERVED: TrParam(
            path=FAPSERVICE2_PATH + 'CellConfig.LTE.EPC.PLMNList.1.CellReservedForOperatorUse',
            is_invasive=True,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        CA_PLMN_ENABLE: TrParam(
            path=FAPSERVICE2_PATH + 'CellConfig.LTE.EPC.PLMNList.1.Enable',
            is_invasive=True,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        CA_PLMN_PRIMARY: TrParam(
            path=FAPSERVICE2_PATH + 'CellConfig.LTE.EPC.PLMNList.1.IsPrimary',
            is_invasive=True,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        CA_PLMN_PLMNID: TrParam(
            path=FAPSERVICE2_PATH + 'CellConfig.LTE.EPC.PLMNList.1.PLMNID',
            is_invasive=True,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
    }


class BaicellsQRTBTrDataModel(DataModel):
    """
    Class to represent relevant data model parameters from TR-196/TR-098/TR-181.
    This class is effectively read-only

    This is for Baicells QRTB based on software BaiBS_QRTB_2.6.2.
    Tested on hw version E01 and A01
    """
    # Parameters to query when reading eNodeB config
    LOAD_PARAMETERS = [ParameterName.DEVICE]
    # Mapping of TR parameter paths to aliases
    DEVICE_PATH = 'Device.'
    FAPSERVICE_PATH = DEVICE_PATH + 'Services.FAPService.1.'
    PARAMETERS = {
        # Top-level objects
        ParameterName.DEVICE: TrParam(
            path=DEVICE_PATH,
            is_invasive=True,
            type=TrParameterType.OBJECT,
            is_optional=False,
        ),
        ParameterName.FAP_SERVICE: TrParam(
            path=FAPSERVICE_PATH,
            is_invasive=True,
            type=TrParameterType.OBJECT,
            is_optional=False,
        ),

        # Device info parameters
        ParameterName.GPS_STATUS: TrParam(
            path=DEVICE_PATH + 'DeviceInfo.X_COM_GPS_Status',
            is_invasive=True,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        ParameterName.PTP_STATUS: TrParam(
            path=DEVICE_PATH + 'DeviceInfo.X_COM_1588_Status',
            is_invasive=True,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        ParameterName.MME_STATUS: TrParam(
            path=DEVICE_PATH + 'DeviceInfo.X_COM_MME_Status',
            is_invasive=True,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        ParameterName.REM_STATUS: TrParam(
            path=DEVICE_PATH + 'DeviceInfo.X_COM_REM_Status',
            is_invasive=True,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        ParameterName.LOCAL_GATEWAY_ENABLE: TrParam(
            path=DEVICE_PATH + 'DeviceInfo.X_COM_LTE_LGW_Switch',
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        ParameterName.GPS_ENABLE: TrParam(
            path=DEVICE_PATH + 'DeviceInfo.X_COM_GpsSyncEnable',
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        ParameterName.GPS_LAT: TrParam(
            path=DEVICE_PATH + 'FAP.GPS.LockedLatitude',
            is_invasive=True,
            type=TrParameterType.INT,
            is_optional=False,
        ),
        ParameterName.GPS_LONG: TrParam(
            path=DEVICE_PATH + 'FAP.GPS.LockedLongitude',
            is_invasive=True,
            type=TrParameterType.INT,
            is_optional=False,
        ),
        ParameterName.SW_VERSION: TrParam(
            path=DEVICE_PATH + 'DeviceInfo.SoftwareVersion',
            is_invasive=True,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        ParameterName.SERIAL_NUMBER: TrParam(
            path=DEVICE_PATH + 'DeviceInfo.SerialNumber',
            is_invasive=True,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        ParameterName.INDOOR_DEPLOYMENT: TrParam(
            path=DEVICE_PATH + 'DeviceInfo.indoorDeployment',
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        ParameterName.ANTENNA_HEIGHT: TrParam(
            path=DEVICE_PATH + 'DeviceInfo.AntennaInfo.Height',
            is_invasive=True,
            type=TrParameterType.INT,
            is_optional=False,
        ),
        ParameterName.ANTENNA_HEIGHT_TYPE: TrParam(
            path=DEVICE_PATH + 'DeviceInfo.AntennaInfo.HeightType',
            is_invasive=True,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        ParameterName.ANTENNA_GAIN: TrParam(
            path=DEVICE_PATH + 'DeviceInfo.AntennaInfo.Gain',
            is_invasive=True,
            type=TrParameterType.INT,
            is_optional=False,
        ),
        ParameterName.CBSD_CATEGORY: TrParam(
            path=DEVICE_PATH + 'DeviceInfo.cbsdCategory',
            is_invasive=True,
            type=TrParameterType.STRING,
            is_optional=False,
        ),

        # Capabilities
        ParameterName.DUPLEX_MODE_CAPABILITY: TrParam(
            path=FAPSERVICE_PATH + 'Capabilities.LTE.DuplexMode',
            is_invasive=True,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        ParameterName.BAND_CAPABILITY: TrParam(
            path=FAPSERVICE_PATH + 'Capabilities.LTE.BandsSupported',
            is_invasive=True,
            type=TrParameterType.STRING,
            is_optional=False,
        ),

        # RF-related parameters
        ParameterName.EARFCNDL: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.EARFCNDL',
            is_invasive=True,
            type=TrParameterType.UNSIGNED_INT,
            is_optional=False,
        ),
        ParameterName.EARFCNUL: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.EARFCNUL',
            is_invasive=True,
            type=TrParameterType.UNSIGNED_INT,
            is_optional=False,
        ),
        ParameterName.BAND: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.FreqBandIndicator',
            is_invasive=True,
            type=TrParameterType.UNSIGNED_INT,
            is_optional=False,
        ),
        ParameterName.PCI: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.PhyCellID',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        ParameterName.DL_BANDWIDTH: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.DLBandwidth',
            is_invasive=True,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        ParameterName.UL_BANDWIDTH: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.ULBandwidth',
            is_invasive=True,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        ParameterName.RADIO_ENABLE: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.X_COM_RadioEnable',
            is_invasive=True,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        ParameterName.SUBFRAME_ASSIGNMENT: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.PHY.TDDFrame.SubFrameAssignment',
            is_invasive=True,
            type=TrParameterType.INT,
            is_optional=False,
        ),
        ParameterName.SPECIAL_SUBFRAME_PATTERN: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.PHY.TDDFrame.SpecialSubframePatterns',
            is_invasive=True,
            type=TrParameterType.INT,
            is_optional=False,
        ),
        ParameterName.CELL_ID: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Common.CellIdentity',
            is_invasive=True,
            type=TrParameterType.UNSIGNED_INT,
            is_optional=False,
        ),
        ParameterName.POWER_SPECTRAL_DENSITY: TrParam(
            path=DEVICE_PATH + 'DeviceInfo.PowerSpectralDensity',
            is_invasive=False,
            type=TrParameterType.UNSIGNED_INT,
            is_optional=False,
        ),

        # Other LTE parameters
        ParameterName.ADMIN_STATE: TrParam(
            path=FAPSERVICE_PATH + 'FAPControl.LTE.AdminState',
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        ParameterName.OP_STATE: TrParam(
            path=FAPSERVICE_PATH + 'FAPControl.LTE.OpState',
            is_invasive=True,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        ParameterName.RF_TX_STATUS: TrParam(
            path=FAPSERVICE_PATH + 'FAPControl.LTE.RFTxStatus',
            is_invasive=True,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),

        # Core network parameters
        ParameterName.MME_IP: TrParam(
            path=FAPSERVICE_PATH + 'FAPControl.LTE.Gateway.S1SigLinkServerList',
            is_invasive=True,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        ParameterName.MME_PORT: TrParam(
            path=FAPSERVICE_PATH + 'FAPControl.LTE.Gateway.S1SigLinkPort',
            is_invasive=True,
            type=TrParameterType.INT,
            is_optional=False,
        ),
        ParameterName.NUM_PLMNS: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNListNumberOfEntries',
            is_invasive=True,
            type=TrParameterType.INT,
            is_optional=False,
        ),

        # PLMN arrays are added below
        ParameterName.PLMN: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNList.',
            is_invasive=True,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        ParameterName.TAC: TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.EPC.TAC',
            is_invasive=True,
            type=TrParameterType.INT,
            is_optional=False,
        ),
        ParameterName.IP_SEC_ENABLE: TrParam(
            path=DEVICE_PATH + 'Services.FAPService.Ipsec.IPSEC_ENABLE',
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        ParameterName.MME_POOL_ENABLE: TrParam(
            path=FAPSERVICE_PATH + 'FAPControl.LTE.Gateway.X_COM_MmePool.Enable',
            is_invasive=True,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),

        # Management server parameters
        ParameterName.PERIODIC_INFORM_ENABLE: TrParam(
            path=DEVICE_PATH + 'ManagementServer.PeriodicInformEnable',
            is_invasive=True,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        ParameterName.PERIODIC_INFORM_INTERVAL: TrParam(
            path=DEVICE_PATH + 'ManagementServer.PeriodicInformInterval',
            is_invasive=True,
            type=TrParameterType.UNSIGNED_INT,
            is_optional=False,
        ),

        # Performance management parameters
        ParameterName.PERF_MGMT_ENABLE: TrParam(
            path=DEVICE_PATH + 'FAP.PerfMgmt.Config.1.Enable',
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        ParameterName.PERF_MGMT_UPLOAD_INTERVAL: TrParam(
            path=DEVICE_PATH + 'FAP.PerfMgmt.Config.1.PeriodicUploadInterval',
            is_invasive=False,
            type=TrParameterType.INT,
            is_optional=False,
        ),
        ParameterName.PERF_MGMT_UPLOAD_URL: TrParam(
            path=DEVICE_PATH + 'FAP.PerfMgmt.Config.1.URL',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),

        # SAS parameters
        ParameterName.SAS_FCC_ID: TrParam(
            path=DEVICE_PATH + 'DeviceInfo.SAS.FccId',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        ParameterName.SAS_USER_ID: TrParam(
            path=DEVICE_PATH + 'DeviceInfo.SAS.UserId',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        ParameterName.SAS_ENABLED: TrParam(
            path=DEVICE_PATH + 'DeviceInfo.SAS.enableMode',
            is_invasive=False,
            type=TrParameterType.INT,
            is_optional=False,
        ),
        ParameterName.SAS_RADIO_ENABLE: TrParam(
            path=DEVICE_PATH + 'DeviceInfo.SAS.RadioEnable',
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
    }

    NUM_PLMNS_IN_CONFIG = 6
    for i in range(1, NUM_PLMNS_IN_CONFIG + 1):  # noqa: WPS604
        PARAMETERS[(ParameterName.PLMN_N) % i] = TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNList.%d.' % i, is_invasive=True, type=TrParameterType.STRING,
            is_optional=False,
        )
        PARAMETERS[ParameterName.PLMN_N_CELL_RESERVED % i] = TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNList.%d.CellReservedForOperatorUse' % i, is_invasive=True,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        )
        PARAMETERS[ParameterName.PLMN_N_ENABLE % i] = TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNList.%d.Enable' % i, is_invasive=True,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        )
        PARAMETERS[ParameterName.PLMN_N_PRIMARY % i] = TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNList.%d.IsPrimary' % i, is_invasive=True,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        )
        PARAMETERS[ParameterName.PLMN_N_PLMNID % i] = TrParam(
            path=FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNList.%d.PLMNID' % i, is_invasive=True,
            type=TrParameterType.STRING,
            is_optional=False,
        )

    TRANSFORMS_FOR_MAGMA = {
        # We don't set GPS, so we don't need transform for enb
        ParameterName.GPS_LAT: transform_for_magma.gps_tr181,
        ParameterName.GPS_LONG: transform_for_magma.gps_tr181,
    }
    PARAMETERS.update(CarrierAggregationParameters.CA_PARAMETERS)  # noqa: WPS604

    @classmethod
    def get_parameter(cls, param_name: ParameterName) -> Optional[TrParam]:
        """
        Retrieve parameter by its name

        Args:
            param_name (ParameterName): parameter name to retrieve

        Returns:
            Optional[TrParam]
        """
        return cls.PARAMETERS.get(param_name)

    @classmethod
    def _get_magma_transforms(
            cls,
    ) -> Dict[ParameterName, Callable[[Any], Any]]:
        return cls.TRANSFORMS_FOR_MAGMA

    @classmethod
    def _get_enb_transforms(cls) -> Dict[ParameterName, Callable[[Any], Any]]:
        return {}

    @classmethod
    def get_load_parameters(cls) -> List[ParameterName]:
        """
        Retrieve all load parameters

        Returns:
             List[ParameterName]
        """
        return cls.LOAD_PARAMETERS

    @classmethod
    def get_num_plmns(cls) -> int:
        """
        Retrieve the number of all PLMN parameters

        Returns:
            int
        """
        return cls.NUM_PLMNS_IN_CONFIG

    @classmethod
    def get_parameter_names(cls) -> List[ParameterName]:
        """
        Retrieve all parameter names

        Returns:
            List[ParameterName]
        """
        excluded_params = [
            str(ParameterName.DEVICE),
            str(ParameterName.FAP_SERVICE),
        ]
        names = list(
            filter(
                lambda x: (not str(x).startswith('PLMN')) and (str(x) not in excluded_params),
                cls.PARAMETERS.keys(),
            ),
        )
        return names

    @classmethod
    def get_numbered_param_names(cls) -> Dict[ParameterName, List[ParameterName]]:
        """
        Retrieve parameter names of all objects

        Returns:
            Dict[ParameterName, List[ParameterName]]
        """
        names = {}
        for i in range(1, cls.NUM_PLMNS_IN_CONFIG + 1):
            params = []
            params.append(ParameterName.PLMN_N_CELL_RESERVED % i)
            params.append(ParameterName.PLMN_N_ENABLE % i)
            params.append(ParameterName.PLMN_N_PRIMARY % i)
            params.append(ParameterName.PLMN_N_PLMNID % i)
            names[ParameterName.PLMN_N % i] = params
        return names


class BaicellsQRTBTrConfigurationInitializer(EnodebConfigurationPostProcessor):
    """
    Overrides desired config on the State Machine
    """

    def postprocess(self, mconfig: Any, service_cfg: Any, desired_cfg: EnodebConfiguration) -> None:
        """
        Add some params to the desired config

        Args:
            mconfig (Any): mconfig
            service_cfg (Any): service config
            desired_cfg (EnodebConfiguration): desired config
        """
        desired_cfg.set_parameter(ParameterName.SAS_ENABLED, 1)
        # Set Cell reservation for both cells
        desired_cfg.set_parameter_for_object(
            ParameterName.PLMN_N_CELL_RESERVED % 1, True,  # noqa: WPS345,WPS425
            ParameterName.PLMN_N % 1,  # noqa: WPS345
        )
        desired_cfg.set_parameter(
            CarrierAggregationParameters.CA_PLMN_CELL_RESERVED, True,
        )

        # Make sure FAPService.1. is Primary
        desired_cfg.set_parameter_for_object(
            ParameterName.PLMN_N_PRIMARY % 1, True,  # noqa: WPS345,WPS425
            ParameterName.PLMN_N % 1,  # noqa: WPS345
        )
        desired_cfg.set_parameter(
            CarrierAggregationParameters.CA_PLMN_PRIMARY, False,
        )

        # Enable both cells
        desired_cfg.set_parameter_for_object(
            ParameterName.PLMN_N_ENABLE % 1, True,  # noqa: WPS345,WPS425
            ParameterName.PLMN_N % 1,  # noqa: WPS345
        )
        desired_cfg.set_parameter(
            CarrierAggregationParameters.CA_PLMN_ENABLE, True,
        )

        parameters_to_delete = [
            ParameterName.RADIO_ENABLE, ParameterName.POWER_SPECTRAL_DENSITY,
            ParameterName.EARFCNDL, ParameterName.EARFCNUL, ParameterName.BAND,
            ParameterName.DL_BANDWIDTH, ParameterName.UL_BANDWIDTH,
            ParameterName.SAS_RADIO_ENABLE,
        ]
        for p in parameters_to_delete:
            if desired_cfg.has_parameter(p):
                desired_cfg.delete_parameter(p)


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
