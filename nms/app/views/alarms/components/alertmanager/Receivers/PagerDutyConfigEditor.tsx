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
import type {ReceiverPagerDutyConfig} from '../../AlarmAPIType';

const useStyles = makeStyles(() => ({
  input: {
    padding: '0 16px',
  },
}));

export default function PagerDutyConfigEditor({
  config,
  onUpdate,
  ...props
}: EditorProps<ReceiverPagerDutyConfig>) {
  const classes = useStyles();

  return (
    <ConfigEditor
      {...props}
      RequiredFields={
        <>
          <AltFormField className={classes.input} dense label="Description">
            <OutlinedInput
              fullWidth={true}
              required
              value={config.description}
              onChange={e => onUpdate({description: e.target.value})}
            />
          </AltFormField>
          <AltFormField className={classes.input} dense label="Severity">
            <OutlinedInput
              fullWidth={true}
              required
              value={config.severity}
              onChange={e => onUpdate({severity: e.target.value})}
              placeholder="Ex: notifications@example.com"
            />
          </AltFormField>
          <AltFormField className={classes.input} dense label="Url">
            <OutlinedInput
              fullWidth={true}
              required
              value={config.url}
              onChange={e => onUpdate({url: e.target.value})}
              placeholder="Ex: webhook.example.com"
            />
          </AltFormField>
          <AltFormField className={classes.input} dense label="Routing Key">
            <OutlinedInput
              fullWidth={true}
              required
              value={config.routing_key}
              onChange={e => onUpdate({routing_key: e.target.value})}
            />
          </AltFormField>
          <AltFormField className={classes.input} dense label="Service Key">
            <OutlinedInput
              fullWidth={true}
              required
              value={config.service_key}
              onChange={e => onUpdate({service_key: e.target.value})}
            />
          </AltFormField>
          <AltFormField className={classes.input} dense label="Client">
            <OutlinedInput
              fullWidth={true}
              value={config.client}
              onChange={e => onUpdate({client: e.target.value})}
            />
          </AltFormField>

          <AltFormField className={classes.input} dense label="Client Url">
            <OutlinedInput
              fullWidth={true}
              required
              value={config.client_url}
              onChange={e => onUpdate({client_url: e.target.value})}
            />
          </AltFormField>
        </>
      }
      data-testid="pager-duty-config-editor"
    />
  );
}
