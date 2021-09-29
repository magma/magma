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

    # HO algorithm parameters
    A1_THRESHOLD_RSRP = 'A1 threshold rsrp'
    LTE_A1_THRESHOLD_RSRQ = 'Lte a1 threshold rsrq'
    HYSTERESIS = 'Hysteresis'
    TIME_TO_TRIGGER = 'Time to trigger'
    A2_THRESHOLD_RSRP = 'A2 threshold rsrp'
    LTE_A2_THRESHOLD_RSRQ = 'Lte a2 threshold rsrq'
    LTE_A2_THRESHOLD_RSRP_IRAT_VOLTE = 'Lte a2 threshold rsrp irat volte'
    LTE_A2_THRESHOLD_RSRQ_IRAT_VOLTE = 'Lte a2 threshold rsrq irat volte'
    A3_OFFSET = 'a3 offset'
    A3_OFFSET_ANR = 'a3 offset anr'
    A4_THRESHOLD_RSRP = 'a4 threshold rsrp'
    LTE_INTRA_A5_THRESHOLD_1_RSRP = 'lte intra a5 threshold1 rsrp'
    LTE_INTRA_A5_THRESHOLD_2_RSRP = 'lte intra a5 threshold2 rsrp'
    LTE_INTER_ANR_A5_THRESHOLD_1_RSRP = 'lte inter anra5 threshold1 rsrp'
    LTE_INTER_ANR_A5_THRESHOLD_2_RSRP = 'lte inter anra5 threshold2 rsrp'
    B2_THRESHOLD1_RSRP = 'b2 threshold1 rsrp'
    B2_THRESHOLD2_RSRP = 'b2 threshold2 rsrp'
    B2_GERAN_IRAT_THRESHOLD = 'b2 geran irat threshold'
    QRXLEVMIN_SIB1 = 'qrxlevmin sib1'
    QRXLEVMINOFFSET = 'qrxlevminoffset'
    S_INTRASEARCH = 's intrasearch'
    S_NONINTRASEARCH = 's nonintrasearch'
    QRXLEVMIN_SIB3 = 'qrxlevmin sib3'
    RESELECTION_PRIORITY = 'reselection priority'
    THRESHSERVINGLOW = 'threshservinglow'
    CIPHERING_ALGORITHM = 'ciphering algorithm'
    INTEGRITY_ALGORITHM = 'integrity algorithm'


class TrParameterType():
    BOOLEAN = 'boolean'
    STRING = 'string'
    INT = 'int'
    UNSIGNED_INT = 'unsignedInt'
    OBJECT = 'object'
