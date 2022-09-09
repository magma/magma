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

import Card from '@mui/material/Card';
import CardHeader from '@mui/material/CardHeader';
import CardTitleRow from '../../components/layout/CardTitleRow';
import DataGrid from '../../components/DataGrid';
import DataUsageIcon from '@mui/icons-material/DataUsage';
import Divider from '@mui/material/Divider';
import Grid from '@mui/material/Grid';
import LoadingFiller from '../../components/LoadingFiller';
import MagmaAPI from '../../api/MagmaAPI';
import React from 'react';
import Text from '../../theme/design-system/Text';
import TextField from '@mui/material/TextField';
import nullthrows from '../../../shared/util/nullthrows';
import {CustomLineChart, DatasetType} from '../../components/CustomMetrics';
import {DateTimePicker} from '@mui/x-date-pickers/DateTimePicker';
import {TimeUnit} from 'chart.js';
import {colors} from '../../theme/default';
import {convertBitToMbit, getPromValue} from './SubscriberUtils';
import {getStep, getStepString} from '../../components/CustomMetrics';
import {makeStyles} from '@mui/styles';
import {subHours} from 'date-fns';
import {useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';
import {useParams} from 'react-router-dom';
import type {DataRows} from '../../components/DataGrid';
import type {Dataset} from '../../components/CustomMetrics';
import type {NetworkId, SubscriberId} from '../../../shared/types/network';
import type {OptionsObject} from 'notistack';

const useStyles = makeStyles({
  dateTimeText: {
    color: colors.primary.comet,
  },
});

type DatasetFetchProps = {
  networkId: NetworkId;
  subscriberId: SubscriberId;
  start: Date;
  end: Date;
  enqueueSnackbar: (
    msg: string,
    cfg: OptionsObject,
  ) => string | number | null | undefined;
};

type DatasetProps = Dataset & {
  datasetKeyProvider: string;
  unit: string;
};

async function getDatasets(props: DatasetFetchProps) {
  const {start, end, networkId, subscriberId} = props;
  const [delta, unit] = getStep(start, end);
  let requestError = false;
  const step = getStepString(delta, unit);
  const toolTipHint = `${delta} ${unit}`;
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
  const allDatasets: Array<DatasetProps> = [];
  const requests = queries.map(async (query, index) => {
    try {
      const resp = (
        await MagmaAPI.metrics.networksNetworkIdPrometheusQueryRangeGet({
          networkId,
          start: start.toISOString(),
          end: end.toISOString(),
          step: getStepString(delta, unit),
          query: query.q,
        })
      ).data;
      const selectedData: Array<DatasetType> = [];
      resp.data.result.forEach(it =>
        it['values']?.map(i => {
          selectedData.push({
            x: parseInt(i[0]) * 1000,
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
        data: selectedData,
        unit: 'MB/s',
      });
    } catch (error) {
      requestError = !!error;
      return [];
    }
  });

  await Promise.all(requests);

  if (requestError) {
    props.enqueueSnackbar('Error getting event counts', {
      variant: 'error',
    });
  }

  return {
    allDatasets,
    unit,
    toolTipHint,
  };
}

function SubscriberDataKPI() {
  const params = useParams();
  const networkId: string = nullthrows(params.networkId);
  const subscriberId: string = nullthrows(params.subscriberId);
  const [kpiRows, setKpiRows] = useState<Array<DataRows>>([]);
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

      const queries = (Object.keys(stepCategoryMap) as Array<
        keyof typeof stepCategoryMap
      >).map(step => {
        const category = stepCategoryMap[step].category;
        const tooltip = stepCategoryMap[step].tooltip;

        return MagmaAPI.metrics
          .networksNetworkIdPrometheusQueryGet({
            networkId,
            query: `avg(rate(ue_reported_usage{IMSI="${subscriberId}",direction="down"}[${step}]))`,
          })
          .then(({data}) => {
            const value = convertBitToMbit(getPromValue(data));
            return {
              value,
              category,
              tooltip,
            };
          })
          .catch(() => {
            enqueueSnackbar('Error getting subscriber KPIs', {
              variant: 'error',
            });
            return {
              value: '-',
              category,
              tooltip,
            };
          });
      });

      await Promise.all(queries).then(allResponses =>
        setKpiRows([allResponses]),
      );
      setIsLoading(false);
    };

    void fetchAllData();
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
  const [unit, setUnit] = useState('' as TimeUnit);
  const [startDate, setStartDate] = useState(subHours(new Date(), 3));
  const [endDate, setEndDate] = useState(new Date());
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

  useEffect(() => {
    // fetch queries
    const fetchAllData = async () => {
      const {allDatasets, unit, toolTipHint} = await getDatasets({
        start: startDate,
        end: endDate,
        networkId,
        subscriberId,
        enqueueSnackbar,
      });
      setDatasets(allDatasets);
      setToolTipHint(toolTipHint);
      setUnit(unit);
      setIsLoading(false);
    };

    void fetchAllData();
  }, [startDate, endDate, enqueueSnackbar, networkId, subscriberId]);

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
              tooltipHandler={tooltipItem => {
                const val = tooltipItem.label;
                return `${tooltipItem.dataset
                  .label!} ${val} ${yLabelUnit} in last ${toolTipHint}s`;
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
