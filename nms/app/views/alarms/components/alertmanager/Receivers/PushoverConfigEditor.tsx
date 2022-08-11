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
 */

import * as React from 'react';
import ConfigEditor from './ConfigEditor';
import Grid from '@mui/material/Grid';
import MenuItem from '@mui/material/MenuItem';
import TextField from '@mui/material/TextField';
import type {EditorProps} from './ConfigEditor';
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
              variant="standard"
              required
              label="User Key"
              value={config.user_key}
              onChange={e => onUpdate({user_key: e.target.value})}
              fullWidth
            />
          </Grid>
          <Grid item>
            <TextField
              variant="standard"
              required
              label="Token"
              value={config.token}
              onChange={e => onUpdate({token: e.target.value})}
              fullWidth
            />
          </Grid>
        </>
      }
      OptionalFields={
        <>
          <Grid item>
            <TextField
              variant="standard"
              label="Title"
              value={config.title}
              onChange={e => onUpdate({title: e.target.value})}
              fullWidth
            />
          </Grid>
          <Grid item>
            <TextField
              variant="standard"
              label="Message"
              value={config.message}
              onChange={e => onUpdate({message: e.target.value})}
              fullWidth
            />
          </Grid>
          <Grid item>
            <TextField
              variant="standard"
              label="Priority"
              value={config.priority || '0'}
              onChange={e => onUpdate({priority: e.target.value})}
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
