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
import ApnEditDialog from './ApnEdit';
import ApnOverview from './ApnOverview';
import ArrowDropDownIcon from '@material-ui/icons/ArrowDropDown';
import BaseNameEditDialog from './BaseNameEdit';
import Button from '@material-ui/core/Button';
import CellWifiIcon from '@material-ui/icons/CellWifi';
import DataPlanEditDialog from './DataPlanEdit';
import DataPlanOverview from './DataPlanOverview';
import LibraryBooksIcon from '@material-ui/icons/LibraryBooks';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import PolicyOverview from './PolicyOverview';
import PolicyRuleEditDialog from './PolicyEdit';
import ProfileEditDialog from './ProfileEdit';
import RatingGroupEditDialog from './RatingGroupEdit';
import React from 'react';
import RssFeedIcon from '@material-ui/icons/RssFeed';
import Text from '../../theme/design-system/Text';
import TopBar from '../../components/TopBar';

import {ApnJsonConfig} from './ApnOverview';
import {PolicyJsonConfig} from './PolicyOverview';
import {Redirect, Route, Switch} from 'react-router-dom';
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '../../../fbc_js_core/ui/hooks';
import {withStyles} from '@material-ui/core/styles';

const useStyles = makeStyles(_ => ({
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

const StyledMenu = withStyles({
  paper: {
    border: '1px solid #d3d4d5',
  },
})(props => (
  <Menu
    data-testid="policy_menu"
    elevation={0}
    getContentAnchorEl={null}
    anchorOrigin={{
      vertical: 'bottom',
      horizontal: 'center',
    }}
    transformOrigin={{
      vertical: 'top',
      horizontal: 'center',
    }}
    {...props}
  />
));

/**
 * Button for creation of policies, rating groups, and profiles
 */
function PolicyMenu() {
  const classes = useStyles();
  const [anchorEl, setAnchorEl] = React.useState(null);
  const [open, setOpen] = React.useState(false);
  const [baseNameDialog, setBaseNameDialog] = React.useState(false);
  const [profileDialog, setProfileDialog] = React.useState(false);
  const [ratingGroupDialog, setRatingGroupDialog] = React.useState(false);

  const handleClick = event => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  return (
    <div>
      <PolicyRuleEditDialog open={open} onClose={() => setOpen(false)} />
      <BaseNameEditDialog
        open={baseNameDialog}
        onClose={() => setBaseNameDialog(false)}
      />
      <ProfileEditDialog
        open={profileDialog}
        onClose={() => setProfileDialog(false)}
      />
      <RatingGroupEditDialog
        open={ratingGroupDialog}
        onClose={() => setRatingGroupDialog(false)}
      />
      <Button
        onClick={handleClick}
        className={classes.appBarBtn}
        endIcon={<ArrowDropDownIcon />}>
        Create New{' '}
      </Button>
      <StyledMenu
        anchorEl={anchorEl}
        keepMounted
        open={Boolean(anchorEl)}
        onClose={handleClose}>
        <MenuItem data-testid="newPolicyMenuItem" onClick={() => setOpen(true)}>
          <Text variant="subtitle2">Policy</Text>
        </MenuItem>
        <MenuItem
          data-testid="newBaseNameMenuItem"
          onClick={() => setBaseNameDialog(true)}>
          <Text variant="subtitle2">Base Name</Text>
        </MenuItem>
        <MenuItem onClick={() => setProfileDialog(true)}>
          <Text variant="subtitle2">Profile</Text>
        </MenuItem>
        <MenuItem
          data-testid="newRatingGroupMenuItem"
          onClick={() => setRatingGroupDialog(true)}>
          <Text variant="subtitle2">Rating Group</Text>
        </MenuItem>
      </StyledMenu>
    </div>
  );
}

/**
 * Wrapper for "Create APN" button
 */
function ApnMenu() {
  const classes = useStyles();
  const [open, setOpen] = React.useState(false);

  return (
    <div>
      <ApnEditDialog open={open} onClose={() => setOpen(false)} />
      <Button
        data-testid="newApnButton"
        onClick={() => setOpen(true)}
        className={classes.appBarBtn}>
        Create New APN
      </Button>
    </div>
  );
}

/**
 * Wrapper for "Create Data Plan" button
 */
function DataPlanMenu() {
  const classes = useStyles();
  const [open, setOpen] = React.useState(false);

  return (
    <div>
      <DataPlanEditDialog
        open={open}
        onClose={() => setOpen(false)}
        dataPlanId={''}
      />
      <Button
        data-testid="newDataPlanButton"
        onClick={() => setOpen(true)}
        className={classes.appBarBtn}>
        Create New Data Plan
      </Button>
    </div>
  );
}

/**
 * Dashboard for "Traffic" related features.
 *
 * Provides a management interface for:
 *  - policies
 *  - rating groups
 *  - profiles
 *  - APNs
 *  - data plans
 *
 * "Read" and "Edit" functionality provided through tables
 * "Create" functiona provided through a header with "Create New"
 * button.
 */
export default function TrafficDashboard() {
  const {relativePath, relativeUrl} = useRouter();

  return (
    <>
      <TopBar
        header="Traffic"
        tabs={[
          {
            label: 'Policies',
            to: '/policy',
            icon: LibraryBooksIcon,
            filters: <PolicyMenu />,
          },
          {
            label: 'APNs',
            to: '/apn',
            icon: RssFeedIcon,
            filters: <ApnMenu />,
          },
          {
            label: 'Data Plans',
            to: '/data_plan',
            icon: CellWifiIcon,
            filters: <DataPlanMenu />,
          },
        ]}
      />

      <Switch>
        <Route
          path={relativePath('/policy/:policyId/json')}
          component={PolicyJsonConfig}
        />
        <Route
          path={relativePath('/apn/:apnId/json')}
          component={ApnJsonConfig}
        />
        <Route path={relativePath('/apn/json')} component={ApnJsonConfig} />
        <Route path={relativePath('/policy')} component={PolicyOverview} />
        <Route path={relativePath('/apn')} component={ApnOverview} />
        <Route path={relativePath('/data_plan')} component={DataPlanOverview} />
        <Redirect to={relativeUrl('/policy')} />
      </Switch>
    </>
  );
}
