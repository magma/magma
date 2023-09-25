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
import {AltFormField} from '../../../../../components/FormField';
import {OutlinedInput} from '@mui/material';
import {makeStyles} from '@mui/styles';

import type {EditorProps} from './ConfigEditor';
import type {ReceiverWebhookConfig} from '../../AlarmAPIType';

const useStyles = makeStyles(() => ({
  input: {
    padding: '0 16px',
  },
}));

export default function WebhookConfigEditor({
  config,
  onUpdate,
  ...props
}: EditorProps<ReceiverWebhookConfig>) {
  const classes = useStyles();
  return (
    <ConfigEditor
      {...props}
      RequiredFields={
        <>
          <AltFormField className={classes.input} dense label="Url">
            <OutlinedInput
              fullWidth={true}
              required
              value={config.url}
              onChange={e => onUpdate({url: e.target.value})}
              placeholder="Ex: webhook.example.com"
            />
          </AltFormField>
        </>
      }
      data-testid="webhook-config-editor"
    />
  );
}
