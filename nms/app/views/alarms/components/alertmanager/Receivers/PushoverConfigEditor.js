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
import * as React from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import ConfigEditor from './ConfigEditor';
import Grid from '@material-ui/core/Grid';
import MenuItem from '@material-ui/core/MenuItem';
import TextField from '@material-ui/core/TextField';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {EditorProps} from './ConfigEditor';
// $FlowFixMe migrated to typescript
import type {ReceiverPushoverConfig} from '../../AlarmAPIType';

export default function PushoverConfigEditor({
  config,
  onUpdate,
  ...props
}: EditorProps<ReceiverPushoverConfig>) {
  return (
    <ConfigEditor
      {...props}
      RequiredFields={
        <>
          <Grid item>
            <TextField
              required
              label="User Key"
              value={config.user_key}
              onChange={e => onUpdate({user_key: (e.target.value: string)})}
              fullWidth
            />
          </Grid>
          <Grid item>
            <TextField
              required
              label="Token"
              value={config.token}
              onChange={e => onUpdate({token: (e.target.value: string)})}
              fullWidth
            />
          </Grid>
        </>
      }
      OptionalFields={
        <>
          <Grid item>
            <TextField
              label="Title"
              value={config.title}
              onChange={e => onUpdate({title: (e.target.value: string)})}
              fullWidth
            />
          </Grid>
          <Grid item>
            <TextField
              label="Message"
              value={config.message}
              onChange={e => onUpdate({message: (e.target.value: string)})}
              fullWidth
            />
          </Grid>
          <Grid item>
            <TextField
              label="Priority"
              value={config.priority || '0'}
              onChange={e => onUpdate({priority: (e.target.value: string)})}
              fullWidth
              select>
              {[
                ['-2', 'Lowest'],
                ['-1', 'Low'],
                ['0', 'Normal'],
                ['1', 'High'],
                ['2', 'Emergency'],
              ].map(([priority, label]) => (
                <MenuItem key={label} value={priority}>
                  {label}
                </MenuItem>
              ))}
            </TextField>
          </Grid>
        </>
      }
      data-testid="email-config-editor"
    />
  );
}
