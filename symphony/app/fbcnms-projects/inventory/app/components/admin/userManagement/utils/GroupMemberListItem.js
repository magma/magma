/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ToggleButtonDisplay} from './ListItem';
import type {User, UserPermissionsGroup} from './UserManagementUtils';

import * as React from 'react';
import MemberListItem from './MemberListItem';
import UserViewer from '../users/UserViewer';
import {useCallback, useEffect, useState} from 'react';
import {useUserManagement} from '../UserManagementContext';

export type AssigenmentButtonProp = $ReadOnly<{|
  assigmentButton?: ?ToggleButtonDisplay,
|}>;

type Props = $ReadOnly<{|
  user: User,
  group?: ?UserPermissionsGroup,
  className?: ?string,
  ...AssigenmentButtonProp,
|}>;

const checkIsMember = (user: User, group?: ?UserPermissionsGroup) =>
  group == null || group.members.find(member => member.id == user.id) != null;

export default function GroupMemberListItem(props: Props) {
  const {user, group, assigmentButton, className} = props;
  const [isMember, setIsMember] = useState(false);
  useEffect(() => setIsMember(checkIsMember(user, group)), [group, user]);
  const userManagement = useUserManagement();

  const toggleAssigment = useCallback(
    (user, shouldAssign) => {
      if (group == null) {
        return Promise.resolve();
      }
      const add = shouldAssign ? [user.id] : [];
      const remove = shouldAssign ? [] : [user.id];
      return userManagement
        .updateGroupMembers(group, add, remove)
        .then(() => setIsMember(shouldAssign));
    },
    [group, userManagement],
  );

  return (
    <MemberListItem
      member={{
        item: user,
        isMember,
      }}
      className={className}
      assigmentButton={assigmentButton}
      onAssignToggle={() => toggleAssigment(user, !isMember)}>
      <UserViewer user={user} showPhoto={true} showRole={true} />
    </MemberListItem>
  );
}
