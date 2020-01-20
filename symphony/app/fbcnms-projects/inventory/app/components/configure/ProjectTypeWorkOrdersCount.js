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
import WorkIcon from '@material-ui/icons/Work';
import classNames from 'classnames';
import symphony from '@fbcnms/ui/theme/symphony';
import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  root: {
    display: 'flex',
    flexDirection: 'row',
    flexGrow: 1,
  },
  text: {
    color: theme.palette.blueGrayDark,
    paddingLeft: '12px',
  },
  enabledIcon: {
    color: theme.palette.primary.main,
  },
  disabledIcon: {
    color: symphony.palette.disabled,
  },
});

type Props = {
  className?: string,
  count: number,
} & WithStyles<typeof styles>;

class ProjectTypeWorkOrdersCount extends React.Component<Props> {
  render() {
    const {classes, className, count} = this.props;
    return (
      <div className={classNames(classes.root, className)}>
        <WorkIcon
          className={count ? classes.enabledIcon : classes.disabledIcon}
        />
        <Text weight="medium" className={classes.text}>
          {this._getCountText(count)}
        </Text>
      </div>
    );
  }

  _getCountText(count: number) {
    switch (count) {
      case 0: {
        return 'No Work Orders';
      }
      case 1: {
        return '1 Work Order';
      }
      default: {
        return `${count} Work Orders`;
      }
    }
  }
}

export default withStyles(styles)(ProjectTypeWorkOrdersCount);
