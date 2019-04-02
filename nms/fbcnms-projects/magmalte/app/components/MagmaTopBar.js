/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import React, {useContext} from 'react';
import AccountCircle from '@material-ui/icons/AccountCircle';
import AppContext from './context/AppContext';
import ArrowDropDownIcon from '@material-ui/icons/ArrowDropDown';
import Divider from '@material-ui/core/Divider';
import MenuItem from '@material-ui/core/MenuItem';
import NetworkContext from './context/NetworkContext';
import TopBar from '@fbcnms/ui/components/layout/TopBar';
import TopBarAnchoredMenu from '@fbcnms/ui/components/layout/TopBarAnchoredMenu';
import {Link} from 'react-router-dom';

import {makeStyles} from '@material-ui/styles';

type Props = {
  children?: any,
};

const useStyles = makeStyles({
  link: {
    textDecoration: 'none',
  },
});

export default function MagmaTopBar(props: Props) {
  const appContext = useContext(AppContext);
  const classes = useStyles();
  const {networkId} = useContext(NetworkContext);
  return (
    <>
      <TopBar>
        {props.children}
        <div>
          <TopBarAnchoredMenu
            id="networks-appbar"
            buttonContent={
              <>
                {networkId} <ArrowDropDownIcon />
              </>
            }>
            {appContext.networkIds.map(id => (
              <Link className={classes.link} to={`/nms/${id}/map/`} key={id}>
                <MenuItem value={id} selected={id === networkId}>
                  {id}
                </MenuItem>
              </Link>
            ))}
            {appContext.user.isSuperUser && (
              <>
                <Divider />
                <Link className={classes.link} to="/nms/network/create">
                  <MenuItem>Create Network</MenuItem>
                </Link>
              </>
            )}
          </TopBarAnchoredMenu>
          <TopBarAnchoredMenu
            id="menu-appbar"
            buttonContent={<AccountCircle />}>
            <MenuItem disabled={true}>{appContext.user.email}</MenuItem>
            <Divider />
            <MenuItem href="/nms/user/logout" component="a">
              Logout
            </MenuItem>
          </TopBarAnchoredMenu>
        </div>
      </TopBar>
    </>
  );
}
