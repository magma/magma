/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import type {ToggleButtonProps} from '../ToggleButton/ToggleButtonGroup';

import * as React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import ToggleButton from '../ToggleButton/ToggleButtonGroup';
import classNames from 'classnames';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  root: {
    display: 'flex',
    flexDirection: 'row',
    padding: '16px 24px',
    paddingBottom: '8px',
  },
  column: {
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'space-between',
    flexShrink: '0',
    '&:not(:last-child)': {
      paddingRight: '8px',
    },
  },
  expandedColumn: {
    flexGrow: '1',
    flexShrink: '1',
  },
  title: {
    paddingTop: '4px',
  },
  collapsablePart: {
    maxHeight: '200px',
    overflow: 'hidden',
    transition: 'max-height 500ms ease-out 0s',
  },
  collapsed: {
    maxHeight: '0px',
  },
  searchBarContainer: {
    paddingTop: '8px',
  },
  viewOptionsContainer: {
    flexGrow: 1,
    display: 'flex',
    justifyContent: 'flex-end',
    paddingBottom: '8px',
  },
  groupButtons: {
    display: 'flex',
    justifyContent: 'flex-end',
  },
  buttonContent: {
    paddingTop: '4px',
  },
  actionButtons: {
    flexGrow: 1,
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'flex-end',
    '&>:not(:first-child)': {
      marginLeft: '12px',
    },
  },
}));

export type ViewHeaderProps = $ReadOnly<{|
  title: React.Node,
  subtitle?: ?React.Node,
  searchBar?: ?React.Node,
  showMinimal?: ?boolean,
  className?: ?string,
|}>;

export type ViewHeaderActionsProps = $ReadOnly<{|
  actionButtons?: $ReadOnlyArray<React.Node>,
|}>;

export type ViewHeaderOptionsProps = $ReadOnly<{|
  viewOptions?: ToggleButtonProps,
|}>;

export type FullViewHeaderProps = $ReadOnly<{|
  ...ViewHeaderProps,
  ...ViewHeaderActionsProps,
  ...ViewHeaderOptionsProps,
|}>;

const ViewHeader = React.forwardRef<FullViewHeaderProps, HTMLElement>(
  (props, ref) => {
    const {
      title,
      subtitle,
      actionButtons,
      viewOptions,
      searchBar,
      showMinimal = false,
      className,
    } = props;
    const classes = useStyles();

    return (
      <div className={classNames(classes.root, className)} ref={ref}>
        <div className={classNames(classes.column, classes.expandedColumn)}>
          <Text variant="h6" className={classes.title}>
            {title}
          </Text>
          <div
            className={classNames(classes.collapsablePart, {
              [classes.collapsed]: showMinimal,
            })}>
            <Text variant="body2" color="gray">
              {subtitle}
            </Text>
            {searchBar != null && (
              <div className={classes.searchBarContainer}>{searchBar}</div>
            )}
          </div>
        </div>
        <div className={classes.column}>
          {viewOptions != null && (
            <div className={classes.viewOptionsContainer}>
              <ToggleButton {...viewOptions} />
            </div>
          )}
          {actionButtons != null && (
            <div
              className={classNames(
                classes.actionButtons,
                classes.collapsablePart,
                {
                  [classes.collapsed]: showMinimal,
                },
              )}>
              {actionButtons}
            </div>
          )}
        </div>
      </div>
    );
  },
);

export default ViewHeader;
