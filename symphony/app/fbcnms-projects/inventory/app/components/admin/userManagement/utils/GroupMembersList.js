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
import type {OptionalRefTypeWrapper} from '../../../../common/EntUtils';
import type {UserManagementUtils_user_base} from './__generated__/UserManagementUtils_user_base.graphql';
import type {UsersGroup} from '../data/UsersGroups';

import * as React from 'react';
import GroupMemberListItem from './GroupMemberListItem';
import List from './List';

type Props = $ReadOnly<{|
  users: $ReadOnlyArray<OptionalRefTypeWrapper<UserManagementUtils_user_base>>,
  group?: ?UsersGroup,
  onChange: UsersGroup => void,
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
