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
import AddCircleOutlineIcon from '@material-ui/icons/AddCircleOutline';
import Dialog from '@material-ui/core/Dialog';
import IconButton from '@material-ui/core/IconButton';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import {makeStyles} from '@material-ui/styles';
import {useRef} from 'react';
import {useState} from 'react';

export type ServiceMenuItem = {
  label: string,
  onClick: () => void,
};

type Props = {
  items: Array<ServiceMenuItem>,
  isOpen: boolean,
  onClose: () => void,
  children: React.Node,
};

const useStyles = makeStyles({
  dialog: {
    width: '80%',
    maxWidth: '1280px',
    height: '90%',
    maxHeight: '800px',
  },
});

const ServiceMenu = (props: Props) => {
  const classes = useStyles();
  const {items, isOpen, onClose, children} = props;
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const anchorElRef = useRef<?HTMLElement>(null);

  return (
    <>
      <IconButton
        buttonRef={anchorElRef}
        className={classes.addButton}
        onClick={() => setIsMenuOpen(true)}>
        <AddCircleOutlineIcon />
      </IconButton>
      <Menu
        anchorEl={anchorElRef?.current}
        keepMounted
        open={isMenuOpen}
        onClose={() => setIsMenuOpen(false)}>
        {items.map(item => (
          <MenuItem
            onClick={() => {
              item.onClick();
              setIsMenuOpen(false);
            }}>
            {item.label}
          </MenuItem>
        ))}
      </Menu>
      <Dialog
        open={isOpen}
        onClose={onClose}
        maxWidth={false}
        fullWidth={true}
        classes={{paperFullWidth: classes.dialog}}>
        {children}
      </Dialog>
    </>
  );
};

export default ServiceMenu;
