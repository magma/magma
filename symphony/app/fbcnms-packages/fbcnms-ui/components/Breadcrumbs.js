/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {BreadcrumbData} from './Breadcrumb';
import type {TextVariant} from '../theme/symphony';

import Breadcrumb from './Breadcrumb';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import Popover from '@material-ui/core/Popover';
import React, {useState} from 'react';
import Text from './design-system/Text';
import classNames from 'classnames';
import {gray8} from '@fbcnms/ui/theme/colors';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  breadcrumbs: {
    display: 'flex',
    alignItems: 'flex-start',
    minWidth: '200px',
  },
  moreIcon: {
    display: 'flex',
    alignItems: 'center',
  },
  moreIconButton: {
    cursor: 'pointer',
    color: gray8,
    '&:hover': {
      color: theme.palette.primary.main,
    },
  },
  collapsedBreadcrumbsList: {
    minWidth: '100px',
  },
  subtext: {
    fontSize: theme.typography.pxToRem(11),
    color: theme.palette.text.secondary,
    marginLeft: '8px',
  },
  largeText: {
    fontSize: '20px',
    lineHeight: '24px',
    fontWeight: 500,
  },
  smallText: {
    fontSize: '14px',
    lineHeight: '24px',
    fontWeight: 500,
  },
  slash: {
    color: gray8,
    margin: '0 6px',
  },
}));

const MAX_NUM_BREADCRUMBS = 3;

type Props = {
  breadcrumbs: Array<BreadcrumbData>,
  className?: string,
  textClassName?: string,
  size?: 'default' | 'small' | 'large',
  variant?: TextVariant,
};

const Breadcrumbs = (props: Props) => {
  const {breadcrumbs, size, className, textClassName, variant} = props;
  const classes = useStyles();

  const [isBreadcrumbsMenuOpen, toggleBreadcrumbsMenuOpen] = useState(false);
  const [anchorEl, setAnchorEl] = React.useState(null);

  let startBreadcrumbs = [];
  let collapsedBreadcrumbs = [];
  const endBreadcrumb = breadcrumbs[breadcrumbs.length - 1];
  if (breadcrumbs.length > MAX_NUM_BREADCRUMBS) {
    startBreadcrumbs = [breadcrumbs[0]];
    collapsedBreadcrumbs = breadcrumbs.slice(1, breadcrumbs.length - 1);
  } else {
    startBreadcrumbs = breadcrumbs.slice(0, breadcrumbs.length - 1);
  }

  return (
    <div className={classNames(classes.breadcrumbs, className)}>
      {startBreadcrumbs.map(b => (
        <Breadcrumb
          key={b.id}
          data={b}
          isLastBreadcrumb={false}
          size={size}
          onClick={b.onClick}
          className={textClassName}
          variant={variant}
        />
      ))}
      {collapsedBreadcrumbs.length > 0 && (
        <div className={classes.moreIcon}>
          <Text
            variant={variant ? variant : size === 'small' ? 'subtitle2' : 'h6'}
            className={classes.moreIconButton}
            onClick={e => {
              toggleBreadcrumbsMenuOpen(true);
              setAnchorEl(e.currentTarget);
            }}>
            {'...'}
          </Text>
          <Text
            variant={variant ? variant : size === 'small' ? 'subtitle2' : 'h6'}
            className={classes.slash}>
            {'/'}
          </Text>
        </div>
      )}
      {endBreadcrumb && (
        <Breadcrumb
          key={endBreadcrumb.id}
          data={endBreadcrumb}
          isLastBreadcrumb={true}
          useEllipsis={false}
          size={size}
          variant={variant}
          className={textClassName}
        />
      )}
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
            <ListItem
              key={`list_item_${b.id}`}
              button
              onClick={() => {
                b.onClick && b.onClick(b.id);
                toggleBreadcrumbsMenuOpen(false);
                setAnchorEl(null);
              }}>
              <Text>{b.name}</Text>
              <Text className={classes.subtext}>{b.subtext}</Text>
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
