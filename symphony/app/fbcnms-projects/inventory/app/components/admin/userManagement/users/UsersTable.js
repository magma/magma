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
import type {User} from '../utils/UserManagementUtils';

import * as React from 'react';
import Table from '@fbcnms/ui/components/design-system/Table/Table';
import Text from '@fbcnms/ui/components/design-system/Text';
import UserDetailsCard from './UserDetailsCard';
import UserViewer from './UserViewer';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import withSuspense from '../../../../common/withSuspense';
import {
  USER_ROLES,
  USER_STATUSES,
  userFullName,
} from '../utils/UserManagementUtils';
import {editUser, useUsers} from '../data/Users';
import {haveDifferentValues} from '../../../../common/EntUtils';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useMemo, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useHistory, useRouteMatch} from 'react-router-dom';

const useStyles = makeStyles(() => ({
  root: {
    flexGrow: 1,
    display: 'flex',
    backgroundColor: symphony.palette.white,
    borderRadius: '4px',
  },
  table: {
    height: 'unset',
  },
  field: {
    margin: '2px',
  },
  nameColumn: {
    width: '200%',
  },
}));

type UserTableRow = TableRowDataType<{|data: User|}>;

const user2UserTableRow: User => UserTableRow = user => ({
  key: user.authID,
  data: user,
});

export const USER_PATH_PARAM = ':id';
export const ALL_USERS_PATH_PARAM = 'all';

function UsersTable() {
  const classes = useStyles();
  const history = useHistory();
  const match = useRouteMatch();

  const users = useUsers();
  const usersTableData = useMemo(() => users.map(user2UserTableRow), [users]);
  const [selectedUserIds, setSelectedUserIds] = useState<
    $ReadOnlyArray<TableRowId>,
  >([]);

  const userRow2UserRole = useCallback(
    userRow =>
      userRow.data.status === USER_STATUSES.DEACTIVATED.key
        ? null
        : USER_ROLES[userRow.data.role].value || userRow.data.role,
    [],
  );

  const activeUserId =
    match.params.id != null && match.params.id != ALL_USERS_PATH_PARAM
      ? match.params.id
      : null;

  const columns = useMemo(() => {
    const isActiveUser = userId =>
      activeUserId != null && activeUserId === userId;
    const returnCols = [
      {
        key: 'name',
        title: <fbt desc="Name column header in users table">Name</fbt>,
        titleClassName: classes.nameColumn,
        className: classes.nameColumn,
        getSortingValue: userRow =>
          `${userFullName(userRow.data)}${userRow.data.authID}`,
        render: userRow => (
          <UserViewer
            user={userRow.data}
            highlightName={isActiveUser(userRow.key)}
            className={classes.field}
          />
        ),
      },
      {
        key: 'role',
        title: <fbt desc="Role column header in users table">Role</fbt>,
        getSortingValue: userRow2UserRole,
        render: userRow2UserRole,
      },
      {
        key: 'status',
        title: <fbt desc="Status column header in users table">Status</fbt>,
        getSortingValue: userRow => USER_STATUSES[userRow.data.status].value,
        render: userRow => (
          <Text
            useEllipsis={true}
            color={
              userRow.data.status === USER_STATUSES.DEACTIVATED.key
                ? 'error'
                : undefined
            }>
            {USER_STATUSES[userRow.data.status].value}
          </Text>
        ),
      },
    ];
    return returnCols;
  }, [classes.nameColumn, classes.field, userRow2UserRole, activeUserId]);

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
  }, [activeUserId, handleError, users]);

  const navigateToUser = useCallback(
    userId => {
      history.push(
        match.path.replace(
          USER_PATH_PARAM,
          `${userId ?? ALL_USERS_PATH_PARAM}`,
        ),
      );
    },
    [history, match.path],
  );

  return (
    <div className={classes.root}>
      <Table
        className={classes.table}
        dataRowsSeparator="border"
        activeRowId={activeUserId}
        onActiveRowIdChanged={navigateToUser}
        selectedIds={selectedUserIds}
        onSelectionChanged={ids => setSelectedUserIds(ids)}
        data={usersTableData}
        columns={columns}
        detailsCard={userDetailsCard}
      />
    </div>
  );
}

export default withSuspense(UsersTable);
