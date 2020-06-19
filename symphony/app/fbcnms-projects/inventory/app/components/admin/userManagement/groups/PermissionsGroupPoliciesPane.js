/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {UsersGroup} from '../data/UsersGroups';

import * as React from 'react';
import Button from '@fbcnms/ui/components/design-system/Button';
import Card from '@fbcnms/ui/components/design-system/Card/Card';
import PermissionsGroupManagePoliciesDialog, {
  DIALOG_TITLE,
} from './PermissionsGroupManagePoliciesDialog';
import PermissionsPoliciesTable from '../policies/PermissionsPoliciesTable';
import ViewContainer from '@fbcnms/ui/components/design-system/View/ViewContainer';
import fbt from 'fbt';
import {ROW_SEPARATOR_TYPES} from '@fbcnms/ui/components/design-system/Table/TableContent';
import {TABLE_VARIANT_TYPES} from '@fbcnms/ui/components/design-system/Table/Table';
import {useMemo, useState} from 'react';
import {wrapRawPermissionsPolicies} from '../data/PermissionsPolicies';

type Props = $ReadOnly<{|
  group: UsersGroup,
  className?: ?string,
  onChange: UsersGroup => void,
|}>;

export default function PermissionsGroupPoliciesPane(props: Props) {
  const {group, className, onChange} = props;
  const policies = useMemo(() => wrapRawPermissionsPolicies(group.policies), [
    group.policies,
  ]);
  const [showManagePoliciesDialog, setShowManagePoliciesDialog] = useState(
    false,
  );

  return (
    <Card className={className} margins="none">
      <ViewContainer
        header={{
          title: <fbt desc="">Policies</fbt>,
          subtitle: (
            <fbt desc="">
              Add policies to apply them on members in this group.
            </fbt>
          ),
          actionButtons: [
            <Button onClick={() => setShowManagePoliciesDialog(true)}>
              {DIALOG_TITLE}
            </Button>,
          ],
        }}>
        {policies.length > 0 ? (
          <PermissionsPoliciesTable
            policies={policies}
            showGroupsColumn={false}
            variant={TABLE_VARIANT_TYPES.embedded}
            dataRowsSeparator={ROW_SEPARATOR_TYPES.border}
          />
        ) : null}
      </ViewContainer>
      <PermissionsGroupManagePoliciesDialog
        selectedPolicies={group.policies}
        onClose={policies => {
          if (policies != null) {
            onChange({
              ...group,
              policies,
            });
          }
          setShowManagePoliciesDialog(false);
        }}
        open={showManagePoliciesDialog}
      />
    </Card>
  );
}
