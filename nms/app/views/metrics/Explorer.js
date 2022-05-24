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
import ActionTable from '../../components/ActionTable';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import CardTitleRow from '../../components/layout/CardTitleRow';
import ExpandLess from '@material-ui/icons/ExpandLess';
import ExpandMore from '@material-ui/icons/ExpandMore';
import ExploreIcon from '@material-ui/icons/Explore';
import Grid from '@material-ui/core/Grid';
// $FlowFixMe migrated to typescript
import LoadingFiller from '../../components/LoadingFiller';
import MagmaV1API from '../../../generated/WebClient';
import React from 'react';
import moment from 'moment';
// $FlowFixMe migrated to typescript
import nullthrows from '../../../shared/util/nullthrows';
import useMagmaAPI from '../../../api/useMagmaAPIFlow';
import type {prometheus_labelset} from '../../../generated/MagmaAPIBindings';

import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useEffect, useMemo, useState} from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';
import {useParams} from 'react-router-dom';

const TITLE = 'Metrics Explorer';

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
  appBarBtnSecondary: {
    color: colors.primary.white,
  },
}));

export type MetricsDetail = {
  MetricName: string,
  PromQL: string,
  Description: string,
  Category: string,
  Service: string,
};

export default function MetricsExplorer() {
  const classes = useStyles();
  const [isLoading, setIsLoading] = useState(true);
  const [LteMetricsTable, setLteMetricsTable] = useState([]);
  const enqueueSnackbar = useEnqueueSnackbar();
  const params = useParams();

  const networkId = nullthrows(params.networkId);
  const startEnd = useMemo(() => {
    return {
      start: moment().subtract(3, 'hours'),
      end: moment(),
    };
  }, []);

  const {
    response: metricSeries,
    isLoading: isMetricSeriesLoading,
  } = useMagmaAPI(MagmaV1API.getNetworksByNetworkIdPrometheusSeries, {
    networkId,
    start: startEnd.start.toISOString(),
    end: startEnd.end.toISOString(),
  });

  const {
    response: tenantMetricSeriesDescription,
    isLoading: isTenantMetricSeriesDescriptionLoading,
  } = useMagmaAPI(MagmaV1API.getTenantsTargetsMetadata, {});

  useEffect(() => {
    fetch('/data/LteMetrics')
      .then(res => res.json())
      .then(
        result => {
          setLteMetricsTable(result);
          setIsLoading(false);
        },
        error => {
          enqueueSnackbar(
            `failed loading metrics data due to following error: ${error}`,
            {
              variant: 'error',
            },
          );
          setIsLoading(false);
        },
      );
  }, [enqueueSnackbar]);

  if (
    isLoading ||
    isMetricSeriesLoading ||
    isTenantMetricSeriesDescriptionLoading
  ) {
    return <LoadingFiller />;
  }

  // filter only those metrics which are relevant to this network
  const metricsMap = {};
  if (metricSeries != null) {
    metricSeries.forEach((labelSet: prometheus_labelset) => {
      metricsMap[labelSet['__name__']] = labelSet;
    });
  }

  const metricsTable = LteMetricsTable.filter((metricEnt: MetricsDetail) => {
    if (metricEnt.MetricName in metricsMap) {
      delete metricsMap[metricEnt.MetricName];
      return true;
    }
    return false;
  });

  const tenantMetricDescription = {};
  tenantMetricSeriesDescription?.forEach(metricDesc => {
    tenantMetricDescription[metricDesc.metric] = metricDesc.help;
  });

  Object.keys(metricsMap).forEach(function (key) {
    let metricDescription = tenantMetricDescription[key];
    if (metricDescription === undefined || metricDescription === '') {
      metricDescription = 'Description unavailable';
    }

    metricsTable.push({
      MetricName: key,
      PromQL: key,
      Category: 'Category unavailable',
      Description: metricDescription,
      Service: metricsMap[key]?.['service'] ?? 'Service name unavailable',
    });
  });

  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={4}>
        <Grid item xs={12}>
          <CardTitleRow icon={ExploreIcon} label={TITLE} />
          <ActionTable
            data={metricsTable}
            columns={[
              {title: 'Name', field: 'MetricName'},
              {title: 'Description', field: 'Description'},
              {title: 'Category', field: 'Category'},
              {title: 'Service', field: 'Service'},
            ]}
            options={{
              actionsColumnIndex: -1,
              pageSizeOptions: [10, 20],
              pageSize: 20,
            }}
            detailPanel={[
              {
                icon: ExpandMore,
                openIcon: ExpandLess,
                render: rowData => (
                  <iframe
                    width="100%"
                    height="400"
                    src={encodeURI(
                      `/grafana/explore?left=["now-1h","now","default",{"expr":"${rowData.PromQL}"}]`,
                    )}
                  />
                ),
              },
            ]}
          />
        </Grid>
      </Grid>
    </div>
  );
}
