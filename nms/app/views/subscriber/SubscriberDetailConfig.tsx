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

import ActionTable from '../../components/ActionTable';
import Button from '@mui/material/Button';
import CardTitleRow from '../../components/layout/CardTitleRow';
import DataGrid from '../../components/DataGrid';
import Grid from '@mui/material/Grid';
import JsonEditor from '../../components/JsonEditor';
import Link from '@mui/material/Link';
import LoadingFiller from '../../components/LoadingFiller';
import MagmaAPI from '../../api/MagmaAPI';
import React from 'react';
import SettingsIcon from '@mui/icons-material/Settings';
import SubscriberContext from '../../context/SubscriberContext';
import nullthrows from '../../../shared/util/nullthrows';
import useMagmaAPI from '../../api/useMagmaAPI';
import {EditSubscriberButton} from './SubscriberEditDialog';
import {
  MutableSubscriber,
  Subscriber,
  SubscriberForbiddenNetworkTypesEnum,
} from '../../../generated';
import {Theme} from '@mui/material/styles';
import {colors, typography} from '../../theme/default';
import {getErrorMessage} from '../../util/ErrorUtils';
import {makeStyles} from '@mui/styles';
import {useCallback, useContext, useState} from 'react';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';
import {useNavigate, useParams, useResolvedPath} from 'react-router-dom';
import type {DataRows} from '../../components/DataGrid';

const useStyles = makeStyles<Theme>(theme => ({
  dashboardRoot: {
    margin: theme.spacing(5),
    flexGrow: 1,
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
}));

export function SubscriberJsonConfig() {
  const params = useParams();
  const [error, setError] = useState('');
  const enqueueSnackbar = useEnqueueSnackbar();
  const subscriberId = nullthrows(params.subscriberId);
  const ctx = useContext(SubscriberContext);
  const subscriberInfo = ctx.state?.[subscriberId];
  const {
    config,
    monitoring: _unused_monitoring,
    state: _unused_state,
    ...subscriberInfoPartial
  } = subscriberInfo;
  const mutableSubscriber: MutableSubscriber = {...subscriberInfoPartial};

  if (config?.static_ips) {
    mutableSubscriber.static_ips = config.static_ips;
  }

  return (
    <JsonEditor
      content={mutableSubscriber}
      error={error}
      onSave={async (subscriber: MutableSubscriber) => {
        try {
          await ctx.setState?.(subscriber.id, subscriber);
          enqueueSnackbar('Subscriber saved successfully', {
            variant: 'success',
          });
          setError('');
        } catch (e) {
          setError(getErrorMessage(e));
        }
      }}
    />
  );
}

export default function SubscriberDetailConfig() {
  const classes = useStyles();
  const params = useParams();
  const navigate = useNavigate();
  const subscriberId = nullthrows(params.subscriberId);
  const networkId: string = nullthrows(params.networkId);
  const ctx = useContext(SubscriberContext);
  const [subscriberConfig, setSubscriberConfig] = useState({} as Subscriber);
  const {isLoading} = useMagmaAPI(
    MagmaAPI.subscribers.lteNetworkIdSubscribersSubscriberIdGet,
    {
      networkId: networkId,
      subscriberId: subscriberId,
    },
    useCallback(
      (response: Subscriber) => {
        setSubscriberConfig(response);

        if (!ctx.state[subscriberId]) {
          void ctx.setState?.('', undefined, {
            ...ctx.state,
            [subscriberId]: response,
          });
        }
      },
      [ctx, subscriberId],
    ),
  );

  if (isLoading) {
    return <LoadingFiller />;
  }

  const subscriberInfo = ctx.state?.[subscriberId] || subscriberConfig;
  function ConfigFilter() {
    return (
      <Button className={classes.appBarBtn} onClick={() => navigate('json')}>
        Edit JSON
      </Button>
    );
  }

  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={4}>
        <Grid item xs={12}>
          <Grid item xs={12}>
            <CardTitleRow
              icon={SettingsIcon}
              label="Config"
              filter={ConfigFilter}
            />
          </Grid>
          <Grid container spacing={4}>
            <Grid item xs={12} md={6}>
              <CardTitleRow
                label="Subscriber"
                filter={() => EditSubscriberButton({editTable: 'subscriber'})}
              />
              <SubscriberInfoConfig subscriberInfo={subscriberInfo} />
            </Grid>

            <Grid item xs={12} md={6}>
              <CardTitleRow
                label="Traffic Policy"
                filter={() =>
                  EditSubscriberButton({editTable: 'trafficPolicy'})
                }
              />
              <SubscriberConfigTrafficPolicy subscriberInfo={subscriberInfo} />
            </Grid>
            <Grid item xs={12} md={6}>
              <CardTitleRow
                label="APN Static IPs"
                filter={() => EditSubscriberButton({editTable: 'staticIps'})}
              />
              <SubscriberApnStaticIpsTable subscriberInfo={subscriberInfo} />
            </Grid>
          </Grid>
        </Grid>
      </Grid>
    </div>
  );
}

function SubscriberConfigTrafficPolicy({
  subscriberInfo,
}: {
  subscriberInfo: Subscriber;
}) {
  function CollapseItems(props: {key: string; data: string}) {
    const data: Array<DataRows> = [
      [
        {
          value: props.data || '-',
        },
      ],
    ];

    return <DataGrid data={data} />;
  }

  const trafficPolicyData: Array<DataRows> = [
    [
      {
        category: 'Active APNs',
        value: subscriberInfo.active_apns?.length || 0,
        collapse: subscriberInfo.active_apns?.map(data => (
          <CollapseItems key={data} data={data} />
        )) || <></>,
      },
    ],
    [
      {
        category: 'Base Names',
        value: subscriberInfo.active_base_names?.length || 0,
        collapse: subscriberInfo.active_base_names?.map(data => (
          <CollapseItems key={data} data={data} />
        )) || <></>,
      },
    ],
    [
      {
        category: 'Active Policies',
        value: subscriberInfo.active_policies?.length || 0,
        collapse: subscriberInfo.active_policies?.map(data => (
          <CollapseItems key={data} data={data} />
        )) || <></>,
      },
    ],
  ];

  return <DataGrid data={trafficPolicyData} />;
}

function SubscriberInfoConfig({subscriberInfo}: {subscriberInfo: Subscriber}) {
  function CollapseItems(props: {
    key: SubscriberForbiddenNetworkTypesEnum;
    data: SubscriberForbiddenNetworkTypesEnum;
  }) {
    const data: Array<DataRows> = [
      [
        {
          value: props.data || '-',
        },
      ],
    ];

    return <DataGrid data={data} />;
  }

  const kpiData: Array<DataRows> = [
    [
      {
        category: 'LTE Network Access',
        value: subscriberInfo.lte.state,
      },
    ],
    [
      {
        category: 'Forbidden Network Types',
        value: subscriberInfo.forbidden_network_types?.length || 0,
        collapse: subscriberInfo.forbidden_network_types?.map(data => (
          <CollapseItems key={data} data={data} />
        )) || <></>,
      },
    ],
    [
      {
        category: 'Data plan',
        value: subscriberInfo.lte.sub_profile,
      },
    ],
    [
      {
        category: 'Auth Key',
        value: subscriberInfo.lte.auth_key,
        obscure: true,
      },
    ],
  ];

  if (subscriberInfo.lte.auth_opc) {
    kpiData.push([
      {
        category: 'Auth OPC',
        value: subscriberInfo.lte.auth_opc,
      },
    ]);
  }

  return <DataGrid data={kpiData} />;
}

function SubscriberApnStaticIpsTable({
  subscriberInfo,
}: {
  subscriberInfo: Subscriber;
}) {
  const resolvedPath = useResolvedPath('');
  const navigate = useNavigate();
  const staticIps = subscriberInfo.config.static_ips || {};
  type SubscriberApnStaticIpsRowType = {
    apnName: string;
    apnStaticIp: string;
  };
  const apnRows: Array<SubscriberApnStaticIpsRowType> = Object.keys(
    staticIps,
  ).map((apnName: string) => {
    return {
      apnName: apnName,
      apnStaticIp: staticIps[apnName],
    };
  });
  const [, setCurrRow] = useState({} as SubscriberApnStaticIpsRowType);
  return (
    <ActionTable
      title=""
      data={apnRows}
      columns={[
        {
          title: 'APN Name',
          field: 'apnName',
          render: currRow => (
            <Link
              variant="body2"
              component="button"
              onClick={() => {
                navigate(
                  resolvedPath.pathname.replace(
                    `subscribers/overview/${subscriberInfo.id}/config`,
                    `traffic/apn`,
                  ),
                );
              }}
              underline="hover">
              {currRow.apnName}
            </Link>
          ),
        },
        {title: 'Static IP', field: 'apnStaticIp'},
      ]}
      handleCurrRow={(row: SubscriberApnStaticIpsRowType) => setCurrRow(row)}
      options={{
        actionsColumnIndex: -1,
        pageSizeOptions: [],
        toolbar: false,
        paging: false,
      }}
    />
  );
}
