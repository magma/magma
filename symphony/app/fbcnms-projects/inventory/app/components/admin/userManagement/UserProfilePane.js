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
import FileUploadArea from '@fbcnms/ui/components/design-system/Experimental/FileUpload/FileUploadArea';
import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import FormFieldTextInput from './FormFieldTextInput';
import Grid from '@material-ui/core/Grid';
import Select from '@fbcnms/ui/components/design-system/Select/Select';
import Text from '@fbcnms/ui/components/design-system/Text';
import UserRoleAndStatusPane from './UserRoleAndStatusPane';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import {makeStyles} from '@material-ui/styles';

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
    display: 'flex',
    flexDirection: 'column',
    marginRight: '24px',
    height: '138px',
  },
  fieldsContainer: {
    display: 'flex',
    flexGrow: '1',
  },
  field: {
    marginRight: '8px',
    flexShrink: '1',
    flexBasis: '240px',
  },
  photoUploadContainer: {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    height: '112px',
    width: '112px',
    backgroundColor: symphony.palette.D10,
    border: `1px dashed ${symphony.palette.D100}`,
    '&:hover': {
      borderColor: symphony.palette.D900,
      cursor: 'pointer',
    },
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
  const {user, onChange} = props;
  const classes = useStyles();

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
            <FormField label={`${fbt('Photo', '')}`}>
              <FileUploadArea onFileChanged={files => alert(files[0].name)} />
            </FormField>
          </div>
          <div className={classes.fieldsContainer}>
            <Grid container spacing={2}>
              <Grid item xs={12} sm={6} lg={6} xl={6}>
                <FormFieldTextInput
                  className={classes.field}
                  label={`${fbt('First Name', '')}`}
                  validationId="first name"
                  value={user.firstName}
                  onValueChanged={newName => {
                    user.firstName = newName;
                    onChange(user);
                  }}
                />
              </Grid>
              <Grid item xs={12} sm={6} lg={6} xl={6}>
                <FormFieldTextInput
                  className={classes.field}
                  label={`${fbt('Last Name', '')}`}
                  validationId="last name"
                  value={user.lastName}
                  onValueChanged={newName => {
                    user.lastName = newName;
                    onChange(user);
                  }}
                />
              </Grid>
              <Grid item xs={12} sm={6} lg={6} xl={6}>
                <FormFieldTextInput
                  className={classes.field}
                  label={`${fbt('Phone Number', '')}`}
                  value={user.phoneNumber || ''}
                  onValueChanged={newPhoneNumber => {
                    user.phoneNumber = newPhoneNumber;
                    onChange(user);
                  }}
                />
              </Grid>
            </Grid>
          </div>
        </div>
      </div>
      <UserRoleAndStatusPane
        className={classes.section}
        user={user}
        onChange={onChange}
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
                className={classes.field}
                label={`${fbt('Job Title', '')}`}
                value={user.jobTitle || ''}
                onValueChanged={newValue => {
                  user.jobTitle = newValue;
                  onChange(user);
                }}
              />
            </Grid>
            <Grid item xs={12} sm={6} lg={4} xl={4}>
              <FormFieldTextInput
                className={classes.field}
                label={`${fbt('Employee ID', '')}`}
                value={user.employeeID || ''}
                onValueChanged={newValue => {
                  user.employeeID = newValue;
                  onChange(user);
                }}
              />
            </Grid>
            <Grid item xs={12} sm={6} lg={4} xl={4}>
              <FormField
                className={classes.field}
                label={`${fbt('Employment Type', '')}`}>
                <Select
                  options={EMPLOYMENT_TYPE_OPTIONS}
                  selectedValue={user.employmentType}
                  onChange={newValue => {
                    user.employmentType = newValue;
                    onChange(user);
                  }}
                />
              </FormField>
            </Grid>
          </Grid>
        </div>
      </div>
    </div>
  );
}
