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
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import CardTitleRow from '../../components/layout/CardTitleRow';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import FEGNetworkContext from '../../components/context/FEGNetworkContext';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import FEGNetworkInfo from './FEGNetworkInfo';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import FEGServicingAccessGatewayTable from './FEGServicingAccessGatewayTable';
import Grid from '@material-ui/core/Grid';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import JsonEditor from '../../components/JsonEditor';
import React from 'react';
import TopBar from '../../components/TopBar';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
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

/**
 * Returns the network page of a federation network. It consists of top
 * bar, which has a button to navigate to the json configuration and a
 * network information section.
 */
export default function NetworkDashboard() {
  const classes = useStyles();
  const navigate = useNavigate();
  const ctx = useContext(FEGNetworkContext);

  return (
    <>
      <TopBar
        header="Network"
        tabs={[
          {
            label: ctx?.state?.id || 'Network',
            to: 'network',
            icon: NetworkCheck,
            filters: (
              <Grid
                container
                justifyContent="flex-end"
                alignItems="center"
                spacing={2}>
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

/**
 * Returns a json config page which allows a user to edit the network
 * information.
 */
export function NetworkJsonConfig() {
  const params = useParams();
  const [error, setError] = useState('');
  const networkId: string = nullthrows(params.networkId);
  const enqueueSnackbar = useEnqueueSnackbar();
  const ctx = useContext(FEGNetworkContext);

  return (
    <JsonEditor
      content={ctx.state}
      error={error}
      onSave={async fegNetwork => {
        try {
          await ctx.updateNetworks({networkId, fegNetwork});
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
