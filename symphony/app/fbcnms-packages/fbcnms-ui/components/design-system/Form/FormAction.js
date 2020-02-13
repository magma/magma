/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import * as React from 'react';
import FormElementContext from './FormElementContext';
import FormValidationContext from '../Form/FormValidationContext';
import {useContext} from 'react';

export type PermissionHandlingProps = {|
  ignorePermissions?: boolean,
  hideWhenDisabled?: boolean,
|};

type Props = {
  children: React.Node,
  disabled?: boolean,
  ...PermissionHandlingProps,
};

const FormAction = (props: Props) => {
  const {
    children,
    disabled: disabledProp = false,
    ignorePermissions = false,
    hideWhenDisabled = true,
  } = props;

  const validationContext = useContext(FormValidationContext);
  const validationDisabling =
    validationContext.editLock.detected && !ignorePermissions;
  const disabled = disabledProp || validationDisabling;
  const shouldHide = validationDisabling && hideWhenDisabled;
  return (
    <FormElementContext.Provider value={{disabled}}>
      {(!shouldHide && children) || null}
    </FormElementContext.Provider>
  );
};

export default FormAction;
