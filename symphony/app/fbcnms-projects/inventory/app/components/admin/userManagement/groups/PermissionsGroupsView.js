/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {TableRowDataType} from '@fbcnms/ui/components/design-system/Table/Table';
import type {UserPermissionsGroup} from '../utils/UserManagementUtils';

import * as React from 'react';
import Table from '@fbcnms/ui/components/design-system/Table/Table';
import fbt from 'fbt';
import {GROUP_STATUSES} from '../utils/UserManagementUtils';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';
import {useUserManagement} from '../UserManagementContext';

export const PERMISSION_GROUPS_VIEW_NAME = fbt(
  'Groups',
  'Header for view showing system permissions groups settings',
);

const useStyles = makeStyles(() => ({
  root: {
    maxHeight: '100%',
  },
  narrowColumn: {
    width: '70%',
  },
  wideColumn: {
    width: '170%',
  },
}));

type GroupTableRow = TableRowDataType<UserPermissionsGroup>;
type GroupTableData = Array<GroupTableRow>;

const group2GroupTableRow: (
  UserPermissionsGroup | GroupTableRow,
) => GroupTableRow = group => ({
  key: group.key || group.id,
  ...group,
});

export default function PermissionsGroupsView() {
  const classes = useStyles();
  const {history} = useRouter();
  const {groups} = useUserManagement();
  const [groupsTable, _setGroupsTable] = useState<GroupTableData>(
    groups.map(group2GroupTableRow),
  );

  const columns = [
    {
      key: 'name',
      title: (
        <fbt desc="Group Name column header in permission groups table">
          Group Name
        </fbt>
      ),
      getSortingValue: groupRow => groupRow.name,
      render: groupRow => groupRow.name,
    },
    {
      key: 'description',
      title: (
        <fbt desc="Description column header in permission groups table">
          Description
        </fbt>
      ),
      getSortingValue: groupRow => groupRow.description,
      render: groupRow => groupRow.description,
      titleClassName: classes.wideColumn,
      className: classes.wideColumn,
    },
    {
      key: 'members',
      title: (
        <fbt desc="Members column header in permission groups table">
          Members
        </fbt>
      ),
      getSortingValue: groupRow => groupRow.members.length,
      render: groupRow => groupRow.members.length,
      titleClassName: classes.narrowColumn,
      className: classes.narrowColumn,
    },
    {
      key: 'status',
      title: (
        <fbt desc="Status column header in permission groups table">Status</fbt>
      ),
      getSortingValue: groupRow => GROUP_STATUSES[groupRow.status].value,
      render: groupRow => GROUP_STATUSES[groupRow.status].value,
      titleClassName: classes.narrowColumn,
      className: classes.narrowColumn,
    },
  ];

  return (
    <div className={classes.root}>
      <Table
        data={groupsTable}
        onActiveRowIdChanged={groupId => {
          if (groupId != null) {
            history.push(`group/${groupId}`);
          }
        }}
        columns={columns}
      />
    </div>
  );
}
