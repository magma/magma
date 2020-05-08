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
import symphony from '@fbcnms/ui/theme/symphony';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  root: {
    display: 'flex',
    flexDirection: 'column',
    width: '100%',
    minWidth: '240px',
  },
  item: {
    borderBottom: `1px solid ${symphony.palette.separatorLight}`,
  },
}));

type Props<T> = $ReadOnly<{|
  items: $ReadOnlyArray<T>,
  emptyState?: ?React.Node,
  className?: ?string,
  children: T => React.Node,
|}>;

export default function List<T>(props: Props<T>) {
  const {items, emptyState, className, children} = props;
  const classes = useStyles();

  return (
    <div className={classNames(classes.root, className)}>
      {items.length == 0 && emptyState != null
        ? emptyState
        : items.map(item => (
            <div className={classes.item}>{children(item)}</div>
          ))}
    </div>
  );
}
