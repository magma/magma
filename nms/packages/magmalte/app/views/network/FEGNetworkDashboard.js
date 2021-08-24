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

import Button from '@material-ui/core/Button';
import CardTitleRow from '../../components/layout/CardTitleRow';
import FEGNetworkContext from '../../components/context/FEGNetworkContext';
import FEGNetworkInfo from './FEGNetworkInfo';
import FEGServicingAccessGatewayTable from './FEGServicingAccessGatewayTable';
import Grid from '@material-ui/core/Grid';
import JsonEditor from '../../components/JsonEditor';
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

/**
 * Returns the network page of a federation network. It consists of top
 * bar, which has a button to navigate to the json configuration and a
 * network information section.
 */
export default function NetworkDashboard() {
  const classes = useStyles();
  const {history, relativePath, relativeUrl} = useRouter();
  const ctx = useContext(FEGNetworkContext);

  return (
    <>
      <TopBar
        header="Network"
        tabs={[
          {
            label: ctx?.state?.id || 'Network',
            to: '/network',
            icon: NetworkCheck,
            filters: (
              <Grid
                container
                justify="flex-end"
                alignItems="center"
                spacing={2}>
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

/**
 * Returns a json config page which allows a user to edit the network
 * information.
 */
export function NetworkJsonConfig() {
  const {match} = useRouter();
  const [error, setError] = useState('');
  const networkId: string = nullthrows(match.params.networkId);
  const enqueueSnackbar = useEnqueueSnackbar();
  const ctx = useContext(FEGNetworkContext);

  return (
    <JsonEditor
      content={ctx.state}
      error={error}
      onSave={async fegNetwork => {
        try {
          ctx.updateNetworks({networkId, fegNetwork});
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

/**
 * Returns information about the federation network and a table of the servicing
 * access gateways alongside the serviced networks they are under.
 */
export function NetworkDashboardInternal() {
  const classes = useStyles();
  const ctx = useContext(FEGNetworkContext);

  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={4}>
        <Grid item xs={12} md={6}>
          <Grid item xs={12}>
            <CardTitleRow label="Network" />
            <FEGNetworkInfo fegNetwork={ctx.state} />
          </Grid>
        </Grid>
        <Grid item xs={12} md={6}>
          <Grid container spacing={4}>
            <Grid item xs={12}>
              <CardTitleRow label="Servicing Access Gateways" />
              <FEGServicingAccessGatewayTable />
            </Grid>
          </Grid>
        </Grid>
      </Grid>
    </div>
  );
}
