/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import AddCircleOutline from '@material-ui/icons/AddCircleOutline';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormControl from '@material-ui/core/FormControl';
import FormGroup from '@material-ui/core/FormGroup';
import IconButton from '@material-ui/core/IconButton';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import RemoveCircleOutline from '@material-ui/icons/RemoveCircleOutline';
import Select from '@material-ui/core/Select';
import TextField from '@material-ui/core/TextField';
import Typography from '@material-ui/core/Typography';

import {DEFAULT_DEVMAND_GATEWAY_CONFIGS} from './DevicesUtils';
import {createDevice} from '@fbcnms/magmalte/app/common/MagmaAPI';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';

const useStyles = makeStyles({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
  title: {
    margin: '15px 0 5px',
  },
  icon: {
    padding: '5px',
  },
});

type Props = {
  onSave: any => void,
  onClose: () => void,
};

export default function DevicesDeviceDialog(props: Props) {
  const classes = useStyles();
  const {match} = useRouter();
  const [hardwareId, setHardwareId] = useState('');
  const [name, setName] = useState('');
  const [deviceConfigs, setDeviceConfigs] = useState<Config[]>([]);

  const onSave = async () => {
    const deviceId = name.replace(/[^a-zA-z0-9]/g, '_').toLowerCase();
    try {
      const device = await createDevice(
        deviceId,
        {
          hardware_id: hardwareId,
          key: {key_type: 'ECHO'},
        },
        'devmand', // type
        DEFAULT_DEVMAND_GATEWAY_CONFIGS,
        {deviceConfigs},
        match,
      );
      props.onSave(device);
    } catch (e) {
      props.onClose();
    }
  };

  return (
    <Dialog open={true} onClose={props.onClose}>
      <DialogTitle>New Controller</DialogTitle>
      <DialogContent>
        <FormGroup>
          <TextField
            required
            className={classes.input}
            label="Controller Hardware ID"
            margin="normal"
            value={hardwareId}
            onChange={event => setHardwareId(event.target.value)}
          />
          <TextField
            required
            className={classes.input}
            label="Controller Name"
            margin="normal"
            value={name}
            onChange={event => setName(event.target.value)}
          />
          <Typography className={classes.title} variant="h6">
            Devices
            <IconButton
              onClick={() =>
                setDeviceConfigs([
                  ...deviceConfigs,
                  {ip: '', name: '', platform: '', channel: ''},
                ])
              }
              className={classes.icon}>
              <AddCircleOutline />
            </IconButton>
          </Typography>
          {deviceConfigs.map((config, index) => (
            <DeviceConfig
              key={index}
              index={index}
              config={config}
              onChange={newConfig => {
                const newDeviceConfigs = [...deviceConfigs];
                newDeviceConfigs[index] = newConfig;
                setDeviceConfigs(newDeviceConfigs);
              }}
              onRemove={() => {
                const newDeviceConfigs = [...deviceConfigs];
                newDeviceConfigs.splice(index, 1);
                setDeviceConfigs(newDeviceConfigs);
              }}
            />
          ))}
        </FormGroup>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose} color="primary">
          Cancel
        </Button>
        <Button onClick={onSave} color="primary" variant="contained">
          Save
        </Button>
      </DialogActions>
    </Dialog>
  );
}

type Config = {ip: string, name: string, platform: string, channel: string};

const DeviceConfig = (props: {
  config: Config,
  index: number,
  onChange: Config => void,
  onRemove: () => void,
}) => {
  const classes = useStyles();
  const onChange = (field: $Keys<Config>) => event => {
    const newConfig = {
      ...props.config,
      [field]: event.target.value,
    };
    props.onChange(newConfig);
  };

  return (
    <FormGroup>
      <Typography variant="subtitle1">
        Device {props.index + 1}
        <IconButton onClick={props.onRemove} className={classes.icon}>
          <RemoveCircleOutline />
        </IconButton>
      </Typography>
      <TextField
        required
        className={classes.input}
        label="Device IP"
        margin="normal"
        value={props.config.ip}
        onChange={onChange('ip')}
      />
      <TextField
        required
        className={classes.input}
        label="Device Name"
        margin="normal"
        value={props.config.name}
        onChange={onChange('name')}
      />
      <FormControl className={classes.input}>
        <InputLabel htmlFor="platform">Platform</InputLabel>
        <Select
          value={props.config.platform}
          onChange={onChange('platform')}
          inputProps={{id: 'platform'}}>
          <MenuItem value="MikroTik">MikroTik Hex</MenuItem>
          <MenuItem value="Ubiquiti M5">Ubiquiti M5</MenuItem>
          <MenuItem value="Ubiquiti Switch">Ubiquiti Switch</MenuItem>
        </Select>
      </FormControl>
      <FormControl className={classes.input}>
        <InputLabel htmlFor="channel">Channel</InputLabel>
        <Select
          value={props.config.channel}
          onChange={onChange('channel')}
          inputProps={{id: 'channel'}}>
          <MenuItem value="snmp">SNMP</MenuItem>
          <MenuItem value="netconf">NetConf</MenuItem>
        </Select>
      </FormControl>
    </FormGroup>
  );
};
