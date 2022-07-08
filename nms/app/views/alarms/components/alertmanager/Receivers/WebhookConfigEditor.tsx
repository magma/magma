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
import Grid from '@material-ui/core/Grid';
import TextField from '@material-ui/core/TextField';
import type {EditorProps} from './ConfigEditor';
import type {ReceiverWebhookConfig} from '../../AlarmAPIType';

export default function WebhookConfigEditor({
  config,
  onUpdate,
  ...props
}: EditorProps<ReceiverWebhookConfig>) {
  return (
    <ConfigEditor
      {...props}
      RequiredFields={
        <>
          <Grid item>
            <TextField
              required
              label="Url"
              placeholder="Ex: webhook.example.com"
              value={config.url}
              onChange={e => onUpdate({url: e.target.value})}
              fullWidth
            />
          </Grid>
        </>
      }
      data-testid="webhook-config-editor"
    />
  );
}
