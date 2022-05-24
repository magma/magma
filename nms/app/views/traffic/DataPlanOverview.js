/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @flow strict-local
 * @format
 */
import type {WithAlert} from '../../components/Alert/withAlert';

import ActionTable from '../../components/ActionTable';
import CardTitleRow from '../../components/layout/CardTitleRow';
import CellWifiIcon from '@material-ui/icons/CellWifi';
import DataPlanEditDialog from './DataPlanEdit';
import Link from '@material-ui/core/Link';
import LteNetworkContext from '../../components/context/LteNetworkContext';
import React from 'react';
import withAlert from '../../components/Alert/withAlert';

// $FlowFixMe migrated to typescript
import nullthrows from '../../../shared/util/nullthrows';
import {makeStyles} from '@material-ui/styles';
import {useContext, useState} from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';
import {useParams} from 'react-router-dom';
import type {UpdateNetworkContextProps} from '../../components/context/LteNetworkContext';

import {
  BITRATE_MULTIPLIER,
  DATA_PLAN_UNLIMITED_RATES,
} from '../../components/network/DataPlanConst';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(3),
    flexGrow: 1,
  },
}));

/**
 * Stores info for a single row of the data plan table.
 *
 * @property {string} dataPlanID
 * @property {string} maxUploadBitRate
 *    - bit rate specified in Mbps
 * @property {string} maxDownloadBitRate
 *    - bit rate specified in Mbps
 */
type DataPlanRowType = {
  dataPlanID: string,
  maxUploadBitRate: string,
  maxDownloadBitRate: string,
};

/**
 * Table displaying all data plans for the current network.
 *
 * Functionality provided also for editing/deleting data plans.
 * Functionality not provided in this table for adding data plans.
 *
 * @param {WithAlert} props
 */
function DataPlanOverview(props: WithAlert) {
  const params = useParams();
  const networkID = nullthrows(params.networkId);
  const classes = useStyles();
  const enqueueSnackbar = useEnqueueSnackbar();
  const [currRow, setCurrRow] = useState<DataPlanRowType>({});
  const [open, setOpen] = React.useState(false);
  const ctx = useContext(LteNetworkContext);
  const epcConfig = ctx.state.cellular.epc;

  const dataPlans = ctx.state.cellular.epc.sub_profiles;
  const dataPlanRows: Array<DataPlanRowType> = dataPlans
    ? Object.keys(dataPlans || {}).map((id: string) => {
        const profile = nullthrows(dataPlans)[id];
        return {
          dataPlanID: id,
          maxUploadBitRate:
            profile.max_dl_bit_rate ===
            DATA_PLAN_UNLIMITED_RATES.max_dl_bit_rate
              ? 'Unlimited'
              : profile.max_dl_bit_rate / BITRATE_MULTIPLIER + ' Mbps',
          maxDownloadBitRate:
            profile.max_ul_bit_rate ===
            DATA_PLAN_UNLIMITED_RATES.max_ul_bit_rate
              ? 'Unlimited'
              : profile.max_ul_bit_rate / BITRATE_MULTIPLIER + ' Mbps',
        };
      })
    : [];

  const onDelete = async (dataPlanId: string) => {
    const subProfiles = epcConfig.sub_profiles || {};
    delete subProfiles[dataPlanId];

    const newConfig = {
      ...epcConfig,
      sub_profiles: subProfiles,
    };
    const updateNetworkProps: UpdateNetworkContextProps = {
      networkId: networkID,
      epcConfigs: newConfig,
    };

    try {
      await ctx.updateNetworks(updateNetworkProps);
      enqueueSnackbar('Data plan deleted successfully', {
        variant: 'success',
      });
    } catch (error) {
      enqueueSnackbar('error.response?.data?.message || error', {
        variant: 'error',
      });
    }
  };

  return (
    <div className={classes.dashboardRoot}>
      <>
        <CardTitleRow key="title" icon={CellWifiIcon} label={'Data Plans'} />
        <DataPlanEditDialog
          open={open}
          onClose={() => setOpen(false)}
          dataPlanId={currRow.dataPlanID}
        />
        <ActionTable
          data={dataPlanRows}
          columns={[
            {
              title: 'Data Plan ID',
              field: 'dataPlanID',
              render: currRow => (
                <Link
                  variant="body2"
                  component="button"
                  onClick={() => {
                    setCurrRow(currRow);
                    setOpen(true);
                  }}>
                  {currRow.dataPlanID}
                </Link>
              ),
            },
            {
              title: 'Max Upload Bit Rate',
              field: 'maxUploadBitRate',
              type: 'numeric',
            },
            {
              title: 'Max Download Bit Rate',
              field: 'maxDownloadBitRate',
              type: 'numeric',
            },
          ]}
          handleCurrRow={(row: DataPlanRowType) => setCurrRow(row)}
          menuItems={[
            {
              name: 'Edit',
              handleFunc: () => {
                setOpen(true);
              },
            },
            {
              name: 'Remove',
              handleFunc: () => {
                props
                  .confirm(
                    `Are you sure you want to delete data plan ${currRow.dataPlanID}?`,
                  )
                  .then(async confirmed => {
                    if (!confirmed) {
                      return;
                    }

                    try {
                      await onDelete(currRow.dataPlanID);
                    } catch (e) {
                      enqueueSnackbar(
                        'failed deleting data plan ' + currRow.dataPlanID,
                        {
                          variant: 'error',
                        },
                      );
                    }
                  });
              },
            },
          ]}
          options={{
            actionsColumnIndex: -1,
            pageSizeOptions: [5, 10],
          }}
        />
      </>
    </div>
  );
}

export default withAlert(DataPlanOverview);
