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
import classNames from 'classnames';
import symphony from '../theme/symphony';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  root: {
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'space-around',
    alignItems: 'center',
  },
  lightBackground: {
    backgroundColor: symphony.palette.white,
    padding: '6px',
  },
  stretchContents: {
    justifyContent: 'flex-start',
  },
}));

type Props = $ReadOnly<{|
  children: React.Node,
  background?: 'regular' | 'light',
  stretchContents?: boolean,
|}>;

const StoryFrame = ({
  children,
  background = 'regular',
  stretchContents = false,
}: Props) => {
  const classes = useStyles();
  return (
    <div
      className={classNames(classes.root, {
        [classes.lightBackground]: background === 'light',
        [classes.stretchContents]: stretchContents === true,
      })}>
      {children}
    </div>
  );
};

export default StoryFrame;
