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
import type {Metrics} from '../../components/context/SubscriberContext';
import type {network_id, subscriber} from '@fbcnms/magma-api';

import ActionTable from '../../components/ActionTable';
import AddSubscriberButton from './SubscriberAddDialog';
import Button from '@material-ui/core/Button';
import CardTitleRow from '../../components/layout/CardTitleRow';
import Grid from '@material-ui/core/Grid';
import Link from '@material-ui/core/Link';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import PeopleIcon from '@material-ui/icons/People';
import React from 'react';
import SubscriberContext from '../../components/context/SubscriberContext';
import SubscriberDetail from './SubscriberDetail';
import TopBar from '../../components/TopBar';
import nullthrows from '@fbcnms/util/nullthrows';

import {Redirect, Route, Switch} from 'react-router-dom';
import {colors, typography} from '../../theme/default';
import {getLabelUnit} from './SubscriberUtils';
import {makeStyles} from '@material-ui/styles';
import {useContext, useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

const TITLE = 'Subscribers';

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

type InitSubscriberStateProps = {
  networkId: network_id,
  setSubscriberMap: ({[string]: subscriber}) => void,
  setSubscriberMetrics: ({[string]: Metrics}) => void,
  enqueueSnackbar: (msg: string, cfg: {}) => ?(string | number),
};

async function InitSubscriberState(props: InitSubscriberStateProps) {
  const {
    networkId,
    setSubscriberMap,
    setSubscriberMetrics,
    enqueueSnackbar,
  } = props;

  let subscribers = {};
  try {
    subscribers = await MagmaV1API.getLteByNetworkIdSubscribers({networkId});
  } catch (e) {
    enqueueSnackbar('failed fetching subscriber information', {
      variant: 'error',
    });
    return [];
  }
  setSubscriberMap(subscribers);

  const subscriberMetrics = {};
  const queries = {
    dailyAvg: 'avg (avg_over_time(ue_traffic[24h])) by (IMSI)',
    currentUsage: 'sum (ue_traffic) by (IMSI)',
  };

  const requests = Object.keys(queries).map(async (queryType: string) => {
    try {
      const resp = await MagmaV1API.getNetworksByNetworkIdPrometheusQuery({
        networkId,
        query: queries[queryType],
      });

      resp?.data?.result?.filter(Boolean).forEach(item => {
        const imsi = Object.values(item?.metric)?.[0];
        if (typeof imsi === 'string') {
          const [value, unit] = getLabelUnit(parseFloat(item?.value?.[1]));
          if (!(imsi in subscriberMetrics)) {
            subscriberMetrics[imsi] = {};
          }
          subscriberMetrics[imsi][queryType] = `${value}${unit}`;
        }
      });
    } catch (e) {
      enqueueSnackbar('failed fetching current usage information', {
        variant: 'error',
      });
    }
  });
  await Promise.all(requests);
  setSubscriberMetrics(subscriberMetrics);
}

export default function SubscriberDashboard() {
  const {match, relativePath, relativeUrl} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);
  const [subscriberMap, setSubscriberMap] = useState({});
  const [subscriberMetrics, setSubscriberMetrics] = useState({});
  const [isLoading, setIsLoading] = useState(true);
  const enqueueSnackbar = useEnqueueSnackbar();

  useEffect(() => {
    const fetchData = async () => {
      await InitSubscriberState({
        networkId,
        setSubscriberMap,
        setSubscriberMetrics,
        enqueueSnackbar,
      });
      setIsLoading(false);
    };
    fetchData();
  }, [networkId, enqueueSnackbar]);

  const updateSubscriberMap = async (key: string, val: subscriber) => {
    if (key in subscriberMap) {
      await MagmaV1API.putLteByNetworkIdSubscribersBySubscriberId({
        networkId,
        subscriber: val,
        subscriberId: key,
      });
    } else {
      await MagmaV1API.postLteByNetworkIdSubscribers({
        networkId,
        subscriber: val,
      });
    }
    setSubscriberMap({...subscriberMap, [key]: val});
  };

  if (isLoading) {
    return <LoadingFiller />;
  }

  const subscriberCtx = {
    state: subscriberMap,
    metrics: subscriberMetrics,
    setState: updateSubscriberMap,
  };

  return (
    <Switch>
      <Route
        path={relativePath('/overview/:subscriberId')}
        render={() => {
          return (
            <SubscriberContext.Provider value={subscriberCtx}>
              <SubscriberDetail />
            </SubscriberContext.Provider>
          );
        }}
      />

      <Route
        path={relativePath('/overview')}
        render={() => {
          return (
            <SubscriberContext.Provider value={subscriberCtx}>
              <SubscriberDashboardInternal />
            </SubscriberContext.Provider>
          );
        }}
      />
      <Redirect to={relativeUrl('/overview')} />
    </Switch>
  );
}

type SubscriberRowType = {
  name: string,
  imsi: string,
  service: string,
  currentUsage: string,
  dailyAvg: string,
  lastReportedTime: Date,
};

function SubscriberDashboardInternal() {
  const classes = useStyles();
  const {history, relativeUrl} = useRouter();
  const [currRow, setCurrRow] = useState<SubscriberRowType>({});
  const ctx = useContext(SubscriberContext);
  const subscriberMap = ctx.state;
  const subscriberMetrics = ctx.metrics;

  return (
    <>
      <TopBar
        header={TITLE}
        tabs={[
          {
            label: 'Subscribers',
            to: '/subscribersv2',
            icon: PeopleIcon,
            filters: (
              <Grid
                container
                justify="flex-end"
                alignItems="center"
                spacing={2}>
                <Grid item>
                  {/* TODO: these button styles need to be localized */}
                  <Button variant="text" className={classes.appBarBtnSecondary}>
                    Secondary Action
                  </Button>
                </Grid>
                <Grid item>
                  <Button variant="contained" className={classes.appBarBtn}>
                    Primary Action
                  </Button>
                </Grid>
              </Grid>
            ),
          },
        ]}
      />

      <div className={classes.dashboardRoot}>
        <Grid container spacing={4}>
          <Grid item xs={12}>
            <CardTitleRow
              icon={PeopleIcon}
              label={TITLE}
              filter={AddSubscriberButton}
            />

            {subscriberMap ? (
              <ActionTable
                data={Object.keys(subscriberMap).map((imsi: string) => {
                  const subscriberInfo = subscriberMap[imsi];
                  const metrics = subscriberMetrics[`${imsi}`];
                  return {
                    name: subscriberInfo.name ?? imsi,
                    imsi: imsi,
                    service: subscriberInfo.lte.state,
                    currentUsage: metrics?.currentUsage ?? 0,
                    dailyAvg: metrics?.dailyAvg ?? 0,
                    lastReportedTime: new Date(
                      subscriberInfo.monitoring?.icmp?.last_reported_time ?? 0,
                    ),
                  };
                })}
                columns={[
                  {title: 'Name', field: 'name'},
                  {
                    title: 'IMSI',
                    field: 'imsi',
                    render: currRow => (
                      <Link
                        variant="body2"
                        component="button"
                        onClick={() =>
                          history.push(relativeUrl('/' + currRow.imsi))
                        }>
                        {currRow.imsi}
                      </Link>
                    ),
                  },
                  {title: 'Service', field: 'service', width: 100},
                  {title: 'Current Usage', field: 'currentUsage', width: 175},
                  {title: 'Daily Average', field: 'dailyAvg', width: 175},
                  {
                    title: 'Last Reported Time',
                    field: 'lastReportedTime',
                    type: 'datetime',
                    width: 200,
                  },
                ]}
                handleCurrRow={(row: SubscriberRowType) => setCurrRow(row)}
                menuItems={[
                  {
                    name: 'View',
                    handleFunc: () => {
                      history.push(relativeUrl('/' + currRow.imsi));
                    },
                  },
                  {
                    name: 'Edit',
                    handleFunc: () => {
                      history.push(relativeUrl('/' + currRow.imsi + '/config'));
                    },
                  },
                  {name: 'Remove'},
                ]}
                options={{
                  actionsColumnIndex: -1,
                  pageSizeOptions: [5, 10],
                }}
              />
            ) : (
              '<Text>No Subscribers Found</Text>'
            )}
          </Grid>
        </Grid>
      </div>
    </>
  );
}
