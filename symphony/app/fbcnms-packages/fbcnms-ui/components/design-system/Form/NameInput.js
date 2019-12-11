/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import FormValidationContext from '@fbcnms/ui/components/design-system/Form/FormValidationContext';
import React, {useContext, useMemo} from 'react';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import shortid from 'shortid';

type Props = {
  value: ?string,
  onChange?: (e: SyntheticInputEvent<HTMLInputElement>) => void,
  inputClass?: string,
  title?: string,
  placeholder?: string,
};

const NameInput = (props: Props) => {
  const {title = 'Name', onChange, value, inputClass, placeholder} = props;
  const onNameChanded = event => {
    if (!onChange) {
      return;
    }
    onChange(event);
  };
  const fieldId = useMemo(() => shortid.generate(), []);
  const validationContext = useContext(FormValidationContext);
  const errorText = validationContext.errorCheck({
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
      hasSpacer={true}>
      <TextInput
        name={fieldId}
        autoFocus={true}
        type="string"
        className={inputClass}
        value={value || ''}
        placeholder={placeholder}
        onChange={onNameChanded}
      />
    </FormField>
  );
};

export default NameInput;
