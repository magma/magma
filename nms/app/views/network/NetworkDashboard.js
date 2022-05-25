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

import AddEditNetworkButton from './NetworkEdit';
import Button from '@material-ui/core/Button';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import CardTitleRow from '../../components/layout/CardTitleRow';
import Grid from '@material-ui/core/Grid';
// $FlowFixMe migrated to typescript
import JsonEditor from '../../components/JsonEditor';
// $FlowFixMe migrated to typescript
import LteNetworkContext from '../../components/context/LteNetworkContext';
import NetworkEpc from './NetworkEpc';
import NetworkInfo from './NetworkInfo';
import NetworkKPI from './NetworkKPIs';
import NetworkRanConfig from './NetworkRanConfig';
import React from 'react';
import TopBar from '../../components/TopBar';
// $FlowFixMe migrated to typescript
import nullthrows from '../../../shared/util/nullthrows';

import {
  Navigate,
  Route,
  Routes,
  useNavigate,
  useParams,
} from 'react-router-dom';
import {NetworkCheck} from '@material-ui/icons';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useContext, useState} from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(5),
  },
  appBarBtn: {
    color: colors.primary.white,
    background: colors.primary.comet,
    fontFamily: typography.button.fontFamily,
    fontWeight: typography.button.fontWeight,
    fontSize: typography.button.fontSize,
    lineHeight: typography.button.lineHeight,
    letterSpacing: typography.button.letterSpacing,

    '&:hover': {
      background: colors.primary.mirage,
    },
  },
}));

export default function NetworkDashboard() {
  const classes = useStyles();
  const navigate = useNavigate();

  return (
    <>
      <TopBar
        header="Network"
        tabs={[
          {
            label: 'Network',
            to: 'network',
            icon: NetworkCheck,
            filters: (
              <Grid
                container
                justifyContent="flex-end"
                alignItems="center"
                spacing={2}>
                <Grid item>
                  <AddEditNetworkButton title={'Add Network'} isLink={false} />
                </Grid>
                <Grid item>
                  <Button
                    className={classes.appBarBtn}
                    onClick={() => navigate('json')}>
                    Edit JSON
                  </Button>
                </Grid>
              </Grid>
            ),
          },
        ]}
      />

      <Routes>
        <Route path="/json" element={<NetworkJsonConfig />} />
        <Route path="/network" element={<NetworkDashboardInternal />} />
        <Route index element={<Navigate to="network" replace />} />
      </Routes>
    </>
  );
}

export function NetworkJsonConfig() {
  const params = useParams();
  const [error, setError] = useState('');
  const networkId: string = nullthrows(params.networkId);
  const enqueueSnackbar = useEnqueueSnackbar();
  const ctx = useContext(LteNetworkContext);

  return (
    <JsonEditor
      content={ctx.state}
      error={error}
      onSave={async lteNetwork => {
        try {
          ctx.updateNetworks({networkId, lteNetwork});
          enqueueSnackbar('Network saved successfully', {
            variant: 'success',
          });
          setError('');
        } catch (e) {
          setError(e.response?.data?.message ?? e.message);
        }
      }}
    />
  );
}

export function NetworkDashboardInternal() {
  const classes = useStyles();
  const ctx = useContext(LteNetworkContext);

  const epcConfigs = ctx.state.cellular?.epc;
  const lteRanConfigs = ctx.state.cellular?.ran;
  const lteDnsConfig = ctx.state?.dns;

  function editNetwork() {
    return (
      <AddEditNetworkButton
        title={'Edit'}
        isLink={true}
        editProps={{
          editTable: 'info',
        }}
      />
    );
  }

  function editRAN() {
    return (
      <AddEditNetworkButton
        title={'Edit'}
        isLink={true}
        editProps={{
          editTable: 'ran',
        }}
      />
    );
  }

  function editEPC() {
    return (
      <AddEditNetworkButton
        title={'Edit'}
        isLink={true}
        editProps={{
          editTable: 'epc',
        }}
      />
    );
  }
  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={4}>
        <Grid item xs={12}>
          <CardTitleRow label="Overview" />
          <NetworkKPI />
        </Grid>
        <Grid item xs={12} md={6}>
          <Grid container spacing={4}>
            <Grid item xs={12}>
              <CardTitleRow label="Network" filter={editNetwork} />
              <NetworkInfo lteNetwork={ctx.state} />
            </Grid>
            <Grid item xs={12}>
              <CardTitleRow label="RAN" filter={editRAN} />
              <NetworkRanConfig
                lteDnsConfig={lteDnsConfig}
                lteRanConfigs={lteRanConfigs}
              />
            </Grid>
          </Grid>
        </Grid>
        <Grid item xs={12} md={6}>
          <Grid container spacing={4}>
            <Grid item xs={12}>
              <CardTitleRow label="EPC" filter={editEPC} />
              <NetworkEpc epcConfigs={epcConfigs} />
            </Grid>
          </Grid>
        </Grid>
      </Grid>
    </div>
  );
}
