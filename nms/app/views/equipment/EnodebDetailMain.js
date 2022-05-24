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

// $FlowFixMe migrated to typescript
import AutorefreshCheckbox from '../../components/AutorefreshCheckbox';
import Button from '@material-ui/core/Button';
import CardTitleRow from '../../components/layout/CardTitleRow';
import DashboardIcon from '@material-ui/icons/Dashboard';
import DataUsageIcon from '@material-ui/icons/DataUsage';
import DateTimeMetricChart from '../../components/DateTimeMetricChart';
import EnodebConfig from './EnodebDetailConfig';
import EnodebContext from '../../components/context/EnodebContext';
import GatewayLogs from './GatewayLogs';
import GraphicEqIcon from '@material-ui/icons/GraphicEq';
import Grid from '@material-ui/core/Grid';
import React from 'react';
import SettingsIcon from '@material-ui/icons/Settings';
import Text from '../../theme/design-system/Text';
import TopBar from '../../components/TopBar';
import moment from 'moment';
// $FlowFixMe migrated to typescript
import nullthrows from '../../../shared/util/nullthrows';
import withAlert from '../../components/Alert/withAlert';

import {DateTimePicker} from '@material-ui/pickers';
import {EnodebJsonConfig} from './EnodebDetailConfig';
import {EnodebStatus, EnodebSummary} from './EnodebDetailSummaryStatus';
import {Navigate, Route, Routes, useParams} from 'react-router-dom';
import {RunGatewayCommands} from '../../state/lte/EquipmentState';
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useContext, useState} from 'react';
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
  dateTimeText: {
    color: colors.primary.comet,
  },
}));
const CHART_TITLE = 'Bandwidth Usage';

export function EnodebDetail() {
  const params = useParams();
  const enodebSerial: string = nullthrows(params.enodebSerial);

  return (
    <>
      <TopBar
        header={`Equipment/${enodebSerial}`}
        tabs={[
          {
            label: 'Overview',
            to: 'overview',
            icon: DashboardIcon,
            filters: <EnodebRebootButton />,
          },
          {
            label: 'Config',
            to: 'config',
            icon: SettingsIcon,
            filters: <EnodebRebootButton />,
          },
        ]}
      />

      <Routes>
        <Route path="/overview" element={<Overview />} />
        <Route path="/config/json" element={<EnodebJsonConfig />} />
        <Route path="/config" element={<EnodebConfig />} />
        <Route path="/logs" element={<GatewayLogs />} />
        <Route index element={<Navigate to="overview" replace />} />
      </Routes>
    </>
  );
}

function EnodebRebootButtonInternal(props: WithAlert) {
  const classes = useStyles();
  const ctx = useContext(EnodebContext);
  const params = useParams();
  const networkId: string = nullthrows(params.networkId);
  const enodebSerial: string = nullthrows(params.enodebSerial);
  const enbInfo = ctx.state.enbInfo[enodebSerial];
  const gatewayId = enbInfo?.enb_state?.reporting_gateway_id;
  const enqueueSnackbar = useEnqueueSnackbar();

  const handleClick = () => {
    if (gatewayId == null) {
      enqueueSnackbar('Unable to trigger reboot, reporting gateway not found', {
        variant: 'error',
      });
      return;
    }

    props
      .confirm(`Are you sure you want to reboot ${enodebSerial}?`)
      .then(async confirmed => {
        if (!confirmed) {
          return;
        }
        const params = {
          command: 'reboot_enodeb',
          params: {shell_params: {[enodebSerial]: {}}},
        };

        try {
          await RunGatewayCommands({
            networkId,
            gatewayId,
            command: 'generic',
            params,
          });
          enqueueSnackbar('eNodeB reboot triggered successfully', {
            variant: 'success',
          });
        } catch (e) {
          enqueueSnackbar(e.response?.data?.message ?? e.message, {
            variant: 'error',
          });
        }
      });
  };

  return (
    <Button
      variant="contained"
      className={classes.appBarBtn}
      onClick={handleClick}>
      Reboot
    </Button>
  );
}
const EnodebRebootButton = withAlert(EnodebRebootButtonInternal);

function Overview() {
  const classes = useStyles();
  const [startDate, setStartDate] = useState(moment().subtract(3, 'hours'));
  const [endDate, setEndDate] = useState(moment());
  const [refresh, setRefresh] = useState(true);

  function MetricChartFilter() {
    return (
      <Grid container justifyContent="flex-end" alignItems="center" spacing={1}>
        <Grid item>
          <Text variant="body3" className={classes.dateTimeText}>
            Filter By Date
          </Text>
        </Grid>
        <Grid item>
          <DateTimePicker
            autoOk
            variant="outlined"
            inputVariant="outlined"
            maxDate={endDate}
            disableFuture
            value={startDate}
            onChange={setStartDate}
          />
        </Grid>
        <Grid item>
          <Text variant="body3" className={classes.dateTimeText}>
            to
          </Text>
        </Grid>
        <Grid item>
          <DateTimePicker
            autoOk
            variant="outlined"
            inputVariant="outlined"
            disableFuture
            value={endDate}
            onChange={setEndDate}
          />
        </Grid>
      </Grid>
    );
  }

  function refreshFilter() {
    return (
      <AutorefreshCheckbox
        autorefreshEnabled={refresh}
        onToggle={() => setRefresh(current => !current)}
      />
    );
  }
  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={4}>
        <Grid item xs={12}>
          <Grid container spacing={4}>
            <Grid item xs={12} md={6} alignItems="center">
              <EnodebSummary />
            </Grid>

            <Grid item xs={12} md={6} alignItems="center">
              <CardTitleRow
                icon={GraphicEqIcon}
                label="Status"
                filter={() => refreshFilter()}
              />
              <EnodebStatus refresh={refresh} />
            </Grid>
          </Grid>
        </Grid>
        <Grid item xs={12}>
          <CardTitleRow
            icon={DataUsageIcon}
            label={CHART_TITLE}
            filter={MetricChartFilter}
          />
          <EnodebMetricChart startDate={startDate} endDate={endDate} />
        </Grid>
      </Grid>
    </div>
  );
}

type Props = {
  startDate: moment$Moment,
  endDate: moment$Moment,
};

function EnodebMetricChart(props: Props) {
  const ctx = useContext(EnodebContext);
  const params = useParams();
  const enodebSerial: string = nullthrows(params.enodebSerial);
  const enbInfo = ctx.state.enbInfo[enodebSerial];
  const enbIpAddress = enbInfo?.enb_state?.ip_address ?? '';

  return (
    <DateTimeMetricChart
      title={CHART_TITLE}
      unit={'Throughput(mb/s)'}
      queries={[
        `rate(gtp_port_user_plane_dl_bytes{service="pipelined", ip_addr="${enbIpAddress}"}[5m])/1000`,
        `rate(gtp_port_user_plane_ul_bytes{service="pipelined", ip_addr="${enbIpAddress}"}[5m])/1000`,
      ]}
      legendLabels={['Download', 'Upload']}
      startDate={props.startDate}
      endDate={props.endDate}
    />
  );
}

export default EnodebDetail;
