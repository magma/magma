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
import type {ReceiverSlackConfig} from '../../AlarmAPIType';

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
              required
              id="apiurl"
              label="Webhook URL"
              placeholder="Ex: https://hooks.slack.com/services/a/b"
              value={config.api_url}
              onChange={e => onUpdate({api_url: (e.target.value: string)})}
              fullWidth
            />
          </Grid>
        </>
      }
      OptionalFields={
        <>
          <Grid item>
            <TextField
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
      data-testid="slack-config-editor"
    />
  );
}
