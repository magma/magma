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
import {TEMP_USERS, USER_STATUSES} from './TempTypes';
import {makeStyles} from '@material-ui/styles';
import {useMemo, useState} from 'react';

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
  key: user.key || user.authId,
  ...user,
});

export default function UsersView() {
  const classes = useStyles();
  const [usersTable, setUsersTable] = useState<UserTableData>(
    TEMP_USERS.map(user2UserTableRow),
  );
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
              {userRow.firstName} {userRow.lastName}
            </Text>
            <Text
              variant="caption"
              color="gray"
              useEllipsis={true}
              className={classes.field}>
              {userRow.authId}
            </Text>
          </>
        ),
      },
      {
        key: 'role',
        title: <fbt desc="Role column header in users table">Role</fbt>,
        render: userRow => userRow.role,
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
            color={
              userRow.status === USER_STATUSES.Deactivated ? 'error' : undefined
            }>
            {userRow.status}
          </Text>
        ),
      },
    ];
  }, [classes.nameColumn, classes.field, activeUserId]);

  const userDetailsCard = useMemo(() => {
    if (!activeUserId) {
      return null;
    }
    const userIndex = TEMP_USERS.findIndex(
      user => user.authId === activeUserId,
    );
    if (userIndex < 0) {
      return null;
    }

    return (
      <UserDetailsCard
        user={TEMP_USERS[userIndex]}
        onChange={changedUser => {
          setUsersTable([
            ...usersTable.slice(0, userIndex),
            user2UserTableRow(changedUser),
            ...usersTable.slice(userIndex + 1),
          ]);
        }}
      />
    );
  }, [activeUserId, usersTable]);
  return (
    <div className={classes.root}>
      <Table
        dataRowsSeparator="border"
        showSelection={true}
        activeRowId={activeUserId}
        onActiveRowIdChanged={setActiveUserId}
        selectedIds={selectedUserIds}
        onSelectionChanged={ids => setSelectedUserIds(ids)}
        data={usersTable}
        columns={columns}
        detailsCard={userDetailsCard}
      />
    </div>
  );
}
