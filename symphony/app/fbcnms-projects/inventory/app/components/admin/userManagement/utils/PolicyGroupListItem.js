/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {PermissionsPolicy} from '../data/PermissionsPolicies';
import type {ToggleButtonDisplay} from './ListItem';
import type {UsersGroup} from '../data/UsersGroups';

import * as React from 'react';
import GroupIcon from '@fbcnms/ui/components/design-system/Icons/Indications/GroupIcon';
import MemberListItem from './MemberListItem';
import Text from '@fbcnms/ui/components/design-system/Text';
import classNames from 'classnames';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import {GROUP_STATUSES} from './UserManagementUtils';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useEffect, useState} from 'react';

const useStyles = makeStyles(() => ({
  policyContainer: {
    display: 'flex',
    height: '100%',
    overflow: 'hidden',
  },
  photoContainer: {
    borderRadius: '50%',
    marginRight: '8px',
    backgroundColor: symphony.palette.D10,
    width: '48px',
    height: '48px',
    display: 'flex',
    flexShrink: 0,
    justifyContent: 'center',
    alignItems: 'center',
  },
  photo: {
    margin: 'auto',
  },
  details: {
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'space-evenly',
    flexShrink: 1,
    overflow: 'hidden',
  },
  metaData: {
    display: 'flex',
    whiteSpace: 'nowrap',
    '& span': {
      marginRight: '2px',
    },
  },
}));

export type AssigenmentButtonProp = $ReadOnly<{|
  assigmentButton?: ?ToggleButtonDisplay,
|}>;

type Props = $ReadOnly<{|
  group: UsersGroup,
  isMember?: ?boolean,
  policy?: ?PermissionsPolicy,
  onChange: PermissionsPolicy => void,
  className?: ?string,
  ...AssigenmentButtonProp,
|}>;

const checkIsGroupInPolicy = (group: UsersGroup, policy?: ?PermissionsPolicy) =>
  policy == null || policy.groups.find(g => g.id === group.id) != null;

export default function PolicyGroupListItem(props: Props) {
  const {group, policy, assigmentButton, className, onChange} = props;
  const classes = useStyles();

  const [isGroupInPolicy, setIsGroupInPolicy] = useState(false);
  useEffect(() => setIsGroupInPolicy(checkIsGroupInPolicy(group, policy)), [
    group,
    policy,
  ]);

  const toggleAssigment = useCallback(
    (group: UsersGroup, shouldAssign) => {
      if (policy == null) {
        return;
      }
      const newGroups = shouldAssign
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
    <MemberListItem
      member={{
        item: group,
        isMember: isGroupInPolicy,
      }}
      className={className}
      assigmentButton={assigmentButton}
      onAssignToggle={() => toggleAssigment(group, !isGroupInPolicy)}>
      <div className={classNames(classes.policyContainer, className)}>
        <div className={classes.photoContainer}>
          <GroupIcon color="gray" />
        </div>
        <div className={classes.details}>
          <Text variant="subtitle2" useEllipsis={true}>
            {group.name}
          </Text>
          <div className={classes.metaData}>
            <Text variant="caption" color="gray" useEllipsis={true}>
              <fbt desc="">
                <fbt:plural count={group.members.length} showCount="yes">
                  member
                </fbt:plural>
              </fbt>
            </Text>
            <Text variant="caption" color="gray" useEllipsis={true}>
              {' • '}
            </Text>
            <Text
              variant="caption"
              color={
                group.status === GROUP_STATUSES.DEACTIVATED.key
                  ? 'error'
                  : 'gray'
              }>
              {GROUP_STATUSES[group.status].value}
            </Text>
          </div>
        </div>
      </div>
    </MemberListItem>
  );
}
