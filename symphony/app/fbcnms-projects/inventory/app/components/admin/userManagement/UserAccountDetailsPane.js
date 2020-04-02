/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {User} from './UserManagementUtils';

import * as React from 'react';
import Button from '@fbcnms/ui/components/design-system/Button';
import FormFieldTextInput from './FormFieldTextInput';
import Grid from '@material-ui/core/Grid';
import Strings from '../../../common/CommonStrings';
import Text from '@fbcnms/ui/components/design-system/Text';
import fbt from 'fbt';
import {makeStyles} from '@material-ui/styles';
import {useEffect, useState} from 'react';
import {useFormContext} from '../../../common/FormContext';

const useStyles = makeStyles(() => ({
  sectionHeader: {
    marginBottom: '16px',
    '&>span': {
      display: 'block',
    },
  },
  sectionBody: {
    '& > div:not($actionButtons)': {
      marginBottom: '0px',
      marginTop: '0px',
    },
  },
  actionButtons: {
    marginTop: '24px',
    '& > button': {
      marginRight: '8px',
    },
  },
}));

export const ACCOUNT_DISPLAY_VARIANTS = {
  newUserDialog: 'newUserDialog',
  userDetailsCard: 'userDetailsCard',
};

type Props = {
  user: User,
  onChange?: (user: User, password: string) => void,
  variant: $Values<typeof ACCOUNT_DISPLAY_VARIANTS>,
  className?: ?string,
};

const UserAccountDetailsPane = (props: Props) => {
  const {user, onChange, className, variant} = props;
  const classes = useStyles();

  const [password, setPassword] = useState<string>('');
  const [passwordVerfication, setPasswordVerification] = useState<string>('');

  const [isEditable, setIsEditable] = useState(
    variant === ACCOUNT_DISPLAY_VARIANTS.newUserDialog,
  );

  const form = useFormContext();
  const passwordRules = form.alerts.error.check({
    fieldId: 'password_rules',
    fieldDisplayName: 'password rules',
    value: password,
    checkCallback: enteredPassword =>
      enteredPassword == null || enteredPassword.length < 10
        ? `${fbt('Password must contain at least 10 characters', '')}`
        : '',
  });
  const passwordMismatch = form.alerts.error.check({
    fieldId: 'password_match',
    fieldDisplayName: 'password match',
    value: !!passwordVerfication && passwordVerfication !== password,
    checkCallback: mismatch =>
      mismatch ? `${fbt("Passwords don't match", '')}` : '',
  });

  useEffect(() => {
    if (
      variant != ACCOUNT_DISPLAY_VARIANTS.newUserDialog ||
      onChange == null ||
      form.alerts.error.detected
    ) {
      return;
    }
    onChange(user, password);
  }, [
    form.alerts.error.detected,
    isEditable,
    onChange,
    password,
    passwordVerfication,
    user,
    variant,
  ]);

  const exitEditMode = () => {
    setIsEditable(false);
    setPasswordVerification('');
    setPassword('');
  };

  const emailField = (
    <FormFieldTextInput
      disabled={variant === ACCOUNT_DISPLAY_VARIANTS.userDetailsCard}
      validationId={
        variant !== ACCOUNT_DISPLAY_VARIANTS.userDetailsCard
          ? 'email'
          : undefined
      }
      label={`${fbt('Email', '')}`}
      value={user.authID}
      onValueChanged={
        onChange == null
          ? undefined
          : newAuthID =>
              onChange(
                {
                  ...user,
                  authID: newAuthID,
                },
                password,
              )
      }
    />
  );
  const passwordField = (
    <FormFieldTextInput
      type="password"
      disabled={!isEditable}
      validationId={isEditable ? 'password' : undefined}
      label={`${fbt('Password', '')}`}
      value={isEditable ? password : '**********'}
      onValueChanged={setPassword}
      hasError={isEditable && !!passwordRules}
      errorText={isEditable ? passwordRules : ''}
      immediateUpdate={true}
    />
  );
  const passwordVerficationField = (
    <FormFieldTextInput
      type="password"
      validationId="password verfication"
      label={`${fbt('Re-type Password', '')}`}
      value={passwordVerfication}
      onValueChanged={setPasswordVerification}
      hasError={!!passwordMismatch}
      errorText={passwordMismatch}
      immediateUpdate={true}
    />
  );
  return (
    <div className={className}>
      <div className={classes.sectionHeader}>
        <Text variant="subtitle1">
          <fbt desc="">Account Details</fbt>
        </Text>
        <Text variant="subtitle2" color="gray">
          <fbt desc="">
            This email will be used to log in to
            <fbt:param name="product name">
              {Strings.common.productName}
            </fbt:param>.
          </fbt>
        </Text>
      </div>
      <div className={classes.sectionBody}>
        {variant === ACCOUNT_DISPLAY_VARIANTS.userDetailsCard ? (
          <>
            <Grid container spacing={2}>
              <Grid key="email" item xs={12} sm={6} lg={4} xl={4}>
                {emailField}
              </Grid>
            </Grid>
            <Grid container spacing={2}>
              <Grid key="password" item xs={12} sm={6} lg={4} xl={4}>
                {passwordField}
              </Grid>
              {isEditable && (
                <Grid
                  key="password_verfication"
                  item
                  xs={12}
                  sm={6}
                  lg={4}
                  xl={4}>
                  {passwordVerficationField}
                </Grid>
              )}
            </Grid>
            <div className={classes.actionButtons}>
              {isEditable ? (
                <>
                  <Button skin="gray" onClick={exitEditMode}>
                    {Strings.common.cancelButton}
                  </Button>
                  <Button
                    onClick={() => {
                      if (onChange) {
                        onChange(user, password);
                      }
                      exitEditMode();
                    }}
                    disabled={form.alerts.error.detected}
                    title={form.alerts.error.message}>
                    <fbt desc="">Save Changes</fbt>
                  </Button>
                </>
              ) : (
                <Button onClick={() => setIsEditable(true)}>
                  <fbt desc="">Change Password</fbt>
                </Button>
              )}
            </div>
          </>
        ) : (
          <Grid container spacing={2}>
            <Grid key="email" item xs={12} sm={6} lg={4} xl={4}>
              {emailField}
            </Grid>
            <Grid key="password" item xs={12} sm={6} lg={4} xl={4}>
              {passwordField}
            </Grid>
            <Grid key="password_verfication" item xs={12} sm={6} lg={4} xl={4}>
              {passwordVerficationField}
            </Grid>
          </Grid>
        )}
      </div>
    </div>
  );
};

export default UserAccountDetailsPane;
