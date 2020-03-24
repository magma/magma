/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {User} from './TempTypes';

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
import {USER_ROLES, USER_STATUSES} from './TempTypes';
import {generateTempId} from '../../../common/EntUtils';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useState} from 'react';
import {useUserManagement} from './UserManagementContext';

const initialUserData: User = {
  id: generateTempId(),
  authID: '',
  firstName: '',
  lastName: '',
  role: USER_ROLES.USER.key,
  status: USER_STATUSES.ACTIVE.key,
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

const NewUserDialog = ({onClose}: Props) => {
  const classes = useStyles();
  const userManagegemtContext = useUserManagement();
  const [creatingUser, setCreatingUser] = useState(false);
  const [user, setUser] = useState<User>({...initialUserData});
  const [password, setPassword] = useState('');

  const enqueueSnackbar = useEnqueueSnackbar();
  const handleError = error => {
    setCreatingUser(false);
    enqueueSnackbar(error.response?.data?.error || error, {variant: 'error'});
  };

  const addUser = () => {
    setCreatingUser(true);
    userManagegemtContext
      .addUser(user, password)
      .finally(() => setCreatingUser(false))
      .then(newUser => {
        onClose(newUser);
      })
      .catch(handleError);
  };

  return (
    <Dialog fullWidth={true} maxWidth="md" open={true}>
      <FormValidationContextProvider>
        <FormValidationContext.Consumer>
          {formValidationContext => {
            formValidationContext.editLock.check({
              fieldId: 'async_save',
              fieldDisplayName: 'Lock while saving',
              value: creatingUser,
              checkCallback: isOnSavingProcess =>
                isOnSavingProcess == true ? 'Saving new user' : '',
            });
            return (
              <>
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
                            currentUser.role = newValue;
                            return currentUser;
                          });
                        },
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
                  <Button onClick={onClose} disabled={creatingUser}>
                    {Strings.common.cancelButton}
                  </Button>
                  <Button
                    onClick={addUser}
                    title={formValidationContext.error.message}
                    disabled={
                      formValidationContext.error.detected || creatingUser
                    }>
                    {Strings.common.saveButton}
                  </Button>
                </DialogActions>
              </>
            );
          }}
        </FormValidationContext.Consumer>
      </FormValidationContextProvider>
    </Dialog>
  );
};

export default NewUserDialog;
