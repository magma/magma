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
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import {TEMP_USERS} from './TempTypes';
import {makeStyles} from '@material-ui/styles';
import {useEffect, useMemo, useState} from 'react';

const useStyles = makeStyles(() => ({
  root: {
    display: 'flex',
    height: '100%',
    backgroundColor: symphony.palette.white,
  },
  field: {
    margin: '2px',
  },
}));

export default function UsersView() {
  const classes = useStyles();
  const [selectedUserIds, setSelectedUserIds] = useState<Array<TableRowId>>([]);
  const [exclusivlySelectedUserId, setExclusivlySelectedUserId] = useState(
    null,
  );
  const tableDate: Array<TableRowDataType<User>> = useMemo(
    () =>
      TEMP_USERS.map(user => ({
        key: user.authId,
        ...user,
      })),
    [],
  );

  useEffect(() => {
    if (selectedUserIds.length === 1) {
      setExclusivlySelectedUserId(selectedUserIds[0]);
    } else {
      setExclusivlySelectedUserId(null);
    }
  }, [selectedUserIds]);

  const isExclusivlySelectedUser = userId =>
    exclusivlySelectedUserId != null && exclusivlySelectedUserId === userId;

  return (
    <div className={classes.root}>
      <Table
        dataRowsSeparator="border"
        showSelection={true}
        selectedIds={selectedUserIds}
        onSelectionChanged={ids => setSelectedUserIds(ids)}
        data={tableDate}
        columns={[
          {
            key: 'name',
            title: <fbt desc="Name column header in users table">Name</fbt>,
            render: userRow => (
              <>
                <Text
                  variant="subtitle2"
                  color={
                    isExclusivlySelectedUser(userRow.key)
                      ? 'primary'
                      : undefined
                  }
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
              <fbt desc="Employment column header in users table">
                Employment
              </fbt>
            ),
            render: userRow => userRow.employmentType ?? '',
          },
          {
            key: 'status',
            title: <fbt desc="Status column header in users table">Status</fbt>,
            render: userRow => userRow.status,
          },
        ]}
      />
    </div>
  );
}
