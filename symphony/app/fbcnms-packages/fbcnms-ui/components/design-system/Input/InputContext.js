/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import React from 'react';

export type InputContextValue = {
  disabled: boolean,
  value: string | number,
};

const InputContext = React.createContext<InputContextValue>({
  disabled: false,
  value: '',
});

export function useInput() {
  return React.useContext(InputContext);
}

export default InputContext;
