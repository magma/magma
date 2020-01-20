/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {WithStyles} from '@material-ui/core';

import * as React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import classNames from 'classnames';
import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  root: {
    color: theme.palette.grey[900],
    fontWeight: 500,
    fontSize: '20px',
    lineHeight: '24px',
  },
});

type Props = {
  title: string,
  children: React.ChildrenArray<null | React.Element<*>>,
  className?: string,
} & WithStyles<typeof styles>;

const CardSection = (props: Props) => {
  const {className, classes, children} = props;
  return (
    <div>
      <Text className={classNames(classes.root, className)}>{props.title}</Text>
      {children}
    </div>
  );
};

export default withStyles(styles)(CardSection);
