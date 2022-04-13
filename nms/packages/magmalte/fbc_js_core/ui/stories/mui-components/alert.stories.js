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

import {action} from '@storybook/addon-actions';
import {storiesOf} from '@storybook/react';

import Alert from '../../components/Alert/Alert';
import Button from '@material-ui/core/Button';
import withAlert from '../../components/Alert/withAlert';
import {STORY_CATEGORIES} from '../storybookUtils';

const DemoButtonWithAlert = withAlert(({alert, label}) => {
  const handleClick = () => {
    alert('This is an alert', label).then(action('dismissed'));
  };
  return (
    <div>
      <Button onClick={handleClick}>Save</Button>
    </div>
  );
});

const DemoButtonWithConfirm = withAlert(({confirm, confirmProps}) => {
  const handleClick = () => {
    confirm(confirmProps).then(action('confirmed')).catch(action('cancelled'));
  };
  return (
    <div>
      <Button onClick={handleClick}>Delete</Button>
    </div>
  );
});

storiesOf(`${STORY_CATEGORIES.MUI_COMPONENTS}/Alert`, module)
  .add('default', () => (
    <Alert
      open={true}
      title="Title"
      message="message"
      onCancel={action('cancelled')}
      onConfirm={action('confirmed')}
    />
  ))
  .add('actions', () => (
    <Alert
      open={true}
      title="Title"
      message="message"
      confirmLabel="Confirm"
      cancelLabel="Cancel"
    />
  ));

storiesOf(`${STORY_CATEGORIES.MUI_COMPONENTS}/Alert/withAlert/alert`, module)
  .add('default', () => <DemoButtonWithAlert />)
  .add('custom label', () => <DemoButtonWithAlert label="Got it" />);

storiesOf(`${STORY_CATEGORIES.MUI_COMPONENTS}/Alert/withAlert/confirm`, module)
  .add('default', () => <DemoButtonWithConfirm confirmProps="Are you sure?" />)
  .add('custom confirm label', () => (
    <DemoButtonWithConfirm
      confirmProps={{message: 'Are you sure?', confirmLabel: 'Delete'}}
    />
  ))
  .add('custom cancel label', () => (
    <DemoButtonWithConfirm
      confirmProps={{message: 'Are you sure?', cancelLabel: 'Abort'}}
    />
  ));
