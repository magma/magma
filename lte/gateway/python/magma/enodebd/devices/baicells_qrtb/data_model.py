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
from typing import Any, Callable, Dict

from magma.enodebd.data_models import transform_for_magma
from magma.enodebd.data_models.data_model import DataModel, TrParam
from magma.enodebd.data_models.data_model_parameters import (
    ParameterName,
    TrParameterType,
)
from magma.enodebd.devices.baicells_qrtb.params import (
    CarrierAggregationParameters,
)


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

    TRANSFORMS_FOR_ENB: Dict[ParameterName, Callable[[Any], Any]] = {}
    TRANSFORMS_FOR_MAGMA = {
        # We don't set GPS, so we don't need transform for enb
        ParameterName.GPS_LAT: transform_for_magma.gps_tr181,
        ParameterName.GPS_LONG: transform_for_magma.gps_tr181,
    }
    PARAMETERS.update(CarrierAggregationParameters.CA_PARAMETERS)  # noqa: WPS604
