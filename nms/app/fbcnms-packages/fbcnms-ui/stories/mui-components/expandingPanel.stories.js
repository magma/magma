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

import AddIcon from '@material-ui/icons/Add';
import ExpandingPanel from '../../components/ExpandingPanel';
import React from 'react';
import Text from '../../components/design-system/Text';
import {STORY_CATEGORIES} from '../storybookUtils';
import {storiesOf} from '@storybook/react';

storiesOf(`${STORY_CATEGORIES.MUI_COMPONENTS}/ExpandingPanel`, module)
  .add('default', () => (
    <ExpandingPanel title="Expanding Panel">
      <Text>This is the content</Text>
    </ExpandingPanel>
  ))
  .add('right button', () => (
    <ExpandingPanel title="Expanding Panel" rightContent={<AddIcon />}>
      <Text>This is the content</Text>
    </ExpandingPanel>
  ));
