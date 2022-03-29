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


class ParameterName():
    # Top-level objects
    DEVICE = 'Device'
    FAP_SERVICE = 'FAPService'

    # Device info parameters
    GPS_STATUS = 'GPS status'
    PTP_STATUS = 'PTP status'
    MME_STATUS = 'MME status'
    REM_STATUS = 'REM status'

    LOCAL_GATEWAY_ENABLE = 'Local gateway enable'
    GPS_ENABLE = 'GPS enable'
    GPS_LAT = 'GPS lat'
    GPS_LONG = 'GPS long'
    GPS_ALTI = 'GPS alti'
    SW_VERSION = 'SW version'
    VENDOR = "VENDOR"
    MODEL_NAME = "MODEL name"
    RF_STATE = "RF state"
    UPTIME = "UPTIME"

    SERIAL_NUMBER = 'Serial number'
    CELL_ID = 'Cell ID'

    # Capabilities
    DUPLEX_MODE_CAPABILITY = 'Duplex mode capability'
    BAND_CAPABILITY = 'Band capability'

    # RF-related parameters
    EARFCNDL = 'EARFCNDL'
    EARFCNUL = 'EARFCNUL'
    BAND = 'Band'
    PCI = 'PCI'
    DL_BANDWIDTH = 'DL bandwidth'
    UL_BANDWIDTH = 'UL bandwidth'
    SUBFRAME_ASSIGNMENT = 'Subframe assignment'
    SPECIAL_SUBFRAME_PATTERN = 'Special subframe pattern'
    POWER_SPECTRAL_DENSITY = 'Power Spectral Density'
    RADIO_ENABLE = "Radio Enable"

    # Other LTE parameters
    ADMIN_STATE = 'Admin state'
    OP_STATE = 'Opstate'
    RF_TX_STATUS = 'RF TX status'

    # RAN parameters
    CELL_RESERVED = 'Cell reserved'
    CELL_BARRED = 'Cell barred'

    # Core network parameters
    MME_IP = 'MME IP'
    MME_PORT = 'MME port'
    NUM_PLMNS = 'Num PLMNs'
    PLMN = 'PLMN'
    PLMN_LIST = 'PLMN List'
    MME_POOL_1 = 'MME Pool 1'
    MME_POOL_2 = 'MME Pool 2'

    # PLMN parameters
    PLMN_N = 'PLMN %d'
    PLMN_N_CELL_RESERVED = 'PLMN %d cell reserved'
    PLMN_N_ENABLE = 'PLMN %d enable'
    PLMN_N_PRIMARY = 'PLMN %d primary'
    PLMN_N_PLMNID = 'PLMN %d PLMNID'

    # PLMN arrays are added below
    TAC = 'TAC'
    IP_SEC_ENABLE = 'IPSec enable'
    MME_POOL_ENABLE = 'MME pool enable'

    # Management server parameters
    PERIODIC_INFORM_ENABLE = 'Periodic inform enable'
    PERIODIC_INFORM_INTERVAL = 'Periodic inform interval'

    # Performance management parameters
    PERF_MGMT_ENABLE = 'Perf mgmt enable'
    PERF_MGMT_UPLOAD_INTERVAL = 'Perf mgmt upload interval'
    PERF_MGMT_UPLOAD_URL = 'Perf mgmt upload URL'
    PERF_MGMT_USER = 'Perf mgmt username'
    PERF_MGMT_PASSWORD = 'Perf mgmt password'

    SAS_ENABLED = 'SAS enabled'
    SAS_FCC_ID = 'SAS FCC ID'
    SAS_USER_ID = 'SAS User ID'
    SAS_RADIO_ENABLE = 'SAS Radio Enable'

    # Download support parameters
    DOWNLOAD_URL = 'Download file server url'
    DOWNLOAD_USER = 'Download user account'
    DOWNLOAD_PASSWORD = 'Download user password'
    DOWNLOAD_FILENAME = 'Download target file name'
    DOWNLOAD_FILESIZE = 'Download file size'
    DOWNLOAD_MD5 = 'Download md check'


class TrParameterType():
    BOOLEAN = 'boolean'
    STRING = 'string'
    INT = 'int'
    UNSIGNED_INT = 'unsignedInt'
    OBJECT = 'object'


class BaicellsParameterName():
    """
        Add the external parameter for Baicells enodeB.
    """
    NUM_LTE_NEIGHBOR_FREQ = 'nums neighbor'
    NEGIH_FREQ_LIST = 'neighbor_freq_list %d'
    NEIGHBOR_FREQ_INDEX_N = 'neighbor freq index %d'
    NEIGHBOR_FREQ_EARFCN_N = 'neighbor freq earfcn %d'
    NEIGHBOR_FREQ_Q_OFFSETRANGE_N = 'neighbor freq qoffsetrange %d'
    NEIGHBOR_FREQ_QRXLEVMINSIB5_N = 'neighbor freq qrxlevminsib5 %d'
    NEIGHBOR_FREQ_PMAX_N = 'neighbor freq pmax %d'
    NEIGHBOR_FREQ_TRESELECTIONEUTRA_N = 'neighbor freq tReselectionEutra %d'
    NEIGHBOR_FREQ_TRESELECTIONEUTRASFMEDIUM_N = 'neighbor freq tReselectionEutraSFMedium %d'
    NEIGHBOR_FREQ_RESELTHRESHHIGH_N = 'neighbor freq ReselThreshHigh %d'
    NEIGHBOR_FREQ_RESELTHRESHLOW_N = 'neighbor freq ReselThreshLow %d'
    NEIGHBOR_FREQ_RESELECTIONPRIORITY_N = 'neighbor freq ReselectionPriority %d'
    NEIGHBOR_FREQ_ENABLE_N = 'neighbor freq enable %d'

    NUM_LTE_NEIGHBOR_CELL = 'nums neighbor cell'
    NEIGHBOR_CELL_LIST_N = 'neighbor_cell_list %d'
    NEIGHBOR_CELL_PLMN_N = 'neighbor_cell_plmn %d'
    NEIGHBOR_CELL_CELL_ID_N = 'neighbor_cell_id %d'
    NEIGHBOR_CELL_EARFCN_N = 'neighbor_cell_earfcn %d'
    NEIGHBOR_CELL_PCI_N = 'neighbor_cell_pci %d'
    NEIGHBOR_CELL_QOFFSET_N = 'neighbor_cell_qoffset %d'
    NEIGHBOR_CELL_CIO_N = 'neighbor_cell_cio %d'
    NEIGHBOR_CELL_TAC_N = 'neighbor_cell_tac %d'
    NEIGHBOR_CELL_ENABLE_N = 'neighbor_cell_enable %d'

    # X2 enable disable
    X2_ENABLE_DISABLE = 'x2 enable disable'

    # Radio Power Control config parameters
    REFERENCE_SIGNAL_POWER = 'Reference Signal Power'
    POWER_CLASS = 'Power Class'
    PA = 'Pa'
    PB = 'Pb'
    # management server
    MANAGEMENT_SERVER = 'Management Server'
    MANAGEMENT_SERVER_PORT = 'Management Server Port'
    MANAGEMENT_SERVER_SSL_ENABLE = 'Management Server SSL Enable'

    # Sync
    SYNC_1588_SWITCH = '1588 sync switch'
    SYNC_1588_DOMAIN = '1588 domain'
    SYNC_1588_SYNC_MSG_INTREVAL = '1588 sync message interval'
    SYNC_1588_DELAY_REQUEST_MSG_INTERVAL = '1588 delay request msg interval'
    SYNC_1588_HOLDOVER = '1588 holdover'
    SYNC_1588_ASYMMETRY = '1588 asymmetry'
    SYNC_1588_UNICAST_ENABLE = '1588 unicast enable'
    SYNC_1588_UNICAST_SERVERIP = '1588 unicast server IP'

    # HO algorithm parameters
    HO_A1_THRESHOLD_RSRP = 'A1 threshold rsrp'
    HO_LTE_A1_THRESHOLD_RSRQ = 'Lte a1 threshold rsrq'
    HO_HYSTERESIS = 'Hysteresis'
    HO_TIME_TO_TRIGGER = 'Time to trigger'
    HO_A2_THRESHOLD_RSRP = 'A2 threshold rsrp'
    HO_LTE_A2_THRESHOLD_RSRQ = 'Lte a2 threshold rsrq'
    HO_A3_OFFSET = 'a3 offset'
    HO_A3_OFFSET_ANR = 'a3 offset anr'
    HO_A4_THRESHOLD_RSRP = 'a4 threshold rsrp'
    HO_LTE_INTRA_A5_THRESHOLD_1_RSRP = 'lte intra a5 threshold1 rsrp'
    HO_LTE_INTRA_A5_THRESHOLD_2_RSRP = 'lte intra a5 threshold2 rsrp'
    HO_LTE_INTER_ANR_A5_THRESHOLD_1_RSRP = 'lte inter anra5 threshold1 rsrp'
    HO_LTE_INTER_ANR_A5_THRESHOLD_2_RSRP = 'lte inter anra5 threshold2 rsrp'
    HO_B2_THRESHOLD1_RSRP = 'b2 threshold1 rsrp'
    HO_B2_THRESHOLD2_RSRP = 'b2 threshold2 rsrp'
    HO_B2_GERAN_IRAT_THRESHOLD = 'b2 geran irat threshold'
    HO_QRXLEVMIN_SELECTION = 'qrxlevmin selection'
    HO_QRXLEVMINOFFSET = 'qrxlevminoffset'
    HO_S_INTRASEARCH = 's intrasearch'
    HO_S_NONINTRASEARCH = 's nonintrasearch'
    HO_QRXLEVMIN_RESELECTION = 'qrxlevmin reselection'
    HO_RESELECTION_PRIORITY = 'reselection priority'
    HO_THRESHSERVINGLOW = 'threshservinglow'
    HO_CIPHERING_ALGORITHM = 'ciphering algorithm'
    HO_INTEGRITY_ALGORITHM = 'integrity algorithm'
