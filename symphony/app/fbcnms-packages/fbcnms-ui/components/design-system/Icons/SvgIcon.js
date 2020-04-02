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
    minWidth: '24px',
    minHeight: '24px',
    width: '24px',
    height: '24px',
  },
  lightColor: {
    fill: symphony.palette.white,
  },
  regularColor: {
    fill: symphony.palette.secondary,
  },
  primaryColor: {
    fill: symphony.palette.primary,
  },
  grayColor: {
    fill: symphony.palette.D500,
  },
  errorColor: {
    fill: symphony.palette.R600,
  },
  inheritColor: {
    fill: 'inherit',
  },
}));

export type SvgIconStyleProps = {
  className?: string,
  color?: 'light' | 'regular' | 'primary' | 'error' | 'gray' | 'inherit',
};

type Props = {
  children: React.Node,
} & SvgIconStyleProps;

const SvgIcon = ({className, children, color = 'regular', ...rest}: Props) => {
  const classes = useStyles();
  return (
    <svg
      viewBox="0 0 24 24"
      width="24px"
      height="24px"
      className={classNames(classes.root, classes[`${color}Color`], className)}
      {...rest}>
      {children}
    </svg>
  );
};

export default SvgIcon;
