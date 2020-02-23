/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

// import type {OptionsToggleProps} from '../Buttons/OptionsToggleButton';
import type {PermissionHandlingProps} from '@fbcnms/ui/components/design-system/Form/FormAction';
import type {ToggleButtonProps} from '../ToggleButton/ToggleButtonGroup';

import * as React from 'react';
import Button from '@fbcnms/ui/components/design-system/Button';
import FormAction from '@fbcnms/ui/components/design-system/Form/FormAction';
import Text from '@fbcnms/ui/components/design-system/Text';
import ToggleButton from '../ToggleButton/ToggleButtonGroup';
import classNames from 'classnames';
import {FormValidationContextProvider} from '@fbcnms/ui/components/design-system/Form/FormValidationContext';
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
    alignItems: 'flex-end',
    justifyContent: 'center',
  },
  actionButton: {
    '&:not(:first-child)': {
      marginLeft: '8px',
    },
  },
}));

export type ActionButtonProps = {
  title: string,
  action: () => void,
  ...PermissionHandlingProps,
};

export type ViewHeaderProps = {
  title: React.Node,
  subtitle?: ?React.Node,
  searchBar?: ?React.Node,
  showMinimal?: ?boolean,
};

export type ViewHeaderActionsProps = {
  actionButtons?: Array<ActionButtonProps>,
};

export type ViewHeaderOptionsProps = {
  viewOptions?: ToggleButtonProps,
};

export type FullViewHeaderProps = ViewHeaderProps &
  ViewHeaderActionsProps &
  ViewHeaderOptionsProps;

const ViewHeader = (props: FullViewHeaderProps) => {
  const {title, subtitle, viewOptions, searchBar, showMinimal = false} = props;
  const actionButtons: Array<ActionButtonProps> = props.actionButtons || [];
  const classes = useStyles();

  return (
    <div className={classes.root}>
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
          <FormValidationContextProvider>
            <div className={classes.actionButtons}>
              {actionButtons.map(actionButton => (
                <FormAction
                  key={actionButton.title}
                  ignorePermissions={actionButton.ignorePermissions}
                  hideWhenDisabled={actionButton.hideWhenDisabled}>
                  <Button
                    className={classes.actionButton}
                    onClick={actionButton.action}>
                    {actionButton.title}
                  </Button>
                </FormAction>
              ))}
            </div>
          </FormValidationContextProvider>
        )}
      </div>
    </div>
  );
};

export default ViewHeader;
