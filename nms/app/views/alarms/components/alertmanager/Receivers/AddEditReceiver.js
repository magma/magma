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
 *
 * @flow
 * @format
 */
import * as React from 'react';
import AddCircleOutlineIcon from '@material-ui/icons/AddCircleOutline';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import Editor from '../../common/Editor';
import EmailConfigEditor from './EmailConfigEditor';
import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import PagerDutyConfigEditor from './PagerDutyConfigEditor';
import PushoverConfigEditor from './PushoverConfigEditor';
import SlackConfigEditor from './SlackConfigEditor';
import TextField from '@material-ui/core/TextField';
import Typography from '@material-ui/core/Typography';
import WebhookConfigEditor from './WebhookConfigEditor';
import useForm from '../../../hooks/useForm';
import {useAlarmContext} from '../../AlarmContext';
import {useParams} from 'react-router-dom';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useSnackbars} from '../../../../../hooks/useSnackbar';

import type {
  AlertReceiver,
  ReceiverConfigListName,
  ReceiverEmailConfig,
  ReceiverPagerDutyConfig,
  ReceiverPushoverConfig,
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
  pagerduty: {
    friendlyName: 'PagerDuty',
    listName: 'pagerduty_configs',
    createConfig: emptyPagerDutyReceiver,
    ConfigEditor: PagerDutyConfigEditor,
  },
  pushover: {
    friendlyName: 'Pushover',
    listName: 'pushover_configs',
    createConfig: emptyPushoverReceiver,
    ConfigEditor: PushoverConfigEditor,
  },
};

export default function AddEditReceiver(props: Props) {
  const {apiUtil} = useAlarmContext();
  const snackbars = useSnackbars();
  const {isNew, receiver, onExit} = props;
  const params = useParams();

  const {
    formState,
    handleInputChange,
    updateListItem,
    removeListItem,
    addListItem,
  } = useForm({
    initialState: receiver,
  });

  const handleAddConfig = React.useCallback(
    (configType: string) => {
      const {listName, createConfig} = CONFIG_TYPES[configType];
      addListItem(listName, createConfig());
    },
    [addListItem],
  );

  const handleSave = React.useCallback(() => {
    async function makeApiCall() {
      try {
        const request = {
          receiver: formState,
          networkId: params.networkId,
        };
        if (isNew) {
          await apiUtil.createReceiver(request);
          onExit();
        } else {
          await apiUtil.editReceiver(request);
        }
        snackbars.success(`Successfully ${isNew ? 'added' : 'saved'} receiver`);
      } catch (error) {
        snackbars.error(
          `Unable to save receiver: ${
            error.response ? error.response.data.message : error.message
          }.`,
        );
      }
    }
    makeApiCall();
  }, [apiUtil, formState, isNew, params.networkId, onExit, snackbars]);

  const configEditorSharedProps = {
    receiver,
    formState,
    updateListItem,
    removeListItem,
  };
  return (
    <Editor
      xs={8}
      isNew={isNew}
      onSave={handleSave}
      onExit={onExit}
      data-testid="add-edit-receiver"
      title={receiver?.name || 'New Receiver'}
      description="Configure channels to notify when an alert fires">
      <Grid item>
        <Card>
          <CardContent>
            <Typography paragraph>Details</Typography>
            <TextField
              required
              id="name"
              label="Name"
              placeholder="Ex: Support Team"
              disabled={!isNew}
              value={formState.name}
              onChange={handleInputChange(val => ({name: val}))}
              fullWidth
            />
          </CardContent>
        </Card>
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
          <ConfigSection
            title={friendlyName}
            onAddConfigClicked={() => handleAddConfig(key)}>
            {list && list.map
              ? list.map((config, idx) => (
                  <ConfigEditor
                    {...getConfigEditorProps({
                      listName: listName,
                      index: idx,
                      createConfig,
                      ...configEditorSharedProps,
                    })}
                  />
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
  onAddConfigClicked,
}: {
  children?: ?React.Node,
  title: string,
  onAddConfigClicked: () => void,
}) {
  return (
    <Grid item>
      <Card>
        <CardContent>
          <Grid container direction="column" wrap="nowrap" spacing={3}>
            <Grid
              container
              item
              xs={12}
              justifyContent="space-between"
              alignItems="center">
              <Grid item>
                <Typography>{title}</Typography>
              </Grid>
              <Grid item>
                <IconButton
                  edge="end"
                  onClick={onAddConfigClicked}
                  data-testid={`add-${title.replace(/\s/g, '')}`}
                  aria-label="add new receiver configuration">
                  <AddCircleOutlineIcon color="primary" />
                </IconButton>
              </Grid>
            </Grid>
            {children}
          </Grid>
        </CardContent>
      </Card>
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

function emptyPagerDutyReceiver(): ReceiverPagerDutyConfig {
  return {};
}

function emptyPushoverReceiver(): ReceiverPushoverConfig {
  return {};
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
