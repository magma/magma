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
import type {
  UserRole,
  UserStatus,
} from './__generated__/UsersView_UsersQuery.graphql';

import * as React from 'react';
import RadioGroup from '@fbcnms/ui/components/design-system/RadioGroup/RadioGroup';
import Text from '@fbcnms/ui/components/design-system/Text';
import classNames from 'classnames';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import {USER_ROLES, USER_STATUSES} from './TempTypes';
import {makeStyles} from '@material-ui/styles';
import {useMemo} from 'react';

const useStyles = makeStyles(() => ({
  sectionHeader: {
    marginBottom: '16px',
    '&>span': {
      display: 'block',
    },
  },
  deactivateOptionSelected: {
    color: symphony.palette.R600,
    '& svg': {
      color: symphony.palette.R600,
    },
  },
}));

type PropValue<T> = {
  value?: ?T,
  onChange: T => void,
};

type Props = {
  role: PropValue<UserRole>,
  status?: PropValue<UserStatus>,
  className?: ?string,
};

const ROLES_OPTIONS: Array<RadioOption> = [
  {
    value: USER_ROLES.USER.key,
    label: USER_ROLES.USER.value,
    details: fbt('Can log in to Symphony desktop and mobile apps', ''),
  },
  {
    value: USER_ROLES.ADMIN.key,
    label: USER_ROLES.ADMIN.value,
    label: fbt('Admin', ''),
    details: fbt(
      'Can log in to desktop and mobile apps, update settings and manage users and permissions',
      '',
    ),
  },
  {
    value: USER_ROLES.OWNER.key,
    label: USER_ROLES.OWNER.value,
    label: fbt('Owner', ''),
    details: fbt(
      'Full access over everything, including inventory and workforce data',
      '',
    ),
  },
];
const STATUS_OPT: RadioOption = {
  value: USER_STATUSES.DEACTIVATED.key,
  label: fbt('Deactivate', ''),
  details: fbt(
    'Temporarely remove all access and permissions for this user',
    '',
  ),
};

const UserRoleAndStatusPane = (props: Props) => {
  const userIsDeactivated =
    props.status?.value === USER_STATUSES.DEACTIVATED.key;
  const classes = useStyles();
  const selectedOptionClass = classNames({
    [classes.deactivateOptionSelected]: userIsDeactivated,
  });
  const handleStatus = props.status != null;
  const options = useMemo(() => {
    if (!handleStatus) {
      return ROLES_OPTIONS;
    }
    const {label, ...partialDeactivateOption} = STATUS_OPT;
    const deactivateOption = {
      ...partialDeactivateOption,
      label: <div className={selectedOptionClass}>{label}</div>,
    };
    return [...ROLES_OPTIONS, deactivateOption];
  }, [handleStatus, selectedOptionClass]);

  const value = userIsDeactivated
    ? USER_STATUSES.DEACTIVATED.key
    : props.role.value ?? USER_ROLES.USER.key;

  return (
    <div className={props.className}>
      <div className={classes.sectionHeader}>
        <Text variant="subtitle1">
          <fbt desc="">Role</fbt>
        </Text>
        <Text variant="subtitle2" color="gray">
          <fbt desc="">
            Roles determine access to key parts of Symphony like Settings and
            User Management. You can view and change what data this user can
            access in Permissions.
          </fbt>
        </Text>
      </div>
      <RadioGroup
        options={options}
        selectedOptionClassName={selectedOptionClass}
        value={value}
        onChange={newValue => {
          if (handleStatus) {
            if (newValue === USER_STATUSES.DEACTIVATED.key) {
              props.status?.onChange(USER_STATUSES.DEACTIVATED.key);
              return;
            }
            if (userIsDeactivated) {
              props.status?.onChange(USER_STATUSES.ACTIVE.key);
            }
          }
          let typedNewValue = null;
          switch (newValue) {
            case USER_ROLES.USER.key:
              typedNewValue = USER_ROLES.USER.key;
              break;
            case USER_ROLES.ADMIN.key:
              typedNewValue = USER_ROLES.ADMIN.key;
              break;
            case USER_ROLES.OWNER.key:
              typedNewValue = USER_ROLES.OWNER.key;
              break;
          }
          if (typedNewValue != null) {
            props.role.onChange(typedNewValue);
          }
        }}
      />
    </div>
  );
};

export default UserRoleAndStatusPane;
