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
import emptyFunction from '@fbcnms/util/emptyFunction';
import fbt from 'fbt';
import {createContext, useCallback, useContext, useMemo, useState} from 'react';
import {Map as immMap} from 'immutable';

type Range = {
  from: number,
  to: number,
};
const isEmpty = value => value == null;

const isEmptyLikeValue = value =>
  Number.isNaN(value) || isEmpty(value) || value === '';

const checkOuterAlertMessage = ruleInfo => ruleInfo?.alert || null;

const checkOuterCallback = ruleInfo =>
  (ruleInfo.checkCallback && ruleInfo.checkCallback(ruleInfo.value)) || null;

const checkRequired = ruleInfo =>
  !!ruleInfo.required && isEmptyLikeValue(ruleInfo.value)
    ? `${fbt(
        fbt.param('field name', ruleInfo.fieldDisplayName) + ' cannot be empty',
        'Required field notation',
      )}`
    : null;

const checkNumberInRange = ruleInfo => {
  if (isEmpty(ruleInfo.value)) {
    return null;
  }
  if (!ruleInfo.range) {
    return null;
  }
  const range: Range = ruleInfo.range;
  const numberValue = Number(ruleInfo.value);
  return Number.isNaN(numberValue) ||
    numberValue < range.from ||
    numberValue > range.to
    ? `${ruleInfo.fieldDisplayName} should be between
      ${range.from} and ${range.to}`
    : null;
};

export type FormRule = $ReadOnly<{|
  fieldId: string,
  fieldDisplayName: string,
  // eslint-disable-next-line flowtype/no-weak-types
  value: ?any,
  alert?: string,
  required?: boolean,
  range?: Range,
  // eslint-disable-next-line flowtype/no-weak-types
  checkCallback?: (value?: ?any) => string,
  notAggregated?: boolean,
|}>;

type FormAlertsContainer = $ReadOnly<{|
  detected: boolean,
  message: string,
  check: (validationInfo: FormRule) => ?string,
  set: (id: string, error: ?string) => ?string,
  clear: (id: string) => void,
|}>;

export type FormAlertsContextType = $ReadOnly<{|
  error: FormAlertsContainer,
  editLock: FormAlertsContainer,
|}>;

const emptyFormAlertsContainer = {
  detected: false,
  message: '',
  check: emptyFunction,
  set: emptyFunction,
  clear: emptyFunction,
};

export const DEFAULT_CONTEXT_VALUE = {
  error: emptyFormAlertsContainer,
  editLock: emptyFormAlertsContainer,
};

const FormAlertsContext = createContext<FormAlertsContextType>(
  DEFAULT_CONTEXT_VALUE,
);

const FormRulesMaintainer = function() {
  const [alertsMap, setAlertsMap] = useState<immMap<string, string>>(
    new immMap<string, string>(),
  );
  const [alertMessage, setAlertMessage] = useState('');
  const [hasAlerts, setHasAlerts] = useState(false);

  const updateContext = useCallback((newAlertsMap: immMap<string, string>) => {
    setAlertsMap(newAlertsMap);
    const aggregatedAlertMessage = newAlertsMap.join();
    setAlertMessage(aggregatedAlertMessage);
    setHasAlerts(aggregatedAlertMessage.length > 0);
  }, []);

  const clearAlert = useCallback(
    id => {
      if (!id) {
        return;
      }
      if (alertsMap.has(id)) {
        updateContext(alertsMap.delete(id));
      }
    },
    [alertsMap, updateContext],
  );

  const setAlert = useCallback(
    (id, alertMessage: ?string) => {
      let returnedAlert = null;
      if (!alertMessage) {
        clearAlert(id);
      } else {
        returnedAlert = alertMessage;
        if (alertsMap.get(id) !== alertMessage) {
          updateContext(alertsMap.set(id, alertMessage));
        }
      }

      return returnedAlert;
    },
    [alertsMap, clearAlert, updateContext],
  );

  const alertsChecks: Array<(FormRule) => ?string> = useMemo(
    () => [
      checkOuterAlertMessage,
      checkRequired,
      checkOuterCallback,
      checkNumberInRange,
    ],
    [],
  );

  const ruleCheck = useCallback(
    (ruleInfo: FormRule) => {
      let alert: ?string = null;
      let checksCount = 0;
      while (alert === null && checksCount < alertsChecks.length) {
        alert = alertsChecks[checksCount](ruleInfo);
        checksCount++;
      }
      return !!ruleInfo.notAggregated
        ? alert
        : setAlert(ruleInfo.fieldId, alert);
    },
    [alertsChecks, setAlert],
  );

  return {
    detected: hasAlerts,
    message: alertMessage,
    check: ruleCheck,
    set: setAlert,
    clear: clearAlert,
  };
};

type Props = {
  children: React.Node,
};

export function FormAlertsContextProvider(props: Props) {
  const errorsContext = FormRulesMaintainer();
  const editLocksContext = FormRulesMaintainer();

  const providerValue = {
    error: errorsContext,
    editLock: editLocksContext,
  };

  return (
    <FormAlertsContext.Provider value={providerValue}>
      {props.children}
    </FormAlertsContext.Provider>
  );
}

export const useFormAlertsContext = () => useContext(FormAlertsContext);

export default FormAlertsContext;
