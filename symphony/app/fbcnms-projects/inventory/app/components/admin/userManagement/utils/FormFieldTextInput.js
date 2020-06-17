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
import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import {useEffect, useState} from 'react';

type FormFieldTextInputProps = {
  validationId?: string,
  label: string,
  type?: string,
  value: string,
  onValueChanged?: ?(string) => void,
  className?: ?string,
  disabled?: ?boolean,
  hasError?: boolean,
  errorText?: ?string,
  immediateUpdate?: boolean,
};

/*
  This hooks helps calling update callback prop.
  It is needed in case we want to call the update callback
  only AFTER context values were recalculated.
  Without it, the callback is called BEFORE the context
  value changes takes effect.
  - SetTimeout is needed for ensuring the used callback is
    surly the one passed AFTER the context values calculations.
  - Using state hook for engaging rendring cycle when triggerred
    (without it, will be using the previous version of given callback).
  - Using effect hook for completing rendering cycle before using callback. 
*/
const useSideEffectCallback = callback => {
  const [shouldTriggerRunCallback, setShouldTriggerRunCallback] = useState(
    false,
  );
  useEffect(() => {
    if (!shouldTriggerRunCallback) {
      return;
    }
    setShouldTriggerRunCallback(false);
    if (callback == null) {
      return;
    }
    callback();
  }, [callback, shouldTriggerRunCallback]);

  return () => setTimeout(() => setShouldTriggerRunCallback(true));
};

const FormFieldTextInput = (props: FormFieldTextInputProps) => {
  const {
    value: propValue,
    onValueChanged,
    validationId,
    label,
    type,
    className,
    hasError,
    errorText,
    disabled,
    immediateUpdate = false,
  } = props;
  const [fieldValue, setFieldValue] = useState<string>('');
  useEffect(() => setFieldValue(propValue), [propValue]);

  const callOnValueChanged = useSideEffectCallback(
    onValueChanged ? () => onValueChanged(fieldValue) : null,
  );
  const updateOnValueChange = updatedValue => {
    const isOnGoingChange = updatedValue != null;
    const currentValue = isOnGoingChange ? updatedValue : fieldValue;
    const trimmedValue = (currentValue && currentValue.trim()) || '';
    if (!isOnGoingChange && trimmedValue != currentValue) {
      setFieldValue(trimmedValue);
    }
    if (trimmedValue == propValue) {
      return;
    }
    callOnValueChanged();
  };

  const isRequired = validationId != null;
  return (
    <FormField
      className={className || undefined}
      label={label}
      required={isRequired}
      validation={
        isRequired
          ? {
              id: validationId || '',
              value: fieldValue,
            }
          : undefined
      }
      hasError={hasError}
      errorText={errorText}>
      <TextInput
        type={type}
        value={fieldValue}
        disabled={disabled ?? false}
        onChange={e => {
          setFieldValue(e.target.value);
          if (immediateUpdate) {
            updateOnValueChange(e.target.value);
          }
        }}
        onBlur={!immediateUpdate ? () => updateOnValueChange() : undefined}
      />
    </FormField>
  );
};

export default FormFieldTextInput;
