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
import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import UserRoleAndStatusPane from './UserRoleAndStatusPane';
import fbt from 'fbt';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  root: {
    display: 'flex',
    flexDirection: 'column',
    height: '100%',
  },
  personalDetails: {
    display: 'flex',
    marginBottom: '16px',
  },
  photoContainer: {
    display: 'flex',
    flexDirection: 'column',
    marginRight: '16px',
    height: '138px',
    width: '112px',
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
}));

type Props = {
  user: User,
  onChange?: ?(User) => void,
};

export default function UserProfilePane(props: Props) {
  const {user, onChange} = props;
  const classes = useStyles();

  return (
    <div className={classes.root}>
      <div className={classes.personalDetails}>
        <div className={classes.photoContainer}>
          <FormField label={`${fbt('Photo', '')}`}>photo placeholder</FormField>
        </div>
        <div className={classes.fieldsContainer}>
          <FormField
            label={`${fbt('First Name', '')}`}
            required={true}
            className={classes.field}>
            <TextInput value={user.firstName} />
          </FormField>
          <FormField
            label={`${fbt('Last Name', '')}`}
            required={true}
            className={classes.field}>
            <TextInput value={user.firstName} />
          </FormField>
        </div>
      </div>
      <UserRoleAndStatusPane
        value={user.role}
        onChange={newRole => {
          user.role = newRole;
          if (onChange) {
            onChange(user);
          }
        }}
      />
    </div>
  );
}
