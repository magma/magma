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
import type {GatewayPoolRecordsType} from '../../components/context/GatewayPoolsContext';
import type {WithAlert} from '../../components/Alert/withAlert';

import ActionTable from '../../components/ActionTable';
import CardTitleRow from '../../components/layout/CardTitleRow';
import GatewayPoolsContext from '../../components/context/GatewayPoolsContext';
import Grid from '@material-ui/core/Grid';
import Link from '@material-ui/core/Link';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import React from 'react';
import SettingsInputAntennaIcon from '@material-ui/icons/SettingsInputAntenna';
import withAlert from '../../components/Alert/withAlert';

import {GatewayPoolEditDialog} from './GatewayPoolEdit';
import {colors} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useContext, useEffect, useState} from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';
import {useNavigate, useResolvedPath} from 'react-router-dom';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(3),
    flexGrow: 1,
  },
  topBar: {
    backgroundColor: colors.primary.mirage,
    padding: '20px 40px 20px 40px',
  },
  tabBar: {
    backgroundColor: colors.primary.brightGray,
    padding: '0 0 0 20px',
  },
  tabs: {
    color: colors.primary.white,
  },
  tab: {
    fontSize: '18px',
    textTransform: 'none',
  },
  tabLabel: {
    padding: '20px 0 20px 0',
  },
  tabIconLabel: {
    verticalAlign: 'middle',
    margin: '0 5px 3px 0',
  },
}));

export default function GatewayPools() {
  const classes = useStyles();
  return (
    <div className={classes.dashboardRoot}>
      <Grid container justifyContent="space-between" spacing={3}>
        <Grid item xs={12}>
          <GatewayPoolsTable />
        </Grid>
      </Grid>
    </div>
  );
}

type GatewayPoolRowType = {
  name: string,
  id: string,
  mmeGroupId: number,
  gatewayPrimary: Array<GatewayPoolRecordsType>,
  gatewaySecondary: Array<GatewayPoolRecordsType>,
};

function GatewayPoolsTableRaw(props: WithAlert) {
  const enqueueSnackbar = useEnqueueSnackbar();
  const ctx = useContext(GatewayPoolsContext);
  const [open, setOpen] = useState(false);
  const [currRow, setCurrRow] = useState<GatewayPoolRowType>({});
  const [gwPool, setGwPool] = useState(ctx.state[currRow.id] || {});
  const resolvedPath = useResolvedPath('');
  const navigate = useNavigate();
  useEffect(() => {
    setGwPool(ctx.state[currRow.id] || {});
  }, [currRow.id, ctx.state]);
  const gwIds = (gateways: Array<GatewayPoolRecordsType>) => {
    return (
      <>
        {gateways.length > 0 ? (
          gateways.map(gw => (
            <List dense>
              <ListItem>
                <Link
                  variant="body2"
                  onClick={() => {
                    navigate(
                      resolvedPath.pathname.replace(
                        `pools`,
                        `gateway/${gw.gateway_id}`,
                      ),
                    );
                  }}>
                  {gw.gateway_id}
                </Link>
              </ListItem>
            </List>
          ))
        ) : (
          <> - </>
        )}
      </>
    );
  };
  const gatewayPoolRows: Array<GatewayPoolRowType> = ctx.state
    ? Object.keys(ctx.state).map((poolId: string) => {
        const gwPoolInfo = ctx.state[poolId];
        const gwPrimary: Array<GatewayPoolRecordsType> = gwPoolInfo.gatewayPoolRecords.filter(
          record =>
            record.mme_relative_capacity === 255 &&
            ctx.state[poolId].gatewayPool.gateway_ids.includes(
              record.gateway_id,
            ),
        );
        const gwSecondary: Array<GatewayPoolRecordsType> = gwPoolInfo.gatewayPoolRecords.filter(
          record =>
            record.mme_relative_capacity === 1 &&
            ctx.state[poolId].gatewayPool.gateway_ids.includes(
              record.gateway_id,
            ),
        );
        return {
          name: gwPoolInfo.gatewayPool.gateway_pool_name || '',
          id: poolId,
          mmeGroupId: gwPoolInfo.gatewayPool.config.mme_group_id,
          gatewayPrimary: gwPrimary,
          gatewaySecondary: gwSecondary,
        };
      })
    : [];

  return (
    <>
      <GatewayPoolEditDialog
        open={open}
        onClose={() => setOpen(false)}
        pool={gwPool}
        isAdd={true}
      />

      <CardTitleRow
        key="title"
        icon={SettingsInputAntennaIcon}
        label={`Gateway Pools (${Object.keys(ctx.state).length})`}
      />
      <ActionTable
        title=""
        data={gatewayPoolRows}
        columns={[
          {
            title: 'Name',
            field: 'name',
          },
          {
            title: 'ID',
            field: 'id',
          },
          {
            title: 'MME Group ID',
            field: 'mmeGroupId',
          },
          {
            title: 'Primary Gateway',
            field: 'gatewayPrimary',
            render: currRow => gwIds(currRow.gatewayPrimary),
          },
          {
            title: 'Secondary Gateway',
            field: 'gatewaySecondary',
            render: currRow => gwIds(currRow.gatewaySecondary),
          },
        ]}
        handleCurrRow={(row: GatewayPoolRowType) => setCurrRow(row)}
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
                .confirm(`Are you sure you want to delete ${currRow.id}?`)
                .then(async confirmed => {
                  if (!confirmed) {
                    return;
                  }

                  try {
                    await ctx.setState(currRow.id);
                  } catch (e) {
                    enqueueSnackbar(e.response?.data?.message ?? e.message, {
                      variant: 'error',
                    });
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
  );
}

const GatewayPoolsTable = withAlert(GatewayPoolsTableRaw);
