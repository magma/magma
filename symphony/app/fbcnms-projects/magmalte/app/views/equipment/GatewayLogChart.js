/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import type {Dataset} from '../../components/CustomHistogram';

import CustomHistogram from '../../components/CustomHistogram';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React from 'react';
import moment from 'moment';
import nullthrows from '@fbcnms/util/nullthrows';
import Text from '../../theme/design-system/Text';

import {Card, CardHeader} from '@material-ui/core/';
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
  setLogCount: number => void,
};

export default function LogChart(props: Props) {
  const {match} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);
  const gatewayId: string = nullthrows(match.params.gatewayId);
  const {start, end, delta, format, unit, setLogCount} = props;
  const enqueueSnackbar = useEnqueueSnackbar();
  const [labels, setLabels] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [dataset, setDataset] = useState<Dataset>({
    label: 'Log Counts',
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
        const response = await MagmaV1API.getNetworksByNetworkIdLogsCount({
          networkId: networkId,
          filters: `gateway_id:${gatewayId}`,
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
          label: 'Log Counts',
          backgroundColor: colors.secondary.dodgerBlue,
          borderColor: colors.secondary.dodgerBlue,
          borderWidth: 1,
          hoverBackgroundColor: colors.secondary.dodgerBlue,
          hoverBorderColor: 'black',
          data: logData,
        };
        setDataset(ds);
        setLogCount(logData.reduce((a, b) => a + b, 0));
        setIsLoading(false);
      })
      .catch(error => {
        requestError = error;
        setIsLoading(false);
      });

    if (requestError) {
      enqueueSnackbar('Error getting log counts', {
        variant: 'error',
      });
    }
  }, [
    setLogCount,
    start,
    end,
    delta,
    format,
    unit,
    enqueueSnackbar,
    gatewayId,
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
