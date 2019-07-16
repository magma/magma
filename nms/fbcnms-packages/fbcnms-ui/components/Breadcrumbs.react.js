/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {BreadcrumbData} from './Breadcrumb.react';

import Breadcrumb from './Breadcrumb.react';
import KeyboardArrowRightIcon from '@material-ui/icons/KeyboardArrowRight';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import MoreHorizIcon from '@material-ui/icons/MoreHoriz';
import Popover from '@material-ui/core/Popover';
import React, {useState} from 'react';
import Typography from '@material-ui/core/Typography';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  breadcrumbs: {
    display: 'flex',
    alignItems: 'flex-start',
  },
  moreIcon: {
    display: 'flex',
    alignItems: 'center',
  },
  moreIconButton: {
    cursor: 'pointer',
    '&:hover': {
      color: theme.palette.primary.main,
    },
  },
  arrowIcon: {
    color: theme.palette.grey[600],
  },
  collapsedBreadcrumbsList: {
    minWidth: '100px',
  },
  subtext: {
    fontSize: theme.typography.pxToRem(11),
    color: theme.palette.text.secondary,
    marginLeft: '8px',
  },
}));

const MAX_NUM_BREADCRUMBS = 3;

type Props = {
  breadcrumbs: Array<BreadcrumbData>,
  size?: 'default' | 'small' | 'large',
};

const Breadcrumbs = (props: Props) => {
  const {breadcrumbs, size} = props;
  const classes = useStyles();

  const [isBreadcrumbsMenuOpen, toggleBreadcrumbsMenuOpen] = useState(false);
  const [anchorEl, setAnchorEl] = React.useState(null);

  let collapsedBreadcrumbs = [];
  if (breadcrumbs.length > MAX_NUM_BREADCRUMBS) {
    collapsedBreadcrumbs = breadcrumbs.slice(1, breadcrumbs.length - 2);
  }
  const hasCollapsedBreadcrumbs = collapsedBreadcrumbs.length > 0;
  const startBreadcrumbs = hasCollapsedBreadcrumbs
    ? breadcrumbs.slice(0, 1)
    : [];
  const endBreadcrumbs = breadcrumbs.slice(
    collapsedBreadcrumbs.length + (hasCollapsedBreadcrumbs ? 1 : 0),
  );

  const arrowStyle = {
    height: size === 'large' || size === 'default' ? '32px' : '21px',
  };

  return (
    <div className={classes.breadcrumbs}>
      {startBreadcrumbs.map(b => (
        <Breadcrumb
          key={b.id}
          data={b}
          isLastBreadcrumb={false}
          size={size}
          onClick={b.onClick}
        />
      ))}
      {hasCollapsedBreadcrumbs && (
        <div className={classes.moreIcon} style={arrowStyle}>
          <MoreHorizIcon
            className={classes.moreIconButton}
            fontSize={size}
            onClick={e => {
              toggleBreadcrumbsMenuOpen(true);
              setAnchorEl(e.currentTarget);
            }}
          />
          <KeyboardArrowRightIcon
            className={classes.arrowIcon}
            style={arrowStyle}
            fontSize={size}
          />
        </div>
      )}
      {endBreadcrumbs.map((b, i) => (
        <Breadcrumb
          key={b.id}
          data={b}
          isLastBreadcrumb={i === endBreadcrumbs.length - 1}
          size={size}
        />
      ))}
      <Popover
        open={isBreadcrumbsMenuOpen}
        anchorEl={anchorEl}
        onClose={() => {
          toggleBreadcrumbsMenuOpen(false);
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
        <List className={classes.collapsedBreadcrumbsList}>
          {collapsedBreadcrumbs.map(b => (
            <ListItem key={`list_item_${b.id}`} button onClick={b.onClick}>
              <Typography>{b.name}</Typography>
              <Typography className={classes.subtext}>{b.subtext}</Typography>
            </ListItem>
          ))}
        </List>
      </Popover>
    </div>
  );
};

Breadcrumbs.defaultProps = {
  size: 'default',
  showTypes: true,
};

export default Breadcrumbs;
