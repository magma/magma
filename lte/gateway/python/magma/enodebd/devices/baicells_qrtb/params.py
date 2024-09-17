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
from magma.enodebd.data_models.data_model import TrParam
from magma.enodebd.data_models.data_model_parameters import TrParameterType


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
