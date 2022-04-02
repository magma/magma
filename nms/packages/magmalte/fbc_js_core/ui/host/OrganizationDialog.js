/**
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

import type {OrganizationPlainAttributes} from '../../../fbc_js_core/sequelize_models/models/organization';

import Button from '../../../fbc_js_core/ui/components/design-system/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogTitle from '@material-ui/core/DialogTitle';
import LoadingFillerBackdrop from '../../../fbc_js_core/ui/components/LoadingFillerBackdrop';
import OrganizationInfoDialog from './OrganizationInfoDialog';
import OrganizationUserDialog from './OrganizationUserDialog';
import React from 'react';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';

import {UserRoles} from '../../../fbc_js_core/auth/types';
import {brightGray, white} from '../../../fbc_js_core/ui/theme/colors';
import {makeStyles} from '@material-ui/styles';
import {useAxios} from '../../../fbc_js_core/ui/hooks';
import {useState} from 'react';

const useStyles = makeStyles(_ => ({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
  tabBar: {
    backgroundColor: brightGray,
    color: white,
  },
}));

export type DialogProps = {
  error: string,
  user: CreateUserType,
  organization: OrganizationPlainAttributes,
  onUserChange: CreateUserType => void,
  onOrganizationChange: OrganizationPlainAttributes => void,
  // Array of networks ids
  allNetworks: Array<string>,
  // If true, enable all networks for an organization
  shouldEnableAllNetworks: boolean,
  setShouldEnableAllNetworks: boolean => void,
};

type Props = {
  onClose: () => void,
  onCreateOrg: (org: CreateOrgType) => void,
  onCreateUser: (user: CreateUserType) => void,
  // flag to display create user tab
  addUser: boolean,
  setAddUser: () => void,
};

type CreateUserType = {
  email: string,
  id?: number,
  networkIDs: Array<string>,
  organization?: string,
  role: ?string,
  tabs?: Array<string>,
  password: string,
  passwordConfirmation?: string,
};

type CreateOrgType = {
  name: string,
  networkIDs: Array<string>,
  customDomains?: Array<string>,
};
/**
 * Create Orgnization Dilaog
 * This component displays a dialog with 2 tabs
 * First tab: OrganizationInfoDialog, to create a new organization
 * Second tab: OrganizationUserDialog, to create a user that belongs to the new organization
 */
export default function (props: Props) {
  const classes = useStyles();
  const {error, isLoading, response} = useAxios({
    method: 'get',
    url: '/host/networks/async',
  });

  const [organization, setOrganization] = useState<OrganizationPlainAttributes>(
    {},
  );
  const [currentTab, setCurrentTab] = useState(props.addUser ? 1 : 0);
  const [shouldEnableAllNetworks, setShouldEnableAllNetworks] = useState(false);
  const [user, setUser] = useState<CreateUserType>({});
  const [createError, setCreateError] = useState('');
  const allNetworks = error || !response ? [] : response.data.sort();

  if (isLoading) {
    return <LoadingFillerBackdrop />;
  }

  const createProps = {
    user,
    organization,
    error: createError,
    onUserChange: (user: CreateUserType) => {
      setUser(user);
    },
    onOrganizationChange: (organization: OrganizationPlainAttributes) => {
      setOrganization(organization);
    },
    allNetworks,
    shouldEnableAllNetworks,
    setShouldEnableAllNetworks,
  };
  const onSave = async () => {
    if (currentTab === 0) {
      if (!organization.name) {
        setCreateError('Name cannot be empty');
        return;
      }
      const payload = {
        name: organization.name,
        networkIDs: shouldEnableAllNetworks
          ? allNetworks
          : Array.from(organization.networkIDs || []),
        customDomains: [], // TODO
        // tabs: Array.from(tabs),
      };

      props.onCreateOrg(payload);
      setCurrentTab(currentTab + 1);
      setCreateError('');
      props.setAddUser();
    } else {
      if (!user.email) {
        setCreateError('Email cannot be empty');
        return;
      }

      if (!user.password) {
        setCreateError('Password cannot be empty');
        return;
      }
      if (user.password != user.passwordConfirmation) {
        setCreateError('Passwords must match');
        return;
      }

      const payload: CreateUserType = {
        email: user.email,
        password: user.password,
        role: user.role,
        networkIDs:
          user.role === UserRoles.SUPERUSER
            ? []
            : Array.from(user.networkIDs || []),
      };
      props.onCreateUser(payload);
    }
  };

  return (
    <Dialog
      open={true}
      onClose={props.onClose}
      maxWidth={'sm'}
      fullWidth={true}>
      <DialogTitle>Add Organization</DialogTitle>
      <Tabs
        indicatorColor="primary"
        value={currentTab}
        className={classes.tabBar}
        onChange={(_, v) => setCurrentTab(v)}>
        <Tab disabled={currentTab === 1} label={'Organization'} />
        <Tab disabled={currentTab === 0} label={'Users'} />
      </Tabs>
      <>
        {currentTab === 0 && <OrganizationInfoDialog {...createProps} />}
        {currentTab === 1 && <OrganizationUserDialog {...createProps} />}
      </>
      <DialogActions>
        <Button onClick={props.onClose} skin="regular">
          Cancel
        </Button>
        <Button skin="comet" onClick={onSave}>
          {currentTab === 0 ? 'Save and Continue' : 'Save'}
        </Button>
      </DialogActions>
    </Dialog>
  );
}
