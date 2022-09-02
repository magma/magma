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
import Checkbox from '@mui/material/Checkbox';
import ConfigEditor from './ConfigEditor';
import FormControlLabel from '@mui/material/FormControlLabel';
import Text from '../../../../../theme/design-system/Text';
import {AltFormField} from '../../../../../components/FormField';
import {ListItem, OutlinedInput} from '@mui/material';
import {makeStyles} from '@mui/styles';

import type {EditorProps} from './ConfigEditor';
import type {ReceiverEmailConfig} from '../../AlarmAPIType';

const useStyles = makeStyles(() => ({
  input: {
    padding: '0 16px',
  },
}));

export default function EmailConfigEditor({
  config,
  onUpdate,
  ...props
}: EditorProps<ReceiverEmailConfig>) {
  const classes = useStyles();

  return (
    <ConfigEditor
      {...props}
      RequiredFields={
        <>
          <AltFormField className={classes.input} dense label="Send To">
            <OutlinedInput
              fullWidth={true}
              value={config.to}
              onChange={e => onUpdate({to: e.target.value})}
              placeholder="Ex: ops@example.com"
            />
          </AltFormField>
          <AltFormField className={classes.input} dense label="From">
            <OutlinedInput
              fullWidth={true}
              value={config.from}
              onChange={e => onUpdate({from: e.target.value})}
              placeholder="Ex: notifications@example.com"
            />
          </AltFormField>
          <AltFormField className={classes.input} label="Host">
            <OutlinedInput
              fullWidth={true}
              value={config.smarthost}
              onChange={e => onUpdate({smarthost: e.target.value})}
              placeholder="Ex: smtp.example.com"
            />
          </AltFormField>
          <AltFormField
            className={classes.input}
            label="Auth Username"
            subLabel="SMTP Auth using CRAM-MD5, LOGIN and PLAIN">
            <OutlinedInput
              fullWidth={true}
              value={config.auth_username}
              onChange={e => onUpdate({auth_username: e.target.value})}
            />
          </AltFormField>
          <AltFormField
            className={classes.input}
            label="Auth Password"
            subLabel="SMTP Auth using LOGIN and PLAIN">
            <OutlinedInput
              fullWidth={true}
              value={config.auth_password}
              onChange={e => onUpdate({auth_password: e.target.value})}
            />
          </AltFormField>
        </>
      }
      OptionalFields={
        <>
          <AltFormField
            className={classes.input}
            label="Auth Secret"
            subLabel="SMTP Auth using CRAM-MD5">
            <OutlinedInput
              fullWidth={true}
              value={config.auth_secret}
              onChange={e => onUpdate({auth_secret: e.target.value})}
            />
          </AltFormField>
          <AltFormField
            className={classes.input}
            label="Auth Identity"
            subLabel="SMTP Auth using PLAIN">
            <OutlinedInput
              fullWidth={true}
              value={config.auth_identity}
              onChange={e => onUpdate({auth_identity: e.target.value})}
            />
          </AltFormField>
          <ListItem>
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
              label={
                <Text weight="medium" variant="subtitle2">
                  Require TLS
                </Text>
              }
            />
          </ListItem>
        </>
      }
      data-testid="email-config-editor"
    />
  );
}
