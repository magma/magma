/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {AssigenmentButtonProp} from './GroupMemberListItem';
import type {
  PermissionsPolicy,
  UserPermissionsGroup,
} from './UserManagementUtils';

import * as React from 'react';
import List from './List';
import PolicyGroupListItem from './PolicyGroupListItem';

type Props = $ReadOnly<{|
  groups: $ReadOnlyArray<UserPermissionsGroup>,
  policy?: ?PermissionsPolicy,
  emptyState?: ?React.Node,
  className?: ?string,
  ...AssigenmentButtonProp,
|}>;

export default function PolicyGroupsList(props: Props) {
  const {groups, policy, assigmentButton, ...rest} = props;

  return (
    <List items={groups} {...rest}>
      {group => (
        <PolicyGroupListItem
          group={group}
          assigmentButton={assigmentButton}
          policy={policy}
        />
      )}
    </List>
  );
}
