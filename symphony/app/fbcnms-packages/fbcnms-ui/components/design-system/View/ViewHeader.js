/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {ButtonProps} from '../Button';
import type {IconButtonProps} from '../IconButton';
import type {PermissionHandlingProps} from '@fbcnms/ui/components/design-system/Form/FormAction';
import type {ToggleButtonProps} from '../ToggleButton/ToggleButtonGroup';

import * as React from 'react';
import Button from '@fbcnms/ui/components/design-system/Button';
import FormAction from '@fbcnms/ui/components/design-system/Form/FormAction';
import IconButton from '../IconButton';
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
  },
  actionButton: {
    '&:not(:first-child)': {
      marginLeft: '12px',
    },
  },
}));

export type ActionRegularButtonProps = {|
  title: React.Node,
  action: () => void,
  className?: ?string,
  ...PermissionHandlingProps,
  ...ButtonProps,
|};

export type ActionIconButtonProps = {|
  action: () => void,
  ...PermissionHandlingProps,
  ...IconButtonProps,
|};

export type ActionButtonProps =
  | ActionRegularButtonProps
  | ActionIconButtonProps;

export type ViewHeaderProps = $ReadOnly<{|
  title: React.Node,
  subtitle?: ?React.Node,
  searchBar?: ?React.Node,
  showMinimal?: ?boolean,
  className?: ?string,
|}>;

export type ViewHeaderActionsProps = $ReadOnly<{|
  actionButtons?: Array<ActionButtonProps>,
|}>;

export type ViewHeaderOptionsProps = $ReadOnly<{|
  viewOptions?: ToggleButtonProps,
|}>;

export type FullViewHeaderProps = $ReadOnly<{|
  ...ViewHeaderProps,
  ...ViewHeaderActionsProps,
  ...ViewHeaderOptionsProps,
|}>;

function FormHeaderAction(props: ActionButtonProps) {
  const {
    ignorePermissions,
    hideOnEditLock,
    disableOnFromError,
    action,
    icon = undefined,
    title = undefined,
    variant = undefined,
    disabled,
    tooltip,
    skin,
    className,
  } = props;
  const classes = useStyles();

  const buttonNode = () => {
    if (icon != null) {
      return (
        <IconButton
          className={classNames(classes.actionButton, className)}
          icon={icon}
          skin={skin}
          onClick={action}
        />
      );
    }
    if (title != null) {
      return (
        <Button
          className={classNames(classes.actionButton, className)}
          skin={skin}
          variant={variant}
          onClick={action}>
          {title}
        </Button>
      );
    }

    return null;
  };

  return (
    <FormAction
      ignorePermissions={ignorePermissions}
      hideOnEditLock={hideOnEditLock}
      disabled={disabled}
      tooltip={tooltip}
      disableOnFromError={disableOnFromError}>
      {buttonNode()}
    </FormAction>
  );
}

const ViewHeader = React.forwardRef<FullViewHeaderProps, HTMLElement>(
  (props, ref) => {
    const {
      title,
      subtitle,
      viewOptions,
      searchBar,
      showMinimal = false,
      className,
    } = props;
    const actionButtons: Array<ActionButtonProps> = props.actionButtons || [];
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
              {actionButtons.map((actionButton, index) => (
                <FormHeaderAction
                  key={`viewHeaderAction${index}`}
                  {...actionButton}
                />
              ))}
            </div>
          )}
        </div>
      </div>
    );
  },
);

export default ViewHeader;
