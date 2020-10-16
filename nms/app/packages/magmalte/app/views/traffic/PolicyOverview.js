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
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';

import ActionTable from '../../components/ActionTable';
import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import LibraryBooksIcon from '@material-ui/icons/LibraryBooks';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import LteNetworkContext from '../../components/context/LteNetworkContext';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import PolicyContext from '../../components/context/PolicyContext';
import PolicyRuleEditDialog from './PolicyEdit';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextField from '@material-ui/core/TextField';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';

import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';
import {Route, Switch} from 'react-router-dom';
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useContext, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

const POLICY_TITLE = 'Policies';

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
  networkWide: string,
};

export function PolicyOverview(props: WithAlert) {
  const [lastRefreshTime, setLastRefreshTime] = useState(new Date().getTime());
  const classes = useStyles();
  const enqueueSnackbar = useEnqueueSnackbar();
  const [currRow, setCurrRow] = useState<PolicyRowType>({});
  const {history, relativeUrl, relativePath} = useRouter();
  const ctx = useContext(PolicyContext);
  const lteNetworkCtx = useContext(LteNetworkContext);
  const lteNetwork = lteNetworkCtx.state;
  const ruleNames = new Set(
    lteNetwork?.subscriber_config?.network_wide_rule_names ?? [],
  );
  const {match} = useRouter();
  const {networkId} = match.params;
  const {isLoading: isPolicyLoading, response: qosProfiles} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdPolicyQosProfiles,
    {
      networkId: networkId,
    },
    undefined,
    lastRefreshTime,
  );
  const {isLoading: isSubsLoading, response: subscribers} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdSubscribers,
    {networkId: networkId},
    undefined,
    lastRefreshTime,
  );

  if (isPolicyLoading || isSubsLoading) {
    return <LoadingFiller />;
  }

  const policies = ctx.state;
  const policyRows: Array<PolicyRowType> = policies
    ? Object.keys(policies).map((policyID: string) => {
        const policyRule = policies[policyID];
        return {
          policyID: policyRule.id,
          numFlows: policyRule.flow_list.length,
          priority: policyRule.priority,
          numSubscribers: policyRule.assigned_subscribers?.length ?? 0,
          monitoringKey: policyRule.monitoring_key ?? '',
          rating: policyRule.rating_group?.toString() ?? 'Not Found',
          trackingType: policyRule.tracking_type ?? 'NO_TRACKING',
          networkWide: ruleNames.has(policyID) ? 'Enabled' : 'Disabled',
        };
      })
    : [];
  const closeDialog = () => {
    setLastRefreshTime(new Date().getTime());
    history.push(relativeUrl(''));
  };
  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={3}>
        <Grid container>
          <Grid item xs={6}>
            <Text key="title" data-testid={`title_${POLICY_TITLE}`}>
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
              {title: 'Network Wide', field: 'networkWide'},
            ]}
            handleCurrRow={(row: PolicyRowType) => setCurrRow(row)}
            menuItems={[
              {
                name: 'Edit',
                handleFunc: () => {
                  history.push(relativeUrl('/' + currRow.policyID + '/json'));
                },
              },
              {name: 'Deactivate'},
              {
                name: 'Remove',
                handleFunc: () => {
                  props
                    .confirm(
                      `Are you sure you want to delete ${currRow.policyID}?`,
                    )
                    .then(async confirmed => {
                      if (!confirmed) {
                        return;
                      }

                      try {
                        // trigger deletion
                        ctx.setState(currRow.policyID);
                      } catch (e) {
                        enqueueSnackbar(
                          'failed deleting policy ' + currRow.policyID,
                          {
                            variant: 'error',
                          },
                        );
                      }
                    });
                },
              },
            ]}
            options={{
              actionsColumnIndex: -1,
              pageSizeOptions: [5, 10],
            }}
          />
        </Grid>
      </Grid>
      <Switch>
        <Route
          path={relativePath('/json')}
          render={() => (
            <PolicyRuleEditDialog
              subscribers={subscribers ?? {}}
              qosProfiles={qosProfiles ?? {}}
              onCancel={closeDialog}
              onSave={async (policy, isNetworkWide) => {
                try {
                  await ctx.setState(policy.id, policy, isNetworkWide);
                  enqueueSnackbar('Policy saved successfully', {
                    variant: 'success',
                  });
                  history.goBack();
                } catch (e) {
                  enqueueSnackbar('Policy failed to save: ' + e.message, {
                    variant: 'error',
                  });
                }
              }}
            />
          )}
        />
        <Route
          path={relativePath('/:policyId/json')}
          render={() => (
            <PolicyRuleEditDialog
              rule={ctx.state[currRow.policyID]}
              subscribers={subscribers ?? {}}
              qosProfiles={qosProfiles ?? {}}
              onCancel={closeDialog}
              onSave={async (policy, isNetworkWide) => {
                try {
                  await ctx.setState(policy.id, policy, isNetworkWide);
                  enqueueSnackbar('Policy saved successfully', {
                    variant: 'success',
                  });
                  history.goBack();
                } catch (e) {
                  enqueueSnackbar('Policy failed to save: ' + e.message, {
                    variant: 'error',
                  });
                }
              }}
            />
          )}
        />
      </Switch>
    </div>
  );
}

export default withAlert(PolicyOverview);
