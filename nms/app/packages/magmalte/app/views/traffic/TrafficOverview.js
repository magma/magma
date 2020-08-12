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
import type {apn, policy_rule} from '@fbcnms/magma-api';

import ApnOverview from './ApnOverview';
import AppBar from '@material-ui/core/AppBar';
import Grid from '@material-ui/core/Grid';
import LibraryBooksIcon from '@material-ui/icons/LibraryBooks';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import PolicyOverview from './PolicyOverview';
import React from 'react';
import RssFeedIcon from '@material-ui/icons/RssFeed';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import Text from '@fbcnms/ui/components/design-system/Text';
import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';

import {ApnJsonConfig} from './ApnOverview';
import {GetCurrentTabPos} from '../../components/TabUtils.js';
import {PolicyJsonConfig} from './PolicyOverview';
import {Redirect, Route, Switch} from 'react-router-dom';
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useState} from 'react';
import {useRouter} from '@fbcnms/ui/hooks';

const POLICY_TITLE = 'Policies';
const APN_TITLE = 'APNs';
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
    color: colors.primary.white,
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

export default function TrafficDashboard() {
  const classes = useStyles();
  const {relativePath, relativeUrl, match} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);
  const [policies, setPolicies] = useState<{[string]: policy_rule}>({});
  const [apns, setApns] = useState<{[string]: apn}>({});
  const {isLoading: policyLoading} = useMagmaAPI(
    MagmaV1API.getNetworksByNetworkIdPoliciesRulesViewFull,
    {
      networkId: networkId,
    },
    useCallback(response => {
      setPolicies(response);
    }, []),
  );

  const {isLoading: apnLoading} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdApns,
    {
      networkId: networkId,
    },
    useCallback(response => {
      setApns(response);
    }, []),
  );
  if (policyLoading || apnLoading) {
    return <LoadingFiller />;
  }
  return (
    <>
      <div className={classes.topBar}>
        <Text color="light" weight="medium">
          Traffic
        </Text>
      </div>

      <AppBar position="static" color="default" className={classes.tabBar}>
        <Grid container>
          <Grid item xs={6}>
            <Tabs
              value={GetCurrentTabPos(match.url, ['policy', 'apn'])}
              indicatorColor="primary"
              TabIndicatorProps={{style: {height: '5px'}}}
              textColor="inherit"
              className={classes.tabs}>
              <Tab
                key="Policy"
                component={NestedRouteLink}
                label={<PolicyDashboardTabLabel />}
                to="/policy"
                className={classes.tab}
              />
              <Tab
                key="APNs"
                component={NestedRouteLink}
                label={<APNDashboardTabLabel />}
                to="/apn"
                className={classes.tab}
              />
            </Tabs>
          </Grid>
        </Grid>
      </AppBar>

      <Switch>
        <Route
          path={relativePath('/policy/:policyId/json')}
          render={() => (
            <PolicyJsonConfig
              policies={policies}
              onSave={policy => setPolicies({...policies, [policy.id]: policy})}
            />
          )}
        />
        <Route
          path={relativePath('/apn/:apnId/json')}
          render={() => (
            <ApnJsonConfig
              apns={apns}
              onSave={apn => setApns({...apns, [apn.apn_name]: apn})}
            />
          )}
        />
        <Route
          path={relativePath('/policy')}
          render={() => <PolicyOverview policies={policies} />}
        />
        <Route
          path={relativePath('/apn')}
          render={() => <ApnOverview apns={apns} />}
        />
        <Redirect to={relativeUrl('/policy')} />
      </Switch>
    </>
  );
}

function PolicyDashboardTabLabel() {
  const classes = useStyles();

  return (
    <div className={classes.tabLabel}>
      <LibraryBooksIcon className={classes.tabIconLabel} /> {POLICY_TITLE}
    </div>
  );
}

function APNDashboardTabLabel() {
  const classes = useStyles();

  return (
    <div className={classes.tabLabel}>
      <RssFeedIcon className={classes.tabIconLabel} /> {APN_TITLE}
    </div>
  );
}
