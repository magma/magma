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
    SW_VERSION = 'SW version'

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


class BaicellsParameterName(object):
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


class TrParameterType():
    BOOLEAN = 'boolean'
    STRING = 'string'
    INT = 'int'
    UNSIGNED_INT = 'unsignedInt'
    OBJECT = 'object'
