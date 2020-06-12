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
import AddIcon from '@material-ui/icons/Add';
import Editor from '../../common/Editor';
import EmailConfigEditor from './EmailConfigEditor';
import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import MenuItem from '@material-ui/core/MenuItem';
import SlackConfigEditor from './SlackConfigEditor';
import TextField from '@material-ui/core/TextField';
import Typography from '@material-ui/core/Typography';
import WebhookConfigEditor from './WebhookConfigEditor';
import useForm from '../../../hooks/useForm';
import useRouter from '../../../hooks/useRouter';
import {useAlarmContext} from '../../AlarmContext';
import {useEnqueueSnackbar} from '../../../hooks/useSnackbar';

import type {
  AlertReceiver,
  ReceiverConfigListName,
  ReceiverEmailConfig,
  ReceiverSlackConfig,
  ReceiverWebhookConfig,
} from '../../AlarmAPIType';

type Props = {
  receiver: AlertReceiver,
  isNew: boolean,
  onExit: () => void,
};

const CONFIG_TYPES: {
  [string]: {
    listName: ReceiverConfigListName,
    friendlyName: string,
    createConfig: () => {},
    ConfigEditor: React.ComponentType<*>,
  },
} = {
  slack: {
    friendlyName: 'Slack Channel',
    listName: 'slack_configs',
    createConfig: emptySlackReceiver,
    ConfigEditor: SlackConfigEditor,
  },
  email: {
    friendlyName: 'Email',
    listName: 'email_configs',
    createConfig: emptyEmailReceiver,
    ConfigEditor: EmailConfigEditor,
  },
  webhook: {
    friendlyName: 'Webhook',
    listName: 'webhook_configs',
    createConfig: emptyWebhookReceiver,
    ConfigEditor: WebhookConfigEditor,
  },
};

export default function AddEditReceiver(props: Props) {
  const {apiUtil} = useAlarmContext();

  const {isNew, receiver, onExit} = props;
  const {match} = useRouter();
  const enqueueSnackbar = useEnqueueSnackbar();

  const {
    formState,
    handleInputChange,
    updateListItem,
    removeListItem,
    addListItem,
  } = useForm({
    initialState: receiver,
  });

  const [addConfigType, setAddConfigType] = React.useState('slack');

  const handleAddConfig = React.useCallback(() => {
    const {listName, createConfig} = CONFIG_TYPES[addConfigType];
    addListItem(listName, createConfig());
  }, [addConfigType, addListItem]);

  const handleSave = React.useCallback(() => {
    async function makeApiCall() {
      try {
        const request = {
          receiver: formState,
          networkId: match.params.networkId,
        };
        if (isNew) {
          await apiUtil.createReceiver(request);
          onExit();
        } else {
          await apiUtil.editReceiver(request);
        }
        enqueueSnackbar(`Successfully ${isNew ? 'added' : 'saved'} receiver`, {
          variant: 'success',
        });
      } catch (error) {
        enqueueSnackbar(
          `Unable to save receiver: ${
            error.response ? error.response.data.message : error.message
          }.`,
          {
            variant: 'error',
          },
        );
      }
    }
    makeApiCall();
  }, [
    apiUtil,
    enqueueSnackbar,
    formState,
    isNew,
    match.params.networkId,
    onExit,
  ]);

  const configEditorSharedProps = {
    receiver,
    formState,
    updateListItem,
    removeListItem,
  };
  return (
    <Editor
      isNew={isNew}
      onSave={handleSave}
      onExit={onExit}
      data-testid="add-edit-receiver"
      title={receiver?.name || 'New Receiver'}
      description="Configure channels to notify when an alert fires">
      <Grid item>
        <TextField
          required
          id="name"
          label="Receiver Name"
          placeholder="Ex: Support Team"
          disabled={!isNew}
          value={formState.name}
          onChange={handleInputChange(val => ({name: val}))}
          fullWidth
        />
      </Grid>
      <Grid container item>
        <Grid item xs={12}>
          <Typography>Add Notification</Typography>
        </Grid>
        <Grid
          item
          container
          xs={12}
          justify="space-between"
          alignItems="center">
          <Grid xs={11} item>
            <TextField
              select
              fullWidth
              value={addConfigType}
              onChange={e => setAddConfigType(e.target.value)}>
              {Object.keys(CONFIG_TYPES).map(key => {
                const {friendlyName} = CONFIG_TYPES[key];
                return (
                  <MenuItem key={key} value={key} data-test-config-type={key}>
                    {friendlyName}
                  </MenuItem>
                );
              })}
            </TextField>
          </Grid>
          <Grid xs={1} item>
            <IconButton
              edge="end"
              onClick={handleAddConfig}
              aria-label="add new receiver configuration">
              <AddIcon />
            </IconButton>
          </Grid>
        </Grid>
      </Grid>

      {Object.keys(CONFIG_TYPES).map(key => {
        const {
          friendlyName,
          createConfig,
          listName,
          ConfigEditor,
        } = CONFIG_TYPES[key];
        const list = formState[listName];
        return (
          <ConfigSection title={friendlyName}>
            {list && list.map
              ? list.map((config, idx) => (
                  <Grid item key={idx}>
                    <ConfigEditor
                      {...getConfigEditorProps({
                        listName: listName,
                        index: idx,
                        createConfig,
                        ...configEditorSharedProps,
                      })}
                    />
                  </Grid>
                ))
              : null}
          </ConfigSection>
        );
      })}
    </Editor>
  );
}

function ConfigSection({
  children,
  title,
}: {
  children?: ?React.Node,
  title: string,
}) {
  return (
    <Grid container item direction="column" wrap="nowrap" spacing={1}>
      <Grid container justify="space-between" alignItems="center" item xs={12}>
        <Grid item>
          <Typography color="textSecondary">{title}</Typography>
        </Grid>
      </Grid>
      {children || null}
    </Grid>
  );
}

function emptySlackReceiver(): ReceiverSlackConfig {
  return {api_url: ''};
}

function emptyEmailReceiver(): ReceiverEmailConfig {
  return {from: '', to: '', smarthost: ''};
}

function emptyWebhookReceiver(): ReceiverWebhookConfig {
  return {
    url: '',
  };
}

/**
 * Creates all the required props for a config editor.
 * Since config editors are rendered in a list and there is no unique
 * identifier, editing is done by list and by index (ie: slack_configs[0]).
 * This binds the callbacks to listname and index so the config editors don't
 * need to worry about their position in the list.
 */
function getConfigEditorProps<TConfig>({
  listName,
  index,
  receiver,
  formState,
  createConfig,
  updateListItem,
  removeListItem,
}: {
  listName: ReceiverConfigListName,
  index: number,
  receiver: AlertReceiver,
  formState: {[string]: Array<$Shape<TConfig>>, name: string},
  createConfig: () => $Shape<TConfig>,
  updateListItem: (
    listName: ReceiverConfigListName,
    index: number,
    update: $Shape<TConfig> | TConfig,
  ) => void,
  removeListItem: (listName: ReceiverConfigListName, index: number) => void,
}): {
  config: TConfig,
  onUpdate: ($Shape<TConfig>) => void,
  onReset: () => void,
  onDelete: () => void,
  isNew: boolean,
} {
  // The instance of a config such as ReceiverSlackConfig or ReceiverEmailConfig
  const config = formState[listName][index];
  const isNew = !receiver[listName] || !receiver[listName][index];

  const onUpdate = (update: $Shape<TConfig> | TConfig) =>
    updateListItem(listName, index, update);
  const onDelete = () => removeListItem(listName, index);
  const onReset = () =>
    updateListItem(
      listName,
      index,
      /**
       * When editing a config, the state of this config will be stored
       * untouched in the receiver object. If the receiver object does not
       * contain a definition for this config, it's new and we can reset it
       * by generating a new instance of the config
       */
      receiver[listName] && receiver[listName][index]
        ? receiver[listName][index]
        : null || createConfig(),
    );

  return {
    config,
    isNew,
    onUpdate,
    onReset,
    onDelete,
  };
}
