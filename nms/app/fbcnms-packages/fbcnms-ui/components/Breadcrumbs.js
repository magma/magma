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
import React, {useMemo, useRef, useState} from 'react';
import Text from './design-system/Text';
import classNames from 'classnames';
// flowlint untyped-import:off
import fbt from 'fbt';
import useResize from './design-system/hooks/useResize';
import {gray8} from '@fbcnms/ui/theme/colors';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  breadcrumbs: {
    display: 'flex',
    alignItems: 'flex-start',
    overflow: 'hidden',
    flexGrow: 1,
    flexShrink: 1,
    flexBasis: 'auto',
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

type Props = {
  breadcrumbs: Array<BreadcrumbData>,
  className?: string,
  textClassName?: string,
  size?: 'default' | 'small' | 'large',
  variant?: TextVariant,
};

const KEEP_UNCOLLAPSED_AT_END = 1;

const Breadcrumbs = (props: Props) => {
  const {
    breadcrumbs,
    size = 'default',
    className,
    textClassName,
    variant,
  } = props;
  const classes = useStyles();
  const [isBreadcrumbsMenuOpen, toggleBreadcrumbsMenuOpen] = useState(false);
  const [anchorEl, setAnchorEl] = React.useState(null);
  const breadcrumbsContainer = useRef(null);
  const [collapsedBreadcrumbsCount, setCollapsedBreadcrumbsCount] = useState(0);
  const maximalCollapsedBreadcrumbsCount =
    breadcrumbs.length - KEEP_UNCOLLAPSED_AT_END;
  const ellipsisBreadcrumb = useMemo(
    () => ({
      id: 'collapsed',
      name: '...',
      subtext: fbt(
        'Click to see hidden parts',
        `tooltip for ellipsis control, showing hidden parts on click`,
      ),
      onClick: (_, clickTarget) => {
        toggleBreadcrumbsMenuOpen(true);
        setAnchorEl(clickTarget);
      },
    }),
    [],
  );
  useResize(breadcrumbsContainer, eventArgs => {
    if (eventArgs.width.expanded) {
      setCollapsedBreadcrumbsCount(0);
    }
    const parentContainer: ?HTMLElement = breadcrumbsContainer.current;
    if (!parentContainer) {
      return;
    }
    let breadcrumbsWidth = 0;
    for (
      let badcrumbIndex = 0;
      badcrumbIndex < parentContainer.children.length;
      badcrumbIndex++
    ) {
      breadcrumbsWidth += parentContainer.children[badcrumbIndex].clientWidth;
    }
    if (
      breadcrumbsWidth > parentContainer.clientWidth &&
      collapsedBreadcrumbsCount < maximalCollapsedBreadcrumbsCount
    ) {
      setCollapsedBreadcrumbsCount(collapsedBreadcrumbsCount + 1);
    }
  });

  const keepUncollapsedAtStart =
    collapsedBreadcrumbsCount > 0 &&
    collapsedBreadcrumbsCount < maximalCollapsedBreadcrumbsCount
      ? 1
      : 0;
  const [
    startBreadcrumbs,
    collapsedBreadcrumbs,
    endBreadcrumbs,
  ] = useMemo(() => {
    return [
      breadcrumbs.slice(0, keepUncollapsedAtStart),
      breadcrumbs.slice(
        keepUncollapsedAtStart,
        collapsedBreadcrumbsCount + keepUncollapsedAtStart,
      ),
      breadcrumbs.slice(
        keepUncollapsedAtStart + collapsedBreadcrumbsCount,
        breadcrumbs.length,
      ),
    ];
  }, [breadcrumbs, collapsedBreadcrumbsCount, keepUncollapsedAtStart]);

  return (
    <div
      aria-id="container"
      ref={breadcrumbsContainer}
      className={classNames(classes.breadcrumbs, className)}>
      {startBreadcrumbs.map(b => {
        return (
          <Breadcrumb
            key={b.id}
            data={b}
            size={size}
            className={textClassName}
            isLastBreadcrumb={false}
            variant={variant}
          />
        );
      })}
      {collapsedBreadcrumbs.length > 0 && (
        <Breadcrumb
          key={'collapsed'}
          data={ellipsisBreadcrumb}
          size={size}
          variant={variant}
          className={textClassName}
          isLastBreadcrumb={false}
        />
      )}
      {endBreadcrumbs.map((b, i) => {
        const isLast = i === endBreadcrumbs.length - 1;
        const isSingle =
          startBreadcrumbs.length === 0 && endBreadcrumbs.length === 1;
        return (
          <Breadcrumb
            key={b.id}
            data={b}
            isLastBreadcrumb={isLast}
            useEllipsis={isSingle}
            size={size}
            className={textClassName}
            variant={variant}
          />
        );
      })}
      <Popover
        open={isBreadcrumbsMenuOpen}
        anchorEl={anchorEl}
        onClose={() => {
          toggleBreadcrumbsMenuOpen(false);
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
              onClick={e => {
                b.onClick && b.onClick(b.id, e.currentTarget);
                toggleBreadcrumbsMenuOpen(false);
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

export default Breadcrumbs;
