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
import symphony from '../../../theme/symphony';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles({
  root: {
    minWidth: '24px',
    minHeight: '24px',
    width: '24px',
    height: '24px',
  },
  regularColor: {
    fill: symphony.palette.secondary,
  },
  lightColor: {
    fill: symphony.palette.white,
  },
  regularColor: {
    fill: symphony.palette.D900,
  },
  primaryColor: {
    fill: symphony.palette.primary,
  },
  grayColor: {
    fill: symphony.palette.D400,
  },
  errorColor: {
    fill: symphony.palette.R600,
  },
});

export type SvgIconStyleProps = {
  className?: string,
  color?: 'light' | 'regular' | 'primary' | 'error' | 'gray',
};

type Props = {
  children: React.Node,
} & SvgIconStyleProps;

const SvgIcon = ({className, children, color, ...rest}: Props) => {
  const classes = useStyles();
  return (
    <svg
      viewBox="0 0 24 24"
      width="24px"
      height="24px"
      className={classNames(
        classes.root,
        classes[`${color ?? 'regular'}Color`],
        className,
      )}
      {...rest}>
      {children}
    </svg>
  );
};

export default SvgIcon;
