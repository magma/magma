/**
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

import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import CardTitleRow from './layout/CardTitleRow';
import DataUsageIcon from '@material-ui/icons/DataUsage';
import LoadingFiller from './LoadingFiller';
import MagmaAPI from '../../api/MagmaAPI';
import React from 'react';
import moment from 'moment';
import nullthrows from '../../shared/util/nullthrows';
import {
  CustomLineChart,
  getQueryRanges,
  getStep,
  getStepString,
} from './CustomMetrics';
import {TimeUnit} from 'chart.js';
import {colors} from '../theme/default';
import {useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '../hooks/useSnackbar';
import {useParams} from 'react-router-dom';
import type {Dataset} from './CustomMetrics';
import type {NetworkId} from '../../shared/types/network';
import type {OptionsObject} from 'notistack';

type Props = {
  startEnd: [moment.Moment, moment.Moment];
};

type DatasetFetchProps = {
  networkId: NetworkId;
  start: moment.Moment;
  end: moment.Moment;
  delta: number;
  unit: TimeUnit;
  enqueueSnackbar: (
    msg: string,
    cfg: OptionsObject,
  ) => (string | number) | null | undefined;
};

async function getEventAlertDataset(props: DatasetFetchProps) {
  const {networkId, start, end, delta, unit} = props;
  let requestError = false;
  const queries = getQueryRanges(start, end, delta, unit);
  const requests = queries.map(async query => {
    try {
      const [s, e] = query;
      const response = (
        await MagmaAPI.events.eventsNetworkIdAboutCountGet({
          networkId: networkId,
          start: s.toISOString(),
          end: e.toISOString(),
        })
      ).data;
      return response;
    } catch (error) {
      requestError = !!error;
    }

    return null;
  });

  // get events data
  const eventData = await Promise.all(requests)
    .then(allResponses => {
      return allResponses.map((r, index) => {
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        const [_, e] = queries[index];

        if (r === null || r === undefined) {
          return {
            t: e.unix() * 1000,
            y: 0,
          };
        }

        return {
          t: e.unix() * 1000,
          y: r,
        };
      });
    })
    .catch(error => {
      requestError = !!error;
      return [];
    });

  const alertsData: Array<{t: number; y: number}> = [];

  try {
    const alertPromResp = (
      await MagmaAPI.metrics.networksNetworkIdPrometheusQueryRangeGet({
        networkId: networkId,
        start: start.toISOString(),
        end: end.toISOString(),
        step: getStepString(delta, unit),
        query: 'sum(ALERTS)',
      })
    ).data;
    alertPromResp.data?.result.forEach(it =>
      it['values']?.map(i => {
        alertsData.push({
          t: parseInt(i[0]) * 1000,
          y: parseFloat(i[1]),
        });
      }),
    );
  } catch (error) {
    requestError = !!error;
  }

  if (requestError) {
    props.enqueueSnackbar('Error getting event counts', {
      variant: 'error',
    });
  }

  return [
    {
      label: 'Alerts',
      fill: false,
      lineTension: 0.2,
      pointHitRadius: 10,
      pointRadius: 0.1,
      borderWidth: 2,
      backgroundColor: colors.data.flamePea,
      borderColor: colors.data.flamePea,
      hoverBackgroundColor: colors.data.flamePea,
      hoverBorderColor: 'black',
      data: alertsData,
    },
    {
      label: 'Events',
      fill: false,
      backgroundColor: colors.secondary.dodgerBlue,
      borderColor: colors.secondary.dodgerBlue,
      borderWidth: 1,
      hoverBackgroundColor: colors.secondary.dodgerBlue,
      hoverBorderColor: 'black',
      data: eventData,
    },
  ];
}

export default function EventAlertChart(props: Props) {
  const params = useParams();
  const networkId: string = nullthrows(params.networkId);
  const [start, end] = props.startEnd;
  const enqueueSnackbar = useEnqueueSnackbar();
  const [isLoading, setIsLoading] = useState(true);

  const [eventDataset, setEventDataset] = useState<Dataset>({
    label: 'Events',
    backgroundColor: colors.secondary.dodgerBlue,
    borderColor: colors.secondary.dodgerBlue,
    borderWidth: 1,
    hoverBackgroundColor: colors.secondary.malibu,
    hoverBorderColor: 'black',
    data: [],
    fill: false,
  });

  const [alertDataset, setAlertDataset] = useState<Dataset>({
    label: 'Alerts',
    backgroundColor: colors.data.flamePea,
    borderColor: colors.data.flamePea,
    borderWidth: 1,
    hoverBackgroundColor: colors.data.flamePea,
    hoverBorderColor: 'black',
    data: [],
    fill: false,
  });

  const [delta, unit] = getStep(start, end);
  useEffect(() => {
    // fetch queries
    const fetchAllData = async () => {
      const [eventDataset, alertDataset] = await getEventAlertDataset({
        start,
        end,
        delta,
        unit,
        networkId,
        enqueueSnackbar,
      });
      setEventDataset(eventDataset);
      setAlertDataset(alertDataset);
      setIsLoading(false);
    };

    void fetchAllData();
  }, [start, end, delta, unit, enqueueSnackbar, networkId]);

  if (isLoading) {
    return <LoadingFiller />;
  }

  return (
    <>
      <CardTitleRow
        icon={DataUsageIcon}
        label="Frequency of Alerts and Events"
      />
      <Card elevation={0}>
        <CardHeader
          subheader={
            <CustomLineChart
              start={start}
              end={end}
              delta={delta}
              unit={unit}
              dataset={[eventDataset, alertDataset]}
              yLabel={'Count'}
            />
          }
        />
      </Card>
    </>
  );
}
