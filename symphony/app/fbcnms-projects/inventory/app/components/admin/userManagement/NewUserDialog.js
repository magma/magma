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
import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import FormValidationContext, {
  FormValidationContextProvider,
} from '@fbcnms/ui/components/design-system/Form/FormValidationContext';
import Grid from '@material-ui/core/Grid';
import Strings from '../../../common/CommonStrings';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
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
  const [password, setPassword] = useState('');
  const [passwordVerfication, setPasswordVerification] = useState('');

  return (
    <Dialog fullWidth={true} maxWidth="md" open={true}>
      <FormValidationContextProvider>
        <FormValidationContext.Consumer>
          {validationContext => {
            const passwordMismatch = validationContext.error.check({
              fieldId: 'password match',
              fieldDisplayName: 'password match',
              value: !!passwordVerfication && passwordVerfication !== password,
              checkCallback: mismatch =>
                mismatch ? `${fbt("Passwords doesn't match", '')}` : '',
            });
            return (
              <>
                <DialogTitle
                  disableTypography={true}
                  className={'classes.dialogTitle'}>
                  <Text variant="h6">
                    <fbt desc="">New User Account</fbt>
                  </Text>
                </DialogTitle>
                <DialogContent>
                  <div className={classes.section}>
                    <div className={classes.sectionHeader}>
                      <Text variant="subtitle1">
                        <fbt desc="">Profile</fbt>
                      </Text>
                    </div>
                    <Grid container spacing={2}>
                      <Grid key="first name" item xs={12} sm={6} lg={4} xl={4}>
                        <FormField
                          className={classes.field}
                          label={`${fbt('First Name', '')}`}
                          required={true}
                          validation={{
                            id: 'first name',
                            value: user.firstName,
                          }}>
                          <TextInput
                            value={user.firstName}
                            autoFocus={true}
                            onChange={e => {
                              const newValue = e.target.value;
                              setUser(currentUser => {
                                return Object.assign({}, currentUser, {
                                  firstName: newValue,
                                });
                              });
                            }}
                          />
                        </FormField>
                      </Grid>
                      <Grid key="last name" item xs={12} sm={6} lg={4} xl={4}>
                        <FormField
                          className={classes.field}
                          label={`${fbt('Last Name', '')}`}
                          required={true}
                          validation={{id: 'last name', value: user.lastName}}>
                          <TextInput
                            value={user.lastName}
                            onChange={e => {
                              const newValue = e.target.value;
                              setUser(currentUser => {
                                return Object.assign({}, currentUser, {
                                  lastName: newValue,
                                });
                              });
                            }}
                          />
                        </FormField>
                      </Grid>
                    </Grid>
                  </div>

                  <div className={classes.section}>
                    <UserRoleAndStatusPane
                      value={user.role}
                      onChange={newValue => {
                        if (USER_ROLES[newValue] == null) {
                          return;
                        }
                        setUser(currentUser => {
                          return Object.assign({}, currentUser, {
                            role: USER_ROLES[newValue],
                          });
                        });
                      }}
                    />
                  </div>
                  <div className={classes.section}>
                    <div className={classes.sectionHeader}>
                      <Text variant="subtitle1">
                        <fbt desc="">Account</fbt>
                      </Text>
                    </div>
                    <Grid container spacing={2}>
                      <Grid key="Email" item xs={12} sm={12} lg={4} xl={4}>
                        <FormField
                          className={classes.field}
                          label={`${fbt('Email', '')}`}
                          required={true}
                          validation={{id: 'email', value: user.authId}}>
                          <TextInput
                            value={user.authId}
                            onChange={e => {
                              const newValue = e.target.value;
                              setUser(currentUser => {
                                return Object.assign({}, currentUser, {
                                  authId: newValue,
                                });
                              });
                            }}
                          />
                        </FormField>
                      </Grid>
                      <Grid key="Password" item xs={12} sm={12} lg={4} xl={4}>
                        <FormField
                          className={classes.field}
                          label={`${fbt('Password', '')}`}
                          required={true}
                          validation={{id: 'password', value: password}}>
                          <TextInput
                            type="password"
                            value={password}
                            onChange={e => {
                              const newValue = e.target.value;
                              setPassword(newValue);
                            }}
                          />
                        </FormField>
                      </Grid>
                      <Grid
                        key="PasswordVerification"
                        item
                        xs={12}
                        sm={12}
                        lg={4}
                        xl={4}>
                        <FormField
                          className={classes.field}
                          label={`${fbt('Re-type Password', '')}`}
                          required={true}
                          validation={{
                            id: 'password verification',
                            value: passwordVerfication,
                          }}
                          hasError={!!passwordMismatch}
                          errorText={passwordMismatch}>
                          <TextInput
                            type="password"
                            value={passwordVerfication}
                            onChange={e => {
                              const newValue = e.target.value;
                              setPasswordVerification(newValue);
                            }}
                          />
                        </FormField>
                      </Grid>
                    </Grid>
                  </div>
                </DialogContent>
                <DialogActions>
                  <Button onClick={() => props.onClose(null)}>
                    {Strings.common.cancelButton}
                  </Button>
                  <Button
                    onClick={() => props.onClose(user)}
                    title={validationContext.error.message}
                    disabled={validationContext.error.detected}>
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
