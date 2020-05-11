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
import type {User, UserPermissionsGroup} from './UserManagementUtils';

import * as React from 'react';
import GroupMemberListItem from './GroupMemberListItem';
import List from './List';

type Props = $ReadOnly<{|
  users: $ReadOnlyArray<User>,
  group?: ?UserPermissionsGroup,
  onChange: UserPermissionsGroup => void,
  emptyState?: ?React.Node,
  className?: ?string,
  ...AssigenmentButtonProp,
|}>;

export default function GroupMembersList(props: Props) {
  const {users, group, assigmentButton, onChange, ...rest} = props;

  return (
    <List items={users} {...rest}>
      {user => (
        <GroupMemberListItem
          user={user}
          assigmentButton={assigmentButton}
          group={group}
          onChange={onChange}
        />
      )}
    </List>
  );
}
