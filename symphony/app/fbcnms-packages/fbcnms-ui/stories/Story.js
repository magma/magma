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
import Text from '../components/design-system/Text';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  root: {
    padding: '16px',
  },
  header: {
    marginBottom: '24px',
  },
}));

type Props = $ReadOnly<{|
  name: React.Node,
  children: React.Node,
|}>;

const Story = ({name, children}: Props) => {
  const classes = useStyles();
  return (
    <div className={classes.root}>
      <div className={classes.header}>
        <Text variant="h3">{name}</Text>
      </div>
      {children}
    </div>
  );
};

export default Story;
