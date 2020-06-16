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
import Grid from '@material-ui/core/Grid';
import ListAltIcon from '@material-ui/icons/ListAlt';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import moment from 'moment';
import nullthrows from '@fbcnms/util/nullthrows';

import {DateTimePicker} from '@material-ui/pickers';
import {getStep} from '../../components/CustomHistogram';
import {useEffect, useMemo, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

const BLUE = '#3984FF';

export default function LogChart() {
  const {match} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);
  const gatewayId: string = nullthrows(match.params.gatewayId);
  const [startDate, setStartDate] = useState(moment().subtract(3, 'hours'));
  const [endDate, setEndDate] = useState(moment());
  const startEnd = useMemo(() => {
    const [delta, unit, format] = getStep(startDate, endDate);
    return {
      start: startDate,
      end: endDate,
      delta: delta,
      unit: unit,
      format: format,
    };
  }, [startDate, endDate]);

  const enqueueSnackbar = useEnqueueSnackbar();
  const [labels, setLabels] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [dataset, setDataset] = useState<Dataset>({
    label: 'Log Counts',
    backgroundColor: BLUE,
    borderColor: BLUE,
    borderWidth: 1,
    hoverBackgroundColor: BLUE,
    hoverBorderColor: 'black',
    data: [],
  });

  useEffect(() => {
    // build queries
    let requestError = '';
    const queries = [];
    const logLabels = [];
    let s = startEnd.start.clone();
    while (startEnd.end.diff(s) >= 0) {
      logLabels.push(s.format(startEnd.format));
      const e = s.clone();
      e.add(startEnd.delta, startEnd.unit);
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
          backgroundColor: BLUE,
          borderColor: BLUE,
          borderWidth: 1,
          hoverBackgroundColor: BLUE,
          hoverBorderColor: 'black',
          data: logData,
        };
        setDataset(ds);
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
  }, [startEnd, enqueueSnackbar, gatewayId, networkId]);

  if (isLoading) {
    return <LoadingFiller />;
  }

  return (
    <Grid container align="top" alignItems="flex-start">
      <Grid item xs={6}>
        <Text>
          <ListAltIcon />
          Logs
        </Text>
      </Grid>
      <Grid item xs={6}>
        <Grid container justify="flex-end" alignItems="center" spacing={1}>
          <Grid item>
            <Text>Filter By Date</Text>
          </Grid>
          <Grid item>
            <DateTimePicker
              autoOk
              variant="inline"
              inputVariant="outlined"
              maxDate={endDate}
              disableFuture
              value={startDate}
              onChange={setStartDate}
            />
          </Grid>
          <Grid item>
            <Text>To</Text>
          </Grid>
          <Grid item>
            <DateTimePicker
              autoOk
              variant="inline"
              inputVariant="outlined"
              disableFuture
              value={endDate}
              onChange={setEndDate}
            />
          </Grid>
        </Grid>
      </Grid>
      <Grid item xs={12}>
        <CustomHistogram dataset={[dataset]} labels={labels} />
      </Grid>
    </Grid>
  );
}
