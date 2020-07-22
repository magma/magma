/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ProjectLink} from './AppDrawerProjectNavigation';

import AppSideBarProjectNavigation from './AppSideBarProjectNavigation';
import ExpandButton from './ExpandButton';
import ProfileButton from '../ProfileButton';
import React, {useState} from 'react';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  root: {
    alignItems: 'center',
    backgroundColor: theme.palette.blueGrayDark,
    boxShadow: '1px 0px 0px 0px rgba(0, 0, 0, 0.1)',
    display: 'flex',
    flexDirection: 'column',
    height: '100vh',
    width: '80px',
    minWidth: '80px',
    padding: '20px 0px 20px 0px',
    position: 'relative',
  },
  mainItems: {
    display: 'flex',
    flexGrow: 1,
    flexDirection: 'column',
    alignItems: 'center',
    width: '100%',
    flexGrow: 1,
    paddingTop: '40px',
  },
  secondaryItems: {
    display: 'flex',
    flexGrow: 1,
    flexDirection: 'column',
    alignItems: 'center',
    width: '100%',
    justifyContent: 'flex-end',
  },
  expandButton: {
    display: 'flex',
    alignItems: 'center',
    width: '100%',
    justifyContent: 'center',
    flexGrow: 1,
  },
  visibleExpandButton: {
    visibility: 'visible',
  },
  hiddenExpandButton: {
    visibility: 'hidden',
  },
}));

type Props = {
  mainItems: any,
  secondaryItems?: any,
  showSettings: boolean,
  user: {
    email: string,
    isSuperUser: boolean,
  },
  projects?: ProjectLink[],
  useExpandButton: boolean,
  showExpandButton: boolean,
  expanded: boolean,
  onExpandClicked?: () => void,
};

const AppSideBar = (props: Props) => {
  const {
    user,
    expanded,
    showExpandButton,
    useExpandButton,
    onExpandClicked,
    mainItems,
    secondaryItems,
  } = props;
  const classes = useStyles();
  const projects = props.projects || [];
  const [hovered, setIsHovered] = useState(false);
  return (
    <div
      className={classes.root}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}>
      <div className={classes.mainItems}>{mainItems}</div>
      {useExpandButton && (
        <div className={classes.expandButton}>
          <div
            className={
              showExpandButton || hovered
                ? classes.visibleExpandButton
                : classes.hiddenExpandButton
            }>
            <ExpandButton
              expanded={expanded}
              onClick={() => onExpandClicked && onExpandClicked()}
            />
          </div>
        </div>
      )}
      <div className={classes.secondaryItems}>
        {secondaryItems}
        <ProfileButton showSettings={props.showSettings} user={user} />
      </div>
      <AppSideBarProjectNavigation projects={projects} />
    </div>
  );
};

AppSideBar.defaultProps = {
  useExpandButton: false,
  showExpandButton: false,
  expanded: false,
};

export default AppSideBar;
