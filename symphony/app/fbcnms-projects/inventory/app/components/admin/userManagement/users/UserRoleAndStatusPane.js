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
import type {User} from '../utils/UserManagementUtils';

import * as React from 'react';
import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import RadioGroup from '@fbcnms/ui/components/design-system/RadioGroup/RadioGroup';
import Text from '@fbcnms/ui/components/design-system/Text';
import classNames from 'classnames';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import {USER_ROLES, USER_STATUSES} from '../utils/UserManagementUtils';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useMemo} from 'react';

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

type Props = $ReadOnly<{
  user: User,
  onChange: User => void,
  canSetDeactivated?: ?boolean,
  className?: ?string,
  disabled?: ?boolean,
}>;

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
  const {
    user,
    className,
    onChange,
    canSetDeactivated = true,
    disabled = false,
  } = props;
  const userIsDeactivated = user.status === USER_STATUSES.DEACTIVATED.key;
  const classes = useStyles();
  const selectedOptionClass = classNames({
    [classes.deactivateOptionSelected]: userIsDeactivated,
  });
  const options = useMemo(() => {
    if (!canSetDeactivated) {
      return ROLES_OPTIONS;
    }
    const {label, ...partialDeactivateOption} = STATUS_OPT;
    const deactivateOption = {
      ...partialDeactivateOption,
      label: <div className={selectedOptionClass}>{label}</div>,
    };
    return [...ROLES_OPTIONS, deactivateOption];
  }, [canSetDeactivated, selectedOptionClass]);

  const onRoleChanged = useCallback(
    newValue => {
      const newUser = {...user};
      if (canSetDeactivated) {
        if (newValue === USER_STATUSES.DEACTIVATED.key) {
          newUser.status = USER_STATUSES.DEACTIVATED.key;
          onChange(newUser);
          return;
        }
        if (userIsDeactivated) {
          newUser.status = USER_STATUSES.ACTIVE.key;
        }
      }
      switch (newValue) {
        case USER_ROLES.USER.key:
          newUser.role = USER_ROLES.USER.key;
          break;
        case USER_ROLES.ADMIN.key:
          newUser.role = USER_ROLES.ADMIN.key;
          break;
        default:
          return;
      }
      onChange(newUser);
    },
    [canSetDeactivated, onChange, user, userIsDeactivated],
  );

  const value = userIsDeactivated
    ? USER_STATUSES.DEACTIVATED.key
    : user.role ?? USER_ROLES.USER.key;

  return (
    <div className={className}>
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
      <FormField disabled={!!disabled}>
        <RadioGroup
          options={options}
          selectedOptionClassName={selectedOptionClass}
          value={value}
          onChange={onRoleChanged}
        />
      </FormField>
    </div>
  );
};

export default UserRoleAndStatusPane;
