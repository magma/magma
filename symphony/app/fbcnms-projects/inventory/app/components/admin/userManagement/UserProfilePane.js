/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {EmploymentType, User} from './TempTypes';
import type {OptionProps} from '@fbcnms/ui/components/design-system/Select/SelectMenu';

import * as React from 'react';
import Button from '@fbcnms/ui/components/design-system/Button';
import FileUploadArea from '@fbcnms/ui/components/design-system/Experimental/FileUpload/FileUploadArea';
import FileUploadButton from '../../FileUpload/FileUploadButton';
import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import FormFieldTextInput from './FormFieldTextInput';
import Grid from '@material-ui/core/Grid';
import Select from '@fbcnms/ui/components/design-system/Select/Select';
import Text from '@fbcnms/ui/components/design-system/Text';
import UserRoleAndStatusPane from './UserRoleAndStatusPane';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import {DocumentAPIUrls} from '../../../common/DocumentAPI';
import {SQUARE_DIMENSION_PX} from '@fbcnms/ui/components/design-system/Experimental/FileUpload/FileUploadArea';
import {makeStyles} from '@material-ui/styles';
import {useEffect, useState} from 'react';
import {useFormContext} from '../../../common/FormContext';

const useStyles = makeStyles(() => ({
  root: {
    display: 'flex',
    flexDirection: 'column',
    height: '100%',
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

const EMPLOYMENT_TYPE_OPTIONS: Array<OptionProps<EmploymentType>> = [
  {
    key: 'FullTime',
    value: 'FullTime',
    label: fbt('Full Time', ''),
  },
  {
    key: 'Contructor',
    value: 'Contructor',
    label: fbt('Contructor', ''),
  },
];

type Props = {
  user: User,
  onChange: User => void,
};

export default function UserProfilePane(props: Props) {
  const {user: propUser, onChange} = props;
  const classes = useStyles();
  const [user, setUser] = useState<?User>(null);
  useEffect(() => setUser(propUser), [propUser]);

  const form = useFormContext();
  const callOnChange = (newUser: User) => {
    setUser(newUser);
    if (form.alerts.error.detected) {
      return;
    }
    onChange(newUser);
  };

  /* temp */
  const [profilePhoto, setProfilePhoto] = useState<?{storeKey: string}>(null);
  useEffect(() => setProfilePhoto(null), [user]);

  if (user == null) {
    return null;
  }

  return (
    <div className={classes.root}>
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
                        src={DocumentAPIUrls.get_url(profilePhoto.storeKey)}
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
                    callOnChange({
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
                    callOnChange({
                      ...user,
                      lastName,
                    })
                  }
                />
              </Grid>
              <Grid item xs={12} sm={6} lg={6} xl={6}>
                <FormFieldTextInput
                  key={`${user.id}_phone`}
                  className={classes.field}
                  label={`${fbt('Phone Number', '')}`}
                  value={user.phoneNumber || ''}
                  onValueChanged={phoneNumber =>
                    callOnChange({
                      ...user,
                      phoneNumber,
                    })
                  }
                />
              </Grid>
            </Grid>
          </div>
        </div>
      </div>
      <UserRoleAndStatusPane
        className={classes.section}
        user={user}
        onChange={callOnChange}
      />
      <div className={classes.section}>
        <div className={classes.sectionHeader}>
          <Text variant="subtitle1">
            <fbt desc="">Employment Information</fbt>
          </Text>
          <Text variant="subtitle2" color="gray">
            <fbt desc="">
              Up-to-date info makes it easier to manage teams and schedule work
              orders.
            </fbt>
          </Text>
        </div>
        <div>
          <Grid container spacing={2}>
            <Grid item xs={12} sm={6} lg={4} xl={4}>
              <FormFieldTextInput
                key={`${user.id}_job`}
                className={classes.field}
                label={`${fbt('Job Title', '')}`}
                value={user.jobTitle || ''}
                onValueChanged={jobTitle =>
                  callOnChange({
                    ...user,
                    jobTitle,
                  })
                }
              />
            </Grid>
            <Grid item xs={12} sm={6} lg={4} xl={4}>
              <FormFieldTextInput
                key={`${user.id}_employee_id`}
                className={classes.field}
                label={`${fbt('Employee ID', '')}`}
                value={user.employeeID || ''}
                onValueChanged={employeeID =>
                  callOnChange({
                    ...user,
                    employeeID,
                  })
                }
              />
            </Grid>
            <Grid item xs={12} sm={6} lg={4} xl={4}>
              <FormField
                className={classes.field}
                label={`${fbt('Employment Type', '')}`}>
                <Select
                  options={EMPLOYMENT_TYPE_OPTIONS}
                  selectedValue={user.employmentType}
                  onChange={employmentType =>
                    callOnChange({
                      ...user,
                      employmentType,
                    })
                  }
                />
              </FormField>
            </Grid>
          </Grid>
        </div>
      </div>
    </div>
  );
}
