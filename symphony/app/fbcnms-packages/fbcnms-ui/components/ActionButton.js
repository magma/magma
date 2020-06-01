/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import AddIcon from './design-system/Icons/Actions/AddIcon';
import IconButton from './design-system/IconButton';
import React from 'react';
import RemoveIcon from './design-system/Icons/Actions/RemoveIcon';

type Props = $ReadOnly<{|
  action: 'add' | 'remove',
  onClick: () => void,
|}>;

const ActionButton = (props: Props) => {
  const {action, onClick} = props;

  switch (action) {
    case 'add':
      return <IconButton icon={AddIcon} onClick={onClick} />;
    case 'remove':
      return <IconButton icon={RemoveIcon} skin="regular" onClick={onClick} />;
    default:
      return null;
  }
};

export default ActionButton;
