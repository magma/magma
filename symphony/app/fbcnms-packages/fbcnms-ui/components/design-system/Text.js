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
import symphony from '../../theme/symphony';
import {makeStyles} from '@material-ui/styles';

const styles = {
  h1: symphony.typography.h1,
  h2: symphony.typography.h2,
  h3: symphony.typography.h3,
  h4: symphony.typography.h4,
  h5: symphony.typography.h5,
  h6: symphony.typography.h6,
  subtitle1: symphony.typography.subtitle1,
  subtitle2: symphony.typography.subtitle2,
  subtitle3: symphony.typography.subtitle3,
  body1: symphony.typography.body1,
  body2: symphony.typography.body2,
  caption: symphony.typography.caption,
  overline: symphony.typography.overline,
  lightColor: {
    color: symphony.palette.white,
  },
  regularColor: {
    color: symphony.palette.D900,
  },
  primaryColor: {
    color: symphony.palette.primary,
  },
  grayColor: {
    color: symphony.palette.D400,
  },
  errorColor: {
    color: symphony.palette.R600,
  },
  warningColor: {
    color: symphony.palette.Y600,
  },
  inheritColor: {
    color: 'inherit',
  },
  lightWeight: {
    fontWeight: 300,
  },
  regularWeight: {
    fontWeight: 400,
  },
  mediumWeight: {
    fontWeight: 500,
  },
  boldWeight: {
    fontWeight: 600,
  },
  truncate: {
    display: 'block',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
  },
};
export const typographyStyles = makeStyles<Props, typeof styles>(() => styles);

type Props = {
  children: ?React.Node,
  variant?:
    | 'h1'
    | 'h2'
    | 'h3'
    | 'h4'
    | 'h5'
    | 'h6'
    | 'subtitle1'
    | 'subtitle2'
    | 'subtitle3'
    | 'body1'
    | 'body2'
    | 'caption'
    | 'overline',
  className?: string,
  useEllipsis?: ?boolean,
  weight?: 'inherit' | 'light' | 'regular' | 'medium' | 'bold',
  color?:
    | 'light'
    | 'regular'
    | 'primary'
    | 'error'
    | 'gray'
    | 'warning'
    | 'inherit',
};

const Text = (props: Props) => {
  const {
    children,
    variant = 'body1',
    className,
    color = 'regular',
    weight = 'inherit',
    useEllipsis = false,
    ...rest
  } = props;
  const classes = typographyStyles();
  return (
    <span
      {...rest}
      className={classNames(
        classes[variant],
        classes[`${color ?? 'regular'}Color`],
        classes[`${weight ? weight : 'regular'}Weight`],
        {[classes.truncate]: useEllipsis},
        className,
      )}>
      {children}
    </span>
  );
};

export default Text;
