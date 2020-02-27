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
import Text from '@fbcnms/ui/components/design-system/Text';
import classNames from 'classnames';
import symphony from '@fbcnms/ui/theme/symphony';
import {makeStyles} from '@material-ui/styles';

const MAJOR_SIZE = '240px';

const useStyles = makeStyles(() => ({
  root: {
    flexBasis: MAJOR_SIZE,
    flexGrow: 0,
    flexShrink: 0,
    display: 'flex',
    flexDirection: 'column',
    overflow: 'hidden',
    flexBasis: MAJOR_SIZE,
    width: MAJOR_SIZE,
    minWidth: MAJOR_SIZE,
    maxWidth: MAJOR_SIZE,
  },
  menuHeader: {
    flexGrow: 0,
    padding: '24px 16px 0 16px',
    backgroundColor: symphony.palette.white,
  },
  menuItemsContainer: {
    flexGrow: 1,
    display: 'flex',
    flexDirection: 'column',
    padding: '16px',
    paddingBottom: '8px',
    backgroundColor: symphony.palette.white,
    overflowX: 'hidden',
    overflowY: 'auto',
  },
  menuItem: {
    flexShrink: 0,
    padding: '8px 16px',
    border: '1px solid transparent',
    borderRadius: '4px',
    '&:hover': {
      backgroundColor: symphony.palette.background,
      cursor: 'pointer',
    },
    '&:not(:last-child)': {
      marginBottom: '2px',
    },
  },
  activeItem: {
    backgroundColor: symphony.palette.background,
  },
}));

export type MenuItem = {
  label: React.Node,
  tooltip?: ?string,
};

type Props = {
  header?: ?React.Node,
  items: Array<MenuItem>,
  activeItemIndex?: number,
  onActiveItemChanged: (activeItem: MenuItem, activeItemIndex: number) => void,
};

export default function SideMenu(props: Props) {
  const {header, items, activeItemIndex, onActiveItemChanged} = props;
  const classes = useStyles();

  return (
    <div className={classes.root}>
      <Text variant="body1" className={classes.menuHeader}>
        {header}
      </Text>
      <div className={classes.menuItemsContainer}>
        {items.map((item, itemIndex) => (
          <div
            title={item.tooltip}
            onClick={() =>
              onActiveItemChanged &&
              onActiveItemChanged(items[itemIndex], itemIndex)
            }
            className={classNames(classes.menuItem, {
              [classes.activeItem]: activeItemIndex === itemIndex,
            })}>
            <Text
              variant="body1"
              useEllipsis={true}
              color={activeItemIndex === itemIndex ? 'primary' : 'gray'}>
              {item.label}
            </Text>
          </div>
        ))}
      </div>
    </div>
  );
}
