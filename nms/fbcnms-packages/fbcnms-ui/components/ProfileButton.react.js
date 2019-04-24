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
import classNames from 'classnames';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import Popout from '@fbcnms/ui/components/Popout.react';
import ProfileIcon from '../icons/ProfileIcon.react';
import React, {useState} from 'react';
import Typography from '@material-ui/core/Typography';

const useStyles = makeStyles(theme => ({
  accountButton: {
    backgroundColor: theme.palette.common.white,
    width: '28px',
    height: '28px',
    fontSize: '28px',
    cursor: 'pointer',
    borderRadius: '100%',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    border: `1px solid ${theme.palette.common.white}`,
    '&:hover, &$openButton': {
      border: `1px solid ${theme.palette.blue60}`,
      '& $accountButtonIcon': {
        fill: theme.palette.blue60,
      },
    },
  },
  openButton: {},
  accountButtonIcon: {
    '&&': {
      fill: theme.palette.grey.A200,
      fontSize: '15px',
    },
  },
  itemGutters: {
    '&&': {
      minWidth: '170px',
      borderRadius: '4px',
      padding: '8px 10px',
      '&:hover': {
        backgroundColor: 'rgba(145, 145, 145, 0.1)',
      },
    },
  },
  profileList: {
    '&&': {
      padding: '10px 5px',
    },
  },
  profileItemText: {
    fontSize: '12px',
    lineHeight: '16px',
  },
}));

type Props = {
  user: {
    email: string,
    isSuperUser: boolean,
  },
};

const ProfileButton = (props: Props) => {
  const {user} = props;
  const {email, isSuperUser} = user;
  const {relativeUrl} = useRouter();
  const classes = useStyles();
  const [isProfileMenuOpen, toggleProfileMenu] = useState(false);

  return (
    <Popout
      content={
        <List component="nav" className={classes.profileList}>
          <ListItem classes={{gutters: classes.itemGutters}} disabled={true}>
            <Typography className={classes.profileItemText}>{email}</Typography>
          </ListItem>
          {isSuperUser && (
            <ListItem
              classes={{gutters: classes.itemGutters}}
              button
              href={relativeUrl('/settings')}
              component="a">
              <Typography className={classes.profileItemText}>
                Settings
              </Typography>
            </ListItem>
          )}
          <ListItem
            classes={{gutters: classes.itemGutters}}
            button
            href="/user/logout"
            component="a">
            <Typography className={classes.profileItemText}>Logout</Typography>
          </ListItem>
        </List>
      }
      onOpen={() => toggleProfileMenu(true)}
      onClose={() => toggleProfileMenu(false)}>
      <div
        className={classNames({
          [classes.accountButton]: true,
          [classes.openButton]: isProfileMenuOpen,
        })}>
        <ProfileIcon className={classes.accountButtonIcon} />
      </div>
    </Popout>
  );
};

export default ProfileButton;
