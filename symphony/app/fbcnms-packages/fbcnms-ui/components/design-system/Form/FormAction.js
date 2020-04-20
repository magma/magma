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
import FormAlertsContext from '../Form/FormAlertsContext';
import FormElementContext from './FormElementContext';
import {joinNullableStrings} from '@fbcnms/util/strings';
import {useContext, useMemo} from 'react';

export type PermissionHandlingProps = {|
  ignorePermissions?: ?boolean,
  hideOnEditLock?: ?boolean,
  disableOnFromError?: ?boolean,
|};

type Props = {
  children: React.Node,
  disabled?: boolean,
  tooltip?: ?string,
  ...PermissionHandlingProps,
};

const FormAction = (props: Props) => {
  const {
    children,
    disabled: disabledProp = false,
    tooltip: tooltipProp,
    ignorePermissions = false,
    hideOnEditLock = true,
    disableOnFromError = false,
  } = props;

  const validationContext = useContext(FormAlertsContext);
  const edittingLocked =
    validationContext.editLock.detected && !ignorePermissions;
  const shouldHide = edittingLocked && hideOnEditLock == true;
  const haveDisablingError =
    validationContext.error.detected && disableOnFromError;
  const disabled: boolean =
    disabledProp || edittingLocked || haveDisablingError == true;
  const tooltip = useMemo(
    () =>
      joinNullableStrings([
        tooltipProp,
        haveDisablingError == true ? validationContext.error.message : null,
        edittingLocked == true ? validationContext.editLock.message : null,
      ]),
    [
      edittingLocked,
      haveDisablingError,
      tooltipProp,
      validationContext.editLock.message,
      validationContext.error.message,
    ],
  );
  return (
    <FormElementContext.Provider value={{disabled, tooltip}}>
      {(!shouldHide && children) || null}
    </FormElementContext.Provider>
  );
};

export default FormAction;
