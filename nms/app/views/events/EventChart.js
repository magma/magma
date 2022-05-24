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
import type {Dataset, DatasetType} from '../../components/CustomMetrics';

import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import CustomHistogram from '../../components/CustomMetrics';
// $FlowFixMe migrated to typescript
import LoadingFiller from '../../components/LoadingFiller';
import MagmaV1API from '../../../generated/WebClient';
import React from 'react';
import moment from 'moment';
// $FlowFixMe migrated to typescript
import nullthrows from '../../../shared/util/nullthrows';

import {colors} from '../../theme/default';
import {getQueryRanges} from '../../components/CustomMetrics';
import {useEffect, useState} from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';
import {useParams} from 'react-router-dom';

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
  const params = useParams();
  const networkId: string = nullthrows(params.networkId);
  const {start, end, delta, format, unit, streams, tags, setEventCount} = props;
  const enqueueSnackbar = useEnqueueSnackbar();
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
    const queries = getQueryRanges(start, end, delta, unit);
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
        const data: Array<DatasetType> = allResponses.map((r, index) => {
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

        const ds: Dataset = {
          label: 'Event Counts',
          backgroundColor: colors.secondary.dodgerBlue,
          borderColor: colors.secondary.dodgerBlue,
          borderWidth: 1,
          hoverBackgroundColor: colors.secondary.dodgerBlue,
          hoverBorderColor: 'black',
          data: data,
        };
        setDataset(ds);
        setEventCount(data.reduce((a, b) => a + b.y, 0));
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
      <CardHeader subheader={<CustomHistogram dataset={[dataset]} />} />
    </Card>
  );
}
