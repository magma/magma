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
 *
 * @flow strict-local
 * @format
 */
import type {Dataset} from '../../components/CustomHistogram';

import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import CustomHistogram from '../../components/CustomHistogram';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React from 'react';
import moment from 'moment';
import nullthrows from '@fbcnms/util/nullthrows';

import {colors} from '../../theme/default';
import {useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

type Props = {
  start: moment,
  end: moment,
  delta: number,
  unit: string,
  format: string,
  streams: string,
  tags: string,
  setEventCount: number => void,
};

export default function EventChart(props: Props) {
  const {match} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);
  const {start, end, delta, format, unit, streams, tags, setEventCount} = props;
  const enqueueSnackbar = useEnqueueSnackbar();
  const [labels, setLabels] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [dataset, setDataset] = useState<Dataset>({
    label: 'Event Counts',
    backgroundColor: colors.secondary.dodgerBlue,
    borderColor: colors.secondary.dodgerBlue,
    borderWidth: 1,
    hoverBackgroundColor: colors.secondary.dodgerBlue,
    hoverBorderColor: 'black',
    data: [],
  });

  useEffect(() => {
    // build queries
    let requestError = '';
    const queries = [];
    const logLabels = [];
    let s = start.clone();
    while (end.diff(s) >= 0) {
      logLabels.push(s.format(format));
      const e = s.clone();
      e.add(delta, unit);
      queries.push([s, e]);
      s = e.clone();
    }
    setLabels(logLabels);

    const requests = queries.map(async (query, _) => {
      try {
        const [s, e] = query;
        const response = await MagmaV1API.getEventsByNetworkIdAboutCount({
          networkId: networkId,
          streams: streams,
          tags: tags,
          start: s.toISOString(),
          end: e.toISOString(),
        });
        return response;
      } catch (error) {
        requestError = error;
      }
      return null;
    });

    Promise.all(requests)
      .then(allResponses => {
        const logData: Array<number> = allResponses.map(r => {
          if (r === null || r === undefined) {
            return 0;
          }
          return r;
        });

        const ds: Dataset = {
          label: 'Event Counts',
          backgroundColor: colors.secondary.dodgerBlue,
          borderColor: colors.secondary.dodgerBlue,
          borderWidth: 1,
          hoverBackgroundColor: colors.secondary.dodgerBlue,
          hoverBorderColor: 'black',
          data: logData,
        };
        setDataset(ds);
        setEventCount(logData.reduce((a, b) => a + b, 0));
        setIsLoading(false);
      })
      .catch(error => {
        requestError = error;
        setIsLoading(false);
      });

    if (requestError) {
      enqueueSnackbar('Error getting event counts', {
        variant: 'error',
      });
    }
  }, [
    setEventCount,
    start,
    end,
    delta,
    format,
    unit,
    enqueueSnackbar,
    streams,
    tags,
    networkId,
  ]);

  if (isLoading) {
    return <LoadingFiller />;
  }

  return (
    <Card elevation={0}>
      <CardHeader
        subheader={<CustomHistogram dataset={[dataset]} labels={labels} />}
      />
    </Card>
  );
}
