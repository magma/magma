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
import useSideEffectCallback from './useSideEffectCallback';
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
