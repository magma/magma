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
import Text from '@fbcnms/ui/components/design-system/Text';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  root: {
    display: 'flex',
    flexDirection: 'column',
    height: '100%',
  },
}));

type Props = {
  user: User,
};

export default function UserPermissionsPane(props: Props) {
  const {user} = props;
  const classes = useStyles();

  return (
    <div className={classes.root}>
      <Text>{`Permission Details for ${user.authID}`}</Text>
    </div>
  );
}
