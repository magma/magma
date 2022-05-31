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
import ConfigEditor from './ConfigEditor';
import Grid from '@material-ui/core/Grid';
import TextField from '@material-ui/core/TextField';
import type {EditorProps} from './ConfigEditor';
// $FlowFixMe migrated to typescript
import type {ReceiverPagerDutyConfig} from '../../AlarmAPIType';

export default function PagerDutyConfigEditor({
  config,
  onUpdate,
  ...props
}: EditorProps<ReceiverPagerDutyConfig>) {
  return (
    <ConfigEditor
      {...props}
      RequiredFields={
        <>
          <Grid item>
            <TextField
              required
              label="Description"
              value={config.description}
              onChange={e => onUpdate({description: (e.target.value: string)})}
              fullWidth
            />
          </Grid>
          <Grid item>
            <TextField
              required
              label="Severity"
              value={config.severity}
              onChange={e => onUpdate({severity: (e.target.value: string)})}
              fullWidth
            />
          </Grid>
          <Grid item>
            <TextField
              required
              label="Url"
              placeholder="Ex: webhook.example.com"
              value={config.url}
              onChange={e => onUpdate({url: (e.target.value: string)})}
              fullWidth
            />
          </Grid>
          <Grid item>
            <TextField
              required
              label="Routing_key"
              value={config.routing_key}
              onChange={e => onUpdate({routing_key: (e.target.value: string)})}
              fullWidth
            />
          </Grid>
          <Grid item>
            <TextField
              required
              label="Service_key"
              value={config.service_key}
              onChange={e => onUpdate({service_key: (e.target.value: string)})}
              fullWidth
            />
          </Grid>
          <Grid item>
            <TextField
              required
              label="Client"
              value={config.client}
              onChange={e => onUpdate({client: (e.target.value: string)})}
              fullWidth
            />
          </Grid>
          <Grid item>
            <TextField
              required
              label="Client Url"
              value={config.client_url}
              onChange={e => onUpdate({client_url: (e.target.value: string)})}
              fullWidth
            />
          </Grid>
        </>
      }
      data-testid="pager-duty-config-editor"
    />
  );
}
