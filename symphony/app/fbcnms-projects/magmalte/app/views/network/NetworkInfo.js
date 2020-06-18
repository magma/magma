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
import type {network} from '../../../../../fbcnms-packages/fbcnms-magma-api';

import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextField from '@material-ui/core/TextField';
import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';

import {useCallback, useState} from 'react';
import {useRouter} from '@fbcnms/ui/hooks';

export default function NetworkInfo({readOnly}: {readOnly: boolean}) {
  const {match} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);
  const [networkInfo, setNetworkInfo] = useState<network>({});
  const [networkType, setNetworkType] = useState('');

  const {isLoading} = useMagmaAPI(
    MagmaV1API.getNetworksByNetworkId,
    {
      networkId: networkId,
    },
    useCallback(networkInfo => {
      setNetworkInfo(networkInfo);
      setNetworkType(networkInfo?.type ?? '');
    }, []),
  );

  if (isLoading) {
    return <LoadingFiller />;
  }

  if (Object.keys(networkInfo).length === 0) {
    return null;
  }

  return (
    <Grid container>
      <Grid container item xs={12}>
        <Grid item>
          <Text weight="medium" variant="h5">
            Network
          </Text>
        </Grid>
        <Grid container item justify="flex-end">
          <Text>Edit</Text>
        </Grid>
      </Grid>
      <Grid item xs={12}>
        <List component={Paper}>
          <ListItem>
            <TextField
              fullWidth={true}
              value={networkInfo.name}
              label="Name"
              onChange={({target}) =>
                setNetworkInfo({...networkInfo, name: target.value})
              }
              InputProps={{disableUnderline: true, readOnly: readOnly}}
            />
          </ListItem>
          <Divider />
          <ListItem>
            <TextField
              fullWidth={true}
              value={networkType}
              label="Network Type"
              onChange={({target}) => {
                setNetworkType(target.value);
                setNetworkInfo({...networkInfo, type: target.value});
              }}
              InputProps={{disableUnderline: true, readOnly: readOnly}}
            />
          </ListItem>
          <Divider />
          <ListItem>
            <TextField
              fullWidth={true}
              value={networkInfo.description}
              label="Description"
              onChange={({target}) =>
                setNetworkInfo({...networkInfo, description: target.value})
              }
              InputProps={{disableUnderline: true, readOnly: readOnly}}
            />
          </ListItem>
        </List>
      </Grid>
    </Grid>
  );
}
