/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';
import * as React from 'react';
import {Link} from 'react-router-dom';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';

const useStyles = makeStyles(theme => ({
  icon: {
    color: theme.palette.grey[400],
  },
  link: {
    textDecoration: 'none',
  },
  selectedIcon: {
    color: theme.palette.common.white,
  },
  selectedRow: {
    '&:hover': {
      backgroundColor: '#2e3c42',
    },
  },
}));

type Props = {
  path: string,
  label: string,
  icon: any,
};

export default function NavListItem(props: Props) {
  const classes = useStyles();
  const router = useRouter();
  const isSelected = router.location.pathname.includes(props.path);
  const iconClass = isSelected ? classes.selectedIcon : classes.icon;

  return (
    <Link to={props.path} className={classes.link}>
      <ListItem button className={classes.selectedRow}>
        <ListItemIcon>
          <div className={iconClass}>{props.icon}</div>
        </ListItemIcon>
        <ListItemText classes={{primary: iconClass}} primary={props.label} />
      </ListItem>
    </Link>
  );
}
