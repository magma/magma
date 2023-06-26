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

import AutorefreshCheckbox from '../../components/AutorefreshCheckbox';
import Button from '@mui/material/Button';
import CardTitleRow from '../../components/layout/CardTitleRow';
import DashboardIcon from '@mui/icons-material/Dashboard';
import DataUsageIcon from '@mui/icons-material/DataUsage';
import DateTimeMetricChart from '../../components/DateTimeMetricChart';
import EnodebConfig from './EnodebDetailConfig';
import EnodebContext from '../../context/EnodebContext';
import GatewayLogs from './GatewayLogs';
import GraphicEqIcon from '@mui/icons-material/GraphicEq';
import Grid from '@mui/material/Grid';
import React from 'react';
import SettingsIcon from '@mui/icons-material/Settings';
import Text from '../../theme/design-system/Text';
import TextField from '@mui/material/TextField';
import TopBar from '../../components/TopBar';
import nullthrows from '../../../shared/util/nullthrows';
import withAlert from '../../components/Alert/withAlert';
import {DateTimePicker} from '@mui/x-date-pickers/DateTimePicker';
import {EnodebJsonConfig} from './EnodebDetailConfig';
import {EnodebStatus, EnodebSummary} from './EnodebDetailSummaryStatus';
import {Navigate, Route, Routes, useParams} from 'react-router-dom';
import {RunGatewayCommands} from './RunGatewayCommands';
import {Theme} from '@mui/material/styles';
import {colors, typography} from '../../theme/default';
import {getErrorMessage} from '../../util/ErrorUtils';
import {makeStyles} from '@mui/styles';
import {subHours} from 'date-fns';
import {useContext, useState} from 'react';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';
import type {WithAlert} from '../../components/Alert/withAlert';

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

    void props
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
          enqueueSnackbar(getErrorMessage(e), {
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
  const [startDate, setStartDate] = useState(subHours(new Date(), 3));
  const [endDate, setEndDate] = useState(new Date());
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
            renderInput={props => <TextField {...props} />}
            maxDate={endDate}
            disableFuture
            value={startDate}
            onChange={date => setStartDate(date!)}
          />
        </Grid>
        <Grid item>
          <Text variant="body3" className={classes.dateTimeText}>
            to
          </Text>
        </Grid>
        <Grid item>
          <DateTimePicker
            renderInput={props => <TextField {...props} />}
            disableFuture
            value={endDate}
            onChange={date => setEndDate(date!)}
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
  startDate: Date;
  endDate: Date;
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
