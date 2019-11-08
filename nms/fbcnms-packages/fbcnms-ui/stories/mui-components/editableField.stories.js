/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import EditableField from '../../components/EditableField';
import React, {useState} from 'react';
import {STORY_CATEGORIES} from '../storybookUtils';
import {storiesOf} from '@storybook/react';

function TestField(props: {type: 'string' | 'date'}) {
  const [value, setValue] = useState(null);
  return (
    <EditableField
      onSave={newValue => {
        setValue(newValue);
        return true;
      }}
      value={value}
      type={props.type}
      editDisabled={false}
    />
  );
}

storiesOf(`${STORY_CATEGORIES.MUI_COMPONENTS}/EditableField`, module)
  .add('string', () => <TestField type="string" />)
  .add('date', () => <TestField type="date" />);
