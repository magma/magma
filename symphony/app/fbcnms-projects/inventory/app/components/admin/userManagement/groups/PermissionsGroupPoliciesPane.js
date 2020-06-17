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
import Card from '@fbcnms/ui/components/design-system/Card/Card';
import PermissionsPoliciesTable from '../policies/PermissionsPoliciesTable';
import ViewContainer from '@fbcnms/ui/components/design-system/View/ViewContainer';
import fbt from 'fbt';
import {ROW_SEPARATOR_TYPES} from '@fbcnms/ui/components/design-system/Table/TableContent';
import {TABLE_VARIANT_TYPES} from '@fbcnms/ui/components/design-system/Table/Table';
import {useMemo} from 'react';
import {wrapRawPermissionsPolicies} from '../data/PermissionsPolicies';

type Props = $ReadOnly<{|
  group: UsersGroup,
  className?: ?string,
|}>;

export default function PermissionsGroupPoliciesPane({
  group,
  className,
}: Props) {
  const policies = useMemo(() => wrapRawPermissionsPolicies(group.policies), [
    group.policies,
  ]);

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
    </Card>
  );
}
