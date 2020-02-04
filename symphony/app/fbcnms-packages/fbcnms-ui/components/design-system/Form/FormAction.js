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
import FormElementContext from './FormElementContext';
import FormValidationContext from '../Form/FormValidationContext';
import {useContext, useMemo} from 'react';

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
  const disabled = useMemo(
    () =>
      disabledProp ||
      (validationContext.editLock.detected && !ignorePermissions),
    [disabledProp, ignorePermissions, validationContext.editLock.detected],
  );
  const shouldHide = useMemo(() => disabled && hideWhenDisabled, [
    disabled,
    hideWhenDisabled,
  ]);
  return (
    <FormElementContext.Provider value={{disabled}}>
      {(!shouldHide && children) || null}
    </FormElementContext.Provider>
  );
};

export default FormAction;
