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
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useUserManagement} from '../UserManagementContext';

const useStyles = makeStyles(() => ({
  root: {
    flexGrow: 1,
    display: 'flex',
    flexDirection: 'column',
    height: '100%',
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
    enqueueSnackbar(error.response?.data?.error || error, {variant: 'error'});
  };

  return (
    <div className={classes.root}>
      <UserAccountDetailsPane
        variant={
          isForCurrentUserSettings
            ? ACCOUNT_DISPLAY_VARIANTS.userSettingsView
            : ACCOUNT_DISPLAY_VARIANTS.userDetailsCard
        }
        user={user}
        onChange={(user, password, currentPassword) => {
          if (isForCurrentUserSettings && currentPassword != null) {
            userManagement
              .changeCurrentUserPassword(currentPassword, password)
              .catch(handleError);
          } else {
            userManagement
              .changeUserPassword(user, password)
              .catch(handleError);
          }
        }}
      />
    </div>
  );
}
