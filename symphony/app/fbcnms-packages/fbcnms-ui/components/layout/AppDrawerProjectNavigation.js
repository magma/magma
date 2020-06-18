/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {Tab} from '@fbcnms/types/tabs';

import ArrowDropDownIcon from '@material-ui/icons/ArrowDropDown';
import IconButton from '@material-ui/core/IconButton';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemSecondaryAction from '@material-ui/core/ListItemSecondaryAction';
import ListItemText from '@material-ui/core/ListItemText';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import React, {useState} from 'react';

import {makeStyles} from '@material-ui/styles';
import {useRouter} from '../../hooks';

export type ProjectLink = {
  id: Tab,
  name: string,
  secondary: string,
  url: string,
};

type Props = {
  projects: ProjectLink[],
};

const useClasses = makeStyles(theme => ({
  menu: {
    minWidth: 300,
  },
  primary: {
    color: theme.palette.common.white,
  },
  secondary: {
    color: theme.palette.common.white,
  },
  selectedRow: {
    borderBottom: '1px solid ' + theme.palette.grey[400],
    borderTop: '1px solid ' + theme.palette.grey[400],
    '&:hover': {
      backgroundColor: '#2e3c42',
    },
  },
}));

export default function AppDrawerProjectNavigation(props: Props) {
  const {projects} = props;
  const classes = useClasses();
  const {history, match} = useRouter();
  const [anchorEl, setAnchorEl] = useState(null);

  const selected = projects.find(item => match.url.startsWith(item.url));

  return (
    <div>
      <List component="nav">
        <ListItem
          button
          className={classes.selectedRow}
          aria-haspopup="true"
          aria-controls="navigation-menu"
          aria-label="When device is locked"
          onClick={({currentTarget}) => setAnchorEl(currentTarget)}>
          <ListItemText
            classes={{
              primary: classes.primary,
              secondary: classes.secondary,
            }}
            primary={selected?.name}
            secondary={selected?.secondary}
          />
          <ListItemSecondaryAction>
            <IconButton
              onClick={({currentTarget}) => setAnchorEl(currentTarget)}>
              <ArrowDropDownIcon />
            </IconButton>
          </ListItemSecondaryAction>
        </ListItem>
      </List>
      <Menu
        className={classes.menu}
        id="navigation-menu"
        anchorEl={anchorEl}
        open={!!anchorEl}
        onClose={() => setAnchorEl(null)}>
        {projects.map(item => (
          <MenuItem
            key={item.url}
            disabled={match.url.startsWith(item.url)}
            selected={match.url.startsWith(item.url)}
            onClick={_event => history.push(item.url)}>
            {item.name} - {item.secondary}
          </MenuItem>
        ))}
      </Menu>
    </div>
  );
}
