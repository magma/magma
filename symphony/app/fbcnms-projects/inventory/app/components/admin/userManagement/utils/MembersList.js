/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {AssigenmentButtonProp, GroupMember} from './GroupMemberViewer';
import type {UserPermissionsGroup} from '../utils/UserManagementUtils';

import * as React from 'react';
import GroupMemberViewer from './GroupMemberViewer';
import classNames from 'classnames';
import symphony from '@fbcnms/ui/theme/symphony';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  root: {
    display: 'flex',
    flexDirection: 'column',
    width: '100%',
    minWidth: '240px',
  },
  user: {
    borderBottom: `1px solid ${symphony.palette.separatorLight}`,
  },
}));

type Props = $ReadOnly<{|
  members: $ReadOnlyArray<GroupMember>,
  group?: ?UserPermissionsGroup,
  emptyState?: ?React.Node,
  className?: ?string,
  ...AssigenmentButtonProp,
|}>;

export default function MembersList(props: Props) {
  const {members, group, emptyState, assigmentButton, className} = props;
  const classes = useStyles();

  return (
    <div className={classNames(classes.root, className)}>
      {members.length == 0 && emptyState != null
        ? emptyState
        : members.map(member => (
            <GroupMemberViewer
              className={classes.user}
              member={member}
              assigmentButton={assigmentButton}
              group={group}
            />
          ))}
    </div>
  );
}
