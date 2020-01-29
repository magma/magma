/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import * as React from 'react';
import ConfigEditor from './ConfigEditor';
import Grid from '@material-ui/core/Grid';
import TextField from '@material-ui/core/TextField';
import type {EditorProps} from './ConfigEditor';
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
        </>
      }
      data-testid="email-config-editor"
    />
  );
}
