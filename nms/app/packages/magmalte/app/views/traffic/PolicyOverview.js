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
import type {policy_rule} from '@fbcnms/magma-api';

import ActionTable from '../../components/ActionTable';
import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import JsonEditor from '../../components/JsonEditor';
import LibraryBooksIcon from '@material-ui/icons/LibraryBooks';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextField from '@material-ui/core/TextField';
import nullthrows from '@fbcnms/util/nullthrows';

import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';

const POLICY_TITLE = 'Policies';
const DEFAULT_POLICY_CONFIG = {
  flow_list: [],
  id: '',
  monitoring_key: '',
  priority: 1,
};

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(3),
    flexGrow: 1,
  },
  topBar: {
    backgroundColor: colors.primary.mirage,
    padding: '20px 40px 20px 40px',
    color: colors.primary.white,
  },
  tabBar: {
    backgroundColor: colors.primary.brightGray,
    padding: `0 ${theme.spacing(5)}px`,
  },
  tabs: {
    color: colors.primary.white,
  },
  tab: {
    fontSize: '18px',
    textTransform: 'none',
  },
  tabLabel: {
    padding: '16px 0 16px 0',
    display: 'flex',
    alignItems: 'center',
  },
  tabIconLabel: {
    marginRight: '8px',
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
    textPrimary: colors.primary.mirage,
    outlined: true,
    contained: true,
    color: colors.primary.mirage,
  },
  // TODO: remove this when we actually fill out the grid sections
  contentPlaceholder: {
    padding: '50px 0',
  },
  paper: {
    height: 100,
    padding: theme.spacing(10),
    textAlign: 'center',
  },
  formControl: {
    margin: theme.spacing(1),
    minWidth: 120,
  },
}));
type PolicyRowType = {
  policyID: string,
  numFlows: number,
  priority: number,
  numSubscribers: number,
  monitoringKey: string,
  rating: string,
  trackingType: string,
};

type policiesType = {
  policies: {[string]: policy_rule},
};
export default function PolicyOverview(props: policiesType) {
  const classes = useStyles();
  // this for enabling edit, deactivate actions
  const [currRow, setCurrRow] = useState<PolicyRowType>({});
  const {history, relativeUrl} = useRouter();
  const policyRows: Array<PolicyRowType> = props.policies
    ? Object.keys(props.policies).map((policyID: string) => {
        const policyRule = props.policies[policyID];
        return {
          policyID: policyRule.id,
          numFlows: policyRule.flow_list.length,
          priority: policyRule.priority,
          numSubscribers: policyRule.assigned_subscribers?.length ?? 0,
          monitoringKey: policyRule.monitoring_key ?? '',
          rating: policyRule.rating_group?.toString() ?? 'not found',
          trackingType: policyRule.tracking_type ?? 'NO_TRACKING',
        };
      })
    : [];
  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={3}>
        <Grid container>
          <Grid item xs={6}>
            <Text key="title">
              <LibraryBooksIcon /> {POLICY_TITLE}
            </Text>
          </Grid>
          <Grid
            container
            item
            xs={6}
            justify="flex-end"
            alignItems="center"
            spacing={2}>
            <Grid item>
              <Button className={classes.appBarBtnSecondary}>
                Download Template
              </Button>
            </Grid>

            <Grid item>
              <Button className={classes.appBarBtnSecondary}>Upload CSV</Button>
            </Grid>

            <Grid item>
              <Button
                className={classes.appBarBtn}
                onClick={() => history.push(relativeUrl('/json'))}>
                Create New Policy
              </Button>
            </Grid>
          </Grid>
        </Grid>

        <Grid item xs={12}>
          <ActionTable
            data={policyRows}
            columns={[
              {title: 'Policy ID', field: 'policyID'},
              {title: 'Flows', field: 'numFlows', type: 'numeric'},
              {title: 'Priority', field: 'priority', type: 'numeric'},
              {title: 'Subscribers', field: 'numSubscribers', type: 'numeric'},
              {
                title: 'Monitoring Key',
                field: 'monitoringKey',
                render: rowData => {
                  return (
                    <TextField
                      type="password"
                      value={rowData.monitoringKey}
                      InputProps={{
                        disableUnderline: true,
                        readOnly: true,
                      }}
                    />
                  );
                },
              },
              {title: 'Rating', field: 'rating'},
              {title: 'Tracking Type', field: 'trackingType'},
            ]}
            handleCurrRow={(row: PolicyRowType) => setCurrRow(row)}
            menuItems={[
              {
                name: 'Edit JSON',
                handleFunc: () => {
                  history.push(relativeUrl('/' + currRow.policyID + '/json'));
                },
              },
              {name: 'Deactivate'},
              {name: 'Remove'},
            ]}
            options={{
              actionsColumnIndex: -1,
              pageSizeOptions: [5, 10],
            }}
          />
        </Grid>
      </Grid>
    </div>
  );
}

type Props = {
  policies: {[string]: policy_rule},
  onSave?: policy_rule => void,
};
export function PolicyJsonConfig(props: Props) {
  const {match, history} = useRouter();
  const [error, setError] = useState('');
  const networkId: string = nullthrows(match.params.networkId);
  const policyID: string = match.params.policyId;
  const enqueueSnackbar = useEnqueueSnackbar();
  const policy: policy_rule = props.policies[policyID] || DEFAULT_POLICY_CONFIG;
  return (
    <JsonEditor
      content={policy}
      error={error}
      onSave={async policy => {
        try {
          if (policyID) {
            await MagmaV1API.putNetworksByNetworkIdPoliciesRulesByRuleId({
              networkId: networkId,
              ruleId: policyID,
              policyRule: (policy: policy_rule),
            });
            enqueueSnackbar('Policy saved successfully', {
              variant: 'success',
            });
          } else {
            await MagmaV1API.postNetworksByNetworkIdPoliciesRules({
              networkId: networkId,
              policyRule: (policy: policy_rule),
            });
            enqueueSnackbar('Policy added successfully', {
              variant: 'success',
            });
          }
          setError('');
          props.onSave?.(policy);
          history.goBack();
        } catch (e) {
          setError(e.response?.data?.message ?? e.message);
        }
      }}
    />
  );
}
