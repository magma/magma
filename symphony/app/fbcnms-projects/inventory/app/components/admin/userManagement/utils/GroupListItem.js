/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {AssigenmentButtonProp} from './MemberListItem';
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
  groupContainer: {
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

type Props = $ReadOnly<{|
  group: UsersGroup,
  isMember?: ?boolean,
  onChange?: boolean => Promise<void> | void,
  className?: ?string,
  ...AssigenmentButtonProp,
|}>;

export default function GroupListItem(props: Props) {
  const {
    group,
    isMember: isMemberProp,
    assigmentButton,
    className,
    onChange,
  } = props;
  const classes = useStyles();

  const [isMember, setIsMember] = useState(false);
  useEffect(() => setIsMember(isMemberProp === true), [isMemberProp]);

  const callOnChange = useCallback(
    newValue => {
      if (onChange == null) {
        return;
      }
      return onChange(newValue);
    },
    [onChange],
  );

  return (
    <MemberListItem
      member={{
        item: group,
        isMember,
      }}
      className={className}
      assigmentButton={assigmentButton}
      onAssignToggle={() => callOnChange(!isMember)}>
      <div className={classNames(classes.groupContainer, className)}>
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
