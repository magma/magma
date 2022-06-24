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
 */

import AddEditNetworkButton from './NetworkEdit';
import Button from '@material-ui/core/Button';
import CardTitleRow from '../../components/layout/CardTitleRow';
import Grid from '@material-ui/core/Grid';
import JsonEditor from '../../components/JsonEditor';
import LteNetworkContext, {
  UpdateNetworkContextProps,
} from '../../components/context/LteNetworkContext';
import NetworkEpc from './NetworkEpc';
import NetworkInfo from './NetworkInfo';
import NetworkKPI from './NetworkKPIs';
import NetworkRanConfig from './NetworkRanConfig';
import React from 'react';
import TopBar from '../../components/TopBar';
import nullthrows from '../../../shared/util/nullthrows';

import {
  Navigate,
  Route,
  Routes,
  useNavigate,
  useParams,
} from 'react-router-dom';
import {NetworkCheck} from '@material-ui/icons';
import {Theme} from '@material-ui/core/styles';
import {colors, typography} from '../../theme/default';
import {getErrorMessage} from '../../util/ErrorUtils';
import {makeStyles} from '@material-ui/styles';
import {useContext, useState} from 'react';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';

const useStyles = makeStyles<Theme>(theme => ({
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
          // TODO[TS-migration] Broken LteNetworkContext type
          await ctx.updateNetworks(({
            networkId,
            lteNetwork,
          } as unknown) as UpdateNetworkContextProps);
          enqueueSnackbar('Network saved successfully', {
            variant: 'success',
          });
          setError('');
        } catch (error) {
          setError(getErrorMessage(error));
        }
      }}
    />
  );
}

export function NetworkDashboardInternal() {
  const classes = useStyles();
  const ctx = useContext(LteNetworkContext);

  // TODO[TS-migration] Broken LteNetworkContext type
  /* eslint-disable @typescript-eslint/no-non-null-asserted-optional-chain */
  const epcConfigs = ctx.state.cellular?.epc!;
  const lteRanConfigs = ctx.state.cellular?.ran!;
  const lteDnsConfig = ctx.state.dns!;
  /* eslint-enable @typescript-eslint/no-non-null-asserted-optional-chain */

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
              <NetworkRanConfig lteRanConfigs={lteRanConfigs} />
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
