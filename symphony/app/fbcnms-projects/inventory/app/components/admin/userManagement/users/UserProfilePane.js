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
import AppContext from '@fbcnms/ui/context/AppContext';
import Button from '@fbcnms/ui/components/design-system/Button';
import FileUploadArea from '@fbcnms/ui/components/design-system/Experimental/FileUpload/FileUploadArea';
import FileUploadButton from '../../../FileUpload/FileUploadButton';
import FormAction from '@fbcnms/ui/components/design-system/Form/FormAction';
import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import FormFieldTextInput from '../utils/FormFieldTextInput';
import Grid from '@material-ui/core/Grid';
import Text from '@fbcnms/ui/components/design-system/Text';
import UserRoleAndStatusPane from './UserRoleAndStatusPane';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import {DocumentAPIUrls} from '../../../../common/DocumentAPI';
import {FormContextProvider} from '../../../../common/FormContext';
import {GROUP_STATUSES, USER_ROLES} from '../utils/UserManagementUtils';
import {Prompt} from 'react-router-dom';
import {SQUARE_DIMENSION_PX} from '@fbcnms/ui/components/design-system/Experimental/FileUpload/FileUploadArea';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useContext, useEffect, useMemo, useState} from 'react';
import {useMessageShowingContext} from '@fbcnms/ui/components/design-system/Dialog/MessageShowingContext';

const useStyles = makeStyles(() => ({
  root: {
    flexGrow: 1,
    display: 'flex',
    flexDirection: 'column',
  },
  form: {
    flexBasis: 'auto',
    flexGrow: '1',
    flexShrink: '0',
    height: '100px',
    overflowY: 'auto',
    overflowX: 'hidden',
    padding: '24px',
    paddingBottom: '0',
    display: 'flex',
    flexDirection: 'column',
  },
  bottomBar: {
    flexGrow: 0,
    padding: '8px 24px',
    borderTop: `1px solid ${symphony.palette.separator}`,
    display: 'flex',
    justifyContent: 'flex-end',
    '&>*': {
      marginLeft: '8px',
    },
  },
  section: {
    display: 'flex',
    flexDirection: 'column',
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
  personalDetails: {
    display: 'flex',
    marginBottom: '16px',
  },
  photoContainer: {
    flexBasis: SQUARE_DIMENSION_PX,
    flexGrow: 0,
    flexShrink: 0,
    overflow: 'hidden',
    display: 'flex',
    flexDirection: 'column',
    marginRight: '24px',
    height: '138px',
  },
  photoWrapper: {
    maxHeight: '100%',
    position: 'relative',
  },
  photoArea: {
    display: 'flex',
    flexDirection: 'column',
    '&:hover $photoRemoveButton': {
      display: 'unset',
    },
  },
  photo: {
    flexBasis: SQUARE_DIMENSION_PX,
    flexGrow: 1,
    overflow: 'hidden',
    width: '100%',
    objectFit: 'cover',
    objectPosition: '50% 50%',
    borderRadius: '4px',
  },
  photoRemoveButton: {
    display: 'none',
    position: 'absolute',
    bottom: '8px',
    left: '50%',
    transform: 'translateX(-50%)',
    opacity: '0.9',
  },
  fieldsContainer: {
    display: 'flex',
    flexGrow: '1',
  },
  field: {
    marginRight: '8px',
    maxHeight: '58px',
    flexShrink: '1',
    flexBasis: '240px',
  },
}));

// const EMPLOYMENT_TYPE_OPTIONS: Array<OptionProps<EmploymentType>> = [
//   {
//     key: 'FullTime',
//     value: 'FullTime',
//     label: fbt('Full Time', ''),
//   },
//   {
//     key: 'Contructor',
//     value: 'Contructor',
//     label: fbt('Contructor', ''),
//   },
// ];

type Props = {
  user: User,
  onChange: User => void,
};

function isSameUserDetails(userA: ?User, userB: ?User) {
  if (userA == null && userB == null) {
    return true;
  }
  if (userA == null || userB == null) {
    return false;
  }
  const aData: string = JSON.stringify(userA);
  const bData: string = JSON.stringify(userB);

  return aData === bData;
}

export default function UserProfilePane(props: Props) {
  const {user: propUser, onChange} = props;
  const classes = useStyles();
  const {isFeatureEnabled} = useContext(AppContext);
  const userManagementDevMode = isFeatureEnabled('user_management_dev');

  const [
    shouldShowVerificationWhenDeactivating,
    setShowVerificationWhenDeactivating,
  ] = useState(true);
  const [
    shouldShowVerificationWhenChangingRole,
    setShowVerificationWhenChangingRole,
  ] = useState(true);

  const [user, setUser] = useState<?User>(null);
  useEffect(() => setUser(propUser), [propUser]);

  const messageShowingContext = useMessageShowingContext();

  const callOnChange = useCallback(() => {
    if (user == null) {
      return;
    }
    onChange(user);
  }, [onChange, user]);

  const revertChanges = useCallback(() => {
    setUser(propUser);
  }, [propUser]);

  /* temp */
  const [profilePhoto, setProfilePhoto] = useState<?{storeKey: string}>(null);
  useEffect(() => setProfilePhoto(null), [user]);

  const userDataChanged = !isSameUserDetails(user, propUser);
  const formButtonsDisablingProps = useMemo(
    () =>
      userDataChanged
        ? null
        : {
            disabled: true,
            tooltip: `${fbt('No changes were made', '')}`,
          },
    [userDataChanged],
  );

  if (user == null) {
    return null;
  }

  return (
    <div className={classes.root}>
      <Prompt
        when={userDataChanged}
        message={`${fbt(
          'You have unsaved data, leaving this page will ignore those changes.',
          '',
        )}`}
      />
      <FormContextProvider permissions={{adminRightsRequired: true}}>
        <div className={classes.form}>
          <div className={classes.section}>
            <div className={classes.sectionHeader}>
              <Text variant="subtitle1">
                <fbt desc="">Personal Details</fbt>
              </Text>
              <Text variant="subtitle2" color="gray">
                <fbt desc="">
                  These details are used when assigning work orders and granting
                  permissions.
                </fbt>
              </Text>
            </div>
            <div className={classes.personalDetails}>
              {userManagementDevMode ? (
                <div className={classes.photoContainer}>
                  <FormField
                    label={`${fbt('Photo', '')}`}
                    className={classes.photoWrapper}>
                    <FileUploadButton
                      useUploadSnackbar={false}
                      multiple={false}
                      onFileUploaded={
                        (_file, storeKey) => setProfilePhoto({storeKey})
                        /*
                    onChange({
                      ...user,
                      profilePhoto: {
                        id: '',
                        storeKey,
                        fileName: file.name,
                      },
                    })
                    */
                      }>
                      {openFileUploadDialog =>
                        profilePhoto?.storeKey != null ? (
                          <div className={classes.photoArea}>
                            <img
                              src={DocumentAPIUrls.get_url(
                                profilePhoto.storeKey,
                              )}
                              className={classes.photo}
                            />
                            <Button
                              className={classes.photoRemoveButton}
                              skin="regular"
                              onClick={() => setProfilePhoto(null)}>
                              <fbt desc="">Remove</fbt>
                            </Button>
                          </div>
                        ) : (
                          <FileUploadArea onClick={openFileUploadDialog} />
                        )
                      }
                    </FileUploadButton>
                  </FormField>
                </div>
              ) : null}
              <div className={classes.fieldsContainer}>
                <Grid container spacing={2}>
                  <Grid item xs={12} sm={6} lg={6} xl={6}>
                    <FormFieldTextInput
                      key={`${user.id}_first_name`}
                      className={classes.field}
                      label={`${fbt('First Name', '')}`}
                      validationId="first name"
                      value={user.firstName}
                      onValueChanged={firstName =>
                        setUser({
                          ...user,
                          firstName,
                        })
                      }
                    />
                  </Grid>
                  <Grid item xs={12} sm={6} lg={6} xl={6}>
                    <FormFieldTextInput
                      key={`${user.id}_last_name`}
                      className={classes.field}
                      label={`${fbt('Last Name', '')}`}
                      validationId="last name"
                      value={user.lastName}
                      onValueChanged={lastName =>
                        setUser({
                          ...user,
                          lastName,
                        })
                      }
                    />
                  </Grid>
                  {/* {userManagementDevMode ? (
                    <Grid item xs={12} sm={6} lg={6} xl={6}>
                      <FormFieldTextInput
                        key={`${user.id}_phone`}
                        className={classes.field}
                        label={`${fbt('Phone Number', '')}`}
                        value={user.phoneNumber || ''}
                        onValueChanged={phoneNumber =>
                          setUser({
                            ...user,
                            phoneNumber,
                          })
                        }
                      />
                    </Grid>
                  ) : null} */}
                </Grid>
              </div>
            </div>
          </div>
          <UserRoleAndStatusPane
            className={classes.section}
            user={user}
            onChange={newUser => {
              if (
                newUser.status !== user?.status &&
                newUser.status === GROUP_STATUSES.DEACTIVATED.key &&
                shouldShowVerificationWhenDeactivating
              ) {
                messageShowingContext.showMessage({
                  title: fbt('Deactivate Account', ''),
                  message: (
                    <>
                      <div>
                        <fbt desc="">
                          Are you sure you want to deactivate this account?
                        </fbt>
                      </div>
                      <div>
                        <fbt desc="">
                          All access and permissions for this user will be
                          disabled until this account is reactivated.
                        </fbt>
                      </div>
                    </>
                  ),
                  verificationCheckbox: {
                    label: fbt("Don't show this again", ''),
                  },
                  confirmLabel: fbt('Deactivate', ''),
                  skin: 'red',
                  onCancel: messageShowingContext.hideMessage,
                  onClose: messageShowingContext.hideMessage,
                  onConfirm: dontShowAgain => {
                    setUser(newUser);
                    setShowVerificationWhenDeactivating(dontShowAgain !== true);
                    messageShowingContext.hideMessage();
                  },
                });
              } else if (
                newUser.role !== user?.role &&
                shouldShowVerificationWhenChangingRole
              ) {
                const useAn =
                  newUser.role == USER_ROLES.ADMIN.key ||
                  newUser.role == USER_ROLES.OWNER.key;
                messageShowingContext.showMessage({
                  title: (
                    <fbt desc="">
                      Change to
                      <fbt:param name="new role type">
                        {USER_ROLES[newUser.role].value}
                      </fbt:param>
                    </fbt>
                  ),
                  message: (
                    <fbt desc="">
                      Are you sure you want to make this user
                      <fbt:enum
                        enum-range={['a', 'an']}
                        value={useAn ? 'an' : 'a'}
                      />
                      <fbt:param name="new role type">
                        {` ${USER_ROLES[newUser.role].value}`}
                      </fbt:param>
                      ?
                    </fbt>
                  ),
                  verificationCheckbox: {
                    label: fbt("Don't show this again", ''),
                  },
                  confirmLabel: fbt('Change Role', ''),
                  onCancel: messageShowingContext.hideMessage,
                  onClose: messageShowingContext.hideMessage,
                  onConfirm: dontShowAgain => {
                    setUser(newUser);
                    setShowVerificationWhenChangingRole(dontShowAgain !== true);
                    messageShowingContext.hideMessage();
                  },
                });
              } else {
                setUser(newUser);
              }
            }}
          />
        </div>
        <div className={classes.bottomBar}>
          <FormAction {...formButtonsDisablingProps}>
            <Button skin="regular" onClick={revertChanges}>
              <fbt desc="">Revert</fbt>
            </Button>
          </FormAction>
          <FormAction disableOnFromError={true} {...formButtonsDisablingProps}>
            <Button onClick={callOnChange}>
              <fbt desc="">Apply</fbt>
            </Button>
          </FormAction>
        </div>
      </FormContextProvider>
    </div>
  );
}
