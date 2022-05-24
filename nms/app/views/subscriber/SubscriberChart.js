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
import type {
  network_id,
  subscriber_id,
} from '../../../generated/MagmaAPIBindings';

import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import CardTitleRow from '../../components/layout/CardTitleRow';
import Divider from '@material-ui/core/Divider';

import DataGrid from '../../components/DataGrid';
import DataUsageIcon from '@material-ui/icons/DataUsage';
import Grid from '@material-ui/core/Grid';
// $FlowFixMe migrated to typescript
import LoadingFiller from '../../components/LoadingFiller';
import MagmaV1API from '../../../generated/WebClient';
import React from 'react';
import Text from '../../theme/design-system/Text';
import moment from 'moment';
// $FlowFixMe migrated to typescript
import nullthrows from '../../../shared/util/nullthrows';

import {CustomLineChart} from '../../components/CustomMetrics';
import {DateTimePicker} from '@material-ui/pickers';
import {colors} from '../../theme/default';
import {convertBitToMbit, getPromValue} from './SubscriberUtils';
import {getStep, getStepString} from '../../components/CustomMetrics';
import {makeStyles} from '@material-ui/styles';
import {useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';
import {useParams} from 'react-router-dom';

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
      q: `avg(rate(ue_reported_usage{IMSI="${subscriberId}",direction="down"}[${step}]))`,
      color: colors.secondary.dodgerBlue,
      label: 'download',
    },
    {
      q: `avg(rate(ue_reported_usage{IMSI="${subscriberId}",direction="up"}[${step}]))`,
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
            y: parseFloat(convertBitToMbit(parseFloat(i[1]))),
          });
        }),
      );

      allDatasets.push({
        datasetKeyProvider: index.toString(),
        label: query.label,
        fill: true,

        borderWidth: 1,
        backgroundColor: query.color,
        borderColor: query.color,
        hoverBackgroundColor: query.color,
        hoverBorderColor: 'black',
        data: data,
        unit: 'MB/s',
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
  const params = useParams();
  const networkId: string = nullthrows(params.networkId);
  const subscriberId: string = nullthrows(params.subscriberId);
  const [kpiRows, setKpiRows] = useState<DataRows[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const enqueueSnackbar = useEnqueueSnackbar();

  useEffect(() => {
    // fetch queries
    const fetchAllData = async () => {
      const stepCategoryMap = {
        '1h': {
          category: 'Hourly Usage MB/s',
          tooltip: 'Average Data Usage in MB/s over the past 1 hour',
        },
        '24h': {
          category: 'Daily Avg MB/s',
          tooltip: 'Average Data Usage in MB/s over the past 1 day',
        },
        '30d': {
          category: 'Monthly Avg Mb/s',
          tooltip: 'Average Data Usage in MB/s over the past 1 month',
        },
        '1y': {
          category: 'Yearly Avg Mb/s',
          tooltip: 'Average Data Usage in MB/s over the past 1 year',
        },
      };

      const queries = Object.keys(stepCategoryMap).map(async step => {
        const category = stepCategoryMap[step].category;
        const tooltip = stepCategoryMap[step].tooltip;
        try {
          const result = await MagmaV1API.getNetworksByNetworkIdPrometheusQuery(
            {
              networkId,
              query: `avg(rate(ue_reported_usage{IMSI="${subscriberId}",direction="down"}[${step}]))`,
            },
          );
          const value = convertBitToMbit(getPromValue(result));
          return {value, category, tooltip};
        } catch (e) {
          enqueueSnackbar('Error getting subscriber KPIs', {variant: 'error'});
          return {value: '-', category, tooltip};
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
  const params = useParams();
  const networkId: string = nullthrows(params.networkId);
  const subscriberId: string = nullthrows(params.subscriberId);
  const enqueueSnackbar = useEnqueueSnackbar();
  const [datasets, setDatasets] = useState<Array<Dataset>>([]);
  const [toolTipHint, setToolTipHint] = useState('');
  const [unit, setUnit] = useState('');
  const [start, setStart] = useState(moment().subtract(3, 'hours'));
  const [end, setEnd] = useState(moment());
  const [isLoading, setIsLoading] = useState(true);
  const yLabelUnit = 'MB/s';

  function Filter() {
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
      setUnit(unit);
      setIsLoading(false);
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
          title={<Text variant="body2">Data Usage {yLabelUnit}</Text>}
          subheader={
            <CustomLineChart
              dataset={datasets}
              unit={unit}
              yLabel={yLabelUnit}
              tooltipHandler={(tooltipItem, data) => {
                const val = tooltipItem.yLabel;
                return (
                  data.datasets[tooltipItem.datasetIndex].label +
                  ` ${val} ${yLabelUnit} in last ${toolTipHint}s`
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
