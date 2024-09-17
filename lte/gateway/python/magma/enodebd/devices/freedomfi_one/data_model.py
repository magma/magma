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
from typing import Any, Callable, Dict, Iterable

from magma.enodebd.data_models import transform_for_magma
from magma.enodebd.data_models.data_model import DataModel, TrParam
from magma.enodebd.data_models.data_model_parameters import (
    ParameterName,
    TrParameterType,
)
from magma.enodebd.devices.freedomfi_one.params import (
    CarrierAggregationParameters,
    FreedomFiOneMiscParameters,
    SASParameters,
    StatusParameters,
)


class FreedomFiOneTrDataModel(DataModel):
    """
    Class to represent relevant data model parameters from TR-196/TR-098.
    This class is effectively read-only.

    These models have these idiosyncrasies (on account of running TR098):

    - Parameter content root is different (InternetGatewayDevice)
    - GetParameter queries with a wildcard e.g. InternetGatewayDevice. do
      not respond with the full tree (we have to query all parameters)
    - MME status is not exposed - we assume the MME is connected if
      the eNodeB is transmitting (OpState=true)
    - Parameters such as band capability/duplex config
      are rooted under `boardconf.` and not the device config root
    - Num PLMNs is not reported by these units
    """
    # Mapping of TR parameter paths to aliases
    DEVICE_PATH = 'Device.'
    FAPSERVICE_PATH = DEVICE_PATH + 'Services.FAPService.1.'
    FAP_CONTROL = FAPSERVICE_PATH + 'FAPControl.'
    BCCH = FAPSERVICE_PATH + 'REM.LTE.Cell.1.BCCH.'

    PARAMETERS = {
        # Top-level objects
        ParameterName.DEVICE: TrParam(
            DEVICE_PATH, is_invasive=False, type=TrParameterType.OBJECT,
            is_optional=False,
        ),
        ParameterName.FAP_SERVICE: TrParam(
            FAP_CONTROL, is_invasive=False,
            type=TrParameterType.OBJECT, is_optional=False,
        ),

        # Device info
        ParameterName.SW_VERSION: TrParam(
            DEVICE_PATH + 'DeviceInfo.SoftwareVersion',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        ParameterName.SERIAL_NUMBER: TrParam(
            DEVICE_PATH + 'DeviceInfo.SerialNumber',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),

        # RF-related parameters
        ParameterName.EARFCNDL: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.EARFCNDL',
            is_invasive=False,
            type=TrParameterType.INT,
            is_optional=False,
        ),
        ParameterName.EARFCNUL: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.EARFCNUL',
            is_invasive=False,
            type=TrParameterType.INT,
            is_optional=False,
        ),
        ParameterName.DL_BANDWIDTH: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.DLBandwidth',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        ParameterName.UL_BANDWIDTH: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.ULBandwidth',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        ParameterName.PCI: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.PhyCellID',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        ParameterName.SUBFRAME_ASSIGNMENT: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.PHY.TDDFrame.SubFrameAssignment',
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        ParameterName.SPECIAL_SUBFRAME_PATTERN: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.PHY.TDDFrame.SpecialSubframePatterns',
            is_invasive=False,
            type=TrParameterType.INT,
            is_optional=False,
        ),
        ParameterName.CELL_ID: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Common.CellIdentity',
            is_invasive=False,
            type=TrParameterType.UNSIGNED_INT,
            is_optional=False,
        ),
        ParameterName.BAND: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.FreqBandIndicator',
            is_invasive=False,
            type=TrParameterType.UNSIGNED_INT,
            is_optional=False,
        ),

        # Readonly LTE state
        ParameterName.ADMIN_STATE: TrParam(
            FAP_CONTROL + 'LTE.AdminState',
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        ParameterName.GPS_ENABLE: TrParam(
            DEVICE_PATH + 'FAP.GPS.ScanOnBoot',
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),

        # Core network parameters
        ParameterName.MME_IP: TrParam(
            FAP_CONTROL + 'LTE.Gateway.S1SigLinkServerList',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        ParameterName.MME_PORT: TrParam(
            FAP_CONTROL + 'LTE.Gateway.S1SigLinkPort',
            is_invasive=False,
            type=TrParameterType.INT,
            is_optional=False,
        ),

        ParameterName.NUM_PLMNS: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNListNumberOfEntries',
            is_invasive=False,
            type=TrParameterType.INT,
            is_optional=False,
        ),

        ParameterName.TAC: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.EPC.TAC', is_invasive=False,
            type=TrParameterType.INT,
            is_optional=False,
        ),
        # Management server parameters
        ParameterName.PERIODIC_INFORM_ENABLE: TrParam(
            DEVICE_PATH + 'ManagementServer.PeriodicInformEnable',
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        ParameterName.PERIODIC_INFORM_INTERVAL: TrParam(
            DEVICE_PATH + 'ManagementServer.PeriodicInformInterval',
            is_invasive=False,
            type=TrParameterType.INT,
            is_optional=False,
        ),

        # Performance management parameters
        ParameterName.PERF_MGMT_ENABLE: TrParam(
            DEVICE_PATH + 'FAP.PerfMgmt.Config.1.Enable',
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        ParameterName.PERF_MGMT_UPLOAD_INTERVAL: TrParam(
            DEVICE_PATH + 'FAP.PerfMgmt.Config.1.PeriodicUploadInterval',
            is_invasive=False,
            type=TrParameterType.INT,
            is_optional=False,
        ),
        ParameterName.PERF_MGMT_UPLOAD_URL: TrParam(
            DEVICE_PATH + 'FAP.PerfMgmt.Config.1.URL',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
    }

    NUM_PLMNS_IN_CONFIG = 1
    for i in range(1, NUM_PLMNS_IN_CONFIG + 1):  # noqa: WPS604
        PARAMETERS[ParameterName.PLMN_N % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNList.%d.' % i,
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        )
        PARAMETERS[ParameterName.PLMN_N_CELL_RESERVED % i] = TrParam(
            FAPSERVICE_PATH
            + 'CellConfig.LTE.EPC.PLMNList.%d.CellReservedForOperatorUse' % i,
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        )
        PARAMETERS[ParameterName.PLMN_N_ENABLE % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNList.%d.Enable' % i,
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        )
        PARAMETERS[ParameterName.PLMN_N_PRIMARY % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNList.%d.IsPrimary' % i,
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        )
        PARAMETERS[ParameterName.PLMN_N_PLMNID % i] = TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.EPC.PLMNList.%d.PLMNID' % i,
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        )

    PARAMETERS.update(SASParameters.SAS_PARAMETERS)  # noqa: WPS604
    PARAMETERS.update(FreedomFiOneMiscParameters.MISC_PARAMETERS)  # noqa: WPS604
    PARAMETERS.update(StatusParameters.STATUS_PARAMETERS)  # noqa: WPS604
    # These are stateful parameters that have no tr-69 representation
    PARAMETERS.update(StatusParameters.DERIVED_STATUS_PARAMETERS)  # noqa: WPS604
    PARAMETERS.update(CarrierAggregationParameters.CA_PARAMETERS)  # noqa: WPS604

    # Parameters to query when reading eNodeB config
    LOAD_PARAMETERS = list(PARAMETERS.keys())

    TRANSFORMS_FOR_ENB: Dict[ParameterName, Callable[[Any], Any]] = {}
    TRANSFORMS_FOR_MAGMA = {
        # We don't set these parameters
        ParameterName.BAND_CAPABILITY: transform_for_magma.band_capability,
        ParameterName.DUPLEX_MODE_CAPABILITY: transform_for_magma.duplex_mode,
        ParameterName.GPS_LAT: transform_for_magma.gps_tr181,
        ParameterName.GPS_LONG: transform_for_magma.gps_tr181,
    }

    @classmethod
    def get_sas_param_names(cls) -> Iterable[ParameterName]:
        """
        Retrieve names of SAS parameters

        Returns:
            List[ParameterName]
        """
        return SASParameters.SAS_PARAMETERS.keys()
