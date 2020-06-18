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
import Text from '../Text';
import classNames from 'classnames';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(_theme => ({
  root: {
    display: 'flex',
    alignItems: 'center',
    marginBottom: '24px',
  },
  titleText: {
    flexGrow: 1,
  },
}));

type Props = {
  className?: string,
  children: string,
  rightContent?: React.Node,
};

const CardHeader = (props: Props) => {
  const {children, className, rightContent} = props;
  const classes = useStyles();
  return (
    <div className={classNames(classes.root, className)}>
      <Text variant="h6" className={classes.titleText}>
        {children}
      </Text>
      {rightContent}
    </div>
  );
};

export default CardHeader;
