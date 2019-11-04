/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {
  federation_gateway,
  gateway_federation_configs,
  gx,
} from '@fbcnms/magma-api';

import AppBar from '@material-ui/core/AppBar';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import LoadingFillerBackdrop from '@fbcnms/ui/components/LoadingFillerBackdrop';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React, {useState} from 'react';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import TextField from '@material-ui/core/TextField';

import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '../../common/useMagmaAPI';
import {
  AddGatewayFields,
  EMPTY_GATEWAY_FIELDS,
  MAGMAD_DEFAULT_CONFIGS,
} from '../AddGatewayDialog';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles({
  appBar: {
    backgroundColor: '#f5f5f5',
    marginBottom: '20px',
  },
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
});

type Props = {|
  onClose: () => void,
  onSave: federation_gateway => void,
  editingGateway?: federation_gateway,
|};

type DiameterValues = {
  address: string,
  dest_host: string,
  dest_realm: string,
  host: string,
  realm: string,
  local_address: string,
  product_name: string,
};

function getDiameterConfigs(cfg: DiameterValues): gx {
  return {
    server: {...cfg},
  };
}

function getInitialDiameterConfigs(cfg: ?gx): DiameterValues {
  return {
    address: cfg?.server?.address || '',
    dest_host: cfg?.server?.dest_host || '',
    dest_realm: cfg?.server?.dest_realm || '',
    host: cfg?.server?.host || '',
    realm: cfg?.server?.realm || '',
    local_address: cfg?.server?.local_address || '',
    product_name: cfg?.server?.product_name || '',
  };
}

export default function FEGGatewayDialog(props: Props) {
  const classes = useStyles();
  const {match} = useRouter();
  const enqueueSnackbar = useEnqueueSnackbar();

  const {editingGateway} = props;
  const [tab, setTab] = useState(editingGateway ? 'gx' : 'general');
  const [generalFields, setGeneralFields] = useState(EMPTY_GATEWAY_FIELDS);
  const [gx, setGx] = useState<DiameterValues>(
    getInitialDiameterConfigs(editingGateway?.federation?.gx),
  );
  const [gy, setGy] = useState<DiameterValues>(
    getInitialDiameterConfigs(editingGateway?.federation?.gy),
  );
  const [s6a, setS6A] = useState<DiameterValues>(
    getInitialDiameterConfigs(editingGateway?.federation?.s6a),
  );

  const networkID = nullthrows(match.params.networkId);
  const {response: tiers, isLoading} = useMagmaAPI(
    MagmaV1API.getNetworksByNetworkIdTiers,
    {networkId: networkID},
  );

  if (isLoading || !tiers) {
    return <LoadingFillerBackdrop />;
  }

  const getFederationConfigs = (): gateway_federation_configs => ({
    aaa_server: {},
    eap_aka: {},
    gx: getDiameterConfigs(gx),
    gy: {server: getDiameterConfigs(gy).server},
    health: {},
    hss: {},
    s6a: getDiameterConfigs(s6a),
    served_network_ids: [],
    swx: {},
  });

  const onSave = async () => {
    try {
      if (editingGateway) {
        await MagmaV1API.putFegByNetworkIdGatewaysByGatewayIdFederation({
          networkId: networkID,
          gatewayId: editingGateway.id,
          config: getFederationConfigs(),
        });
      } else {
        await MagmaV1API.postFegByNetworkIdGateways({
          networkId: networkID,
          gateway: {
            device: {
              hardware_id: generalFields.hardwareID,
              key: {
                key: generalFields.challengeKey,
                key_type: 'ECHO',
              },
            },
            federation: getFederationConfigs(),
            magmad: MAGMAD_DEFAULT_CONFIGS,
            id: generalFields.gatewayID,
            description: generalFields.description,
            name: generalFields.name,
            tier: generalFields.tier,
          },
        });
      }

      const gateway = await MagmaV1API.getFegByNetworkIdGatewaysByGatewayId({
        networkId: networkID,
        gatewayId: editingGateway?.id || generalFields.gatewayID,
      });
      props.onSave(gateway);
    } catch (e) {
      enqueueSnackbar(e?.response?.data?.message || e?.message || e, {
        variant: 'error',
      });
    }
  };

  let content;
  switch (tab) {
    case 'general':
      content = (
        <AddGatewayFields
          onChange={setGeneralFields}
          values={generalFields}
          tiers={tiers}
        />
      );
      break;
    case 'gx':
      content = <DiameterFields onChange={setGx} values={gx} />;
      break;
    case 'gy':
      content = <DiameterFields onChange={setGy} values={gy} />;
      break;
    case 's6a':
      content = <DiameterFields onChange={setS6A} values={s6a} />;
      break;
  }

  return (
    <Dialog open={true} onClose={props.onClose} maxWidth="md" scroll="body">
      <AppBar position="static" className={classes.appBar}>
        <Tabs
          indicatorColor="primary"
          textColor="primary"
          value={tab}
          onChange={(event, tab) => setTab(tab)}>
          {!editingGateway && <Tab label="General" values="general" />}
          <Tab label="Gx" value="gx" />
          <Tab label="Gy" value="gy" />
          <Tab label="S6A" value="s6a" />
        </Tabs>
      </AppBar>
      <DialogContent>{content}</DialogContent>
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

function DiameterFields(props: {
  values: DiameterValues,
  onChange: DiameterValues => void,
}) {
  const classes = useStyles();
  const {values} = props;
  const onChange = field => event =>
    props.onChange({...values, [field]: event.target.value});

  return (
    <>
      <TextField
        label="Address"
        className={classes.input}
        value={values.address}
        onChange={onChange('address')}
        placeholder="example.magma.com:5555"
      />
      <TextField
        label="Destination Host"
        className={classes.input}
        value={values.dest_host}
        onChange={onChange('dest_host')}
        placeholder="magma-fedgw.magma.com"
      />
      <TextField
        label="Dest Realm"
        className={classes.input}
        value={values.dest_realm}
        onChange={onChange('dest_realm')}
        placeholder="magma.com"
      />
      <TextField
        label="Host"
        className={classes.input}
        value={values.host}
        onChange={onChange('host')}
        placeholder="magma.com"
      />
      <TextField
        label="Realm"
        className={classes.input}
        value={values.realm}
        onChange={onChange('realm')}
        placeholder="realm"
      />
      <TextField
        label="Local Address"
        className={classes.input}
        value={values.local_address}
        onChange={onChange('local_address')}
        placeholder=":56789"
      />
      <TextField
        label="Product Name"
        className={classes.input}
        value={values.product_name}
        onChange={onChange('product_name')}
        placeholder="Magma"
      />
    </>
  );
}
