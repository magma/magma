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
import type {
  mutable_subscriber,
  subscriber,
} from '../../../generated/MagmaAPIBindings';

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import ActionTable from '../../components/ActionTable';
import Button from '@material-ui/core/Button';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import CardTitleRow from '../../components/layout/CardTitleRow';
import DataGrid from '../../components/DataGrid';
import Grid from '@material-ui/core/Grid';
import JsonEditor from '../../components/JsonEditor';
import Link from '@material-ui/core/Link';
import React from 'react';
import SettingsIcon from '@material-ui/icons/Settings';
import SubscriberContext from '../../components/context/SubscriberContext';
// $FlowFixMe migrated to typescript
import nullthrows from '../../../shared/util/nullthrows';

import {EditSubscriberButton} from './SubscriberEditDialog';
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useContext, useState} from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';
import {useNavigate, useParams, useResolvedPath} from 'react-router-dom';

const useStyles = makeStyles(theme => ({
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
  const mutableSubscriber: mutable_subscriber = {
    ...subscriberInfoPartial,
  };

  if (config?.static_ips) {
    mutableSubscriber.static_ips = config.static_ips;
  }

  return (
    <JsonEditor
      content={mutableSubscriber}
      error={error}
      onSave={async (subscriber: mutable_subscriber) => {
        try {
          await ctx.setState?.(subscriber.id, subscriber);
          enqueueSnackbar('Subscriber saved successfully', {
            variant: 'success',
          });
          setError('');
        } catch (e) {
          setError(e.response?.data?.message ?? e.message);
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
  const ctx = useContext(SubscriberContext);
  const subscriberInfo = ctx.state?.[subscriberId];

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
  subscriberInfo: subscriber,
}) {
  function CollapseItems(props) {
    const data: DataRows[] = [
      [
        {
          value: props.data || '-',
        },
      ],
    ];

    return <DataGrid data={data} />;
  }

  const trafficPolicyData: DataRows[] = [
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

function SubscriberInfoConfig({subscriberInfo}: {subscriberInfo: subscriber}) {
  const [authKey, _setAuthKey] = useState(subscriberInfo.lte.auth_key);
  const [authOPC, _setAuthOPC] = useState(subscriberInfo.lte.auth_opc ?? false);
  const [dataPlan, _setDataPlan] = useState(subscriberInfo.lte.sub_profile);

  function CollapseItems(props) {
    const data: DataRows[] = [
      [
        {
          value: props.data || '-',
        },
      ],
    ];

    return <DataGrid data={data} />;
  }

  const kpiData: DataRows[] = [
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
        value: dataPlan,
      },
    ],
    [
      {
        category: 'Auth Key',
        value: authKey,
        obscure: true,
      },
    ],
  ];

  if (authOPC) {
    kpiData.push([
      {
        category: 'Auth OPC',
        value: authOPC,
      },
    ]);
  }

  return <DataGrid data={kpiData} />;
}

function SubscriberApnStaticIpsTable({
  subscriberInfo,
}: {
  subscriberInfo: subscriber,
}) {
  const resolvedPath = useResolvedPath('');
  const navigate = useNavigate();
  const staticIps = subscriberInfo.config.static_ips || {};
  type SubscriberApnStaticIpsRowType = {
    apnName: string,
    apnStaticIp: string,
  };
  const apnRows: Array<SubscriberApnStaticIpsRowType> = Object.keys(
    staticIps,
  ).map((apnName: string) => {
    return {
      apnName: apnName,
      apnStaticIp: staticIps[apnName],
    };
  });
  const [_currRow, setCurrRow] = useState<SubscriberApnStaticIpsRowType>({});
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
              }}>
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
