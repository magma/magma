/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {User, UserRole} from './TempTypes';

import * as React from 'react';
import Button from '@fbcnms/ui/components/design-system/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormFieldTextInput from './FormFieldTextInput';
import FormValidationContext, {
  FormValidationContextProvider,
} from '@fbcnms/ui/components/design-system/Form/FormValidationContext';
import Grid from '@material-ui/core/Grid';
import Strings from '../../../common/CommonStrings';
import Text from '@fbcnms/ui/components/design-system/Text';
import UserAccountDetailsPane, {
  ACCOUNT_DISPLAY_VARIANTS,
} from './UserAccountDetailsPane';
import UserRoleAndStatusPane from './UserRoleAndStatusPane';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import {makeStyles} from '@material-ui/styles';
import {useState} from 'react';

export const USER_ROLES = {
  User: 'User',
  Admin: 'Admin',
  Owner: 'Owner',
};

export const USER_STATUSES = {
  Active: 'Active',
  Deactivated: 'Deactivated',
  Deleted: 'Deleted',
};

const initialUserData: User = {
  authId: '',
  firstName: '',
  lastName: '',
  role: USER_ROLES.User,
  status: USER_STATUSES.Active,
};

const useStyles = makeStyles(() => ({
  field: {},
  section: {
    '&:not(:last-child)': {
      paddingBottom: '16px',
      borderBottom: `1px solid ${symphony.palette.separator}`,
    },
    marginBottom: '16px',
  },
  sectionHeader: {
    marginBottom: '16px',
    '&>span': {
      display: 'block',
    },
  },
}));

type Props = {
  onClose: (?User) => void,
};

const NewUserDialog = (props: Props) => {
  const classes = useStyles();
  const [user, setUser] = useState<User>(initialUserData);
  const [_password, setPassword] = useState('');

  return (
    <Dialog fullWidth={true} maxWidth="md" open={true}>
      <FormValidationContextProvider>
        <DialogTitle disableTypography={true}>
          <Text variant="h6">
            <fbt desc="">New User Account</fbt>
          </Text>
        </DialogTitle>
        <DialogContent>
          <div className={classes.section}>
            <div className={classes.sectionHeader}>
              <Text variant="subtitle1">
                <fbt desc="">Personal Details</fbt>
              </Text>
            </div>
            <Grid container spacing={2}>
              <Grid key="first_name" item xs={12} sm={6} lg={4} xl={4}>
                <FormFieldTextInput
                  validationId="first_name"
                  label={`${fbt('First Name', '')}`}
                  value={user.firstName || ''}
                  onValueChanged={newValue =>
                    setUser(currentUser => {
                      currentUser.firstName = newValue;
                      return currentUser;
                    })
                  }
                />
              </Grid>
              <Grid key="last_name" item xs={12} sm={6} lg={4} xl={4}>
                <FormFieldTextInput
                  validationId="last_name"
                  label={`${fbt('Last Name', '')}`}
                  value={user.lastName || ''}
                  onValueChanged={newValue =>
                    setUser(currentUser => {
                      currentUser.lastName = newValue;
                      return currentUser;
                    })
                  }
                />
              </Grid>
            </Grid>
          </div>
          <div className={classes.section}>
            <UserRoleAndStatusPane
              role={{
                value: user.role,
                onChange: newValue => {
                  if (USER_ROLES[newValue] == null) {
                    return;
                  }
                  setUser(currentUser => {
                    currentUser.role = USER_ROLES[newValue];
                    return currentUser;
                  });
                },
              }}
              onChange={(newValue: UserRole) => {
                if (USER_ROLES[newValue] == null) {
                  return;
                }
                setUser(currentUser => {
                  currentUser.role = USER_ROLES[newValue];
                  return currentUser;
                });
              }}
            />
          </div>
          <UserAccountDetailsPane
            variant={ACCOUNT_DISPLAY_VARIANTS.newUserDialog}
            className={classes.section}
            user={user}
            onChange={(updatedUser, updatedPassword) => {
              setUser(updatedUser);
              setPassword(updatedPassword);
            }}
          />
        </DialogContent>
        <DialogActions>
          <FormValidationContext.Consumer>
            {formValidationContext => (
              <>
                <Button onClick={() => props.onClose(null)}>
                  {Strings.common.cancelButton}
                </Button>
                <Button
                  onClick={() => props.onClose(user)}
                  title={formValidationContext.error.message}
                  disabled={formValidationContext.error.detected}>
                  {Strings.common.saveButton}
                </Button>
              </>
            )}
          </FormValidationContext.Consumer>
        </DialogActions>
      </FormValidationContextProvider>
    </Dialog>
  );
};

export default NewUserDialog;
