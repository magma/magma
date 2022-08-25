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
import TextField from '@mui/material/TextField';
import {AltFormField} from '../../../../../components/FormField';
import {Button, OutlinedInput} from '@mui/material';
import {ExpandMore} from '@mui/icons-material';
import {makeStyles} from '@mui/styles';

import type {EditorProps} from './ConfigEditor';
import type {ReceiverSlackConfig} from '../../AlarmAPIType';

const useStyles = makeStyles(() => ({
  expandMoreRotation: {
    transform: 'rotate(-180deg)',
    transition: '.3s',
  },
  expandLessRotation: {
    transition: '.3s',
  },
}));

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
              variant="standard"
              required
              data-testid="slack-config-editor"
              id="apiurl"
              label="Webhook URL"
              placeholder="Ex: https://hooks.slack.com/services/a/b"
              value={config.api_url}
              onChange={e => onUpdate({api_url: e.target.value})}
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

export function SlackConfig(isNew: boolean) {
  const classes = useStyles();
  const [advancedSettings, setAdvancedSettings] = React.useState<boolean>(
    false,
  );

  return (
    <Grid container spacing={2}>
      <AltFormField label={'Webhook URL'}>
        <OutlinedInput
          disabled={!isNew}
          required
          data-testid="slack-config-editor"
          id="apiurl"
          placeholder="Ex: https://hooks.slack.com/services/a/b"
          value={'config.api_url'}
          onChange={e => console.log('api_url:', e.target.value)}
          fullWidth
        />
      </AltFormField>
      <Button
        endIcon={
          <ExpandMore
            className={
              advancedSettings
                ? classes.expandMoreRotation
                : classes.expandLessRotation
            }
          />
        }
        variant="text"
        onClick={() => setAdvancedSettings(settings => !settings)}>
        Advanced Settings
      </Button>
      {advancedSettings && (
        <AltFormField label={'Message Title'}>
          <OutlinedInput
            disabled={!isNew}
            id="title"
            placeholder="Ex: Urgent"
            value={'config.title'}
            onChange={e => console.log('title:', e.target.value)}
            fullWidth
          />
        </AltFormField>
      )}
    </Grid>
  );
}
