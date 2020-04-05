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
import AppContext from '@fbcnms/ui/context/AppContext';
import Table from '@fbcnms/ui/components/design-system/Table/Table';
import Text from '@fbcnms/ui/components/design-system/Text';
import UserDetailsCard from './UserDetailsCard';
import UserViewer from './UserViewer';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import {USER_ROLES, USER_STATUSES} from '../utils/UserManagementUtils';
import {haveDifferentValues} from '../../../../common/EntUtils';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useEffect, useMemo, useState} from 'react';
import {useContext} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useUserManagement} from '../UserManagementContext';

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

type UserTableRow = TableRowDataType<{|data: User|}>;
type UserTableData = Array<UserTableRow>;

const user2UserTableRow: User => UserTableRow = user => ({
  key: user.authID,
  data: user,
});

export default function UsersTable() {
  const classes = useStyles();

  const {isFeatureEnabled} = useContext(AppContext);
  const userManagementDevMode = isFeatureEnabled('user_management_dev');

  const [usersTableData, setUsersTableData] = useState<UserTableData>([]);
  const {users, editUser} = useUserManagement();
  useEffect(() => setUsersTableData(users.map(user2UserTableRow)), [users]);
  const [selectedUserIds, setSelectedUserIds] = useState<Array<TableRowId>>([]);
  const [activeUserId, setActiveUserId] = useState(null);

  const columns = useMemo(() => {
    const isActiveUser = userId =>
      activeUserId != null && activeUserId === userId;
    const returnCols = [
      {
        key: 'name',
        title: <fbt desc="Name column header in users table">Name</fbt>,
        titleClassName: classes.nameColumn,
        className: classes.nameColumn,
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
        render: userRow =>
          userRow.data.status === USER_STATUSES.DEACTIVATED.key
            ? null
            : USER_ROLES[userRow.data.role].value || userRow.data.role,
      },
      {
        key: 'status',
        title: <fbt desc="Status column header in users table">Status</fbt>,
        render: userRow => (
          <Text
            useEllipsis={true}
            color={
              userRow.data.status === USER_STATUSES.DEACTIVATED.key
                ? 'error'
                : undefined
            }>
            {USER_STATUSES[userRow.data.status].value || userRow.data.status}
          </Text>
        ),
      },
    ];
    if (userManagementDevMode) {
      returnCols.push(
        ...[
          {
            key: 'job_title',
            title: (
              <fbt desc="Job Title column header in users table">Job Title</fbt>
            ),
            render: userRow => userRow.data.jobTitle ?? '',
          },
          {
            key: 'employment',
            title: (
              <fbt desc="Employment column header in users table">
                Employment
              </fbt>
            ),
            render: userRow => userRow.data.employmentType ?? '',
          },
        ],
      );
    }
    return returnCols;
  }, [classes.nameColumn, classes.field, userManagementDevMode, activeUserId]);

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
