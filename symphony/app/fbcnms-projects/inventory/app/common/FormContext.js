/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {FormAlertsContextType} from '@fbcnms/ui/components/design-system/Form/FormAlertsContext';

import * as React from 'react';
import AppContext from '@fbcnms/ui/context/AppContext';
import FormAlertsContext, {
  DEFAULT_CONTEXT_VALUE as DEFAULT_ALERTS,
  FormAlertsContextProvider,
} from '@fbcnms/ui/components/design-system/Form/FormAlertsContext';
import fbt from 'fbt';
import {createContext, useContext} from 'react';

type FromContextType = $ReadOnly<{|
  alerts: FormAlertsContextType,
|}>;

const DEFAULT_CONTEXT_VALUE = {
  alerts: DEFAULT_ALERTS,
};

const FormContext = createContext<FromContextType>(DEFAULT_CONTEXT_VALUE);

type Props = {
  children: React.Node,
};

export function FormContextProvider(props: Props) {
  const {user} = useContext(AppContext);
  return (
    <FormAlertsContextProvider>
      <FormAlertsContext.Consumer>
        {alerts => {
          alerts.editLock.check({
            fieldId: 'System Rules',
            fieldDisplayName: 'Read Only User',
            value: user?.isReadOnlyUser,
            checkCallback: isReadOnlyUser =>
              isReadOnlyUser
                ? `${fbt(
                    'Writing permissions are required. Contact your system administrator.',
                    '',
                  )}`
                : '',
          });
          return (
            <FormContext.Provider value={{alerts}}>
              {props.children}
            </FormContext.Provider>
          );
        }}
      </FormAlertsContext.Consumer>
    </FormAlertsContextProvider>
  );
}

export const useFormContext = () => useContext(FormContext);

export default FormContext;
