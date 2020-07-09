/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import type {network_epc_configs} from '@fbcnms/magma-api';

import Button from '@fbcnms/ui/components/design-system/Button';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import Divider from '@material-ui/core/Divider';
import FormControl from '@material-ui/core/FormControl';
import FormLabel from '@material-ui/core/FormLabel';
import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import Input from '@material-ui/core/Input';
import InputAdornment from '@material-ui/core/InputAdornment';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import MenuItem from '@material-ui/core/MenuItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Select from '@material-ui/core/Select';
import Text from '@fbcnms/ui/components/design-system/Text';
import Visibility from '@material-ui/icons/Visibility';
import VisibilityOff from '@material-ui/icons/VisibilityOff';

import {colors} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useState} from 'react';

const useStyles = makeStyles(() => ({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '50%',
    fullWidth: true,
  },
  itemTitle: {
    color: colors.primary.comet,
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
  },
  itemValue: {
    color: colors.primary.brightGray,
  },
}));

type Props = {
  epcConfigs: network_epc_configs,
};

export default function NetworkEpc(props: Props) {
  const classes = useStyles();
  const typographyProps = {
    primaryTypographyProps: {
      variant: 'caption',
      className: classes.itemTitle,
    },
    secondaryTypographyProps: {
      variant: 'h6',
      className: classes.itemValue,
    },
  };
  const [showPassword, setShowPassword] = React.useState(false);
  return (
    <List component={Paper} data-testid="epc">
      <ListItem>
        <ListItemText
          primary={'Policy Enforcement Enabled'}
          secondary={props.epcConfigs.relay_enabled ? 'Enabled' : 'Disabled'}
          {...typographyProps}
        />
      </ListItem>
      <Divider />
      <ListItem>
        <ListItemText
          disableTypography={true}
          primary={<Text variant="caption">LTE Auth AMF</Text>}
          secondary={
            <Input
              type={showPassword ? 'text' : 'password'}
              data-testid="epcPassword"
              fullWidth={true}
              value={props.epcConfigs.lte_auth_amf}
              disableUnderline={true}
              readOnly={true}
              endAdornment={
                <InputAdornment position="end">
                  <IconButton
                    aria-label="toggle password visibility"
                    onClick={() => setShowPassword(!showPassword)}
                    onMouseDown={event => event.preventDefault()}>
                    {showPassword ? <Visibility /> : <VisibilityOff />}
                  </IconButton>
                </InputAdornment>
              }
            />
          }
        />
      </ListItem>
      <Divider />
      <ListItem>
        <ListItemText
          primary={'MCC'}
          secondary={props.epcConfigs.mcc}
          {...typographyProps}
        />
      </ListItem>
      <Divider />
      <ListItem>
        <ListItemText
          primary={'MNC'}
          secondary={props.epcConfigs.mnc}
          {...typographyProps}
        />
      </ListItem>
      <Divider />
      <ListItem>
        <ListItemText
          primary={'TAC'}
          secondary={props.epcConfigs.tac}
          {...typographyProps}
        />
      </ListItem>
    </List>
  );
}

type EditProps = {
  saveButtonTitle: string,
  networkId: string,
  epcConfigs: ?network_epc_configs,
  onClose: () => void,
  onSave: network_epc_configs => void,
};

export function NetworkEpcEdit(props: EditProps) {
  const classes = useStyles();
  const [showPassword, setShowPassword] = React.useState(false);
  const enqueueSnackbar = useEnqueueSnackbar();
  const [error, setError] = useState('');
  const [epcConfigs, setEpcConfigs] = useState<network_epc_configs>(
    props.epcConfigs || {
      cloud_subscriberdb_enabled: false,
      default_rule_id: 'default_rule_1',
      lte_auth_amf: 'gAA=',
      lte_auth_op: 'EREREREREREREREREREREQ==',
      mcc: '001',
      mnc: '01',
      network_services: ['policy_enforcement'],
      relay_enabled: false,
      sub_profiles: {},
      tac: 1,
    },
  );

  const onSave = async () => {
    try {
      MagmaV1API.putLteByNetworkIdCellularEpc({
        networkId: props.networkId,
        config: epcConfigs,
      });
      enqueueSnackbar('EPC configs saved successfully', {variant: 'success'});
      props.onSave(epcConfigs);
    } catch (e) {
      setError(e.data?.message ?? e.message);
    }
  };

  return (
    <>
      <DialogContent data-testid="networkEpcEdit">
        {error !== '' && <FormLabel error>{error}</FormLabel>}
        <List component={Paper}>
          <ListItem>
            <Grid container>
              <Grid item xs={12}>
                Policy Enforcement Enabled
              </Grid>
              <Grid item xs={12}>
                <FormControl variant={'outlined'} className={classes.input}>
                  <Select
                    value={epcConfigs.relay_enabled ? 1 : 0}
                    onChange={({target}) => {
                      setEpcConfigs({
                        ...epcConfigs,
                        relay_enabled: target.value === 1,
                      });
                    }}
                    input={<OutlinedInput id="relayEnabled" />}>
                    <MenuItem value={0}>Disabled</MenuItem>
                    <MenuItem value={1}>Enabled</MenuItem>
                  </Select>
                </FormControl>
              </Grid>
            </Grid>
          </ListItem>
          <ListItem>
            <Grid container>
              <Grid item xs={12}>
                LTE Auth AMF
              </Grid>
              <Grid item xs={12}>
                <OutlinedInput
                  data-testid="password"
                  className={classes.input}
                  type={showPassword ? 'text' : 'password'}
                  fullWidth={true}
                  value={epcConfigs.lte_auth_amf}
                  onChange={({target}) => {
                    setEpcConfigs({...epcConfigs, lte_auth_amf: target.value});
                  }}
                  endAdornment={
                    <InputAdornment position="end">
                      <IconButton
                        aria-label="toggle password visibility"
                        onClick={() => setShowPassword(!showPassword)}
                        onMouseDown={event => event.preventDefault()}>
                        {showPassword ? <Visibility /> : <VisibilityOff />}
                      </IconButton>
                    </InputAdornment>
                  }
                />
              </Grid>
            </Grid>
          </ListItem>
          <ListItem>
            <Grid container>
              <Grid item xs={12}>
                MCC
              </Grid>
              <Grid item xs={12}>
                <OutlinedInput
                  data-testid="mcc"
                  className={classes.input}
                  fullWidth={true}
                  value={epcConfigs.mcc}
                  onChange={({target}) =>
                    setEpcConfigs({...epcConfigs, mcc: target.value})
                  }
                />
              </Grid>
            </Grid>
          </ListItem>
          <ListItem>
            <Grid container>
              <Grid item xs={12}>
                MNC
              </Grid>
              <Grid item xs={12}>
                <OutlinedInput
                  data-testid="mnc"
                  className={classes.input}
                  fullWidth={true}
                  value={epcConfigs.mnc}
                  onChange={({target}) =>
                    setEpcConfigs({...epcConfigs, mnc: target.value})
                  }
                />
              </Grid>
            </Grid>
          </ListItem>
          <ListItem>
            <Grid container>
              <Grid item xs={12}>
                TAC
              </Grid>
              <Grid item xs={12}>
                <OutlinedInput
                  data-testid="tac"
                  className={classes.input}
                  type="number"
                  fullWidth={true}
                  value={epcConfigs.tac}
                  onChange={({target}) =>
                    setEpcConfigs({...epcConfigs, tac: parseInt(target.value)})
                  }
                />
              </Grid>
            </Grid>
          </ListItem>
        </List>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose} skin="regular">
          Cancel
        </Button>
        <Button onClick={onSave}>{props.saveButtonTitle}</Button>
      </DialogActions>
    </>
  );
}
