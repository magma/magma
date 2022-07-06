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
 * @flow
 * @format
 */

import * as React from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import ConfigEditor from './ConfigEditor';
import Grid from '@material-ui/core/Grid';
import TextField from '@material-ui/core/TextField';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {EditorProps} from './ConfigEditor';
// $FlowFixMe migrated to typescript
import type {ReceiverSlackConfig} from '../../AlarmAPIType';

export default function SlackConfigEditor({
  config,
  onUpdate,
  ...props
}: EditorProps<ReceiverSlackConfig>) {
  return (
    <ConfigEditor
      {...props}
      RequiredFields={
        <>
          <Grid item>
            <TextField
              required
              data-testid="slack-config-editor"
              id="apiurl"
              label="Webhook URL"
              placeholder="Ex: https://hooks.slack.com/services/a/b"
              value={config.api_url}
              onChange={e => onUpdate({api_url: (e.target.value: string)})}
              fullWidth
            />
          </Grid>
        </>
      }
      OptionalFields={
        <>
          <Grid item>
            <TextField
              id="title"
              label="Message Title"
              placeholder="Ex: Urgent"
              value={config.title}
              onChange={e => onUpdate({title: e.target.value})}
              fullWidth
            />
          </Grid>
        </>
      }
    />
  );
}
