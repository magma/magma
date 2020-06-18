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
import type {network_epc_configs} from '../../../../../fbcnms-packages/fbcnms-magma-api';

import Divider from '@material-ui/core/Divider';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Grid from '@material-ui/core/Grid';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Switch from '@material-ui/core/Switch';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextField from '@material-ui/core/TextField';
import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';

import {useCallback, useState} from 'react';
import {useRouter} from '@fbcnms/ui/hooks';

export default function NetworkEpc(props: {readOnly: boolean}) {
  const {match} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);
  const [networkEpc, setNetworkEpc] = useState<network_epc_configs>({});

  const {isLoading} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdCellularEpc,
    {
      networkId: networkId,
    },
    useCallback(epc => setNetworkEpc(epc), []),
  );

  if (isLoading) {
    return <LoadingFiller />;
  }

  if (Object.keys(networkEpc).length === 0) {
    return null;
  }

  //(PolicyEnforcement, LTE Auth AMF, MCC, MNC, TAC)
  return (
    <Grid container>
      <Grid container item xs={12}>
        <Grid item>
          <Text weight="medium" variant="h5">
            EPC
          </Text>
        </Grid>
        <Grid container item justify="flex-end">
          <Text>Edit</Text>
        </Grid>
      </Grid>
      <Grid item xs={12}>
        <List component={Paper}>
          <ListItem>
            <FormControlLabel
              control={
                <Switch
                  disabled={props.readOnly}
                  color="primary"
                  checked={networkEpc.relay_enabled}
                  onChange={({target}) =>
                    setNetworkEpc({
                      ...networkEpc,
                      relay_enabled: target.checked,
                    })
                  }
                  name="checkedA"
                />
              }
              label="Policy Enforcement Enabled"
            />
          </ListItem>
          <Divider />
          <ListItem>
            <TextField
              type="password"
              fullWidth={true}
              value={networkEpc.lte_auth_amf}
              label="LTE Auth AMF"
              onChange={({target}) => {
                setNetworkEpc({...networkEpc, lte_auth_amf: target.value});
              }}
              InputProps={{disableUnderline: true, readOnly: props.readOnly}}
            />
          </ListItem>
          <Divider />
          <ListItem>
            <TextField
              fullWidth={true}
              value={networkEpc.mcc}
              label="MCC"
              onChange={({target}) =>
                setNetworkEpc({...networkEpc, mcc: target.value})
              }
              InputProps={{disableUnderline: true, readOnly: props.readOnly}}
            />
          </ListItem>
          <Divider />
          <ListItem>
            <TextField
              fullWidth={true}
              value={networkEpc.mnc}
              label="MNC"
              onChange={({target}) =>
                setNetworkEpc({...networkEpc, mnc: target.value})
              }
              InputProps={{disableUnderline: true, readOnly: props.readOnly}}
            />
          </ListItem>
          <Divider />
          <ListItem>
            <TextField
              type="number"
              fullWidth={true}
              value={networkEpc.tac}
              label="TAC"
              onChange={({target}) =>
                setNetworkEpc({...networkEpc, tac: parseInt(target.value)})
              }
              InputProps={{disableUnderline: true, readOnly: props.readOnly}}
            />
          </ListItem>
        </List>
      </Grid>
    </Grid>
  );
}
