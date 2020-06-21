/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {OptionalRefTypeWrapper} from '../../../../common/EntUtils';
import type {ToggleButtonDisplay} from './ListItem';
import type {UserManagementUtils_user_base} from './__generated__/UserManagementUtils_user_base.graphql';
import type {UsersGroup} from '../data/UsersGroups';

import * as React from 'react';
import MemberListItem from './MemberListItem';
import UserViewer from '../users/UserViewer';
import {useCallback, useEffect, useState} from 'react';

export type AssigenmentButtonProp = $ReadOnly<{|
  assigmentButton?: ?ToggleButtonDisplay,
|}>;

type Props = $ReadOnly<{|
  user: OptionalRefTypeWrapper<UserManagementUtils_user_base>,
  group?: ?UsersGroup,
  onChange: UsersGroup => void,
  className?: ?string,
  ...AssigenmentButtonProp,
|}>;

const checkIsMember = (
  user: OptionalRefTypeWrapper<UserManagementUtils_user_base>,
  group?: ?UsersGroup,
) =>
  group == null || group.members.find(member => member.id == user.id) != null;

export default function GroupMemberListItem(props: Props) {
  const {user, group, onChange, assigmentButton, className} = props;
  const [isMember, setIsMember] = useState(false);
  useEffect(() => setIsMember(checkIsMember(user, group)), [group, user]);

  const toggleAssigment = useCallback(
    (user, shouldAssign) => {
      if (group == null) {
        return;
      }
      const {id, authID, firstName, lastName, email, status, role} = user;
      const newMembers = shouldAssign
        ? [
            ...group.members,
            {
              id,
              authID,
              firstName,
              lastName,
              email,
              status,
              role,
            },
          ]
        : group.members.filter(m => m.id != user.id);
      onChange({
        ...group,
        members: newMembers,
      });
    },
    [group, onChange],
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
