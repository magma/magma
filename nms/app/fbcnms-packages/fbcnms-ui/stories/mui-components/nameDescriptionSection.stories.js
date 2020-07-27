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

import NameDescriptionSection from '../../components/NameDescriptionSection';
import React from 'react';
import {STORY_CATEGORIES} from '../storybookUtils';
import {storiesOf} from '@storybook/react';

storiesOf(`${STORY_CATEGORIES.MUI_COMPONENTS}/NameDescriptionSection`, module)
  .add('default', () => (
    <div>
      <NameDescriptionSection />
    </div>
  ))
  .add('custom title', () => (
    <div>
      <NameDescriptionSection title="Ttile" />
    </div>
  ))
  .add('placeholders', () => (
    <div>
      <NameDescriptionSection
        namePlaceholder="Add a name"
        descriptionPlaceholder="Add some details"
      />
    </div>
  ))
  .add('populated', () => (
    <div>
      <NameDescriptionSection
        name="Foo"
        description="Lorem ipsum dolor sit amet, consectetur adipiscing elit,
        sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
        Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi
        ut aliquip ex ea commodo consequat. Duis aute irure dolor in
        reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla
        pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa
         qui officia deserunt mollit anim id est laborum."
      />
    </div>
  ));
