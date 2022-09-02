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
import MenuItem from '@mui/material/MenuItem';
import {AltFormField} from '../../../../../components/FormField';
import {FormControl, OutlinedInput, Select} from '@mui/material';
import {makeStyles} from '@mui/styles';

import type {EditorProps} from './ConfigEditor';
import type {ReceiverPushoverConfig} from '../../AlarmAPIType';

const useStyles = makeStyles(() => ({
  input: {
    padding: '0 16px',
  },
}));

export default function PushoverConfigEditor({
  config,
  onUpdate,
  ...props
}: EditorProps<ReceiverPushoverConfig>) {
  const classes = useStyles();
  const priorityList = [
    ['-2', 'Lowest'],
    ['-1', 'Low'],
    ['0', 'Normal'],
    ['1', 'High'],
    ['2', 'Emergency'],
  ];
  return (
    <ConfigEditor
      {...props}
      RequiredFields={
        <>
          <AltFormField className={classes.input} dense label="User Key">
            <OutlinedInput
              fullWidth={true}
              required
              value={config.user_key}
              onChange={e => onUpdate({user_key: e.target.value})}
            />
          </AltFormField>
          <AltFormField className={classes.input} dense label="Token">
            <OutlinedInput
              fullWidth={true}
              required
              value={config.token}
              onChange={e => onUpdate({token: e.target.value})}
            />
          </AltFormField>
        </>
      }
      OptionalFields={
        <>
          <AltFormField className={classes.input} dense label="Title">
            <OutlinedInput
              fullWidth={true}
              value={config.title}
              onChange={e => onUpdate({title: e.target.value})}
            />
          </AltFormField>
          <AltFormField className={classes.input} dense label="Message">
            <OutlinedInput
              fullWidth={true}
              value={config.message}
              onChange={e => onUpdate({message: e.target.value})}
            />
          </AltFormField>
          <AltFormField label={'Priority'}>
            <FormControl fullWidth>
              <Select
                value={config.priority || '0'}
                onChange={e => onUpdate({priority: e.target.value})}
                input={<OutlinedInput id="deviceClass" />}>
                {priorityList.map(([priority, label]) => (
                  <MenuItem key={label} value={priority}>
                    {label}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </AltFormField>
        </>
      }
      data-testid="email-config-editor"
    />
  );
}
