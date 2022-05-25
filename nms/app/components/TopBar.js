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
import typeof SvgIcon from '@material-ui/core/@@SvgIcon';

import AppBar from '@material-ui/core/AppBar';
import Grid from '@material-ui/core/Grid';
// $FlowFixMe migrated to typescript
import NestedRouteLink from './NestedRouteLink';
import React from 'react';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import Text from '../theme/design-system/Text';

import NetworkSelector from './NetworkSelector';
// $FlowFixMe migrated to typescript
import {GetCurrentTabPos} from './TabUtils';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {colors} from '../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useResolvedPath} from 'react-router-dom';

const useStyles = makeStyles(theme => ({
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
  dateTimeText: {
    color: colors.primary.selago,
  },
  versionText: {
    textAlign: 'right',
  },
}));

type BarLabel = {
  icon?: SvgIcon,
  label: string,
  to: string,
  key?: string,
  filters?: React$Node,
};

type Props = {header: string, tabs: BarLabel[]};

export default function TopBar(props: Props) {
  const classes = useStyles();
  const {pathname} = useResolvedPath('');
  const currentTab = GetCurrentTabPos(
    pathname,
    props.tabs.map(tab => tab.to),
  );
  function tabLabel(label, icon) {
    const Icon = icon;

    return (
      <div className={classes.tabLabel}>
        {Icon ? <Icon className={classes.tabIconLabel} /> : null}
        {label}
      </div>
    );
  }
  return (
    <>
      <div className={classes.topBar}>
        <Grid
          container
          direction="row"
          justifyContent="space-between"
          alignItems="center">
          <Grid item xs>
            <Text variant="body2">{props.header}</Text>
          </Grid>
          <NetworkSelector />
        </Grid>
      </div>
      {props.tabs.length > 0 && (
        <AppBar position="static" color="default" className={classes.tabBar}>
          <Grid
            container
            direction="row"
            justifyContent="space-between"
            alignItems="center">
            <Grid item xs>
              <Tabs
                value={currentTab}
                indicatorColor="primary"
                TabIndicatorProps={{style: {height: '5px'}}}
                textColor="inherit"
                className={classes.tabs}>
                {props.tabs.map(tab => (
                  <Tab
                    key={tab.key ?? tab.label}
                    component={NestedRouteLink}
                    label={tabLabel(tab.label, tab.icon)}
                    to={tab.to}
                    className={classes.tab}
                    data-testid={tab.label}
                  />
                ))}
              </Tabs>
            </Grid>
            {props.tabs.map((tab, i) => (
              <React.Fragment key={`fragment-${i}`}>
                {currentTab === i ? (
                  <Grid key={i} item>
                    {tab.filters}
                  </Grid>
                ) : null}
              </React.Fragment>
            ))}
          </Grid>
        </AppBar>
      )}
    </>
  );
}
