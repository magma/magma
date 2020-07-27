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

import MoreActionsButton from '../../components/MoreActionsButton';
import React from 'react';
import {STORY_CATEGORIES} from '../storybookUtils';
import {storiesOf} from '@storybook/react';

storiesOf(`${STORY_CATEGORIES.MUI_COMPONENTS}/MoreActionsButton`, module).add(
  'string',
  () => (
    <MoreActionsButton
      variant="primary"
      items={[
        {name: 'Item 1', onClick: () => window.alert('clicked item #1')},
        {name: 'Item 2', onClick: () => window.alert('clicked item #2')},
        {name: 'Item 3', onClick: () => window.alert('clicked item #3')},
      ]}
    />
  ),
);
