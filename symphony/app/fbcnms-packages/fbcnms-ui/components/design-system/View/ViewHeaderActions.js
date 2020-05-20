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
import type {
  ErrorHandlingProps,
  PermissionHandlingProps,
} from '@fbcnms/ui/components/design-system/Form/FormAction';
import type {IconButtonProps} from '../IconButton';
import type {OptionProps} from '../Select/SelectMenu';

import * as React from 'react';
import Button from '@fbcnms/ui/components/design-system/Button';
import FormAction from '@fbcnms/ui/components/design-system/Form/FormAction';
import IconButton from '../IconButton';
import PopoverMenu from '../Select/PopoverMenu';
import classNames from 'classnames';

import {withStyles} from '@material-ui/core/styles';
import type {WithStyles} from '@material-ui/core';

const styles = {
  actionButton: {
    '&:not(:first-child)': {
      marginLeft: '12px',
    },
  },
};

type BaseHeaderActionProps = $ReadOnly<{|
  tooltip?: string,
  disabled?: boolean,
  ...PermissionHandlingProps,
  ...ErrorHandlingProps,
|}>;

export type ButtonActionProps = $ReadOnly<{|
  action: () => void,
  className?: ?string,
  ...BaseHeaderActionProps,
  ...ButtonProps,
  children: React.Node,
|}> &
  WithStyles<typeof styles>;

class ButtonActionComponent extends React.Component<ButtonActionProps> {
  constructor(props: ButtonActionProps) {
    super(props);
  }

  render() {
    const {
      variant,
      action,
      skin = 'primary',
      classes,
      className,
      useEllipsis,
      children,
      ...formActionProps
    } = this.props;

    return (
      <FormAction {...formActionProps}>
        <Button
          className={classNames(classes.actionButton, className)}
          skin={skin}
          variant={variant}
          onClick={action}
          useEllipsis={useEllipsis}>
          {children}
        </Button>
      </FormAction>
    );
  }
}
export const ButtonAction = withStyles(styles)(ButtonActionComponent);

export type IconActionProps = $ReadOnly<{|
  action: () => void,
  ...BaseHeaderActionProps,
  ...IconButtonProps,
|}> &
  WithStyles<typeof styles>;

class IconActionComponent extends React.Component<IconActionProps> {
  constructor(props: IconActionProps) {
    super(props);
  }

  render() {
    const {
      icon,
      action,
      skin = 'primary',
      classes,
      className,
      ...formActionProps
    } = this.props;

    return (
      <FormAction {...formActionProps}>
        <IconButton
          className={classNames(classes.actionButton, className)}
          icon={icon}
          skin={skin}
          onClick={action}
        />
      </FormAction>
    );
  }
}
export const IconAction = withStyles(styles)(IconActionComponent);

export type OptionsActionProps = $ReadOnly<{|
  children: React.Node,
  options: Array<OptionProps<string>>,
  optionAction: (option: string) => void,
  className?: ?string,
  ...BaseHeaderActionProps,
|}> &
  WithStyles<typeof styles>;

class OptionsActionComponent extends React.Component<OptionsActionProps> {
  constructor(props: OptionsActionProps) {
    super(props);
  }

  render() {
    const {
      options,
      optionAction,
      skin = 'primary',
      classes,
      className,
      children,
      ...formActionProps
    } = this.props;

    return (
      <FormAction {...formActionProps}>
        <PopoverMenu
          className={classNames(classes.actionButton, className)}
          skin={skin}
          menuDockRight={true}
          options={options}
          onChange={optionAction}>
          {children}
        </PopoverMenu>
      </FormAction>
    );
  }
}
export const OptionsAction = withStyles(styles)(OptionsActionComponent);
