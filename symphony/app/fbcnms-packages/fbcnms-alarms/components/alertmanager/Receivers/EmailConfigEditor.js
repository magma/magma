/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
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
        </>
      }
      data-testid="email-config-editor"
    />
  );
}
