/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import MoreHorizIcon from '@material-ui/icons/MoreHoriz';
import MoreVertIcon from '@material-ui/icons/MoreVert';
import Popover from '@material-ui/core/Popover';
import React, {useState} from 'react';
import Text from './design-system/Text';
import classNames from 'classnames';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  moreButton: {
    cursor: 'pointer',
    fill: theme.palette.primary.main,
  },
  itemText: {
    textTransform: 'capitalize',
  },
}));

type Props = {
  icon: 'horizontal' | 'vertical',
  items: Array<{
    name: string,
    onClick: () => void | Promise<void>,
  }>,
  iconClassName?: string,
};

const MoreActionsButton = (props: Props) => {
  const {items, icon, iconClassName} = props;
  const [anchorEl, setAnchorEl] = useState(null);
  const [isMenuOpen, toggleMenuOpen] = useState(false);
  const classes = useStyles();
  const MoreIcon = icon === 'horizontal' ? MoreHorizIcon : MoreVertIcon;
  return (
    <>
      <MoreIcon
        className={classNames(classes.moreButton, iconClassName)}
        onClick={e => {
          toggleMenuOpen(true);
          setAnchorEl(e.currentTarget);
        }}
      />
      <Popover
        open={isMenuOpen}
        anchorEl={anchorEl}
        onClose={() => {
          toggleMenuOpen(false);
          setAnchorEl(null);
        }}
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'center',
        }}
        transformOrigin={{
          vertical: 'top',
          horizontal: 'center',
        }}>
        <List>
          {items.map(item => (
            <ListItem
              key={`list_item_${item.name}`}
              button
              onClick={() => {
                item.onClick();
                toggleMenuOpen(false);
                setAnchorEl(null);
              }}>
              <Text className={classes.itemText}>{item.name}</Text>
            </ListItem>
          ))}
        </List>
      </Popover>
    </>
  );
};

MoreActionsButton.defaultProps = {
  icon: 'horizontal',
};

export default MoreActionsButton;
