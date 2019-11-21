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
import Text from '@fbcnms/ui/components/design-system/Text';
import {withStyles} from '@material-ui/core/styles';

type Props = WithStyles<typeof styles> & {
  text: ?string,
};

const styles = theme => ({
  root: {
    minWidth: '200px',
    padding: '10px',
  },
  label: {
    color: theme.palette.dark,
  },
});

class WorkOrderDetailsPaneItem extends React.Component<Props> {
  render() {
    const {classes} = this.props;
    return (
      <div className={classes.root}>
        <Text variant="body2" className={classes.label}>
          {this.props.text ?? ''}
        </Text>
      </div>
    );
  }
}

export default withStyles(styles)(WorkOrderDetailsPaneItem);
