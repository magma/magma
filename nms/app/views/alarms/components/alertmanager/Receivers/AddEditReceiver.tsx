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
import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import Editor from '../../common/Editor';
import EmailConfigEditor from './EmailConfigEditor';
import Grid from '@mui/material/Grid';
import IconButton from '@mui/material/IconButton';
import PagerDutyConfigEditor from './PagerDutyConfigEditor';
import PushoverConfigEditor from './PushoverConfigEditor';
import SlackConfigEditor from './SlackConfigEditor';
import TextField from '@mui/material/TextField';
import Typography from '@mui/material/Typography';
import WebhookConfigEditor from './WebhookConfigEditor';
import useForm from '../../../hooks/useForm';
import {useAlarmContext} from '../../AlarmContext';
import {useParams} from 'react-router-dom';
import {useSnackbars} from '../../../../../hooks/useSnackbar';

import {getErrorMessage} from '../../../../../util/ErrorUtils';
import type {
  AlertReceiver,
  ReceiverConfigListName,
  ReceiverEmailConfig,
  ReceiverPagerDutyConfig,
  ReceiverPushoverConfig,
  ReceiverSlackConfig,
  ReceiverWebhookConfig,
} from '../../AlarmAPIType';
import {Theme} from '@mui/material/styles';

import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  FormLabel,
  List,
  ListItem,
  ListItemText,
  MenuItem,
  OutlinedInput,
  Select,
  Tab,
  Tabs,
  ToggleButton,
  ToggleButtonGroup,
} from '@mui/material';
import {makeStyles} from '@mui/styles';
import {
  AltFormField,
  AltFormFieldSubheading,
} from '../../../../../components/FormField';
// import DialogTitle from '../../../../../theme/design-system/DialogTitle';
import {colors} from '../../../../../theme/default';
import {DEFAULT_ATTRIBUTE} from '@mui/system/cssVars/getInitColorSchemeScript';
import Text from '../../../../../theme/design-system/Text';
import {ArrowDropDown, DeleteOutline, ExpandMore} from '@mui/icons-material';
import {
  FlowMatchDirectionEnum,
  FlowMatchIpProtoEnum,
} from '../../../../../../generated';
import MenuButton from '../../../../../components/MenuButton';

const useStyles = makeStyles(() => ({
  tabBar: {
    backgroundColor: colors.primary.brightGray,
  },
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
  root: {
    '&$expanded': {
      minHeight: 'auto',
    },
    marginTop: '0px',
    marginBottom: '0px',
  },
  expanded: {marginTop: '-8px', marginBottom: '-8px'},
  block: {
    display: 'block',
  },
  flex: {display: 'flex'},
  panel: {flexGrow: 1},
  removeIcon: {alignSelf: 'baseline'},
  dialog: {height: '640px'},
  title: {textAlign: 'center', margin: 'auto', marginLeft: '0px'},
  description: {
    color: colors.primary.mirage,
  },
  switch: { margin: 'auto 0px' },
  expandMoreRotation: {
    transform: 'rotate(-180deg)',
    transition: '.3s',
  },
  expandLessRotation: {
    transition: '.3s',
  },
}));

type Props = {
  receiver: AlertReceiver;
  isNew: boolean;
  error: string;
  onChange: (receiver: AlertReceiver) => void;
};

const DEFAULT_RECEIVER: AlertReceiver = {name: ''};

const CONFIG_TYPES = {
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
} as const;

export default function AddEditReceiver(props: Props) {
  const {apiUtil} = useAlarmContext();
  const snackbars = useSnackbars();
  const {isNew, receiver} = props;
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
      const {listName, createConfig} = CONFIG_TYPES[
        configType as keyof typeof CONFIG_TYPES
      ];
      addListItem(listName, createConfig());
    },
    [addListItem],
  );

  const handleSave = React.useCallback(() => {
    async function makeApiCall() {
      try {
        const request = {
          receiver: formState,
          networkId: params.networkId!,
        };
        if (isNew) {
          await apiUtil.createReceiver(request);
          // onExit();
        } else {
          await apiUtil.editReceiver(request);
        }
        snackbars.success(`Successfully ${isNew ? 'added' : 'saved'} receiver`);
      } catch (error) {
        snackbars.error(`Unable to save receiver: ${getErrorMessage(error)}.`);
      }
    }
    void makeApiCall();
  }, [apiUtil, formState, isNew, params.networkId, snackbars]);

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
      // onExit={}
      data-testid="add-edit-receiver"
      title={receiver?.name || 'New Receiver'}
      description="Configure channels to notify when an alert fires">
      <Grid item>
        <Card>
          <CardContent>
            <Typography paragraph>Details</Typography>
            <TextField
              variant="standard"
              required
              id="name"
              label="Name"
              placeholder="Ex: Support Team"
              disabled={!isNew}
              value={formState.name}
              onChange={handleInputChange((val: string) => ({name: val}))}
              fullWidth
            />
          </CardContent>
        </Card>
      </Grid>

      {(Object.keys(CONFIG_TYPES) as Array<keyof typeof CONFIG_TYPES>).map(
        key => {
          const {
            friendlyName,
            createConfig,
            listName,
            ConfigEditor,
          } = CONFIG_TYPES[key];
          const list = formState[listName];
          return (
            // TODO[TS-migration] Typing of the ConfigEditor is weak here because the association between field name and editor is lost here
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
        },
      )}
    </Editor>
  );
}

function ConfigSection({
  children,
  title,
  onAddConfigClicked,
}: {
  children?: React.ReactNode;
  title: string;
  onAddConfigClicked: () => void;
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
                  aria-label="add new receiver configuration"
                  size="large">
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
  return ({} as unknown) as ReceiverPagerDutyConfig;
}

function emptyPushoverReceiver(): ReceiverPushoverConfig {
  return ({} as unknown) as ReceiverPushoverConfig;
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
  listName: ReceiverConfigListName;
  index: number;
  receiver: AlertReceiver;
  formState: Record<ReceiverConfigListName, Array<Partial<TConfig>>> & {
    name: string;
  };
  createConfig: () => Partial<TConfig>;
  updateListItem: (
    listName: ReceiverConfigListName,
    index: number,
    update: Partial<TConfig> | TConfig,
  ) => void;
  removeListItem: (listName: ReceiverConfigListName, index: number) => void;
}): {
  config: any;
  onUpdate: (update: Partial<TConfig>) => void;
  onReset: () => void;
  onDelete: () => void;
  isNew: boolean;
} {
  // The instance of a config such as ReceiverSlackConfig or ReceiverEmailConfig
  const config = formState[listName][index];
  const isNew = !receiver[listName]?.[index];

  const onUpdate = (update: Partial<TConfig> | TConfig) =>
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
      (receiver[listName]?.[index] || createConfig()) as
        | Partial<TConfig>
        | TConfig,
    );

  return {
    config,
    isNew,
    onUpdate,
    onReset,
    onDelete,
  };
}

export function ReceiverDialog(props: {
  onClose: () => void;
  open: boolean;
  receiver: AlertReceiver | null;
  isNew: boolean;
}) {
  const classes = useStyles();
  const [currentTab, setCurrentTab] = React.useState(1);
  const [receiver, setReceiver] = React.useState(
    props.receiver || DEFAULT_RECEIVER,
  );
  const [error, setError] = React.useState<string>('');
  const onSave = () => console.log('save : ', receiver);

  React.useEffect(() => {
    setReceiver(props.receiver || DEFAULT_RECEIVER);
    setError('');
    setCurrentTab(0);
  }, [props.open, props.receiver]);

  return (
    <Dialog
      open={props.open}
      onClose={() => props.onClose()}
      fullWidth={true}
      maxWidth="md"
      scroll="body"
      data-testid="">
      <DialogTitle>{false ? 'Edit Receiver' : 'Add New Receiver'}</DialogTitle>

      <Tabs
        value={currentTab}
        className={classes.tabBar}
        indicatorColor="primary"
        onChange={(_, tab) => setCurrentTab(tab as number)}>
        <Tab label="Receiver" />
        <Tab label="Channels" />
      </Tabs>
      <DialogContent>
        {currentTab === 0 && (
          <AddEditReceiverInfos
            receiver={receiver}
            onChange={newReceiver => setReceiver(newReceiver)}
            isNew={props.isNew}
            error={error}
          />
        )}
        {currentTab === 1 && (
          <AddEditReceiverChannels
            receiver={receiver}
            onChange={newReceiver => setReceiver(newReceiver)}
            isNew={props.isNew}
            error={error}
          />
        )}
      </DialogContent>
      <DialogActions>
        <Button onClick={() => props.onClose()}>Cancel</Button>
        <Button
          onClick={() => (currentTab === 0 ? setCurrentTab(1) : void onSave())}
          color="primary"
          variant="contained">
          {currentTab === 0 ? 'Next' : 'Save and Add Receiver'}
        </Button>
      </DialogActions>
    </Dialog>
  );
}

export function AddEditReceiverInfos(props: Props) {
  const {apiUtil} = useAlarmContext();
  const snackbars = useSnackbars();
  const {isNew, receiver} = props;
  const params = useParams();
  const [error, setError] = React.useState(props.error);


  return (
    <List>
      {error !== '' && (
        <AltFormField label={''}>
          <FormLabel data-testid="" error>
            {error}
          </FormLabel>
        </AltFormField>
      )}
      <AltFormField label={'Name'}>
        <OutlinedInput
          disabled={!isNew}
          data-testid="name"
          placeholder="Enter Name"
          fullWidth={true}
          value={receiver.name}
          onChange={({target}) =>
            props.onChange({...receiver, name: target.value})
          }
        />
      </AltFormField>
    </List>
  );
}

export function AddEditReceiverChannels(props: Props) {
  const classes = useStyles();
  const {apiUtil} = useAlarmContext();
  const snackbars = useSnackbars();
  const {isNew, receiver} = props;
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
      const {listName, createConfig} = CONFIG_TYPES[
        configType as keyof typeof CONFIG_TYPES
      ];
      addListItem(listName, createConfig());
    },
    [addListItem],
  );

  const handleSave = React.useCallback(() => {
    async function makeApiCall() {
      try {
        const request = {
          receiver: formState,
          networkId: params.networkId!,
        };
        if (isNew) {
          await apiUtil.createReceiver(request);
          // onExit();
        } else {
          await apiUtil.editReceiver(request);
        }
        snackbars.success(`Successfully ${isNew ? 'added' : 'saved'} receiver`);
      } catch (error) {
        snackbars.error(`Unable to save receiver: ${getErrorMessage(error)}.`);
      }
    }
    void makeApiCall();
  }, [apiUtil, formState, isNew, params.networkId, snackbars]);

  const configEditorSharedProps = {
    receiver,
    formState,
    updateListItem,
    removeListItem,
  };
  return (
    <List>
      <ListItem>
        <MenuButton label="Create New" size="small">
          <MenuItem
            data-testid="newBaseNameMenuItem"
            onClick={() => console.log('Stack Channel')}>
            <Text variant="body2">Stack Channel</Text>
          </MenuItem>
          <MenuItem onClick={() => console.log('Email')}>
            <Text variant="body2">Email</Text>
          </MenuItem>
          <MenuItem
            data-testid="newRatingGroupMenuItem"
            onClick={() => console.log('Webhook')}>
            <Text variant="body2">Webhook</Text>
          </MenuItem>
          <MenuItem
            data-testid="newRatingGroupMenuItem"
            onClick={() => console.log('Pager Duty')}>
            <Text variant="body2">Pager Duty</Text>
          </MenuItem>
          <MenuItem
            data-testid="newRatingGroupMenuItem"
            onClick={() => console.log('Pushover')}>
            <Text variant="body2">Pushover</Text>
          </MenuItem>
        </MenuButton>
      </ListItem>
      <ListItem>
        <div className={classes.flex}>

        </div>
      </ListItem>
    </List>
  );
}

function ChannelItem(props: {channelTitle: string}) {
  const classes = useStyles();
  return (
    <Accordion defaultExpanded className={classes.panel}>
    <AccordionSummary
      classes={{
        root: classes.root,
        expanded: classes.expanded,
      }}
      expandIcon={<ExpandMore />}>
      <Grid container justifyContent="space-between">
        <Grid item className={classes.title}>
            <Text weight="medium" variant="body2">
              {props.channelTitle}
          
          </Text>
        </Grid>
        <Grid item>
          <IconButton
            className={classes.removeIcon}
            onClick={() => console.log('delete')}
            size="large">
            <DeleteOutline />
          </IconButton>
        </Grid>
      </Grid>
    </AccordionSummary>
    <AccordionDetails
      classes={{
        root: classes.block,
      }}>
      <div className={classes.flex}>
      </div>
    </AccordionDetails>
  </Accordion>
  )
}

function SlackConfig(isNew: boolean) {
  const classes = useStyles();
  const [advancedSettings, setAdvancedSettings] = React.useState(false)

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
    <Button endIcon={<ExpandMore className={advancedSettings ? classes.expandMoreRotation : classes.expandLessRotation } />} variant='text' onClick={() => setAdvancedSettings(settings => !settings)}>
      Advanced Settings
    </Button>
    {advancedSettings &&
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
    }
  </Grid>
  )
}