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
import type {enodeb, network_ran_configs} from '@fbcnms/magma-api';

import AddEditEnodeButton from './EnodebDetailConfigEdit';
import Button from '@material-ui/core/Button';
import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';
import JsonEditor from '../../components/JsonEditor';
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
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
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
  appBarBtn: {
    color: colors.primary.white,
    background: colors.primary.comet,
    fontFamily: typography.button.fontFamily,
    fontWeight: typography.button.fontWeight,
    fontSize: typography.button.fontSize,
    lineHeight: typography.button.lineHeight,
    letterSpacing: typography.button.letterSpacing,

    '&:hover': {
      background: colors.primary.mirage,
    },
  },
}));

type Props = {
  enbInfo: EnodebInfo,
  lteRanConfigs?: ?network_ran_configs,
  onSave?: enodeb => void,
};

export function EnodebJsonConfig(props: Props) {
  const {match} = useRouter();
  const [error, setError] = useState('');
  const networkId: string = nullthrows(match.params.networkId);
  const enqueueSnackbar = useEnqueueSnackbar();

  return (
    <JsonEditor
      content={props.enbInfo.enb}
      error={error}
      onSave={async enb => {
        try {
          await MagmaV1API.putLteByNetworkIdEnodebsByEnodebSerial({
            networkId: networkId,
            enodebSerial: props.enbInfo.enb.serial,
            enodeb: (enb: enodeb),
          });
          enqueueSnackbar('eNodeb saved successfully', {
            variant: 'success',
          });
          setError('');
          props.onSave?.(enb);
        } catch (e) {
          setError(e.response?.data?.message ?? e.message);
        }
      }}
    />
  );
}

export default function EnodebConfig(props: Props) {
  const classes = useStyles();
  const {enbInfo, onSave} = props;
  const {history, relativeUrl, match} = useRouter();
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
    onSave: onSave,
  };

  if (isLoading) {
    return <LoadingFiller />;
  }

  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={3} alignItems="stretch">
        <Grid container spacing={3} alignItems="stretch" item xs={12}>
          <Grid container item xs={12}>
            <Grid item xs={6}>
              <Text>
                <SettingsIcon /> Config
              </Text>
            </Grid>
            <Grid container item xs={6} justify="flex-end">
              <Text>
                <Button
                  className={classes.appBarBtn}
                  onClick={() => {
                    history.push(relativeUrl('/json'));
                  }}>
                  Edit JSON
                </Button>
              </Text>
            </Grid>
          </Grid>

          <Grid item xs={6}>
            <Grid container>
              <Grid item xs={6}>
                <Text>eNodeB</Text>
              </Grid>
              <Grid container item xs={6} justify="flex-end">
                <AddEditEnodeButton
                  title={'Edit'}
                  isLink={true}
                  editProps={{
                    ...editProps,
                    editTable: 'config',
                    onSave: enb => editProps.onSave?.(enb),
                  }}
                />
              </Grid>
            </Grid>
            <EnodebInfoConfig enbInfo={enbInfo} />
          </Grid>

          <Grid item xs={6}>
            <Grid container>
              <Grid item xs={6}>
                <Text>RAN</Text>
              </Grid>
              <Grid container item xs={6} justify="flex-end">
                <AddEditEnodeButton
                  title={'Edit'}
                  isLink={true}
                  editProps={{
                    ...editProps,
                    editTable: 'ran',
                    onSave: enb => editProps.onSave?.(enb),
                  }}
                />
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
