/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import React from 'react';
import Text from '../Text';
import classNames from 'classnames';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(_theme => ({
  root: {
    marginBottom: '24px',
  },
}));

type Props = {
  className?: string,
  children: string,
};

const CardHeader = (props: Props) => {
  const {children, className} = props;
  const classes = useStyles();
  return (
    <div className={classNames(classes.root, className)}>
      <Text variant="h6">{children}</Text>
    </div>
  );
};

export default CardHeader;
