/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ButtonSkin} from './Button';
import type {SvgIconStyleProps} from './Icons/SvgIcon';

import * as React from 'react';
import Button from './Button';

export type IconComponent = React.ComponentType<SvgIconStyleProps>;

export type IconButtonProps = $ReadOnly<{|
  className?: string,
  icon: IconComponent,
  skin?: ButtonSkin,
  disabled?: boolean,
  tooltip?: string,
|}>;

type Props = $ReadOnly<{|
  onClick?:
    | void
    | (void | ((SyntheticMouseEvent<HTMLElement>) => void | Promise<void>)),
  ...IconButtonProps,
|}>;

const IconButton = ({icon: Icon, ...buttonProps}: Props) => {
  return (
    <Button variant="text" {...buttonProps}>
      <Icon color="inherit" />
    </Button>
  );
};

export default IconButton;
