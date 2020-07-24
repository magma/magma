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

import React, {useState} from 'react';
import SideNavigationPanel from '../../components/SideNavigationPanel';
import {STORY_CATEGORIES} from '../storybookUtils';
import {storiesOf} from '@storybook/react';

const TemplatesPanel = () => {
  const [selectedItem, setSelectedItem] = useState('0');
  return (
    <SideNavigationPanel
      title="Templates"
      items={[
        {key: '0', label: 'Work Orders'},
        {key: '1', label: 'Projects'},
      ]}
      selectedItemId={selectedItem}
      onItemClicked={({key}) => setSelectedItem(key)}
    />
  );
};

storiesOf(
  `${STORY_CATEGORIES.MUI_COMPONENTS}/SideNavigationPanel`,
  module,
).add('default', () => <TemplatesPanel />);
