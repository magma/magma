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

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import ApnEditDialog from './ApnEdit';
import ApnOverview, {ApnJsonConfig} from './ApnOverview';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import BaseNameEditDialog from './BaseNameEdit';
import Button from '@material-ui/core/Button';
import CellWifiIcon from '@material-ui/icons/CellWifi';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import DataPlanEditDialog from './DataPlanEdit';
import DataPlanOverview from './DataPlanOverview';
import LibraryBooksIcon from '@material-ui/icons/LibraryBooks';
// $FlowFixMe migrated to typescript
import MenuButton from '../../components/MenuButton';
import MenuItem from '@material-ui/core/MenuItem';
import PolicyOverview, {PolicyJsonConfig} from './PolicyOverview';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import PolicyRuleEditDialog from './PolicyEdit';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import ProfileEditDialog from './ProfileEdit';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import RatingGroupEditDialog from './RatingGroupEdit';
import React from 'react';
import RssFeedIcon from '@material-ui/icons/RssFeed';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import Text from '../../theme/design-system/Text';
// $FlowFixMe migrated to typescript
import TopBar from '../../components/TopBar';
import {Navigate, Route, Routes} from 'react-router-dom';

/**
 * Button for creation of policies, rating groups, and profiles
 */
function PolicyMenu() {
  const [open, setOpen] = React.useState(false);
  const [baseNameDialog, setBaseNameDialog] = React.useState(false);
  const [profileDialog, setProfileDialog] = React.useState(false);
  const [ratingGroupDialog, setRatingGroupDialog] = React.useState(false);

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
      <MenuButton label="Create New" size="small">
        <MenuItem data-testid="newPolicyMenuItem" onClick={() => setOpen(true)}>
          <Text variant="body2">Policy</Text>
        </MenuItem>
        <MenuItem
          data-testid="newBaseNameMenuItem"
          onClick={() => setBaseNameDialog(true)}>
          <Text variant="body2">Base Name</Text>
        </MenuItem>
        <MenuItem onClick={() => setProfileDialog(true)}>
          <Text variant="body2">Profile</Text>
        </MenuItem>
        <MenuItem
          data-testid="newRatingGroupMenuItem"
          onClick={() => setRatingGroupDialog(true)}>
          <Text variant="body2">Rating Group</Text>
        </MenuItem>
      </MenuButton>
    </div>
  );
}

/**
 * Wrapper for "Create APN" button
 */
function ApnMenu() {
  const [open, setOpen] = React.useState(false);

  return (
    <div>
      <ApnEditDialog open={open} onClose={() => setOpen(false)} />
      <Button
        data-testid="newApnButton"
        variant="contained"
        color="primary"
        size="small"
        onClick={() => setOpen(true)}>
        Create New APN
      </Button>
    </div>
  );
}

/**
 * Wrapper for "Create Data Plan" button
 */
function DataPlanMenu() {
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
        variant="contained"
        color="primary"
        size="small"
        onClick={() => setOpen(true)}>
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
  return (
    <>
      <TopBar
        header="Traffic"
        tabs={[
          {
            label: 'Policies',
            to: 'policy',
            icon: LibraryBooksIcon,
            filters: <PolicyMenu />,
          },
          {
            label: 'APNs',
            to: 'apn',
            icon: RssFeedIcon,
            filters: <ApnMenu />,
          },
          {
            label: 'Data Plans',
            to: 'data_plan',
            icon: CellWifiIcon,
            filters: <DataPlanMenu />,
          },
        ]}
      />

      <Routes>
        <Route path="/policy/:policyId/json" element={<PolicyJsonConfig />} />
        <Route path="/apn/:apnId/json" element={<ApnJsonConfig />} />
        <Route path="/apn/json" element={<ApnJsonConfig />} />
        <Route path="/policy" element={<PolicyOverview />} />
        <Route path="/apn" element={<ApnOverview />} />
        <Route path="/data_plan" element={<DataPlanOverview />} />
        <Route index element={<Navigate to="policy" replace />} />
      </Routes>
    </>
  );
}
