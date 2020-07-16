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
import classNames from 'classnames';
import symphony from '../../../theme/symphony';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  root: {
    padding: '8px 0px',
    backgroundColor: symphony.palette.white,
    boxShadow: symphony.shadows.DP3,
    borderRadius: '4px',
    maxHeight: '322px',
    overflowY: 'auto',
    minWidth: '112px',
    maxWidth: '360px',
  },
  fullWidth: {
    width: '100%',
  },
  normalWidth: {
    width: '236px',
  },
}));

export type MenuBaseProps = $ReadOnly<{|
  className?: string,
  size?: 'normal' | 'full',
|}>;

type Props = $ReadOnly<{|
  children: React.Node,
  ...MenuBaseProps,
|}>;

const MenuBase = ({children, className, size = 'normal'}: Props) => {
  const classes = useStyles();
  return (
    <div
      className={classNames(classes.root, className, {
        [classes.fullWidth]: size === 'full',
        [classes.normalWidth]: size === 'normal',
      })}>
      {children}
    </div>
  );
};

export default MenuBase;
