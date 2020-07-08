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
import type {UsersGroup} from '../data/UsersGroups';

import * as React from 'react';
import GroupListItem from './GroupListItem';
import List from './List';

export type GroupsListProps = $ReadOnly<{|
  groups: $ReadOnlyArray<UsersGroup>,
  emptyState?: ?React.Node,
  className?: ?string,
  groupClassName?: ?string,
|}>;

type Props = $ReadOnly<{|
  ...GroupsListProps,
  assignment?: {
    ...AssigenmentButtonProp,
    isGroupAssigned: UsersGroup => boolean,
    onGroupAssignmentChange: (UsersGroup, boolean) => Promise<void> | void,
  },
|}>;

export default function GroupsList(props: Props) {
  const {groups, assignment, groupClassName, ...rest} = props;

  return (
    <List items={groups} {...rest}>
      {group => (
        <GroupListItem
          key={group.id}
          className={groupClassName}
          group={group}
          isMember={assignment?.isGroupAssigned(group)}
          onChange={
            assignment != null
              ? shouldAssign =>
                  assignment.onGroupAssignmentChange(
                    group,
                    shouldAssign === true,
                  )
              : undefined
          }
          assigmentButton={assignment?.assigmentButton}
        />
      )}
    </List>
  );
}
