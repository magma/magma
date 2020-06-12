/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {User} from '../utils/UserManagementUtils';

import * as React from 'react';
import UserAccountDetailsPane, {
  ACCOUNT_DISPLAY_VARIANTS,
} from './UserAccountDetailsPane';
import {FormContextProvider} from '../../../../common/FormContext';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useUserManagement} from '../UserManagementContext';

const useStyles = makeStyles(() => ({
  root: {
    flexGrow: 1,
    display: 'flex',
    flexDirection: 'column',
    height: '100%',
    padding: '24px',
  },
}));

type Props = {
  user: User,
  isForCurrentUserSettings?: ?boolean,
};

export default function UserAccountPane(props: Props) {
  const {user, isForCurrentUserSettings = false} = props;
  const classes = useStyles();
  const userManagement = useUserManagement();
  const enqueueSnackbar = useEnqueueSnackbar();

  const handleError = error => {
    enqueueSnackbar(error.response?.data?.error || error.message || error, {
      variant: 'error',
    });
    throw error;
  };

  return (
    <div className={classes.root}>
      <FormContextProvider
        permissions={{ignorePermissions: isForCurrentUserSettings}}>
        <UserAccountDetailsPane
          variant={
            isForCurrentUserSettings
              ? ACCOUNT_DISPLAY_VARIANTS.userSettingsView
              : ACCOUNT_DISPLAY_VARIANTS.userDetailsCard
          }
          user={user}
          onChange={(user, password, currentPassword) => {
            if (isForCurrentUserSettings && currentPassword != null) {
              return userManagement
                .changeCurrentUserPassword(currentPassword, password)
                .catch(handleError);
            }
            return userManagement
              .changeUserPassword(user, password)
              .then(() => undefined)
              .catch(handleError);
          }}
        />
      </FormContextProvider>
    </div>
  );
}
