/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @flow
 * @format
 */

import Button from '@fbcnms/ui/components/design-system/Button';
import FormControl from '@material-ui/core/FormControl';
import FormHelperText from '@material-ui/core/FormHelperText';
import Grid from '@material-ui/core/Grid';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import MenuItem from '@material-ui/core/MenuItem';
import React, {useCallback} from 'react';
import Select from '@material-ui/core/Select';
import axios from 'axios';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';
import {makeStyles} from '@material-ui/styles';
import {sortBy} from 'lodash';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useState} from 'react';

const useStyles = makeStyles(() => ({
  header: {
    margin: '10px',
    display: 'flex',
    justifyContent: 'space-between',
  },
  paper: {
    margin: '10px',
  },
  noNetworks: {
    height: '70vh',
  },
}));

function Alerts() {
  const classes = useStyles();
  const enqueueSnackbar = useEnqueueSnackbar();

  const [network, setNetwork] = useState<string>('');
  const [networkIDs, setNetworkIDs] = useState(null);

  const {error, isLoading} = useMagmaAPI(
    MagmaV1API.getNetworks,
    {},
    useCallback(res => setNetworkIDs(sortBy(res, [n => n.toLowerCase()])), []),
  );

  if (error || isLoading || !networkIDs) {
    return <LoadingFiller />;
  }

  return (
    <div className={classes.paper}>
      <div className={classes.header}>
        <h1>Predefined Alerts</h1>
      </div>
      <p>
        Use this to create and/or sync predefined alerts for the selected
        network.
      </p>
      <Grid container>
        <Grid item xs={1}>
          <FormControl>
            <Select value={network} onChange={e => setNetwork(e.target.value)}>
              {networkIDs.map(networkID => (
                <MenuItem key={networkID} value={networkID}>
                  {networkID}
                </MenuItem>
              ))}
            </Select>
            <FormHelperText>Network to Sync</FormHelperText>
          </FormControl>
        </Grid>
        <Grid item xs={3}>
          <Button onClick={() => triggerAlertSync(network, enqueueSnackbar)}>
            Sync Alerts
          </Button>
        </Grid>
      </Grid>
    </div>
  );
}

export async function triggerAlertSync(
  networkID: string,
  enqueueSnackbar: any,
) {
  try {
    await axios.post(`/sync_alerts/${networkID}`);
    enqueueSnackbar(`Successfully synced alerts for ${networkID}`, {
      variant: 'success',
    });
  } catch (e) {
    enqueueSnackbar(`Error syncing alerts: ${e?.response?.data}`, {
      variant: 'error',
    });
  }
}

export default Alerts;
