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
import DashboardIcon from '@material-ui/icons/Dashboard';
import DateTimeMetricChart from '../../components/DateTimeMetricChart';
import EnodebConfig from './EnodebDetailConfig';
import EnodebContext from '../../components/context/EnodebContext';
import GatewayLogs from './GatewayLogs';
import GraphicEqIcon from '@material-ui/icons/GraphicEq';
import Grid from '@material-ui/core/Grid';
import React from 'react';
import SettingsIcon from '@material-ui/icons/Settings';
import SettingsInputAntennaIcon from '@material-ui/icons/SettingsInputAntenna';
import TopBar from '../../components/TopBar';
import nullthrows from '@fbcnms/util/nullthrows';

import {EnodebJsonConfig} from './EnodebDetailConfig';
import {EnodebStatus, EnodebSummary} from './EnodebDetailSummaryStatus';
import {Redirect, Route, Switch} from 'react-router-dom';
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useContext} from 'react';
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
  appBarBtnSecondary: {
    color: colors.primary.white,
  },
}));
const CHART_TITLE = 'Bandwidth Usage';

export function EnodebDetail() {
  const ctx = useContext(EnodebContext);
  const classes = useStyles();
  const {relativePath, relativeUrl, match} = useRouter();
  const enodebSerial: string = nullthrows(match.params.enodebSerial);
  const enbInfo = ctx.state.enbInfo[enodebSerial];

  return (
    <>
      <TopBar
        header={`Equipment/${enbInfo.enb.name}`}
        tabs={[
          {
            label: 'Overview',
            to: '/overview',
            icon: DashboardIcon,
            filters: (
              <Grid
                container
                justify="flex-end"
                alignItems="center"
                spacing={2}>
                <Grid item>
                  <Button className={classes.appBarBtnSecondary}>
                    Secondary Action
                  </Button>
                </Grid>
                <Grid item>
                  <Button className={classes.appBarBtn} variant="contained">
                    Reboot
                  </Button>
                </Grid>
              </Grid>
            ),
          },
          {
            label: 'Config',
            to: '/config',
            icon: SettingsIcon,
            filters: (
              <Grid
                container
                justify="flex-end"
                alignItems="center"
                spacing={2}>
                <Grid item>
                  <Button className={classes.appBarBtnSecondary}>
                    Secondary Action
                  </Button>
                </Grid>
                <Grid item>
                  <Button className={classes.appBarBtn} variant="contained">
                    Reboot
                  </Button>
                </Grid>
              </Grid>
            ),
          },
        ]}
      />

      <Switch>
        <Route path={relativePath('/overview')} component={Overview} />
        <Route
          path={relativePath('/config/json')}
          component={EnodebJsonConfig}
        />
        <Route path={relativePath('/config')} component={EnodebConfig} />
        <Route path={relativePath('/logs')} component={GatewayLogs} />
        <Redirect to={relativeUrl('/overview')} />
      </Switch>
    </>
  );
}

function Overview() {
  const ctx = useContext(EnodebContext);
  const classes = useStyles();
  const {match} = useRouter();
  const enodebSerial: string = nullthrows(match.params.enodebSerial);
  const enbInfo = ctx.state.enbInfo[enodebSerial];
  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={4}>
        <Grid item xs={12}>
          <Grid container spacing={4}>
            <Grid item xs={12} md={6} alignItems="center">
              <CardTitleRow
                icon={SettingsInputAntennaIcon}
                label={enbInfo.enb.name}
              />
              <EnodebSummary />
            </Grid>

            <Grid item xs={12} md={6} alignItems="center">
              <CardTitleRow icon={GraphicEqIcon} label="Status" />
              <EnodebStatus />
            </Grid>
          </Grid>
        </Grid>
        <Grid item xs={12}>
          <DateTimeMetricChart
            title={CHART_TITLE}
            unit={'Throughput(mb/s)'}
            queries={[
              `sum(pdcp_user_plane_bytes_dl{service="enodebd", enodeb="${enodebSerial}"})/1000`,
              `sum(pdcp_user_plane_bytes_ul{service="enodebd", enodeb="${enodebSerial}"})/1000`,
            ]}
            legendLabels={['Download', 'Upload']}
          />
        </Grid>
      </Grid>
    </div>
  );
}

export default EnodebDetail;
