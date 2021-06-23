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
import type {DataRows} from '../../components/DataGrid';
import type {Dataset} from '../../components/CustomMetrics';
import type {EnqueueSnackbarOptions} from 'notistack';
import type {network_id, subscriber_id} from '@fbcnms/magma-api';

import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import CardTitleRow from '../../components/layout/CardTitleRow';
import CustomHistogram from '../../components/CustomMetrics';
import Divider from '@material-ui/core/Divider';

import DataGrid from '../../components/DataGrid';
import DataUsageIcon from '@material-ui/icons/DataUsage';
import Grid from '@material-ui/core/Grid';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React from 'react';
import Text from '../../theme/design-system/Text';
import moment from 'moment';
import nullthrows from '@fbcnms/util/nullthrows';

import {DateTimePicker} from '@material-ui/pickers';
import {colors} from '../../theme/default';
import {getLabelUnit, getPromValue} from './SubscriberUtils';
import {getStep, getStepString} from '../../components/CustomMetrics';
import {makeStyles} from '@material-ui/styles';
import {useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

export type DateTimeMetricChartProps = {
  title: string,
  queries: Array<string>,
  legendLabels: Array<string>,
  unit?: string,
};

const useStyles = makeStyles(_ => ({
  dateTimeText: {
    color: colors.primary.comet,
  },
}));

type DatasetFetchProps = {
  networkId: network_id,
  subscriberId: subscriber_id,
  start: moment,
  end: moment,
  enqueueSnackbar: (
    msg: string,
    cfg: EnqueueSnackbarOptions,
  ) => ?(string | number),
};

async function getDatasets(props: DatasetFetchProps) {
  const {start, end, networkId, subscriberId} = props;
  const [delta, unit] = getStep(start, end);
  let requestError = '';
  const step = getStepString(delta, unit);
  const toolTipHint = delta + ' ' + unit;

  const queries = [
    {
      q: `sum(sum_over_time(ue_reported_usage{IMSI="${subscriberId}",direction="down"}[${step}]))`,
      color: colors.secondary.dodgerBlue,
      label: 'download',
    },
    {
      q: `sum(sum_over_time(ue_reported_usage{IMSI="${subscriberId}",direction="up"}[${step}]))`,
      color: colors.data.flamePea,
      label: 'upload',
    },
  ];
  const allDatasets = [];
  const requests = queries.map(async (query, index) => {
    try {
      const resp = await MagmaV1API.getNetworksByNetworkIdPrometheusQueryRange({
        networkId,
        start: start.toISOString(),
        end: end.toISOString(),
        step: getStepString(delta, unit),
        query: query.q,
      });

      const data = [];
      resp.data.result.forEach(it =>
        it['values']?.map(i => {
          data.push({
            t: parseInt(i[0]) * 1000,
            y: parseFloat(i[1]),
          });
        }),
      );

      allDatasets.push({
        datasetKeyProvider: index.toString(),
        label: query.label,
        fill: false,
        barPercentage: 0.7,

        borderWidth: 2,
        backgroundColor: query.color,
        borderColor: query.color,
        hoverBackgroundColor: query.color,
        hoverBorderColor: 'black',
        data: data,
        unit: 'bytes',
      });
    } catch (error) {
      requestError = error;
      return [];
    }
  });

  await Promise.all(requests);
  if (requestError) {
    props.enqueueSnackbar('Error getting event counts', {
      variant: 'error',
    });
  }
  return {allDatasets, unit, toolTipHint};
}

function SubscriberDataKPI() {
  const {match} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);
  const subscriberId: string = nullthrows(match.params.subscriberId);
  const [kpiRows, setKpiRows] = useState<DataRows[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const enqueueSnackbar = useEnqueueSnackbar();

  useEffect(() => {
    // fetch queries
    const fetchAllData = async () => {
      const stepCategoryMap = {
        '1h': 'Last 1 Hour',
        '24h': 'Last 1 Day',
        '7d': 'Last 1 Week',
        '30d': 'Last 1 Month',
      };

      const queries = Object.keys(stepCategoryMap).map(async step => {
        const category = stepCategoryMap[step];
        try {
          const result = await MagmaV1API.getNetworksByNetworkIdPrometheusQuery(
            {
              networkId,
              query: `sum(increase(ue_reported_usage{IMSI="${subscriberId}",direction="down"}[${step}]))`,
            },
          );
          const [value, unit] = getLabelUnit(getPromValue(result));
          return {value, unit, category};
        } catch (e) {
          enqueueSnackbar('Error getting subscriber KPIs', {variant: 'error'});
          return {value: '-', unit: '', category};
        }
      });

      Promise.all(queries).then(allResponses => setKpiRows([allResponses]));
      // setKpiRows(kpiRows);
      setIsLoading(false);
    };

    fetchAllData();
  }, [enqueueSnackbar, networkId, subscriberId]);

  if (isLoading) {
    return <LoadingFiller />;
  }
  return <DataGrid data={kpiRows} />;
}

export default function SubscriberChart() {
  const classes = useStyles();
  const {match} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);
  const subscriberId: string = nullthrows(match.params.subscriberId);
  const enqueueSnackbar = useEnqueueSnackbar();
  const [datasets, setDatasets] = useState<Array<Dataset>>([]);
  const [toolTipHint, setToolTipHint] = useState('');
  const [unit, setUnit] = useState('');
  const [start, setStart] = useState(moment().subtract(3, 'hours'));
  const [end, setEnd] = useState(moment());
  const [isLoading, setIsLoading] = useState(true);

  function Filter() {
    return (
      <Grid container justify="flex-end" alignItems="center" spacing={1}>
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
            maxDate={end}
            disableFuture
            value={start}
            onChange={setStart}
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
            value={end}
            onChange={setEnd}
          />
        </Grid>
      </Grid>
    );
  }
  useEffect(() => {
    // fetch queries
    const fetchAllData = async () => {
      const {allDatasets, unit, toolTipHint} = await getDatasets({
        start,
        end,
        networkId,
        subscriberId,
        enqueueSnackbar,
      });
      setDatasets(allDatasets);
      setToolTipHint(toolTipHint);
      setIsLoading(false);
      setUnit(unit);
    };

    fetchAllData();
  }, [start, end, enqueueSnackbar, networkId, subscriberId]);

  if (isLoading) {
    return <LoadingFiller />;
  }

  return (
    <>
      <CardTitleRow icon={DataUsageIcon} label={'Data Usage'} filter={Filter} />
      <Card elevation={0}>
        <CardHeader
          title={<Text variant="body2">Data Usage Pattern</Text>}
          subheader={
            <CustomHistogram
              dataset={datasets}
              unit={unit}
              yLabel={'Bytes'}
              tooltipHandler={(tooltipItem, data) => {
                const [val, units] = getLabelUnit(tooltipItem.yLabel);
                return (
                  data.datasets[tooltipItem.datasetIndex].label +
                  ` ${val}${units} in last ${toolTipHint}s`
                );
              }}
            />
          }
        />
      </Card>
      <Divider />
      <SubscriberDataKPI />
    </>
  );
}
