/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {User} from './TempTypes';

import * as React from 'react';
import EmojiIcon from '@fbcnms/ui/components/design-system/Icons/Communication/EmojiIcon';
import Text from '@fbcnms/ui/components/design-system/Text';
import classNames from 'classnames';
import symphony from '@fbcnms/ui/theme/symphony';
import {USER_ROLES} from './TempTypes';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  root: {
    display: 'flex',
    height: '100%',
    overflow: 'hidden',
  },
  photoContainer: {
    borderRadius: '50%',
    marginRight: '8px',
    backgroundColor: symphony.palette.D200,
    width: '48px',
    height: '48px',
    display: 'flex',
    flexShrink: 0,
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

type Props = $ReadOnly<{
  user: User,
  highlightName?: ?boolean,
  showPhoto?: ?boolean,
  showRole?: ?boolean,
  className?: ?string,
}>;

export default function UserViewer(props: Props) {
  const {
    user,
    highlightName = false,
    showPhoto = false,
    showRole = false,
    className,
  } = props;
  const classes = useStyles();

  return (
    <div className={classNames(classes.root, className)}>
      {showPhoto ? (
        <div className={classes.photoContainer}>
          {user.photoId ? (
            user.photoId
          ) : (
            <EmojiIcon color="light" className={classes.photo} />
          )}
        </div>
      ) : null}
      <div className={classes.details}>
        <Text
          variant="subtitle2"
          useEllipsis={true}
          color={highlightName ? 'primary' : undefined}>
          {`${user.firstName} ${user.lastName}`.trim() || '_'}
        </Text>
        <div className={classes.metaData}>
          <Text variant="caption" color="gray" useEllipsis={true}>
            {user.authID}
          </Text>
          {showRole ? (
            <Text variant="caption" color="gray">
              {` • ${USER_ROLES[user.role].value}`}
            </Text>
          ) : null}
        </div>
      </div>
    </div>
  );
}
