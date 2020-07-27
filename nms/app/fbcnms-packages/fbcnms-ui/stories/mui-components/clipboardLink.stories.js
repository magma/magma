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

import React from 'react';

import {storiesOf} from '@storybook/react';

import Button from '@material-ui/core/Button';
import ClipboardLink from '../../components/ClipboardLink';
import IconButton from '@material-ui/core/IconButton';
import LinkIcon from '@material-ui/icons/Link';
import {STORY_CATEGORIES} from '../storybookUtils';

storiesOf(`${STORY_CATEGORIES.MUI_COMPONENTS}/ClipboardLink`, module)
  .add('default', () => (
    <ClipboardLink>
      {({copyString}) => (
        <Button onClick={() => copyString('hi')}>Copy 'hi' to clipboard</Button>
      )}
    </ClipboardLink>
  ))
  .add('custom tooltip options', () => (
    <ClipboardLink
      title="Copy 'hi' to clipboard"
      placement="right"
      leaveDelay={400}>
      {({copyString}) => <Button onClick={() => copyString('hi')}>Copy</Button>}
    </ClipboardLink>
  ))
  .add('with IconButton', () => (
    <ClipboardLink title="Copy url to clipboard">
      {({copyString}) => (
        <IconButton onClick={() => copyString(window.location.href)}>
          <LinkIcon />
        </IconButton>
      )}
    </ClipboardLink>
  ));
