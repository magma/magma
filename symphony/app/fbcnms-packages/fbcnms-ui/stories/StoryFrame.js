/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import * as React from 'react';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  root: {
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'space-around',
    alignItems: 'center',
  },
}));

type Props = $ReadOnly<{|
  children: React.Node,
|}>;

const StoryFrame = ({children}: Props) => {
  const classes = useStyles();
  return <div className={classes.root}>{children}</div>;
};

export default StoryFrame;
