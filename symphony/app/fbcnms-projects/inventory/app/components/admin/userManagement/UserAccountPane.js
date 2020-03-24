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
import UserAccountDetailsPane, {
  ACCOUNT_DISPLAY_VARIANTS,
} from './UserAccountDetailsPane';
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
  onChange: User => void,
};

export default function UserAccountPane(props: Props) {
  const {user, onChange} = props;
  const classes = useStyles();

  return (
    <div className={classes.root}>
      <UserAccountDetailsPane
        variant={ACCOUNT_DISPLAY_VARIANTS.userDetailsCard}
        user={user}
        onChange={(user, _password) => onChange(user)}
      />
    </div>
  );
}
