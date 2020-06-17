/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {CheckListItemFillingProps} from './CheckListItemFilling';

import * as React from 'react';
import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  textInput: {
    width: '100%',
  },
}));

const FreeTextCheckListItemFilling = ({
  item,
  onChange,
}: CheckListItemFillingProps): React.Node => {
  const classes = useStyles();

  const _updateOnChange = newValue => {
    if (!onChange) {
      return;
    }
    const updatedItem = {
      ...item,
      stringValue: newValue,
      checked: !!newValue && newValue.trim().length > 0,
    };
    onChange(updatedItem);
  };

  return (
    <FormField>
      <TextInput
        className={classes.textInput}
        type="multiline"
        rows={5}
        placeholder={item.helpText || ''}
        value={item.stringValue || ''}
        onChange={event => _updateOnChange(event.target.value)}
      />
    </FormField>
  );
};

export default FreeTextCheckListItemFilling;
