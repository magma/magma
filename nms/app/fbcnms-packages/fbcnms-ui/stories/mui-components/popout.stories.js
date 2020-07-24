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

import Button from '@material-ui/core/Button';
import Popout from '../../components/Popout';
import React from 'react';
import Text from '../../components/design-system/Text';
import {STORY_CATEGORIES} from '../storybookUtils';
import {storiesOf} from '@storybook/react';

storiesOf(`${STORY_CATEGORIES.MUI_COMPONENTS}/Popout`, module).add(
  'default',
  () => (
    <div style={{padding: 100}}>
      <Popout
        content={
          <div style={{padding: 20}}>
            <Text variant="body2">Content</Text>
          </div>
        }>
        <Button variant="contained" color="primary">
          Click me!
        </Button>
      </Popout>
    </div>
  ),
);
