/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import type {Property} from '../../common/Property';
import type {PropertyType} from '../../common/PropertyType';

import React, {useState} from 'react';
import Tokenizer from '@fbcnms/ui/components/Tokenizer';
import update from 'immutability-helper';
import {isJSON} from '@fbcnms/ui/utils/displayUtils';

type Props<T: Property | PropertyType> = {|
  property: T,
  onBlur?: ?() => void,
  onChange: T => void,
  disabled?: ?boolean,
|};

function EnumPropertyValueInput<T: Property | PropertyType>(props: Props<T>) {
  const {property, onBlur, onChange, disabled} = props;
  const jsonStr = property.stringValue || '';
  const options = isJSON(jsonStr) ? JSON.parse(jsonStr) : [];
  const optionsArr = Array.isArray(options) ? options : [];
  const [tokens, setTokens] = useState(
    optionsArr.map(option => ({
      id: option,
      label: option,
    })),
  );
  return (
    <Tokenizer
      searchSource="UserInput"
      tokens={tokens}
      disabled={disabled}
      onEntriesRequested={() => {}}
      placeholder="Press Enter after each value"
      onChange={newEntries => {
        setTokens(newEntries);
        onChange(
          update(property, {
            stringValue: {
              $set: JSON.stringify(newEntries.map(t => t.label)),
            },
          }),
        );
      }}
      onBlur={() => {
        onBlur && onBlur();
      }}
    />
  );
}

export default EnumPropertyValueInput;
