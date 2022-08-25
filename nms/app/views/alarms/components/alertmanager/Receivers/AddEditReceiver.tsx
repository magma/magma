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
import Grid from '@mui/material/Grid';
import IconButton from '@mui/material/IconButton';
// import useForm from '../../../hooks/useForm';

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
  MenuItem,
  OutlinedInput,
  Tab,
  Tabs,
} from '@mui/material';
import {AltFormField} from '../../../../../components/FormField';
// import {getErrorMessage} from '../../../../../util/ErrorUtils';
import {makeStyles} from '@mui/styles';
// import {useAlarmContext} from '../../AlarmContext';
// import {useParams} from 'react-router-dom';
// import {useSnackbars} from '../../../../../hooks/useSnackbar';
import type {AlertReceiver} from '../../AlarmAPIType';
// import DialogTitle from '../../../../../theme/design-system/DialogTitle';
import MenuButton from '../../../../../components/MenuButton';
import Text from '../../../../../theme/design-system/Text';
import {DeleteOutline, ExpandMore} from '@mui/icons-material';
import {colors} from '../../../../../theme/default';

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
  switch: {margin: 'auto 0px'},
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
const configTypes: {[key: string]: string} = {
  slack_configs: 'Slack Channel',
  email_configs: 'Email',
  webhook_configs: 'Webhook',
  pagerduty_configs: 'PagerDuty',
  pushover_configs: 'Pushover',
} as const;

export default function ReceiverDialog(props: {
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
  const {isNew, receiver} = props;
  const [error, setError] = React.useState(props.error);

  React.useEffect(() => setError(props.error), [props.error]);

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
  // const {apiUtil} = useAlarmContext();
  // const snackbars = useSnackbars();
  // const {isNew, receiver} = props;
  // const params = useParams();
  const [receiver, setReceiver] = React.useState();

  const handleAddNewChannel = React.useCallback((configType: string) => {}, []);

  // const handleSave = React.useCallback(() => {
  //   async function makeApiCall() {
  //     try {
  //       const request = {
  //         receiver: formState,
  //         networkId: params.networkId!,
  //       };
  //       if (isNew) {
  //         await apiUtil.createReceiver(request);
  //         // onExit();
  //       } else {
  //         await apiUtil.editReceiver(request);
  //       }
  //       snackbars.success(`Successfully ${isNew ? 'added' : 'saved'} receiver`);
  //     } catch (error) {
  //       snackbars.error(`Unable to save receiver: ${getErrorMessage(error)}.`);
  //     }
  //   }
  //   void makeApiCall();
  // }, [apiUtil, formState, isNew, params.networkId, snackbars]);

  return (
    <List>
      <ListItem>
        <MenuButton label="Create New" size="small">
          {Object.keys(configTypes).map((key: string) => (
            <MenuItem onClick={() => handleAddNewChannel(key)}>
              <Text variant="body2">{configTypes[key]}</Text>
            </MenuItem>
          ))}
        </MenuButton>
      </ListItem>
      <ListItem>
        {Object.keys(props.receiver).map((key: keyof AlertReceiver) => {
          if (key !== 'name') {
            (props.receiver[key] || []).map(config => (
              <ChannelItem channelTitle={configTypes[key]} />
            ));
          }
        })}
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
        <div className={classes.flex} />
      </AccordionDetails>
    </Accordion>
  );
}
