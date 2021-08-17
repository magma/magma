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
import CardTitleRow from '../../components/layout/CardTitleRow';
import Grid from '@material-ui/core/Grid';
import JsonEditor from '../../components/JsonEditor';
import LteNetworkContext from '../../components/context/LteNetworkContext';
import NetworkEpc from './NetworkEpc';
import NetworkInfo from './NetworkInfo';
import NetworkKPI from './NetworkKPIs';
import NetworkRanConfig from './NetworkRanConfig';
import React from 'react';
import TopBar from '../../components/TopBar';
import nullthrows from '@fbcnms/util/nullthrows';

import {NetworkCheck} from '@material-ui/icons';
import {Redirect, Route, Switch} from 'react-router-dom';
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useContext, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

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
  const {history, relativePath, relativeUrl} = useRouter();

  return (
    <>
      <TopBar
        header="Network"
        tabs={[
          {
            label: 'Network',
            to: '/network',
            icon: NetworkCheck,
            filters: (
              <Grid
                container
                justify="flex-end"
                alignItems="center"
                spacing={2}>
                <Grid item>
                  <AddEditNetworkButton title={'Add Network'} isLink={false} />
                </Grid>
                <Grid item>
                  <Button
                    className={classes.appBarBtn}
                    onClick={() => {
                      history.push(relativeUrl('/json'));
                    }}>
                    Edit JSON
                  </Button>
                </Grid>
              </Grid>
            ),
          },
        ]}
      />

      <Switch>
        <Route path={relativePath('/json')} component={NetworkJsonConfig} />
        <Route
          path={relativePath('/network')}
          component={NetworkDashboardInternal}
        />
        <Redirect to={relativeUrl('/network')} />
      </Switch>
    </>
  );
}

export function NetworkJsonConfig() {
  const {match} = useRouter();
  const [error, setError] = useState('');
  const networkId: string = nullthrows(match.params.networkId);
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
