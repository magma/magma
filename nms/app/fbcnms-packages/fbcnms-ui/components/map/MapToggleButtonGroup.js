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
import ToggleButtonGroup from '@material-ui/core/Typography';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  toggleGroup: {
    boxShadow: '0px 0px 0.5px 0.5px grey',
    borderRadius: '4px',
  },
}));

type Props = {children: ?React.Node};

const MapToggleButtonGroup = (props: Props) => {
  const classes = useStyles();

  return (
    <ToggleButtonGroup className={classes.toggleGroup}>
      {props.children}
    </ToggleButtonGroup>
  );
};

export default MapToggleButtonGroup;
