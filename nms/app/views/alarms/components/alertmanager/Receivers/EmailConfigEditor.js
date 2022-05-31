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
import Checkbox from '@material-ui/core/Checkbox';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import ConfigEditor from './ConfigEditor';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Grid from '@material-ui/core/Grid';
import TextField from '@material-ui/core/TextField';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {EditorProps} from './ConfigEditor';
// $FlowFixMe migrated to typescript
import type {ReceiverEmailConfig} from '../../AlarmAPIType';

export default function EmailConfigEditor({
  config,
  onUpdate,
  ...props
}: EditorProps<ReceiverEmailConfig>) {
  return (
    <ConfigEditor
      {...props}
      RequiredFields={
        <>
          <Grid item>
            <TextField
              required
              label="Send To"
              placeholder="Ex: ops@example.com"
              value={config.to}
              onChange={e => onUpdate({to: e.target.value})}
              fullWidth
            />
          </Grid>
          <Grid item>
            <TextField
              required
              label="From"
              placeholder="Ex: notifications@example.com"
              value={config.from}
              onChange={e => onUpdate({from: e.target.value})}
              fullWidth
            />
          </Grid>
          <Grid item>
            <TextField
              required
              label="Host"
              placeholder="Ex: smtp.example.com"
              value={config.smarthost}
              onChange={e => onUpdate({smarthost: e.target.value})}
              fullWidth
            />
          </Grid>
          <Grid item>
            <TextField
              label="Auth Username"
              value={config.auth_username}
              onChange={e => onUpdate({auth_username: e.target.value})}
              fullWidth
              helperText="SMTP Auth using CRAM-MD5, LOGIN and PLAIN"
            />
          </Grid>
          <Grid item>
            <TextField
              label="Auth Password"
              value={config.auth_password}
              onChange={e => onUpdate({auth_password: e.target.value})}
              fullWidth
              helperText="SMTP Auth using LOGIN and PLAIN"
            />
          </Grid>
        </>
      }
      OptionalFields={
        <>
          <Grid item>
            <TextField
              label="Auth Secret"
              value={config.auth_secret}
              onChange={e => onUpdate({auth_secret: e.target.value})}
              fullWidth
              helperText="SMTP Auth using CRAM-MD5"
            />
          </Grid>
          <Grid item>
            <TextField
              label="Auth Identity"
              value={config.auth_identity}
              onChange={e => onUpdate({auth_identity: e.target.value})}
              fullWidth
              helperText="SMTP Auth using PLAIN"
            />
          </Grid>
          <Grid item>
            <FormControlLabel
              control={
                <Checkbox
                  checked={config.require_tls}
                  onChange={e => onUpdate({require_tls: e.target.checked})}
                  name="require_tls"
                  color="primary"
                  indeterminate={config.require_tls == null}
                />
              }
              label="Require TLS"
            />
          </Grid>
        </>
      }
      data-testid="email-config-editor"
    />
  );
}
