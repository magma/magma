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
import type {EnodebInfo} from '../../components/lte/EnodebUtils';

import Collapse from '@material-ui/core/Collapse';
import Divider from '@material-ui/core/Divider';
import ExpandLess from '@material-ui/icons/ExpandLess';
import ExpandMore from '@material-ui/icons/ExpandMore';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import GraphicEqIcon from '@material-ui/icons/GraphicEq';
import Grid from '@material-ui/core/Grid';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import SettingsIcon from '@material-ui/icons/Settings';
import Switch from '@material-ui/core/Switch';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextField from '@material-ui/core/TextField';
import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';

import {makeStyles} from '@material-ui/styles';
import {useCallback, useState} from 'react';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(3),
    flexGrow: 1,
  },
  topBar: {
    backgroundColor: theme.palette.magmalte.background,
    padding: '20px 40px 20px 40px',
  },
  tabBar: {
    backgroundColor: theme.palette.magmalte.appbar,
    padding: '0 0 0 20px',
  },
  tabs: {
    color: 'white',
  },
  tab: {
    fontSize: '18px',
    textTransform: 'none',
  },
  tabLabel: {
    padding: '20px 0 20px 0',
  },
  tabIconLabel: {
    verticalAlign: 'middle',
    margin: '0 5px 3px 0',
  },
  // TODO: remove this when we actually fill out the grid sections
  contentPlaceholder: {
    padding: '50px 0',
  },
  paper: {
    height: 100,
    padding: theme.spacing(10),
    textAlign: 'center',
    color: theme.palette.text.secondary,
  },
  card: {
    variant: 'outlined',
  },
  root: {
    '& > *': {
      margin: theme.spacing(1),
      width: '25ch',
    },
  },
}));

export default function EnodebConfig({enbInfo}: {enbInfo: EnodebInfo}) {
  const classes = useStyles();
  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={3} alignItems="stretch">
        <Grid container spacing={3} alignItems="stretch" item xs={12}>
          <Grid item xs={6}>
            <Grid container>
              <Grid item xs={6}>
                <Text>
                  <SettingsIcon /> Config
                </Text>
              </Grid>
              <Grid container item xs={6} justify="flex-end">
                <Text>Edit</Text>
              </Grid>
            </Grid>
            <EnodebInfoConfig readOnly={true} enbInfo={enbInfo} />
          </Grid>

          <Grid item xs={6}>
            <Grid container>
              <Grid item xs={6}>
                <Text>
                  <GraphicEqIcon />
                  RAN
                </Text>
              </Grid>
              <Grid container item xs={6} justify="flex-end">
                <Text>Edit</Text>
              </Grid>
            </Grid>
            <EnodebRanConfig readOnly={true} enbInfo={enbInfo} />
          </Grid>
        </Grid>
      </Grid>
    </div>
  );
}

function EnodebRanConfig({
  enbInfo,
  readOnly,
}: {
  enbInfo: EnodebInfo,
  readOnly: boolean,
}) {
  const [open, setOpen] = React.useState(true);
  const [bandwidth, setBandwidth] = useState(enbInfo.enb.config.bandwidth_mhz);
  const [cellID, setCellID] = useState(enbInfo.enb.config.cell_id);
  const [pci, setPci] = useState(enbInfo.enb.config.pci);
  const [specialSubframePattern, setSpecialSubframePattern] = useState(
    enbInfo.enb.config.special_subframe_pattern,
  );
  const [subframeAssignment, setSubFrameAssignment] = useState(
    enbInfo.enb.config.special_subframe_pattern,
  );
  const [tac, setTac] = useState(enbInfo.enb.config.tac);
  const [transmit, setTransmit] = useState(enbInfo.enb.config.transmit_enabled);
  const [earfcndl, setEarFcnDl] = useState(0);
  const [earfcnul, setEarFcnUl] = useState(0);

  const {match} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);

  const {response: lteRanConfigs, isLoading} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdCellularRan,
    {
      networkId: networkId,
    },
    useCallback(lteRanConfigs => {
      if (lteRanConfigs && lteRanConfigs?.tdd_config) {
        setEarFcnDl(lteRanConfigs?.tdd_config?.earfcndl);
      }
      if (lteRanConfigs && lteRanConfigs?.fdd_config) {
        setEarFcnUl(lteRanConfigs?.fdd_config?.earfcnul);
        setEarFcnDl(lteRanConfigs?.fdd_config?.earfcndl);
      }
    }, []),
  );

  if (isLoading) {
    return <LoadingFiller />;
  }

  const tddConfigJSX = [];
  const fddConfigJSX = [];

  if (lteRanConfigs && lteRanConfigs?.tdd_config) {
    tddConfigJSX.push(
      <List key="tddConfigs">
        <ListItem button onClick={() => setOpen(!open)}>
          <ListItemText primary="TDD" />
          {open ? <ExpandLess /> : <ExpandMore />}
        </ListItem>
        <Collapse key="tdd" in={open} timeout="auto" unmountOnExit>
          <ListItem>
            <TextField
              fullWidth={true}
              value={earfcndl}
              label="EARFCNDL"
              onChange={({target}) => setEarFcnDl(target.value)}
              InputProps={{disableUnderline: true, readOnly: readOnly}}
            />
          </ListItem>
          <ListItem>
            <TextField
              fullWidth={true}
              value={specialSubframePattern}
              label="Special Subframe Pattern"
              onChange={({target}) => setSpecialSubframePattern(target.value)}
              InputProps={{disableUnderline: true, readOnly: readOnly}}
            />
          </ListItem>
          <ListItem>
            <TextField
              fullWidth={true}
              value={subframeAssignment}
              label="Subframe Assignment"
              onChange={({target}) => setSubFrameAssignment(target.value)}
              InputProps={{disableUnderline: true, readOnly: readOnly}}
            />
          </ListItem>
        </Collapse>
      </List>,
    );
  }

  if (lteRanConfigs && lteRanConfigs?.fdd_config) {
    fddConfigJSX.push(
      <List key="fddConfigs">
        <ListItem button onClick={() => setOpen(!open)}>
          <ListItemText primary="FDD" />
          {open ? <ExpandLess /> : <ExpandMore />}
        </ListItem>
        <Divider />
        <Collapse key="fdd" in={open} timeout="auto" unmountOnExit>
          <ListItem>
            <Grid container>
              <Grid item xs={6}>
                <TextField
                  fullWidth={true}
                  value={earfcndl}
                  label="EARFCNDL"
                  onChange={({target}) => setEarFcnDl(target.value)}
                  InputProps={{disableUnderline: true, readOnly: readOnly}}
                />
              </Grid>
              <Grid item xs={6}>
                <TextField
                  fullWidth={true}
                  value={earfcnul}
                  label="EARFCNUL"
                  onChange={({target}) => setEarFcnUl(target.value)}
                  InputProps={{disableUnderline: true, readOnly: readOnly}}
                />
              </Grid>
            </Grid>
          </ListItem>
        </Collapse>
      </List>,
    );
  }
  return (
    <List component={Paper}>
      <ListItem>
        <TextField
          fullWidth={true}
          value={bandwidth}
          label="Bandwidth"
          onChange={({target}) => setBandwidth(target.value)}
          InputProps={{disableUnderline: true, readOnly: readOnly}}
        />
      </ListItem>
      <Divider />
      <ListItem>
        <TextField
          fullWidth={true}
          value={cellID}
          label="Cell ID"
          onChange={({target}) => setCellID(target.value)}
          InputProps={{disableUnderline: true, readOnly: readOnly}}
        />
      </ListItem>
      <Divider />
      {fddConfigJSX}
      {tddConfigJSX}
      <Divider />
      <ListItem>
        <TextField
          fullWidth={true}
          value={pci}
          label="PCI"
          onChange={({target}) => setPci(target.value)}
          InputProps={{disableUnderline: true, readOnly: readOnly}}
        />
      </ListItem>
      <Divider />
      <ListItem>
        <TextField
          fullWidth={true}
          value={tac}
          label="TAC"
          onChange={({target}) => setTac(target.value)}
          InputProps={{disableUnderline: true, readOnly: readOnly}}
        />
      </ListItem>
      <Divider />
      <ListItem>
        <FormControlLabel
          label="Transmit"
          control={
            <Switch
              disabled={readOnly}
              checked={transmit}
              onChange={({target}) => setTransmit(target.checked)}
              color="primary"
            />
          }
        />
      </ListItem>
    </List>
  );
}

function EnodebInfoConfig({
  enbInfo,
  readOnly,
}: {
  enbInfo: EnodebInfo,
  readOnly: boolean,
}) {
  const [name, setName] = useState(enbInfo.enb.name);
  const [serial, setSerial] = useState(enbInfo.enb.serial);
  const [description, setDescription] = useState(enbInfo.enb.description);
  return (
    <List component={Paper}>
      <ListItem>
        <TextField
          fullWidth={true}
          value={name}
          label="Name"
          onChange={({target}) => setName(target.value)}
          InputProps={{disableUnderline: true, readOnly: readOnly}}
        />
      </ListItem>
      <Divider />
      <ListItem>
        <TextField
          fullWidth={true}
          value={serial}
          label="Serial Number"
          onChange={({target}) => setSerial(target.value)}
          InputProps={{disableUnderline: true, readOnly: readOnly}}
        />
      </ListItem>
      <Divider />
      <ListItem>
        <TextField
          fullWidth={true}
          value={description}
          label="Description"
          onChange={({target}) => setDescription(target.value)}
          InputProps={{disableUnderline: true, readOnly: readOnly}}
        />
      </ListItem>
    </List>
  );
}
