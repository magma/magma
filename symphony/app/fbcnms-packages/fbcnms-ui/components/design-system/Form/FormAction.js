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
import {joinNullableStrings} from '@fbcnms/util/strings';
import {useFormAlertsContext} from '../Form/FormAlertsContext';
import {useMemo} from 'react';

export type PermissionHandlingProps = $ReadOnly<{|
  ignorePermissions?: ?boolean,
  hideOnMissingPermissions?: ?boolean,
  disableOnFromError?: ?boolean,
|}>;

export type FormActionProps = $ReadOnly<{|
  children: React.Node,
  disabled?: boolean,
  tooltip?: ?string,
|}>;

type Props = $ReadOnly<{|
  ...FormActionProps,
  ...PermissionHandlingProps,
|}>;

const FormAction = (props: Props) => {
  const {
    children,
    disabled: disabledProp = false,
    tooltip: tooltipProp,
    ignorePermissions = false,
    hideOnMissingPermissions = true,
    disableOnFromError = false,
  } = props;

  const validationContext = useFormAlertsContext();
  const missingPermissions =
    ignorePermissions !== true && validationContext.missingPermissions.detected;
  const edittingLocked =
    missingPermissions || validationContext.editLock.detected;
  const shouldHide = missingPermissions && hideOnMissingPermissions == true;
  const haveDisablingError =
    validationContext.error.detected && disableOnFromError;
  const disabled: boolean =
    disabledProp || edittingLocked || haveDisablingError == true;
  const tooltip = useMemo(
    () =>
      joinNullableStrings([
        tooltipProp,
        haveDisablingError == true ? validationContext.error.message : null,
        edittingLocked == true
          ? validationContext.missingPermissions.message
          : null,
      ]),
    [
      edittingLocked,
      haveDisablingError,
      tooltipProp,
      validationContext.missingPermissions.message,
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
