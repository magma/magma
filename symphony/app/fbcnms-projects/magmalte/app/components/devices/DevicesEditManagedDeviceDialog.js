/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {symphony_device_config} from '@fbcnms/magma-api';

import Button from '@fbcnms/ui/components/design-system/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormGroup from '@material-ui/core/FormGroup';
import FormLabel from '@material-ui/core/FormLabel';
import Input from '@material-ui/core/Input';
import ListItemText from '@material-ui/core/ListItemText';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import Select from '@material-ui/core/Select';
import TextField from '@material-ui/core/TextField';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';

import nullthrows from '@fbcnms/util/nullthrows';
import {makeStyles} from '@material-ui/styles';
import {useEffect, useState} from 'react';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
  formContainer: {
    paddingBottom: theme.spacing(2),
  },
  formGroup: {
    marginLeft: theme.spacing(2),
    marginBottom: theme.spacing(2),
  },
  title: {
    margin: '15px 0 5px',
  },
}));

type Props = {
  title: string,
  onCancel: () => void,
  onSave: string => void,
};

export default function DevicesEditManagedDeviceDialog(props: Props) {
  const classes = useStyles();
  const {match} = useRouter();

  const {deviceID: initialDeviceID} = match.params;

  const [deviceID, setDeviceID] = useState(initialDeviceID || '');
  const [error, setError] = useState('');

  const standardPlatforms = {
    snmp: 'SNMP + ping',
    ping: 'ping',
  };
  const [availablePlatforms, setAvailablePlatforms] = useState<{
    [string]: string,
  }>(standardPlatforms);

  // sectioned device configs
  const [hostTextbox, setHostTextbox] = useState('0.0.0.0');
  const [agentTextbox, setAgentTextbox] = useState('');
  const [platformTextbox, setPlatformTextbox] = useState('Snmp');
  const [typeTextbox, setTypeTextbox] = useState('snmp');
  const [frinxTextbox, setFrinxTextbox] = useState(
    JSON.stringify(
      {
        authorization: 'Basic YWRtaW46YWRtaW4=', // admin:admin
        device_type: 'ios',
        device_version: '15.2',
        frinx_port: 8181,
        host: 'frinx',
        password: 'frinx',
        port: 23,
        transport_type: 'telnet',
        username: 'username',
      },
      null,
      2,
    ),
  );

  const [snmpCommunityTextbox, setSnmpCommunityTextbox] = useState('public');
  const [snmpVersionTextbox, setSnmpVersionTextbox] = useState('v1');

  const [cambiumTextbox, setCambiumTextbox] = useState(
    JSON.stringify(
      {
        client_id: 'Ya4dNAAYSUoFMUSs',
        client_ip: '10.0.0.1',
        client_mac: '58:C1:7A:90:36:50',
        client_secret: '2kBgGnMr69NgNdeHB3s7x4GzUODLkc',
      },
      null,
      2,
    ),
  );
  const [otherChannelTextbox, setOtherChannelTextbox] = useState(
    JSON.stringify({}, null, 2),
  );
  const [configTextbox, setConfigTextbox] = useState('{}'); // device configs

  const genManagedDeviceConfig = (): symphony_device_config => {
    const config: symphony_device_config = {
      host: hostTextbox,
      platform: platformTextbox,
      device_type: typeTextbox.length > 0 ? [typeTextbox] : [],
      channels: {
        snmp_channel: {
          community: snmpCommunityTextbox,
          version: snmpVersionTextbox,
        },
      },
    };

    function setConfigIfExist(fieldName, textbox, setter) {
      if (!textbox || textbox.length === 0) {
        return;
      }
      try {
        setter(JSON.parse(textbox));
      } catch (err) {
        throw {message: `${fieldName}: ${err.message}`};
      }
    }

    setConfigIfExist('Frinx Channel', frinxTextbox, jsonValue => {
      if (config.channels) {
        config.channels.frinx_channel = jsonValue;
      }
    });

    setConfigIfExist('Cambium Channel', cambiumTextbox, jsonValue => {
      if (config.channels) {
        config.channels.cambium_channel = jsonValue;
      }
    });

    setConfigIfExist('Other Channel', otherChannelTextbox, jsonValue => {
      if (config.channels) {
        config.channels.other_channel = {
          channel_props: jsonValue,
        };
      }
    });

    setConfigIfExist(
      'Device Config',
      configTextbox,
      jsonValue => (config.device_config = JSON.stringify(jsonValue)),
    );

    return config;
  };

  const onEdit = async () => {
    try {
      await MagmaV1API.putSymphonyByNetworkIdDevicesByDeviceId({
        networkId: nullthrows(match.params.networkId),
        deviceId: deviceID,
        symphonyDevice: {
          id: deviceID,
          name: deviceID,
          managing_agent: agentTextbox,
          config: genManagedDeviceConfig(),
        },
      });
      props.onSave(deviceID);
    } catch (error) {
      setError(`${error.message} ${error.response?.data?.message || ''}`);
      return;
    }
  };

  const onCreate = async () => {
    try {
      await MagmaV1API.postSymphonyByNetworkIdDevices({
        networkId: nullthrows(match.params.networkId),
        symphonyDevice: {
          id: deviceID,
          name: deviceID,
          managing_agent: agentTextbox,
          config: genManagedDeviceConfig(),
        },
      });
      setDeviceID(deviceID);
      props.onSave(deviceID);
    } catch (error) {
      setError(`${error.message} ${error.response?.data?.message || ''}`);
      return;
    }
  };

  // TODO: separate out create from edit flow so we don't have extra api call
  const {isLoading, error: responseError, response} = useMagmaAPI(
    MagmaV1API.getSymphonyByNetworkIdDevicesByDeviceId,
    {
      networkId: nullthrows(match.params.networkId),
      deviceId: initialDeviceID,
    },
  );

  useEffect(() => {
    // TODO: separate out create from edit flow so we don't have extra api call
    if (initialDeviceID) {
      const initialDeviceConfig = response?.config || {};

      setHostTextbox(initialDeviceConfig.host || '');
      setAgentTextbox(response?.managing_agent || '');
      setPlatformTextbox(initialDeviceConfig.platform || '');
      // TODO: support more than 1 device_type in the list
      setTypeTextbox(initialDeviceConfig.device_type?.[0] || '');
      setFrinxTextbox(
        JSON.stringify(initialDeviceConfig.channels?.frinx_channel, null, 2),
      );
      setSnmpCommunityTextbox(
        initialDeviceConfig.channels?.snmp_channel?.community || '',
      );
      setSnmpVersionTextbox(
        initialDeviceConfig.channels?.snmp_channel?.version || '',
      );
      setCambiumTextbox(
        JSON.stringify(initialDeviceConfig.channels?.cambium_channel, null, 2),
      );
      setOtherChannelTextbox(
        JSON.stringify(
          initialDeviceConfig.channels?.other_channel?.channel_props,
          null,
          2,
        ),
      );

      // special case config where config should be JSON.
      // This attempts to format as JSON, failing that, fill in text box as-is.
      try {
        setConfigTextbox(
          JSON.stringify(
            JSON.parse(initialDeviceConfig.device_config || ''),
            null,
            2,
          ),
        );
      } catch (err) {
        setConfigTextbox(initialDeviceConfig.device_config);
      }
    }
  }, [initialDeviceID, response]);

  useEffect(() => {
    // TODO: separate out create from edit flow so we don't have extra api call
    if (initialDeviceID && responseError) {
      if (responseError.response.status === 404) {
        setError(
          'Warning! Missing config - please enter a new device config below.',
        );
      } else {
        setError(responseError.message);
      }
    }
  }, [initialDeviceID, responseError]);

  // if initialDeviceID is set, then don't allow modifying deviceID
  const deviceIDcontent = initialDeviceID ? (
    <>
      <FormLabel>Device ID: </FormLabel>
      <div>{initialDeviceID}</div>
    </>
  ) : (
    <TextField
      required
      className={classes.input}
      label="Device ID"
      margin="normal"
      onChange={({target}) => {
        setDeviceID(target.value);
      }}
      value={deviceID}
    />
  );

  const content = (
    <div className={classes.formContainer}>
      {deviceIDcontent}

      <FormLabel>Host</FormLabel>
      <FormGroup row className={classes.formGroup}>
        <TextField
          required
          label="IP"
          className={classes.input}
          onChange={({target}) => setHostTextbox(target.value)}
          value={hostTextbox}
        />
      </FormGroup>

      <FormLabel>Agent</FormLabel>
      <FormGroup row className={classes.formGroup}>
        <TextField
          label="Agent"
          className={classes.input}
          onChange={({target}) => setAgentTextbox(target.value)}
          value={agentTextbox}
        />
      </FormGroup>

      <FormLabel>Platform</FormLabel>
      <FormGroup row className={classes.formGroup}>
        <Select
          value={platformTextbox}
          onChange={({target}) => setPlatformTextbox(target.value)}
          input={<Input id="types" />}>
          {Object.keys(availablePlatforms).map(key => (
            <MenuItem key={key} value={key}>
              <ListItemText primary={availablePlatforms[key]} />
            </MenuItem>
          ))}
        </Select>

        <TextField
          required
          label="Platform value (or custom)"
          className={classes.input}
          onChange={({target}) => {
            const targetString = target.value;
            if (!(targetString in standardPlatforms)) {
              setAvailablePlatforms({
                ...standardPlatforms,
                [targetString]: '<Custom Platform>',
              });
              setPlatformTextbox(targetString);
            } else {
              setPlatformTextbox(targetString);
            }
          }}
          value={platformTextbox}
        />
      </FormGroup>

      <TextField
        label="Device Type"
        style={{display: 'none'}} // TODO: show after implemented in agent
        className={classes.input}
        onChange={({target}) => setTypeTextbox(target.value)}
        value={typeTextbox}
      />

      <FormLabel>SNMP Channel Config</FormLabel>
      <FormGroup row className={classes.formGroup}>
        <TextField
          label="Community"
          className={classes.input}
          onChange={({target}) => setSnmpCommunityTextbox(target.value)}
          value={snmpCommunityTextbox}
        />

        <TextField
          label="Version"
          className={classes.input}
          onChange={({target}) => setSnmpVersionTextbox(target.value)}
          value={snmpVersionTextbox}
        />
      </FormGroup>

      <FormLabel>Additional Channel Configs</FormLabel>
      <FormGroup row className={classes.formGroup}>
        <TextField
          label="Frinx Channel Config"
          style={{display: 'none'}} // TODO: show after implemented in agent
          multiline={true}
          lines={Infinity}
          onChange={({target}) => setFrinxTextbox(target.value)}
          value={frinxTextbox}
        />

        <TextField
          label="Cambium Channel Config"
          style={{display: 'none'}} // TODO: show after implemented in agent
          className={classes.input}
          multiline={true}
          lines={Infinity}
          onChange={({target}) => setCambiumTextbox(target.value)}
          value={cambiumTextbox}
        />

        <TextField
          label="Other Channel Config Props"
          className={classes.input}
          multiline={true}
          lines={Infinity}
          onChange={({target}) => setOtherChannelTextbox(target.value)}
          value={otherChannelTextbox}
        />
      </FormGroup>

      <FormLabel>Device Configs</FormLabel>
      <FormGroup row className={classes.formGroup}>
        <TextField
          className={classes.input}
          multiline={true}
          lines={Infinity}
          onChange={({target}) => setConfigTextbox(target.value)}
          value={configTextbox}
        />
      </FormGroup>
    </div>
  );

  return (
    <Dialog open={true} onClose={props.onCancel} fullWidth={true} scroll="body">
      <DialogTitle>{props.title}</DialogTitle>
      <DialogContent>
        {error ? <FormLabel error>{error}</FormLabel> : null}
        {isLoading ? <LoadingFiller /> : content}
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onCancel} skin="regular">
          Cancel
        </Button>
        <Button onClick={initialDeviceID ? onEdit : onCreate}>Save</Button>
      </DialogActions>
    </Dialog>
  );
}
