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
from typing import Any, Dict

from magma.common.service import MagmaService
from magma.enodebd.data_models.data_model import DataModel
from magma.enodebd.data_models.data_model_parameters import ParameterName
from magma.enodebd.device_config.enodeb_config_postprocessor import (
    EnodebConfigurationPostProcessor,
)
from magma.enodebd.device_config.enodeb_configuration import EnodebConfiguration
from magma.enodebd.devices.baicells_qrtb.data_model import (
    BaicellsQRTBTrDataModel,
)
from magma.enodebd.devices.baicells_qrtb.params import (
    CarrierAggregationParameters,
)
from magma.enodebd.devices.baicells_qrtb.states import (
    BaicellsQRTBEndSessionState,
    BaicellsQRTBNotifyDPState,
    BaicellsQRTBQueuedEventsWaitState,
    BaicellsQRTBWaitInformRebootState,
)
from magma.enodebd.devices.device_utils import EnodebDeviceName
from magma.enodebd.state_machines.enb_acs_impl import BasicEnodebAcsStateMachine
from magma.enodebd.state_machines.enb_acs_states import (
    AddObjectsState,
    CheckFirmwareUpgradeDownloadState,
    DeleteObjectsState,
    EnbSendRebootState,
    EnodebAcsState,
    ErrorState,
    FirmwareUpgradeDownloadState,
    GetObjectParametersState,
    GetParametersState,
    SendGetTransientParametersState,
    SetParameterValuesState,
    WaitEmptyMessageState,
    WaitForFirmwareUpgradeDownloadResponse,
    WaitGetObjectParametersState,
    WaitGetParametersState,
    WaitGetTransientParametersState,
    WaitInformState,
    WaitRebootResponseState,
    WaitSetParameterValuesState,
)


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
