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
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {DataRows} from './DataGrid';
import type {generic_command_response} from '../../generated/MagmaAPIBindings';

import Button from '@material-ui/core/Button';
import Check from '@material-ui/icons/Check';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import DataGrid from './DataGrid';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import Divider from '@material-ui/core/Divider';
import Fade from '@material-ui/core/Fade';
// $FlowFixMe migrated to typescript
import FormField from './FormField';
import FormLabel from '@material-ui/core/FormLabel';
import Input from '@material-ui/core/Input';
import LinearProgress from '@material-ui/core/LinearProgress';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
// $FlowFixMe migrated to typescript
import LoadingFiller from './LoadingFiller';
import MagmaV1API from '../../generated/WebClient';
import React from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import Text from '../theme/design-system/Text';
import grey from '@material-ui/core/colors/grey';
// $FlowFixMe migrated to typescript
import nullthrows from '../../shared/util/nullthrows';
import useMagmaAPI from '../../api/useMagmaAPIFlow';

// $FlowFixMe migrated to typescript
import {AltFormField} from './FormField';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useState} from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../app/hooks/useSnackbar';
import {useParams} from 'react-router-dom';

const useStyles = makeStyles(theme => ({
  input: {
    margin: '10px 0',
    width: '100%',
  },
  divider: {
    margin: '10px 0',
  },
  jsonTextarea: {
    fontFamily: 'monospace',
    height: '95%',
    border: 'none',
    margin: theme.spacing(2),
    width: '100%',
  },
}));

type Props = {
  onClose?: () => void,
  gatewayID: string,
  showRestartCommand: boolean,
  showRebootEnodebCommand: boolean,
  showPingCommand: boolean,
  showGenericCommand: boolean,
};

function CommandResponse(props) {
  return (
    <pre
      style={{
        backgroundColor: grey[100],
        fontSize: '12px',
        color: grey[900],
      }}>
      {props.showProgressBar && <LinearProgress />}
      <code>{props.response}</code>
    </pre>
  );
}

export default function GatewayCommandFields(props: Props) {
  return (
    <>
      <DialogContent>
        <RebootButton gatewayID={props.gatewayID} />
        {props.showRestartCommand && (
          <RestartServicesButton gatewayID={props.gatewayID} />
        )}
        {props.showRebootEnodebCommand && (
          <RebootEnodebControls gatewayID={props.gatewayID} />
        )}
        {props.showPingCommand && (
          <PingCommandControls gatewayID={props.gatewayID} />
        )}
        {props.showGenericCommand && (
          <GenericCommandControls gatewayID={props.gatewayID} />
        )}
      </DialogContent>
      {props.onClose && (
        <DialogActions>
          <Button variant="outlined" onClick={props.onClose} color="primary">
            Close
          </Button>
        </DialogActions>
      )}
    </>
  );
}

type ChildProps = {gatewayID: string};

function RebootButton(props: ChildProps) {
  const params = useParams();
  const enqueueSnackbar = useEnqueueSnackbar();
  const [showCheck, setShowCheck] = useState(false);

  const onClick = () => {
    const {gatewayID} = props;
    MagmaV1API.postNetworksByNetworkIdGatewaysByGatewayIdCommandReboot({
      networkId: nullthrows(params.networkId),
      gatewayId: gatewayID,
    })
      .then(_resp => {
        enqueueSnackbar('Successfully initiated reboot', {variant: 'success'});
        setShowCheck(true);
        setTimeout(() => setShowCheck(false), 5000);
      })
      .catch(error =>
        enqueueSnackbar('Reboot failed: ' + error.response.data.message, {
          variant: 'error',
        }),
      );
  };

  return (
    <>
      <Text variant="subtitle1">Reboot</Text>
      <FormField
        label="Reboot Device"
        tooltip="Reboot the Magma gateway server">
        <Button variant="outlined" onClick={onClick} color="primary">
          Reboot
        </Button>
        <Fade in={showCheck} timeout={500}>
          <Check style={{verticalAlign: 'middle'}} htmlColor="green" />
        </Fade>
      </FormField>
    </>
  );
}

function RestartServicesButton(props: ChildProps) {
  const params = useParams();
  const enqueueSnackbar = useEnqueueSnackbar();
  const [showCheck, setShowCheck] = useState(false);

  const onClick = () => {
    const {gatewayID} = props;
    MagmaV1API.postNetworksByNetworkIdGatewaysByGatewayIdCommandRestartServices(
      {
        networkId: nullthrows(params.networkId),
        gatewayId: gatewayID,
        services: [],
      },
    )
      .then(_resp => {
        enqueueSnackbar('Successfully initiated service restart', {
          variant: 'success',
        });
        setShowCheck(true);
        setTimeout(() => setShowCheck(false), 5000);
      })
      .catch(error =>
        enqueueSnackbar(
          'Restart services failed: ' + error.response.data.message,
          {variant: 'error'},
        ),
      );
  };

  return (
    <>
      <FormField
        label="Restart Services"
        tooltip="Restart all MagmaD services on this gateway">
        <Button variant="outlined" onClick={onClick} color="primary">
          Restart Services
        </Button>
        <Fade in={showCheck} timeout={500}>
          <Check style={{verticalAlign: 'middle'}} htmlColor="green" />
        </Fade>
      </FormField>
    </>
  );
}

function RebootEnodebControls(props: ChildProps) {
  const classes = useStyles();
  const {networkId} = useParams();
  const enqueueSnackbar = useEnqueueSnackbar();
  const [showProgress, setShowProgress] = useState(false);
  const [rebootResponse, setRebootResponse] = useState();
  const [enodebSerial, setEnodebSerial] = useState('');

  const onClick = () => {
    const {gatewayID} = props;
    const params =
      enodebSerial.length > 0
        ? {
            command: 'reboot_enodeb',
            params: {shell_params: ([enodebSerial]: any)},
          }
        : {
            command: 'reboot_all_enodeb',
            params: {},
          };

    setShowProgress(true);
    MagmaV1API.postNetworksByNetworkIdGatewaysByGatewayIdCommandGeneric({
      networkId: nullthrows(networkId),
      gatewayId: gatewayID,
      parameters: params,
    })
      .then(resp => setRebootResponse(JSON.stringify(resp, null, 2)))
      .catch(error =>
        enqueueSnackbar(
          'Reboot eNodeB failed: ' + error.response.data.message,
          {variant: 'error'},
        ),
      )
      .finally(() => setShowProgress(false));
  };

  return (
    <div>
      <Divider className={classes.divider} />
      <Text variant="subtitle1">Reboot eNodeB</Text>
      <FormField label="eNodeB Serial ID">
        <Input
          className={classes.input}
          value={enodebSerial}
          onChange={({target}) => setEnodebSerial(target.value)}
          placeholder="Leave empty to reboot every connected eNodeB"
        />
      </FormField>
      <FormField label="">
        <Button variant="outlined" onClick={onClick} color="primary">
          Reboot
        </Button>
      </FormField>
      <FormField label="">
        <CommandResponse
          response={rebootResponse}
          showProgressBar={showProgress}
        />
      </FormField>
    </div>
  );
}

export function PingCommandControls(props: ChildProps) {
  const classes = useStyles();
  const params = useParams();
  const enqueueSnackbar = useEnqueueSnackbar();
  const [pingHosts, setPingHosts] = useState('');
  const [pingPackets, setPingPackets] = useState('');
  const [pingResponse, setPingResponse] = useState();
  const [showProgress, setShowProgress] = useState();

  const onClick = () => {
    const {gatewayID} = props;
    const hosts = pingHosts.split('\n').filter(host => host);
    const packets = parseInt(pingPackets);

    setShowProgress(true);
    MagmaV1API.postNetworksByNetworkIdGatewaysByGatewayIdCommandPing({
      networkId: nullthrows(params.networkId),
      gatewayId: gatewayID,
      pingRequest: {
        hosts,
        packets,
      },
    })
      .then(resp => setPingResponse(JSON.stringify(resp, null, 2)))
      .catch(error =>
        enqueueSnackbar('Ping failed: ' + error.response.data.message, {
          variant: 'error',
        }),
      )
      .finally(() => setShowProgress(false));
  };

  return (
    <div>
      <Divider className={classes.divider} />
      <Text variant="subtitle1">Ping</Text>
      <FormField label="Host(s) (one per line)">
        <Input
          className={classes.input}
          value={pingHosts}
          onChange={({target}) => setPingHosts(target.value)}
          placeholder="E.g. example.com"
          multiline={true}
        />
      </FormField>
      <FormField label="Packets (default 4)">
        <Input
          className={classes.input}
          value={pingPackets}
          onChange={({target}) => setPingPackets(target.value)}
          placeholder="E.g. 4"
          type="number"
        />
      </FormField>
      <FormField label="">
        <Button variant="outlined" onClick={onClick} color="primary">
          Ping
        </Button>
      </FormField>
      <FormField label="">
        <CommandResponse
          response={pingResponse}
          showProgressBar={showProgress}
        />
      </FormField>
    </div>
  );
}

export function GenericCommandControls(props: ChildProps) {
  const classes = useStyles();
  const {networkId} = useParams();
  const enqueueSnackbar = useEnqueueSnackbar();
  const [commandName, setCommandName] = useState('');
  const [commandParams, setCommandParams] = useState('{\n}');
  const [genericResponse, setGenericResponse] = useState();
  const [showProgress, setShowProgress] = useState();

  const onClick = () => {
    const {gatewayID} = props;
    let params = {};
    try {
      params = JSON.parse(commandParams);
    } catch (e) {
      enqueueSnackbar('Error parsing params: ' + e, {variant: 'error'});
      return;
    }
    const parameters = {
      command: commandName,
      params,
    };

    setShowProgress(true);
    MagmaV1API.postNetworksByNetworkIdGatewaysByGatewayIdCommandGeneric({
      networkId: nullthrows(networkId),
      gatewayId: gatewayID,
      parameters,
    })
      .then(resp => setGenericResponse(JSON.stringify(resp, null, 2)))
      .catch(error =>
        enqueueSnackbar(
          'Generic command failed: ' + error.response.data.message,
          {variant: 'error'},
        ),
      )
      .finally(() => setShowProgress(false));
  };

  return (
    <div>
      <Divider className={classes.divider} />
      <Text variant="subtitle1">Generic</Text>
      <FormField label="Command">
        <Input
          className={classes.input}
          value={commandName}
          onChange={({target}) => setCommandName(target.value)}
          placeholder="Command name"
        />
      </FormField>
      <FormField label="Parameters">
        <Input
          className={classes.input}
          value={commandParams}
          onChange={({target}) => setCommandParams(target.value)}
          multiline={true}
          style={{fontFamily: 'monospace', fontSize: '14px'}}
        />
      </FormField>
      <FormField label="">
        <Button variant="outlined" onClick={onClick} color="primary">
          Execute
        </Button>
      </FormField>
      <FormField label="">
        <CommandResponse
          response={genericResponse}
          showProgressBar={showProgress}
        />
      </FormField>
    </div>
  );
}

type FileComponentProps = {
  title?: string,
  content: string,
  error?: string,
};
function FileComponent(props: FileComponentProps) {
  const classes = useStyles();
  let content = props.content.replace(/\\n/g, '\n');
  content = content.slice(1, -1);

  return (
    <List>
      {props.title ?? <ListItemText> {props.title} </ListItemText>}
      {props.error !== '""' && (
        <ListItemText>
          <AltFormField label={''}>
            <FormLabel data-testid="fileError" error>
              {props.error}
            </FormLabel>
          </AltFormField>
        </ListItemText>
      )}
      <ListItem>
        <textarea
          data-testid="fileContent"
          rows="8"
          className={classes.jsonTextarea}
          autoCapitalize="none"
          autoComplete="none"
          autoCorrect="none"
          spellCheck={false}
          value={content}
        />
      </ListItem>
    </List>
  );
}

const TROUBLESHOOTING_HINTS = {
  FLUENTD_SUCCESS:
    'Gateway contains fluentd parameters, Verify if \
control proxy config contains the right fluentd address and port information, \
typically fluend address is fluentd.<orc8r_domain_name> and fluend port is 24224',
  FLUENTD_MISSING:
    'Gateway is missing fluentd parameters, Add the \
right fluentd address and port information to the control proxy config on the \
gateway, typically fluend address is fluentd.<orc8r_domain_name> and fluend port\
 is 24224',
  AGG_API_SUCCESS:
    'event and log aggregation API are returning successful \
responses',
};

const CONTROL_PROXY_CONTENT = 'cat /var/opt/magma/configs/control_proxy.yml';
const FLUENT_BIT_LOGS = 'journalctl -u magma@td-agent-bit  -n 10';
export function TroubleshootingControl(props: ChildProps) {
  const params = useParams();
  const [
    controlProxyContent,
    setControlProxyContent,
  ] = useState<generic_command_response>({});
  const [
    tdAgentLogsContent,
    setTdAgentLogsContent,
  ] = useState<generic_command_response>({});
  const networkId = nullthrows(params.networkId);
  const controlProxyParams = {
    command: 'bash',
    params: {
      shell_params: [`-c '${CONTROL_PROXY_CONTENT}'`],
    },
  };
  const {isLoading: isProxyFileLoading} = useMagmaAPI(
    MagmaV1API.postNetworksByNetworkIdGatewaysByGatewayIdCommandGeneric,
    // $FlowIssue[incompatible-call]
    {networkId, gatewayId: props.gatewayID, parameters: controlProxyParams},
    useCallback(response => setControlProxyContent(response), []),
  );
  const tdAgentBitLogs = {
    command: 'bash',
    params: {
      shell_params: [`-c '${FLUENT_BIT_LOGS}'`],
    },
  };
  const {isLoading: isTdAgentBitLogsLoading} = useMagmaAPI(
    MagmaV1API.postNetworksByNetworkIdGatewaysByGatewayIdCommandGeneric,
    // $FlowIssue[incompatible-call]
    {networkId, gatewayId: props.gatewayID, parameters: tdAgentBitLogs},
    useCallback(response => setTdAgentLogsContent(response), []),
  );
  const {
    isLoading: isEventAPILoading,
    error,
  } = useMagmaAPI(MagmaV1API.getEventsByNetworkIdAboutCount, {networkId});

  if (isProxyFileLoading || isEventAPILoading || isTdAgentBitLogsLoading) {
    return <LoadingFiller />;
  }

  const errContent = JSON.stringify(
    controlProxyContent?.response?.['stderr'] ?? {},
  );
  const fileContent = JSON.stringify(
    controlProxyContent?.response?.['stdout'] ?? {},
    null,
    2,
  );
  const tdErrContent = JSON.stringify(
    controlProxyContent?.response?.['stderr'] ?? {},
  );
  const tdAgentLogsFileContent = JSON.stringify(
    tdAgentLogsContent?.response?.['stdout'] ?? {},
    null,
    2,
  );

  const containsFluentdParams =
    fileContent.includes('fluentd_address') &&
    fileContent.includes('fluentd_port');

  const kpiData: DataRows[] = [
    [
      {
        category: 'Control Proxy Config Validation',
        value: containsFluentdParams ? 'Good' : 'Bad',
        status: containsFluentdParams,
        statusCircle: true,
        tooltip: containsFluentdParams
          ? TROUBLESHOOTING_HINTS.FLUENTD_SUCCESS
          : TROUBLESHOOTING_HINTS.FLUENTD_MISSING,
        collapse: <FileComponent content={fileContent} error={errContent} />,
      },
    ],
    [
      {
        category: 'API validation',
        value: error == null ? 'Good' : 'Bad',
        status: error == null,
        statusCircle: true,
        tooltip:
          error == null
            ? TROUBLESHOOTING_HINTS.AGG_API_SUCCESS
            : `event and log aggregation api is failing,  ${
                error ?? 'internal server error'
              }`,
      },
    ],
    [
      {
        category: 'Fluent Bit Logs',
        value: '',
        tooltip: 'fluend bit logs',
        collapse: (
          <FileComponent
            content={tdAgentLogsFileContent}
            error={tdErrContent}
          />
        ),
      },
    ],
  ];
  return <DataGrid data={kpiData} />;
}
