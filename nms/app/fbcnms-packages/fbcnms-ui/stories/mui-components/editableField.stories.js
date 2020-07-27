/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @flow strict-local
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
