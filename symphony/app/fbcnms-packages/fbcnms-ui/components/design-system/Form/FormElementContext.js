/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import React from 'react';

export type FormElementContextValue = {
  disabled: boolean,
  hasError?: boolean,
};

const FormElementContext = React.createContext<FormElementContextValue>({
  disabled: false,
  hasError: false,
});

export function useFormElementContext() {
  return React.useContext(FormElementContext);
}

export default FormElementContext;
