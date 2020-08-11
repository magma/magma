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
import type {subscriber} from '@fbcnms/magma-api';

import Button from '@material-ui/core/Button';
import CardTitleRow from '../../components/layout/CardTitleRow';
import DataGrid from '../../components/DataGrid';
import Grid from '@material-ui/core/Grid';
import JsonEditor from '../../components/JsonEditor';
import React from 'react';
import SettingsIcon from '@material-ui/icons/Settings';
import SubscriberContext from '../../components/context/SubscriberContext';
import nullthrows from '@fbcnms/util/nullthrows';

import {EditSubscriberButton} from './SubscriberAddDialog';
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useContext, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

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
  const {match} = useRouter();
  const [error, setError] = useState('');
  const enqueueSnackbar = useEnqueueSnackbar();
  const subscriberId = nullthrows(match.params.subscriberId);
  const ctx = useContext(SubscriberContext);
  const subscriberInfo = ctx.state?.[subscriberId];

  return (
    <JsonEditor
      content={subscriberInfo}
      error={error}
      onSave={async subscriber => {
        try {
          await ctx.setState(subscriber.id, {...subscriber});
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

export default function SubscriberDetailConfig({
  subscriberInfo,
}: {
  subscriberInfo: subscriber,
}) {
  const classes = useStyles();
  const {history, relativeUrl} = useRouter();
  function TrafficFilter() {
    return <Button variant="text">Edit</Button>;
  }

  function ConfigFilter() {
    return (
      <Button
        className={classes.appBarBtn}
        onClick={() => {
          history.push(relativeUrl('/json'));
        }}>
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
              <CardTitleRow label="Subscriber" filter={EditSubscriberButton} />
              <SubscriberInfoConfig
                readOnly={true}
                subscriberInfo={subscriberInfo}
              />
            </Grid>

            <Grid item xs={12} md={6}>
              <CardTitleRow label="Traffic Policy" filter={TrafficFilter} />
              <SubscriberConfigTrafficPolicy
                readOnly={true}
                subscriberInfo={subscriberInfo}
              />
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
        collapse:
          subscriberInfo.active_apns?.map(data => (
            <CollapseItems data={data} />
          )) || false,
      },
    ],
    [
      {
        category: 'Base Names',
        value: subscriberInfo.active_base_names?.length || 0,
        collapse:
          subscriberInfo.active_base_names?.map(data => (
            <CollapseItems data={data} />
          )) || false,
      },
    ],
    [
      {
        category: 'Active Policies',
        value: subscriberInfo.active_policies?.length || 0,
        collapse:
          subscriberInfo.active_policies?.map(data => (
            <CollapseItems data={data} />
          )) || false,
      },
    ],
  ];

  return <DataGrid data={trafficPolicyData} />;
}

function SubscriberInfoConfig({subscriberInfo}: {subscriberInfo: subscriber}) {
  const [authKey, _setAuthKey] = useState(subscriberInfo.lte.auth_key);
  const [authOPC, _setAuthOPC] = useState(subscriberInfo.lte.auth_opc ?? false);
  const [dataPlan, _setDataPlan] = useState(subscriberInfo.lte.sub_profile);

  const kpiData: DataRows[] = [
    [
      {
        category: 'LTE Network Access',
        value: subscriberInfo.lte.state,
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
