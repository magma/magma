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
import type {network_ran_configs} from '@fbcnms/magma-api';

import Button from '@fbcnms/ui/components/design-system/Button';
import Collapse from '@material-ui/core/Collapse';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import Divider from '@material-ui/core/Divider';
import ExpandLess from '@material-ui/icons/ExpandLess';
import ExpandMore from '@material-ui/icons/ExpandMore';
import FddConfig from './NetworkRanFddConfig';
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
import TddConfig from './NetworkRanTddConfig';

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
  lteRanConfigs: network_ran_configs,
};

export default function NetworkRan(props: Props) {
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
  const [open, setOpen] = React.useState(true);

  return (
    <List component={Paper} data-testid="ran">
      <ListItem>
        <ListItemText
          primary={'Bandwidth'}
          secondary={props.lteRanConfigs?.bandwidth_mhz}
          {...typographyProps}
        />
      </ListItem>
      <Divider />
      {props.lteRanConfigs?.tdd_config && (
        <List key="tddConfigs">
          <ListItem button onClick={() => setOpen(!open)}>
            <ListItemText
              primary="RAN Config"
              secondary="TDD"
              {...typographyProps}
            />
            {open ? <ExpandLess /> : <ExpandMore />}
          </ListItem>
          <Collapse key="tdd" in={open} timeout="auto" unmountOnExit>
            <ListItem>
              <ListItemText
                primary={'EARFCNDL'}
                secondary={props.lteRanConfigs?.tdd_config?.earfcndl}
                {...typographyProps}
              />
            </ListItem>
            <ListItem>
              <ListItemText
                primary={'Special Subframe Pattern'}
                secondary={
                  props.lteRanConfigs.tdd_config?.special_subframe_pattern
                }
                {...typographyProps}
              />
            </ListItem>
            <ListItem>
              <ListItemText
                primary={'Subframe Assignment'}
                secondary={props.lteRanConfigs?.tdd_config?.subframe_assignment}
                {...typographyProps}
              />
            </ListItem>
          </Collapse>
        </List>
      )}
      {props.lteRanConfigs?.fdd_config && (
        <List key="fddConfigs">
          <ListItem button onClick={() => setOpen(!open)}>
            <ListItemText
              primary="RAN Config"
              secondary="FDD"
              {...typographyProps}
            />
            {open ? <ExpandLess /> : <ExpandMore />}
          </ListItem>
          <Divider />
          <Collapse key="fdd" in={open} timeout="auto" unmountOnExit>
            <ListItem>
              <Grid container>
                <Grid item xs={6}>
                  <ListItemText
                    primary={'EARFCNDL'}
                    secondary={props.lteRanConfigs?.fdd_config?.earfcndl}
                    {...typographyProps}
                  />
                </Grid>
                <Grid item xs={6}>
                  <ListItemText
                    primary={'EARFCNUL'}
                    secondary={props.lteRanConfigs?.fdd_config?.earfcnul}
                    {...typographyProps}
                  />
                </Grid>
              </Grid>
            </ListItem>
          </Collapse>
        </List>
      )}
    </List>
  );
}

type EditProps = {
  saveButtonTitle: string,
  networkId: string,
  lteRanConfigs: ?network_ran_configs,
  onClose: () => void,
  onSave: network_ran_configs => void,
};
type BandType = 'tdd' | 'fdd';
const ValidBandwidths = [3, 5, 10, 15, 20];

export function NetworkRanEdit(props: EditProps) {
  const classes = useStyles();
  const enqueueSnackbar = useEnqueueSnackbar();
  const [error, setError] = useState('');
  const [bandType, setBandType] = useState<BandType>('tdd');
  const defaultTddConfig = {
    earfcndl: 0,
    special_subframe_pattern: 0,
    subframe_assignment: 0,
  };
  const defaulFddConfig = {
    earfcndl: 0,
    earfcnul: 0,
  };
  const [lteRanConfigs, setLteRanConfigs] = useState(
    props?.lteRanConfigs || {
      bandwidth_mhz: 20,
      fdd_config: undefined,
      tdd_config: defaultTddConfig,
    },
  );

  const onSave = async () => {
    const config: network_ran_configs = {
      ...lteRanConfigs,
    };
    if (bandType === 'tdd') {
      config.fdd_config = undefined;
    } else {
      config.tdd_config = undefined;
    }
    try {
      await MagmaV1API.putLteByNetworkIdCellularRan({
        networkId: props.networkId,
        config: config,
      });
      enqueueSnackbar('RAN configs saved successfully', {variant: 'success'});
      props.onSave(config);
    } catch (e) {
      setError(e.response.data?.message ?? e.message);
    }
  };

  return (
    <>
      <DialogContent data-testid="networkRanEdit">
        {error !== '' && <FormLabel error>{error}</FormLabel>}
        <List component={Paper}>
          <ListItem>
            <Grid container>
              <Grid item xs={12}>
                Bandwidth
              </Grid>
              <Grid item xs={12}>
                <FormControl className={classes.input}>
                  <Select
                    value={lteRanConfigs.bandwidth_mhz}
                    onChange={({target}) => {
                      if (
                        target.value === 3 ||
                        target.value === 5 ||
                        target.value === 10 ||
                        target.value === 15 ||
                        target.value === 20
                      ) {
                        setLteRanConfigs({
                          ...lteRanConfigs,
                          bandwidth_mhz: target.value,
                        });
                      }
                    }}
                    input={<OutlinedInput id="bandwidth" />}>
                    {ValidBandwidths.map((k: number, idx: number) => (
                      <MenuItem key={idx} value={k}>
                        {k}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </Grid>
            </Grid>
          </ListItem>
          <ListItem>
            <Grid container>
              <Grid item xs={12}>
                Band Type
              </Grid>
              <Grid item xs={12}>
                <FormControl className={classes.input}>
                  <Select
                    value={bandType}
                    onChange={({target}) => {
                      if (target.value === 'fdd') {
                        setLteRanConfigs({
                          fdd_config: defaulFddConfig,
                          ...lteRanConfigs,
                        });
                        setBandType('fdd');
                      } else {
                        setLteRanConfigs({
                          tdd_config: defaultTddConfig,
                          ...lteRanConfigs,
                        });
                        setBandType(target.value === 'tdd' ? 'tdd' : 'fdd');
                      }
                    }}
                    input={<Input id="bandType" />}>
                    <MenuItem value={'tdd'}>TDD</MenuItem>
                    <MenuItem value={'fdd'}>FDD</MenuItem>
                  </Select>
                </FormControl>
              </Grid>
            </Grid>
          </ListItem>
          {bandType === 'tdd' && (
            <TddConfig
              lteRanConfigs={lteRanConfigs}
              setLteRanConfigs={setLteRanConfigs}
            />
          )}
          {bandType === 'fdd' && (
            <FddConfig
              lteRanConfigs={lteRanConfigs}
              setLteRanConfigs={setLteRanConfigs}
            />
          )}
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
