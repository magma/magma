/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import * as React from 'react';
import symphony from '../theme/symphony';

import {makeStyles} from '@material-ui/styles';

type Props = {
  onClick: () => void,
  children: React.Node,
};

const useStyles = makeStyles(() => ({
  root: {
    cursor: 'pointer',
    textDecoration: 'underline',
    color: symphony.palette.primary,
  },
}));

// TODO(T38660666) - style according to design
export default function Link(props: Props) {
  const classes = useStyles();
  const {onClick, children} = props;
  return (
    <a className={classes.root} onClick={onClick}>
      {children}
    </a>
  );
}
