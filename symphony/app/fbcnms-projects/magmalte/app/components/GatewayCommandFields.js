/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import Button from '@fbcnms/ui/components/design-system/Button';
import Check from '@material-ui/icons/Check';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import Divider from '@material-ui/core/Divider';
import Fade from '@material-ui/core/Fade';
import FormField from './FormField';
import Input from '@material-ui/core/Input';
import LinearProgress from '@material-ui/core/LinearProgress';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import grey from '@material-ui/core/colors/grey';

import nullthrows from '@fbcnms/util/nullthrows';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';

const useStyles = makeStyles(() => ({
  input: {
    margin: '10px 0',
    width: '100%',
  },
  divider: {
    margin: '10px 0',
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
          <Button variant="text" onClick={props.onClose} skin="primary">
            Close
          </Button>
        </DialogActions>
      )}
    </>
  );
}

type ChildProps = {gatewayID: string};

function RebootButton(props: ChildProps) {
  const {match} = useRouter();
  const enqueueSnackbar = useEnqueueSnackbar();
  const [showCheck, setShowCheck] = useState(false);

  const onClick = () => {
    const {gatewayID} = props;
    MagmaV1API.postNetworksByNetworkIdGatewaysByGatewayIdCommandReboot({
      networkId: nullthrows(match.params.networkId),
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
        <Button variant="text" onClick={onClick} skin="primary">
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
  const {match} = useRouter();
  const enqueueSnackbar = useEnqueueSnackbar();
  const [showCheck, setShowCheck] = useState(false);

  const onClick = () => {
    const {gatewayID} = props;
    MagmaV1API.postNetworksByNetworkIdGatewaysByGatewayIdCommandRestartServices(
      {
        networkId: nullthrows(match.params.networkId),
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
        <Button variant="text" onClick={onClick} skin="primary">
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
  const {match} = useRouter();
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
      networkId: nullthrows(match.params.networkId),
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
        <Button variant="text" onClick={onClick} skin="primary">
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

function PingCommandControls(props: ChildProps) {
  const classes = useStyles();
  const {match} = useRouter();
  const enqueueSnackbar = useEnqueueSnackbar();
  const [pingHosts, setPingHosts] = useState('');
  const [pingPackets, setPingPackets] = useState('');
  const [pingResponse, setPingResponse] = useState();
  const [showProgress, setShowProgress] = useState();

  const onClick = () => {
    const {gatewayID} = props;
    const hosts = pingHosts.split('\n').filter(host => host);
    const packets = parseInt(pingPackets);
    const params = {
      hosts,
      packets,
    };

    setShowProgress(true);
    MagmaV1API.postNetworksByNetworkIdGatewaysByGatewayIdCommandPing({
      networkId: nullthrows(match.params.networkId),
      gatewayId: gatewayID,
      pingRequest: params,
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
        <Button variant="text" onClick={onClick} skin="primary">
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

function GenericCommandControls(props: ChildProps) {
  const classes = useStyles();
  const {match} = useRouter();
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
      networkId: nullthrows(match.params.networkId),
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
        <Button variant="text" onClick={onClick} skin="primary">
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
