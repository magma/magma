/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {User, UserPermissionsGroup} from '../utils/UserManagementUtils';

import * as React from 'react';
import Button from '@fbcnms/ui/components/design-system/Button';
import CheckIcon from '@fbcnms/ui/components/design-system/Icons/Indications/CheckIcon';
import PlusIcon from '@fbcnms/ui/components/design-system/Icons/Actions/PlusIcon';
import Strings from '@fbcnms/strings/Strings';
import UserViewer from '../users/UserViewer';
import classNames from 'classnames';
import fbt from 'fbt';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useEffect, useState} from 'react';
import {useUserManagement} from '../UserManagementContext';

const useStyles = makeStyles(() => ({
  root: {
    padding: '8px 0',
    display: 'flex',
    overflow: 'hidden',
    alignItems: 'center',
    '&:not(:hover):not($alwaysShowAssigmentButton) $userAssignButton': {
      display: 'none',
    },
    flexShrink: 0,
  },
  alwaysShowAssigmentButton: {},
  userDetails: {
    flexBasis: '10px',
    flexGrow: 1,
    flexShrink: 1,
  },
  userAssignButton: {
    maxWidth: '88px',
    '& $removeText': {
      display: 'none',
    },
    '&:hover, &$togglingAssignment': {
      '& $addedIcon, $addedText': {
        display: 'none',
      },
      '& $removeText': {
        display: 'unset',
      },
    },
  },
  togglingAssignment: {},
  addedIcon: {},
  addedText: {},
  removeText: {},
}));

export const ASSIGNMENT_BUTTON_VIEWS = {
  always: 'always',
  onHover: 'onHover',
};

type AssigenmentButtonView = $Keys<typeof ASSIGNMENT_BUTTON_VIEWS>;

export type GroupMember = $ReadOnly<{|
  user: User,
  isMember: boolean,
|}>;

export type AssigenmentButtonProp = $ReadOnly<{|
  assigmentButton: AssigenmentButtonView,
|}>;

type Props = $ReadOnly<{|
  member: GroupMember,
  group?: ?UserPermissionsGroup,
  className?: ?string,
  ...AssigenmentButtonProp,
|}>;

export default function GroupMemberViewer(props: Props) {
  const [member, setMember] = useState(props.member);
  useEffect(() => setMember(props.member), [props.member]);
  const {group, assigmentButton, className} = props;
  const classes = useStyles();
  const userManagement = useUserManagement();
  const [isProcessed, setIsProcessed] = useState(false);

  const toggleAssigment = useCallback(
    (memberUser, shouldAssign) => {
      if (group == null) {
        return;
      }

      setIsProcessed(true);
      const add = shouldAssign ? [memberUser.user.id] : [];
      const remove = shouldAssign ? [] : [memberUser.user.id];
      userManagement
        .updateGroupMembers(group, add, remove)
        .then(() =>
          setMember({
            isMember: shouldAssign,
            user: {
              ...memberUser.user,
              groups: shouldAssign
                ? [...memberUser.user.groups, {id: group.id, name: group.name}]
                : memberUser.user.groups.filter(
                    userGroup => userGroup?.id != group.id,
                  ),
            },
          }),
        )
        .finally(() => setIsProcessed(false));
    },
    [group, userManagement],
  );

  return (
    <div
      className={classNames(classes.root, className, {
        [classes.alwaysShowAssigmentButton]:
          assigmentButton == ASSIGNMENT_BUTTON_VIEWS.always || isProcessed,
      })}>
      <UserViewer
        className={classes.userDetails}
        user={member.user}
        showPhoto={true}
        showRole={true}
      />
      {group == null ? null : (
        <Button
          className={classNames(classes.userAssignButton, {
            [classes.togglingAssignment]: isProcessed,
          })}
          disabled={isProcessed}
          onClick={() => toggleAssigment(member, !member.isMember)}
          skin={member.isMember ? 'gray' : 'primary'}
          leftIcon={member.isMember ? CheckIcon : PlusIcon}
          leftIconClass={member.isMember ? classes.addedIcon : undefined}>
          {member.isMember ? (
            <>
              <div className={classes.removeText}>
                {isProcessed ? (
                  <fbt desc="">Removing</fbt>
                ) : (
                  Strings.common.removeButton
                )}
              </div>
              <div className={classes.addedText}>
                <fbt desc="">Added</fbt>
              </div>
            </>
          ) : isProcessed ? (
            <fbt desc="">Adding</fbt>
          ) : (
            Strings.common.addButton
          )}
        </Button>
      )}
    </div>
  );
}
