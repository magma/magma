/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {PermissionsPolicy} from '../utils/UserManagementUtils';
import type {TableRowDataType} from '@fbcnms/ui/components/design-system/Table/Table';

import * as React from 'react';
import Table from '@fbcnms/ui/components/design-system/Table/Table';
import fbt from 'fbt';
import {POLICY_TYPES} from '../utils/UserManagementUtils';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';
import {useUserManagement} from '../UserManagementContext';

export const PERMISSION_POLICIES_VIEW_NAME = fbt(
  'Polices',
  'Header for view showing system permissions policies settings',
);

const ALL_USERS = `${fbt('All Users', '')}`;

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

type PolicyTableRow = TableRowDataType<PermissionsPolicy>;
type PolicyTableData = Array<PolicyTableRow>;

const policy2PolicyTableRow: PermissionsPolicy => PolicyTableRow = policy => ({
  key: policy.id,
  ...policy,
});

const getPolicyUsersCount = (PolicyRow: PolicyTableRow) =>
  PolicyRow.isGlobal ? ALL_USERS : PolicyRow.groups.length;

const getPolicyType = (PolicyRow: PolicyTableRow) => {
  switch (PolicyRow.type) {
    case POLICY_TYPES.InventoryPolicy.key:
      return POLICY_TYPES.InventoryPolicy.value;
    case POLICY_TYPES.WorkforcePolicy.key:
      return POLICY_TYPES.WorkforcePolicy.value;
    default:
      return null;
  }
};

export default function PermissionsPoliciesView() {
  const classes = useStyles();
  const {history} = useRouter();
  const {policies} = useUserManagement();
  const [policiesTable, _setPoliciesTable] = useState<PolicyTableData>(
    policies.map(policy2PolicyTableRow),
  );

  const columns = [
    {
      key: 'name',
      title: (
        <fbt desc="Policy Name column header in permission policies table">
          Policy Name
        </fbt>
      ),
      getSortingValue: PolicyRow => PolicyRow.name,
      render: PolicyRow => PolicyRow.name,
    },
    {
      key: 'description',
      title: (
        <fbt desc="Description column header in permission policies table">
          Description
        </fbt>
      ),
      getSortingValue: PolicyRow => PolicyRow.description,
      render: PolicyRow => PolicyRow.description,
      titleClassName: classes.wideColumn,
      className: classes.wideColumn,
    },
    {
      key: 'type',
      title: (
        <fbt desc="Policy Type column header in permission policies table">
          Policy Type
        </fbt>
      ),
      getSortingValue: getPolicyType,
      render: getPolicyType,
      titleClassName: classes.narrowColumn,
      className: classes.narrowColumn,
    },
    {
      key: 'groups',
      title: (
        <fbt desc="Gropus Applied column header in permission groups table">
          Gropus Applied
        </fbt>
      ),
      getSortingValue: getPolicyUsersCount,
      render: getPolicyUsersCount,
      titleClassName: classes.narrowColumn,
      className: classes.narrowColumn,
    },
  ];

  return (
    <div className={classes.root}>
      <Table
        data={policiesTable}
        onActiveRowIdChanged={policyId => {
          if (policyId != null) {
            history.push(`policy/${policyId}`);
          }
        }}
        columns={columns}
      />
    </div>
  );
}
