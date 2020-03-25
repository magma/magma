/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {Organization} from './Organizations';

import Button from '@fbcnms/ui/components/design-system/Button';
import Checkbox from '@material-ui/core/Checkbox';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormControl from '@material-ui/core/FormControl';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import FormLabel from '@material-ui/core/FormLabel';
import Input from '@material-ui/core/Input';
import InputLabel from '@material-ui/core/InputLabel';
import ListItemText from '@material-ui/core/ListItemText';
import LoadingFillerBackdrop from '@fbcnms/ui/components/LoadingFillerBackdrop';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import Select from '@material-ui/core/Select';
import TextField from '@material-ui/core/TextField';
import axios from 'axios';

import renderList from '@fbcnms/util/renderList';
import {getProjectTabs} from '@fbcnms/magmalte/app/common/projects';
import {makeStyles} from '@material-ui/styles';
import {useAxios} from '@fbcnms/ui/hooks';
import {useState} from 'react';

const useStyles = makeStyles(_ => ({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
}));

type Props = {
  onClose: () => void,
  onSave: Organization => void,
};

export default function (props: Props) {
  const classes = useStyles();
  const {error, isLoading, response} = useAxios({
    method: 'get',
    url: '/master/networks/async',
  });

  const [name, setName] = useState('');
  const [networkIds, setNetworkIds] = useState(new Set());
  const [tabs, setTabs] = useState(new Set());
  const [shouldEnableAllNetworks, setShouldEnableAllNetworks] = useState(false);

  if (isLoading) {
    return <LoadingFillerBackdrop />;
  }

  const allNetworks = error || !response ? [] : response.data.sort();
  const onSave = async () => {
    const response = await axios.post('/master/organization/async', {
      name,
      networkIDs: shouldEnableAllNetworks
        ? allNetworks
        : Array.from(networkIds),
      customDomains: [], // TODO
      tabs: Array.from(tabs),
    });
    props.onSave(response.data.organization);
  };

  return (
    <Dialog open={true} onClose={props.onClose}>
      <DialogTitle>Add Organization</DialogTitle>
      <DialogContent>
        {error && <FormLabel error>{error}</FormLabel>}
        <TextField
          name="name"
          label="Name"
          className={classes.input}
          value={name}
          onChange={({target}) => setName(target.value)}
        />
        <FormControl className={classes.input}>
          <InputLabel htmlFor="tabs">Accessible Tabs</InputLabel>
          <Select
            multiple
            value={Array.from(tabs)}
            onChange={({target}) => setTabs(new Set(target.value))}
            renderValue={renderList}
            input={<Input id="tabs" />}>
            {getProjectTabs().map(tab => (
              <MenuItem key={tab.id} value={tab.id}>
                <Checkbox checked={tabs.has(tab.id)} />
                <ListItemText primary={tab.name} />
              </MenuItem>
            ))}
          </Select>
        </FormControl>
        <FormControlLabel
          control={
            <Checkbox
              checked={shouldEnableAllNetworks}
              onChange={({target}) =>
                setShouldEnableAllNetworks(target.checked)
              }
              color="primary"
            />
          }
          label="Enable All Networks"
        />
        {!shouldEnableAllNetworks && (
          <FormControl className={classes.input}>
            <InputLabel htmlFor="network_ids">Accessible Networks</InputLabel>
            <Select
              multiple
              value={Array.from(networkIds)}
              onChange={({target}) => setNetworkIds(new Set(target.value))}
              renderValue={renderList}
              input={<Input id="network_ids" />}>
              {allNetworks.map(network => (
                <MenuItem key={network} value={network}>
                  <Checkbox checked={networkIds.has(network)} />
                  <ListItemText primary={network} />
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        )}
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose} skin="regular">
          Cancel
        </Button>
        <Button onClick={onSave}>Save</Button>
      </DialogActions>
    </Dialog>
  );
}
