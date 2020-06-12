/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {ActivityField} from './__generated__/ActivityPost_activity.graphql';
import type {SvgIconStyleProps} from '@fbcnms/ui/components/design-system/Icons/SvgIcon';

import * as React from 'react';
import symphony from '@fbcnms/ui/theme/symphony';
import {
  EditIcon,
  MessageIcon,
  ProfileIcon,
} from '@fbcnms/ui/components/design-system/Icons';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  activityIndicator: {
    marginRight: '8px',
  },
  activityTypeIcon: {
    fontSize: '16px',
    fill: symphony.palette.D100,
  },
}));

type Props = {
  field: ActivityField | 'COMMENT',
};

const ActivityCommentsIcon = ({field}: Props) => {
  const classes = useStyles();
  let Icon: ?React.ComponentType<SvgIconStyleProps> = null;
  switch (field) {
    case 'COMMENT':
      Icon = MessageIcon;
      break;
    case 'ASSIGNEE':
    case 'OWNER':
      Icon = ProfileIcon;
      break;
    default:
      Icon = EditIcon;
  }
  return (
    <div className={classes.activityIndicator}>
      <Icon className={classes.activityTypeIcon} />
    </div>
  );
};

export default ActivityCommentsIcon;
