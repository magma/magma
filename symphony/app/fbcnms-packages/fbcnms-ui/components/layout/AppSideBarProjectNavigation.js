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

import AppsIcon from '@material-ui/icons/Apps';
import Popout from '../Popout';
import React from 'react';
import Text from '../design-system/Text';
import classNames from 'classnames';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '../../hooks';

type Props = {
  projects: ProjectLink[],
};

const useClasses = makeStyles(theme => ({
  root: {
    marginTop: '24px',
  },
  menuPaper: {
    outline: 'none',
    overflowX: 'visible',
    overflowY: 'visible',
    padding: '10px 5px',
    position: 'absolute',
    '&:before': {
      borderBottom: '6px solid transparent',
      borderLeft: '6px solid transparent',
      borderRight: '6px solid transparent',
      borderTop: '6px solid white',
      bottom: '-12px',
      content: '""',
      left: '18px',
      position: 'absolute',
      zIndex: 10,
    },
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
      backgroundColor: theme.palette.grey[100],
    },
  },
  contentRoot: {
    padding: '10px 5px',
  },
  appsButton: {
    backgroundColor: theme.palette.common.white,
    borderRadius: '100%',
    color: theme.palette.blueDarkGray,
    cursor: 'pointer',
    display: 'flex',
    width: '36px',
    height: '36px',
    alignItems: 'center',
    justifyContent: 'center',
    '&:hover': {
      color: theme.palette.primary.main,
    },
  },
  popover: {
    '& div:not(:last-child)': {
      marginBottom: '8px',
    },
    '& $menuPaper': {
      boxShadow: '0px 0px 4px 0px rgba(0, 0, 0, 0.15)',
    },
  },
  menuItem: {
    '&:not(:last-child)': {
      marginBottom: '8px',
    },
    minWidth: '170px',
    borderRadius: '4px',
    cursor: 'pointer',
    padding: '8px 10px',
    '&:hover': {
      backgroundColor: 'rgba(145, 145, 145, 0.1)',
    },
    '&$selected': {
      backgroundColor: theme.palette.grey.A100,
    },
  },
  menuItemText: {
    fontSize: '12px',
    lineHeight: '16px',
  },
  selected: {
    '& $menuItemText': {
      color: theme.palette.primary.main,
    },
  },
}));

export default function AppSideBarProjectNavigation(props: Props) {
  const {projects} = props;
  const classes = useClasses();
  const {history, match} = useRouter();

  if (projects.length === 0) {
    return null;
  }

  const selected = projects.find(item => match.url.startsWith(item.url));

  return (
    <div className={classes.root}>
      <Popout
        content={
          <div className={classes.contentRoot}>
            {projects.map(item => (
              <div
                key={item.url}
                className={classNames({
                  [classes.menuItem]: true,
                  [classes.selected]: item.id === selected?.id,
                })}
                disabled={match.url.startsWith(item.url)}
                onClick={_event => history.push(item.url)}>
                <Text className={classes.menuItemText} variant="body2">
                  {item.secondary}
                </Text>
              </div>
            ))}
          </div>
        }>
        <div className={classes.appsButton}>
          <AppsIcon />
        </div>
      </Popout>
    </div>
  );
}
