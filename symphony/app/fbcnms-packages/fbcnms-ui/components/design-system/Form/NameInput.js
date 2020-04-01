/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import FormAlertsContext from '@fbcnms/ui/components/design-system/Form/FormAlertsContext';
import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import React, {useContext, useMemo} from 'react';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import shortid from 'shortid';

type Props = {
  value: ?string,
  onChange?: (e: SyntheticInputEvent<HTMLInputElement>) => void,
  inputClass?: string,
  title?: string,
  placeholder?: string,
  disabled?: boolean,
  onBlur?: () => void,
  hasSpacer?: boolean,
};

const NameInput = (props: Props) => {
  const {
    title = 'Name',
    onChange,
    value,
    inputClass,
    placeholder,
    disabled,
    onBlur,
    hasSpacer,
  } = props;
  const onNameChanded = event => {
    if (!onChange) {
      return;
    }
    onChange(event);
  };
  const fieldId = useMemo(() => shortid.generate(), []);
  const validationContext = useContext(FormAlertsContext);
  const errorText = validationContext.error.check({
    fieldId,
    fieldDisplayName: title,
    value: value,
    required: true,
  });
  return (
    <FormField
      label={title}
      required={true}
      hasError={!!errorText}
      errorText={errorText}
      hasSpacer={hasSpacer ?? true}>
      <TextInput
        name={fieldId}
        autoFocus={true}
        type="string"
        className={inputClass}
        value={value || ''}
        placeholder={placeholder}
        onChange={onNameChanded}
        disabled={disabled}
        onBlur={onBlur}
      />
    </FormField>
  );
};

export default NameInput;
