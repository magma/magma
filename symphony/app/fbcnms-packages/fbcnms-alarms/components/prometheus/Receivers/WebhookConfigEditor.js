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
              onChange={e => onUpdate({url: (e.target.value: string)})}
              fullWidth
            />
          </Grid>
        </>
      }
      data-testid="webhook-config-editor"
    />
  );
}
