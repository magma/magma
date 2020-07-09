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
import type {GroupsListProps} from './GroupsList';
import type {PermissionsPolicy} from '../data/PermissionsPolicies';
import type {UsersGroup} from '../data/UsersGroups';

import * as React from 'react';
import GroupsList from './GroupsList';
import {useCallback} from 'react';

type Props = $ReadOnly<{|
  policy?: ?PermissionsPolicy,
  onChange: PermissionsPolicy => void,
  ...GroupsListProps,
  ...AssigenmentButtonProp,
|}>;

const checkIsGroupInPolicy = (policy?: ?PermissionsPolicy) => (
  group: UsersGroup,
) => policy == null || policy.groups.find(g => g.id === group.id) != null;

export default function PolicyGroupsList(props: Props) {
  const {policy, onChange, assigmentButton, ...rest} = props;

  const toggleAssigment = useCallback(
    (group, shouldAssign) => {
      if (policy == null) {
        return;
      }
      const newGroups =
        shouldAssign === true
          ? [...policy.groups, group]
          : policy.groups.filter(g => g.id != group.id);
      onChange({
        ...policy,
        groups: newGroups,
      });
    },
    [onChange, policy],
  );

  return (
    <GroupsList
      {...rest}
      assignment={{
        assigmentButton,
        isGroupAssigned: checkIsGroupInPolicy(policy),
        onGroupAssignmentChange: toggleAssigment,
      }}
    />
  );
}
