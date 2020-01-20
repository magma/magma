/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import * as React from 'react';
import ReactDOM from 'react-dom';

type Props = {
  children: React.Node,
  target: ?HTMLElement,
};

const Portal = ({children, target}: Props) => {
  return target != null ? ReactDOM.createPortal(children, target) : null;
};

export default Portal;
