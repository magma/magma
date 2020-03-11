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

const useStyles = makeStyles(_theme => ({
  root: {
    padding: '24px',
    backgroundColor: symphony.palette.white,
    boxShadow: symphony.shadows.DP1,
    borderRadius: '4px',
  },
}));

type Props = {
  className?: string,
  children: React.Node,
};

const Card = (props: Props) => {
  const {children, className} = props;
  const classes = useStyles();
  return <div className={classNames(classes.root, className)}>{children}</div>;
};

export default Card;
