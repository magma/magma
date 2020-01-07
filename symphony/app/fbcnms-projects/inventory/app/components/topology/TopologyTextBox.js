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

import React from 'react';
import symphony from '@fbcnms/ui/theme/symphony';
import {withStyles} from '@material-ui/core/styles';

const TEXT_FONT_SIZE = 12;

type Props = {
  transform: string,
  text: ?string,
} & WithStyles<typeof styles>;

const styles = _ => ({
  nodeText: {
    ...symphony.typography.caption,
    fontWeight: 500,
    fill: symphony.palette.D900,
    fontSize: TEXT_FONT_SIZE,
    cursor: 'pointer',
    pointerEvents: 'none',
    stroke: 'none',
    textAnchor: 'middle',
  },
  nodeRect: {
    fill: symphony.palette.D10,
  },
});

const TopologyTextBox = (props: Props) => {
  const {transform, text, classes} = props;

  return (
    <g id="textBox" transform={transform}>
      <rect rx={10} ry={10} className={classes.nodeRect} />
      <text className={classes.nodeText}>{text ?? ''}</text>
    </g>
  );
};

export default withStyles(styles)(TopologyTextBox);
