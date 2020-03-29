/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {
  TableRowDataType,
  TableRowId,
} from '@fbcnms/ui/components/design-system/Table/Table';
import type {User} from './TempTypes';

import * as React from 'react';
import Table from '@fbcnms/ui/components/design-system/Table/Table';
import Text from '@fbcnms/ui/components/design-system/Text';
import UserDetailsCard from './UserDetailsCard';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import {USER_ROLES, USER_STATUSES} from './TempTypes';
import {haveDifferentValues} from '../../../common/EntUtils';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useEffect, useMemo, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useUserManagement} from './UserManagementContext';

const useStyles = makeStyles(() => ({
  root: {
    flexGrow: 1,
    display: 'flex',
    backgroundColor: symphony.palette.white,
    borderRadius: '4px',
  },
  field: {
    margin: '2px',
  },
  nameColumn: {
    width: '200%',
  },
}));

type UserTableRow = TableRowDataType<User>;
type UserTableData = Array<UserTableRow>;

const user2UserTableRow: (User | UserTableRow) => UserTableRow = user => ({
  key: user.key || user.authID,
  ...user,
});

export default function UsersTable() {
  const classes = useStyles();
  const [usersTableData, setUsersTableData] = useState<UserTableData>([]);
  const {users, editUser} = useUserManagement();
  useEffect(() => setUsersTableData(users.map(user2UserTableRow)), [users]);
  const [selectedUserIds, setSelectedUserIds] = useState<Array<TableRowId>>([]);
  const [activeUserId, setActiveUserId] = useState(null);

  const columns = useMemo(() => {
    const isActiveUser = userId =>
      activeUserId != null && activeUserId === userId;
    return [
      {
        key: 'name',
        title: <fbt desc="Name column header in users table">Name</fbt>,
        titleClassName: classes.nameColumn,
        className: classes.nameColumn,
        render: userRow => (
          <>
            <Text
              variant="subtitle2"
              color={isActiveUser(userRow.key) ? 'primary' : undefined}
              useEllipsis={true}
              className={classes.field}>
              {userRow.firstName || userRow.lastName
                ? `${userRow.firstName} ${userRow.lastName}`
                : '_'}
            </Text>
            <Text
              variant="caption"
              color="gray"
              useEllipsis={true}
              className={classes.field}>
              {userRow.authID}
            </Text>
          </>
        ),
      },
      {
        key: 'role',
        title: <fbt desc="Role column header in users table">Role</fbt>,
        render: userRow =>
          userRow.status === USER_STATUSES.DEACTIVATED.key
            ? null
            : USER_ROLES[userRow.role].value || userRow.role,
      },
      {
        key: 'job_title',
        title: (
          <fbt desc="Job Title column header in users table">Job Title</fbt>
        ),
        render: userRow => userRow.jobTitle ?? '',
      },
      {
        key: 'employment',
        title: (
          <fbt desc="Employment column header in users table">Employment</fbt>
        ),
        render: userRow => userRow.employmentType ?? '',
      },
      {
        key: 'status',
        title: <fbt desc="Status column header in users table">Status</fbt>,
        render: userRow => (
          <Text
            useEllipsis={true}
            color={
              userRow.status === USER_STATUSES.DEACTIVATED.key
                ? 'error'
                : undefined
            }>
            {USER_STATUSES[userRow.status].value || userRow.status}
          </Text>
        ),
      },
    ];
  }, [classes.nameColumn, classes.field, activeUserId]);

  const enqueueSnackbar = useEnqueueSnackbar();
  const handleError = useCallback(
    error => {
      enqueueSnackbar(error.response?.data?.error || error, {variant: 'error'});
    },
    [enqueueSnackbar],
  );

  const userDetailsCard = useMemo(() => {
    if (!activeUserId) {
      return null;
    }
    const userIndex = users.findIndex(user => user.authID === activeUserId);
    if (userIndex < 0) {
      return null;
    }

    return (
      <UserDetailsCard
        user={users[userIndex]}
        onChange={user => {
          if (haveDifferentValues(users[userIndex], user)) {
            editUser(user).catch(handleError);
          }
        }}
      />
    );
  }, [activeUserId, editUser, handleError, users]);
  return (
    <div className={classes.root}>
      <Table
        dataRowsSeparator="border"
        activeRowId={activeUserId}
        onActiveRowIdChanged={setActiveUserId}
        selectedIds={selectedUserIds}
        onSelectionChanged={ids => setSelectedUserIds(ids)}
        data={usersTableData}
        columns={columns}
        detailsCard={userDetailsCard}
      />
    </div>
  );
}
