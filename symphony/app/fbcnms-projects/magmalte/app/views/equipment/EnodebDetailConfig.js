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
import type {network_ran_configs} from '@fbcnms/magma-api';

import AddEditEnodeButton from './EnodebDetailConfigEdit';
import Divider from '@material-ui/core/Divider';
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
import Text from '@fbcnms/ui/components/design-system/Text';
import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';

import {EnodeConfigFdd} from './EnodebDetailConfigFdd';
import {EnodeConfigTdd} from './EnodebDetailConfigTdd';
import {colors} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(3),
    flexGrow: 1,
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
  enbInfo: EnodebInfo,
  lteRanConfigs?: ?network_ran_configs,
};

export default function EnodebConfig(props: Props) {
  const classes = useStyles();
  const [enbInfo, setEnbInfo] = useState<EnodebInfo>(props.enbInfo);
  const {match} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);

  const {response: lteRanConfigs, isLoading} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdCellularRan,
    {
      networkId: networkId,
    },
  );

  const editProps = {
    enb: enbInfo.enb,
    lteRanConfigs: lteRanConfigs,
    onSave: enb =>
      setEnbInfo({
        ...enbInfo,
        enb: enb,
      }),
  };

  if (isLoading) {
    return <LoadingFiller />;
  }

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
                <AddEditEnodeButton
                  title={'Edit'}
                  isLink={true}
                  editProps={{
                    editTable: 'config',
                    ...editProps,
                  }}
                />
              </Grid>
            </Grid>
            <EnodebInfoConfig enbInfo={enbInfo} />
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
                <AddEditEnodeButton
                  title={'Edit'}
                  isLink={true}
                  editProps={{
                    editTable: 'ran',
                    ...editProps,
                  }}
                />{' '}
              </Grid>
            </Grid>
            <EnodebRanConfig lteRanConfigs={lteRanConfigs} enbInfo={enbInfo} />
          </Grid>
        </Grid>
      </Grid>
    </div>
  );
}

function EnodebRanConfig(props: Props) {
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
    <List component={Paper} data-testid="ran">
      <ListItem>
        <ListItemText
          primary="Bandwidth"
          secondary={props.enbInfo.enb.config.bandwidth_mhz}
          {...typographyProps}
        />
      </ListItem>
      <Divider />
      <ListItem>
        <ListItemText
          secondary={props.enbInfo.enb.config.cell_id}
          primary="Cell ID"
          {...typographyProps}
        />
      </ListItem>
      <Divider />
      {props.lteRanConfigs?.tdd_config && (
        <EnodeConfigTdd
          earfcndl={props.enbInfo.enb.config.earfcndl ?? 0}
          specialSubframePattern={
            props.enbInfo.enb.config.special_subframe_pattern ?? 0
          }
          subframeAssignment={props.enbInfo.enb.config.subframe_assignment ?? 0}
        />
      )}
      {props.lteRanConfigs?.fdd_config && (
        <EnodeConfigFdd
          earfcndl={props.enbInfo.enb.config.earfcndl ?? 0}
          earfcnul={props.lteRanConfigs.fdd_config.earfcnul}
        />
      )}
      <Divider />
      <ListItem>
        <ListItemText
          secondary={props.enbInfo.enb.config.pci}
          primary="PCI"
          {...typographyProps}
        />
      </ListItem>
      <Divider />
      <ListItem>
        <ListItemText
          secondary={props.enbInfo.enb.config.tac}
          primary="TAC"
          {...typographyProps}
        />
      </ListItem>
      <Divider />
      <ListItem>
        <ListItemText
          secondary={
            props.enbInfo.enb.config.transmit_enabled ? 'Enabled' : 'Disabled'
          }
          primary="Transmit"
          {...typographyProps}
        />
      </ListItem>
    </List>
  );
}

function EnodebInfoConfig(props: Props) {
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
    <List component={Paper} data-testid="config">
      <ListItem>
        <ListItemText
          primary="Name"
          secondary={props.enbInfo.enb.name}
          {...typographyProps}
        />
      </ListItem>
      <Divider />
      <ListItem>
        <ListItemText
          primary="Serial Number"
          secondary={props.enbInfo.enb.serial}
          {...typographyProps}
        />
      </ListItem>
      <Divider />
      <ListItem>
        <ListItemText
          primary="Description"
          secondary={props.enbInfo.enb.description}
          {...typographyProps}
        />
      </ListItem>
    </List>
  );
}
