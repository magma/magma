/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import MoreHorizIcon from '@material-ui/icons/MoreHoriz';
import Popover from '@material-ui/core/Popover';
import React, {useState} from 'react';
import Typography from '@material-ui/core/Typography';
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
  items: Array<{
    name: string,
    onClick: () => void,
  }>,
};

export default function MoreActionsButton(props: Props) {
  const {items} = props;
  const [anchorEl, setAnchorEl] = useState(null);
  const [isMenuOpen, toggleMenuOpen] = useState(false);
  const classes = useStyles();
  return (
    <>
      <MoreHorizIcon
        className={classes.moreButton}
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
              <Typography className={classes.itemText}>{item.name}</Typography>
            </ListItem>
          ))}
        </List>
      </Popover>
    </>
  );
}
