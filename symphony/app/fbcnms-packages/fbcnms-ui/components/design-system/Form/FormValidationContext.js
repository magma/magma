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
import * as imm from 'immutable';
import emptyFunction from '@fbcnms/util/emptyFunction';
import fbt from 'fbt';
import {useCallback, useMemo, useState} from 'react';

type Range = {
  from: number,
  to: number,
};

export type FormInputValueValidation = {
  fieldId: string,
  fieldDisplayName: string,
  // eslint-disable-next-line flowtype/no-weak-types
  value: ?any,
  errorMessage?: string,
  required?: boolean,
  range?: Range,
  // eslint-disable-next-line flowtype/no-weak-types
  checkCallbalck?: (value?: ?any) => string,
};

type FormIssuesContainer = {
  detected: boolean,
  message: string,
  check: (validationInfo: FormInputValueValidation) => ?string,
  set: (id: string, error: ?string) => ?string,
  clear: (id: string) => void,
};

type FormValidationContextType = {
  error: FormIssuesContainer,
  editLock: FormIssuesContainer,
};

const emptyFormIssuesContainer = {
  detected: false,
  message: '',
  check: emptyFunction,
  set: emptyFunction,
  clear: emptyFunction,
};

const FormValidationContext = React.createContext<FormValidationContextType>({
  error: emptyFormIssuesContainer,
  editLock: emptyFormIssuesContainer,
});

type Props = {
  children: React.Node,
};

type ErrorsMap = imm.Map<string, string>;

const FormValidationMaintainer = function() {
  const [errorsMap, setErrorsMap] = useState<ErrorsMap>(
    new imm.Map<string, string>(),
  );
  const [errorMessage, setErrorMessage] = useState('');
  const [hasErrors, setHasErrors] = useState(false);

  const updateContext = useCallback((newErrorsMap: ErrorsMap) => {
    setErrorsMap(newErrorsMap);
    const aggregatedErrorMessage = newErrorsMap.join();
    setErrorMessage(aggregatedErrorMessage);
    setHasErrors(aggregatedErrorMessage.length > 0);
  }, []);

  const clearError = useCallback(
    id => {
      if (!id) {
        return;
      }
      if (errorsMap.has(id)) {
        updateContext(errorsMap.delete(id));
      }
    },
    [errorsMap, updateContext],
  );

  const setError = useCallback(
    (id, errorMessage: ?string) => {
      let returnedError = null;
      if (!errorMessage) {
        clearError(id);
      } else {
        returnedError = errorMessage;
        if (errorsMap.get(id) !== errorMessage) {
          updateContext(errorsMap.set(id, errorMessage));
        }
      }

      return returnedError;
    },
    [errorsMap, clearError, updateContext],
  );
  const isEmpty = useCallback(value => value == null, []);
  const isEmptyLikeValue = useCallback(
    value => Number.isNaN(value) || isEmpty(value) || value === '',
    [isEmpty],
  );
  const checkOuterErrorMessage = useCallback(
    validationInfo => validationInfo?.errorMessage || null,
    [],
  );
  const checkOuterCallback = useCallback(
    validationInfo =>
      (validationInfo.checkCallbalck &&
        validationInfo.checkCallbalck(validationInfo.value)) ||
      null,
    [],
  );
  const checkRequired = useCallback(
    validationInfo =>
      !!validationInfo.required && isEmptyLikeValue(validationInfo.value)
        ? `${fbt(
            fbt.param('field name', validationInfo.fieldDisplayName) +
              ' cannot be empty',
            'Required field notation',
          )}`
        : null,
    [isEmptyLikeValue],
  );
  const checkNumberInRange = useCallback(
    validationInfo => {
      if (isEmpty(validationInfo.value)) {
        return null;
      }
      if (!validationInfo.range) {
        return null;
      }
      const range: Range = validationInfo.range;
      const numberValue = Number(validationInfo.value);
      return Number.isNaN(numberValue) ||
        numberValue < range.from ||
        numberValue > range.to
        ? `${validationInfo.fieldDisplayName} should be between
         ${range.from} and ${range.to}`
        : null;
    },
    [isEmpty],
  );

  const errorChecks: Array<(v: FormInputValueValidation) => ?string> = useMemo(
    () => [
      checkOuterErrorMessage,
      checkRequired,
      checkOuterCallback,
      checkNumberInRange,
    ],
    [
      checkOuterErrorMessage,
      checkRequired,
      checkOuterCallback,
      checkNumberInRange,
    ],
  );

  const errorCheck = useCallback(
    (validationInfo: FormInputValueValidation) => {
      let errorMessage: ?string = null;
      let checksCount = 0;
      while (errorMessage === null && checksCount < errorChecks.length) {
        errorMessage = errorChecks[checksCount](validationInfo);
        checksCount++;
      }
      return setError(validationInfo.fieldId, errorMessage);
    },
    [errorChecks, setError],
  );

  return {
    detected: hasErrors,
    message: errorMessage,
    check: errorCheck,
    set: setError,
    clear: clearError,
  };
};

export function FormValidationContextProvider(props: Props) {
  const errorsContext = FormValidationMaintainer();
  const editLocksContext = FormValidationMaintainer();

  const providerValue = {
    error: errorsContext,
    editLock: editLocksContext,
  };

  return (
    <FormValidationContext.Provider value={providerValue}>
      {props.children}
    </FormValidationContext.Provider>
  );
}

export default FormValidationContext;
