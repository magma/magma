/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {RadioOption} from '@fbcnms/ui/components/design-system/RadioGroup/RadioGroup';
import type {UserRole} from './TempTypes';

import * as React from 'react';
import RadioGroup from '@fbcnms/ui/components/design-system/RadioGroup/RadioGroup';
import Text from '@fbcnms/ui/components/design-system/Text';
import fbt from 'fbt';
import {USER_ROLES} from './TempTypes';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  sectionHeader: {
    marginBottom: '16px',
    '&>span': {
      display: 'block',
    },
  },
}));

type Props = {
  value?: ?UserRole,
  onChange: UserRole => void,
};

const ROLES_OPTIONS: Array<RadioOption> = [
  {
    value: USER_ROLES.User,
    label: fbt('User', ''),
    details: fbt('Can log in to Symphony desktop and mobile apps', ''),
  },
  {
    value: USER_ROLES.Admin,
    label: fbt('Admin', ''),
    details: fbt(
      'Can log in to desktop and mobile apps, update settings and manage users and permissions',
      '',
    ),
  },
  {
    value: USER_ROLES.Owner,
    label: fbt('Owner', ''),
    details: fbt(
      'Full access over everything, including inventory and workforce data',
      '',
    ),
  },
];

const UserRoleAndStatusPane = (props: Props) => {
  const onChange = props.onChange;
  const value = props.value || USER_ROLES.User;
  const classes = useStyles();

  return (
    <>
      <div className={classes.sectionHeader}>
        <Text variant="subtitle1">
          <fbt desc="">User Role and Status</fbt>
        </Text>
        <Text variant="subtitle2" color="gray">
          <fbt desc="">
            Description of what roles are. To view the this user's permissions,
            go to Permissions.
          </fbt>
        </Text>
      </div>
      <RadioGroup
        options={ROLES_OPTIONS}
        value={value}
        onChange={newValue => {
          if (USER_ROLES[newValue] == null) {
            return;
          }
          onChange(USER_ROLES[newValue]);
        }}
      />
    </>
  );
};

export default UserRoleAndStatusPane;
