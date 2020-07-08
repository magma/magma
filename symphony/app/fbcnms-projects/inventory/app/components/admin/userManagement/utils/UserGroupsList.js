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
import type {User} from '../data/Users';
import type {UsersGroup} from '../data/UsersGroups';

import * as React from 'react';
import GroupsList from './GroupsList';
import {updateUserGroups} from '../data/Users';
import {useCallback} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';

type Props = $ReadOnly<{|
  user?: ?User,
  ...GroupsListProps,
  ...AssigenmentButtonProp,
|}>;

const checkIsGroupAppliedToUser = (user?: ?User) => (group: UsersGroup) =>
  user == null || user.groups.find(g => g?.id === group.id) != null;

export default function UserGroupsList(props: Props) {
  const {user, assigmentButton, ...rest} = props;

  const enqueueSnackbar = useEnqueueSnackbar();
  const handleError = useCallback(
    error => {
      enqueueSnackbar(error.response?.data?.error || error, {variant: 'error'});
    },
    [enqueueSnackbar],
  );

  const toggleAssigment = useCallback(
    (group, shouldAssign) => {
      if (user == null) {
        return;
      }
      const editPromise = shouldAssign
        ? updateUserGroups(user, [group.id], [])
        : updateUserGroups(user, [], [group.id]);
      return editPromise.catch(handleError).then(() => undefined);
    },
    [handleError, user],
  );

  return (
    <GroupsList
      {...rest}
      assignment={{
        assigmentButton,
        isGroupAssigned: checkIsGroupAppliedToUser(user),
        onGroupAssignmentChange: toggleAssigment,
      }}
    />
  );
}
