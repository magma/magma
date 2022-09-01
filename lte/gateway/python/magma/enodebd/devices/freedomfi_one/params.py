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
from typing import Dict

from magma.enodebd.data_models.data_model import InvalidTrParamPath, TrParam
from magma.enodebd.data_models.data_model_parameters import (
    ParameterName,
    TrParameterType,
)
from magma.enodebd.device_config.enodeb_configuration import EnodebConfiguration
from magma.enodebd.logger import EnodebdLogger


class SASParameters(object):
    """ Class modeling the SAS parameters and their TR path"""
    # SAS parameters for FreedomFiOne
    FAP_CONTROL = 'Device.Services.FAPService.1.FAPControl.'
    FAPSERVICE_PATH = 'Device.Services.FAPService.1.'

    # Sas management parameters
    # TODO move param definitions to ParameterNames class or think of something to make them more generic across devices
    SAS_ENABLE = "sas_enabled"
    SAS_METHOD = "sas_method"  # 0 = SAS client, 1 = DP mode
    SAS_SERVER_URL = "sas_server_url"
    SAS_UID = "sas_uid"
    SAS_CATEGORY = "sas_category"
    SAS_CHANNEL_TYPE = "sas_channel_type"
    SAS_CERT_SUBJECT = "sas_cert_subject"
    SAS_IC_GROUP_ID = "sas_icg_group_id"
    SAS_LOCATION = "sas_location"
    SAS_HEIGHT_TYPE = "sas_height_type"
    SAS_CPI_ENABLE = "sas_cpi_enable"
    SAS_CPI_IPE = "sas_cpi_ipe"  # Install param supplied enable
    SAS_USER_ID = "sas_uid"  # TODO this should be set to a constant value in config
    SAS_FCC_ID = "sas_fcc_id"
    # For CBRS radios we set this to value returned by SAS, eNB can reduce the
    # power if needed.
    SAS_MAX_EIRP = "sas_max_eirp"
    SAS_ANTENNA_GAIN = "sas_antenna_gain"

    SAS_PARAMETERS = {
        SAS_ENABLE: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.Enable',
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        SAS_METHOD: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.Method',
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        SAS_SERVER_URL: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.Server',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        SAS_USER_ID: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.UserContactInformation',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        SAS_FCC_ID: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.FCCIdentificationNumber',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        SAS_CATEGORY: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.Category',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        SAS_CHANNEL_TYPE: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.ProtectionLevel',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        SAS_CERT_SUBJECT: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.CertSubject',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        # SAS_IC_GROUP_ID: TrParam(
        #     FAP_CONTROL + 'LTE.X_000E8F_SAS.ICGGroupId', is_invasive=False,
        #     type=TrParameterType.STRING, False),
        SAS_LOCATION: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.Location',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        SAS_HEIGHT_TYPE: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.HeightType',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        SAS_CPI_ENABLE: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.CPIEnable',
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        SAS_CPI_IPE: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.CPIInstallParamSuppliedEnable',
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        SAS_MAX_EIRP: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.MaxEirpMHz_Carrier1',
            is_invasive=False,
            type=TrParameterType.INT,
            is_optional=False,
        ),
        SAS_ANTENNA_GAIN: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.AntennaGain',
            is_invasive=True,
            type=TrParameterType.INT,
            is_optional=False,
        ),
    }

    # Hardcoded defaults
    defaults = {
        SAS_ENABLE: True,
        SAS_METHOD: False,
    }


class StatusParameters(object):
    """
    Stateful class that converts eNB status to Magma understood status.
    FreedomFiOne has many status params that interact with each other.
    This class maintains the business logic of parsing these status params
    and converting it to Magma understood fields.
    """
    STATUS_PATH = "Device.X_000E8F_DeviceFeature.X_000E8F_NEStatus."

    # Status parameters
    DEFAULT_GW = 'defaultGW'
    SYNC_STATUS = 'syncStatus'
    SAS_STATUS = 'sasStatus'
    ENB_STATUS = 'enbStatus'
    GPS_SCAN_STATUS = 'gpsScanStatus'

    STATUS_PARAMETERS = {
        # Status nodes
        DEFAULT_GW: TrParam(
            STATUS_PATH + 'X_000E8F_DEFGW_Status',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        SYNC_STATUS: TrParam(
            STATUS_PATH + 'X_000E8F_Sync_Status',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        SAS_STATUS: TrParam(
            STATUS_PATH + 'X_000E8F_SAS_Status',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        ENB_STATUS: TrParam(
            STATUS_PATH + 'X_000E8F_eNB_Status',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),

        # GPS status, lat, long
        GPS_SCAN_STATUS: TrParam(
            'Device.FAP.GPS.ScanStatus',
            is_invasive=False, type=TrParameterType.STRING, is_optional=False,
        ),
        ParameterName.GPS_LAT: TrParam(
            'Device.FAP.GPS.LockedLatitude',
            is_invasive=False, type=TrParameterType.STRING, is_optional=False,
        ),
        ParameterName.GPS_LONG: TrParam(
            'Device.FAP.GPS.LockedLongitude',
            is_invasive=False, type=TrParameterType.STRING, is_optional=False,
        ),
    }

    # Derived status params that don't have tr69 representation.
    DERIVED_STATUS_PARAMETERS = {
        ParameterName.RF_TX_STATUS: TrParam(
            InvalidTrParamPath,
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        ParameterName.GPS_STATUS: TrParam(
            InvalidTrParamPath,
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        ParameterName.PTP_STATUS: TrParam(
            InvalidTrParamPath,
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        ParameterName.MME_STATUS: TrParam(
            InvalidTrParamPath,
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        ParameterName.OP_STATE: TrParam(
            InvalidTrParamPath,
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
    }

    @classmethod
    def set_magma_device_cfg(
            cls, name_to_val: Dict,
            device_cfg: EnodebConfiguration,
    ):
        """
        Convert FreedomFiOne name_to_val representation to magma device_cfg
        """
        success_str = "SUCCESS"  # String constant returned by radio
        insync_str = "INSYNC"

        if (
            name_to_val.get(cls.DEFAULT_GW)
            and name_to_val[cls.DEFAULT_GW].upper() != success_str
        ):
            # Nothing will proceed if the eNB doesn't have an IP on the WAN
            serial_num = "unknown"
            if device_cfg.has_parameter(ParameterName.SERIAL_NUMBER):
                serial_num = device_cfg.get_parameter(
                    ParameterName.SERIAL_NUMBER,
                )
            EnodebdLogger.error(
                "Radio with serial number %s doesn't have IP address "
                "on WAN", serial_num,
            )
            device_cfg.set_parameter(
                param_name=ParameterName.RF_TX_STATUS,
                value=False,
            )
            device_cfg.set_parameter(
                param_name=ParameterName.GPS_STATUS,
                value=False,
            )
            device_cfg.set_parameter(
                param_name=ParameterName.PTP_STATUS,
                value=False,
            )
            device_cfg.set_parameter(
                param_name=ParameterName.MME_STATUS,
                value=False,
            )
            device_cfg.set_parameter(
                param_name=ParameterName.OP_STATE,
                value=False,
            )
            return

        if (
            name_to_val.get(cls.SAS_STATUS)
            and name_to_val[cls.SAS_STATUS].upper() == success_str
        ):
            device_cfg.set_parameter(
                param_name=ParameterName.RF_TX_STATUS,
                value=True,
            )
        else:
            # No sas grant so not transmitting. There is no explicit node for Tx
            # in FreedomFiOne
            device_cfg.set_parameter(
                param_name=ParameterName.RF_TX_STATUS,
                value=False,
            )

        if (
            name_to_val.get(cls.GPS_SCAN_STATUS)
            and name_to_val[cls.GPS_SCAN_STATUS].upper() == success_str
        ):
            device_cfg.set_parameter(
                param_name=ParameterName.GPS_STATUS,
                value=True,
            )
            # Time comes through GPS so can only be insync with GPS is
            # in sync, we use PTP_STATUS field to overload timer is in Sync.
            if (
                name_to_val.get(cls.SYNC_STATUS)
                and name_to_val[cls.SYNC_STATUS].upper() == insync_str
            ):
                device_cfg.set_parameter(
                    param_name=ParameterName.PTP_STATUS,
                    value=True,
                )
            else:
                device_cfg.set_parameter(
                    param_name=ParameterName.PTP_STATUS, value=False,
                )
        else:
            device_cfg.set_parameter(
                param_name=ParameterName.GPS_STATUS,
                value=False,
            )
            device_cfg.set_parameter(
                param_name=ParameterName.PTP_STATUS,
                value=False,
            )

        if (
            name_to_val.get(cls.DEFAULT_GW)
            and name_to_val[cls.DEFAULT_GW].upper() == success_str
        ):
            device_cfg.set_parameter(
                param_name=ParameterName.MME_STATUS,
                value=True,
            )
        else:
            device_cfg.set_parameter(
                param_name=ParameterName.MME_STATUS,
                value=False,
            )

        if (
            name_to_val.get(cls.ENB_STATUS)
            and name_to_val[cls.ENB_STATUS].upper() == success_str
        ):
            device_cfg.set_parameter(
                param_name=ParameterName.OP_STATE,
                value=True,
            )
        else:
            device_cfg.set_parameter(
                param_name=ParameterName.OP_STATE,
                value=False,
            )

        pass_through_params = [ParameterName.GPS_LAT, ParameterName.GPS_LONG]
        for name in pass_through_params:
            device_cfg.set_parameter(name, name_to_val[name])


class FreedomFiOneMiscParameters(object):
    """
    Default set of parameters that enable carrier aggregation and other
    miscellaneous properties
    """
    FAP_CONTROL = 'Device.Services.FAPService.1.FAPControl.'
    FAPSERVICE_PATH = 'Device.Services.FAPService.1.'

    # Tunnel ref format clobber it to non IPSEC as we don't support
    # IPSEC
    TUNNEL_REF = "tunnel_ref"
    PRIM_SOURCE = "prim_src"

    CONTIGUOUS_CC = "contiguous_cc"
    WEB_UI_ENABLE = "web_ui_enable"  # Enable or disable local enb UI

    MISC_PARAMETERS = {
        WEB_UI_ENABLE: TrParam(
            'Device.X_000E8F_DeviceFeature.X_000E8F_WebServerEnable',
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        CONTIGUOUS_CC: TrParam(
            FAP_CONTROL
            + 'LTE.X_000E8F_RRMConfig.X_000E8F_CELL_Freq_Contiguous',
            is_invasive=False,
            type=TrParameterType.INT,
            is_optional=False,
        ),
        TUNNEL_REF: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.Tunnel.1.TunnelRef',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        PRIM_SOURCE: TrParam(
            FAPSERVICE_PATH + 'REM.X_000E8F_tfcsManagerConfig.primSrc',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
    }

    # Hardcoded defaults
    defaults = {
        # Use IPV4 only
        TUNNEL_REF: "Device.IP.Interface.1.IPv4Address.1.",
        CONTIGUOUS_CC: 0,  # Its not contiguous carrier
        WEB_UI_ENABLE: False,  # Disable WebUI by default
    }


class CarrierAggregationParameters(object):
    """
    eNB parameters related to Carrier Aggregation
    """
    FAP_CONTROL = 'Device.Services.FAPService.1.FAPControl.'
    FAPSERVICE_PATH = 'Device.Services.FAPService.1.'

    CA_ENABLE = "CA Enable"
    CA_CARRIER_NUMBER = "CA Carrier Number"
    CA_EARFCNDL = "CA EARFCNDL"
    CA_EARFCNUL = "CA EARFCNUL"
    CA_CELL_ID = "CA Cell ID"
    CA_TAC = "CA TAC"
    CA_MAX_EIRP_MHZ = "CA Max EIRP"
    CA_BAND = "CA Band"

    CA_PARAMETERS = {
        CA_ENABLE: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_RRMConfig.X_000E8F_CA_Enable',
            is_invasive=False,
            type=TrParameterType.BOOLEAN,
            is_optional=False,
        ),
        CA_CARRIER_NUMBER: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_RRMConfig.X_000E8F_Cell_Number',
            is_invasive=False,
            type=TrParameterType.INT,
            is_optional=False,
        ),
        CA_EARFCNDL: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.X_000E8F_EARFCNDL2',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        CA_EARFCNUL: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.X_000E8F_EARFCNUL2',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        CA_CELL_ID: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.Common.X_000E8F_CellIdentity2',
            is_invasive=False,
            type=TrParameterType.STRING,
            is_optional=False,
        ),
        CA_TAC: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.EPC.X_000E8F_TAC2',
            is_invasive=False,
            type=TrParameterType.INT,
            is_optional=False,
        ),
        CA_MAX_EIRP_MHZ: TrParam(
            FAP_CONTROL + 'LTE.X_000E8F_SAS.MaxEirpMHz_Carrier2',
            is_invasive=False,
            type=TrParameterType.UNSIGNED_INT,
            is_optional=False,
        ),
        CA_BAND: TrParam(
            FAPSERVICE_PATH + 'CellConfig.LTE.RAN.RF.X_000E8F_FreqBandIndicator2',
            is_invasive=False,
            type=TrParameterType.UNSIGNED_INT,
            is_optional=False,
        ),
    }

    # By default, disable Carrier Aggregation and switch to Single Carrier mode
    defaults = {
        CA_ENABLE: False,
        CA_CARRIER_NUMBER: 1,
    }
