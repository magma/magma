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
import type {network, network_dns_config} from '@fbcnms/magma-api';

import Button from '@fbcnms/ui/components/design-system/Button';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import Divider from '@material-ui/core/Divider';
import FormControl from '@material-ui/core/FormControl';
import FormLabel from '@material-ui/core/FormLabel';
import Grid from '@material-ui/core/Grid';
import Input from '@material-ui/core/Input';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import MenuItem from '@material-ui/core/MenuItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Select from '@material-ui/core/Select';
import axios from 'axios';

import {AllNetworkTypes} from '@fbcnms/types/network';
import {CWF, FEG, LTE} from '@fbcnms/types/network';
import {colors} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useState} from 'react';

const useStyles = makeStyles(_ => ({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '50%',
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
  networkInfo: network,
};

export default function NetworkInfo(props: Props) {
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
  return (
    <List component={Paper} data-testid="info">
      <ListItem>
        <ListItemText
          primary="ID"
          secondary={props.networkInfo.id}
          {...typographyProps}
        />
      </ListItem>
      <Divider />
      <ListItem>
        <ListItemText
          primary="Name"
          secondary={props.networkInfo.name}
          {...typographyProps}
        />
      </ListItem>
      <Divider />
      <ListItem>
        <ListItemText
          primary="Network Type"
          secondary={props.networkInfo.type}
          {...typographyProps}
        />
      </ListItem>
      <Divider />
      <ListItem>
        <ListItemText
          primary="Description"
          secondary={props.networkInfo.description}
          {...typographyProps}
        />
      </ListItem>
    </List>
  );
}

type EditProps = {
  saveButtonTitle: string,
  networkInfo: ?network,
  onClose: () => void,
  onSave: network => void,
};

const DEFAULT_DNS_CONFIG: network_dns_config = {
  enable_caching: false,
  local_ttl: 0,
  records: [],
};

export function NetworkInfoEdit(props: EditProps) {
  const classes = useStyles();
  const [error, setError] = useState('');
  const [fegNetworkID, setFegNetworkID] = useState('');
  const [servedNetworkIDs, setServedNetworkIDs] = useState('');
  const enqueueSnackbar = useEnqueueSnackbar();
  const [networkType, setNetworkType] = useState(
    props.networkInfo?.type || LTE,
  );
  const [networkInfo, setNetworkInfo] = useState<network>(
    props.networkInfo || {
      name: '',
      id: '',
      description: '',
      dns: DEFAULT_DNS_CONFIG,
    },
  );

  const onSave = async () => {
    const payload = {
      networkID: networkInfo.id,
      data: {
        name: networkInfo.name,
        description: networkInfo.description,
        networkType,
        fegNetworkID,
        servedNetworkIDs,
      },
    };
    if (props.networkInfo) {
      // edit
      try {
        await MagmaV1API.putNetworksByNetworkId({
          networkId: networkInfo.id,
          network: {
            ...networkInfo,
            type: networkType,
          },
        });
        enqueueSnackbar('Network configs saved successfully', {
          variant: 'success',
        });
        props.onSave(networkInfo);
      } catch (e) {
        setError(e.data?.message ?? e.message);
      }
    } else {
      try {
        const response = await axios.post('/nms/network/create', payload);
        if (response.data.success) {
          enqueueSnackbar(`Network $networkInfo.name} successfully created`, {
            variant: 'success',
          });
          props.onSave(networkInfo);
        } else {
          setError(response.data.message);
        }
      } catch (e) {
        setError(e.data?.message ?? e.message);
      }
    }
  };

  return (
    <>
      <DialogContent data-testid="networkInfoEdit">
        {error !== '' && <FormLabel error>{error}</FormLabel>}
        <List component={Paper}>
          <ListItem>
            <Grid container>
              <Grid item xs={12}>
                Network ID
              </Grid>
              <Grid item xs={12}>
                <OutlinedInput
                  data-testid="networkID"
                  fullWidth={true}
                  className={classes.input}
                  value={networkInfo.id}
                  onChange={({target}) =>
                    setNetworkInfo({...networkInfo, id: target.value})
                  }
                  readOnly={props.networkInfo ? true : false}
                />
              </Grid>
            </Grid>
          </ListItem>
          <ListItem>
            <Grid container>
              <Grid item xs={12}>
                Network Name
              </Grid>
              <Grid item xs={12}>
                <OutlinedInput
                  data-testid="networkName"
                  fullWidth={true}
                  className={classes.input}
                  value={networkInfo.name}
                  onChange={({target}) =>
                    setNetworkInfo({...networkInfo, name: target.value})
                  }
                />
              </Grid>
            </Grid>
          </ListItem>
          <ListItem>
            <Grid container>
              <Grid item xs={12}>
                Network Type
              </Grid>
              <Grid item xs={12}>
                <FormControl className={classes.input}>
                  <Select
                    variant={'outlined'}
                    value={networkType}
                    onChange={({target}) => {
                      setNetworkType(target.value);
                    }}
                    data-testid="networkType"
                    input={<Input id="networkType" />}>
                    {AllNetworkTypes.map(type => (
                      <MenuItem key={type} value={type}>
                        <ListItemText primary={type} />
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </Grid>
            </Grid>
          </ListItem>

          {networkType === CWF && (
            <ListItem>
              <Grid container>
                <Grid item xs={12}>
                  Federation Network ID
                </Grid>
                <Grid item xs={12}>
                  <OutlinedInput
                    className={classes.input}
                    value={fegNetworkID}
                    onChange={({target}) => setFegNetworkID(target.value)}
                  />
                </Grid>
              </Grid>
            </ListItem>
          )}
          {networkType === FEG && (
            <ListItem>
              <Grid container>
                <Grid item xs={12}>
                  Served Network IDs
                </Grid>
                <Grid item xs={12}>
                  <OutlinedInput
                    placeholder="network1,network2"
                    className={classes.input}
                    value={servedNetworkIDs}
                    onChange={({target}) => setServedNetworkIDs(target.value)}
                  />
                </Grid>
              </Grid>
            </ListItem>
          )}
          <ListItem>
            <Grid container>
              <Grid item xs={12}>
                Add Description
              </Grid>
              <Grid item xs={12}>
                <OutlinedInput
                  data-testid="networkDescription"
                  fullWidth={true}
                  multiline
                  rows={4}
                  value={networkInfo.description}
                  onChange={({target}) =>
                    setNetworkInfo({...networkInfo, description: target.value})
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
