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
 */

import type {GenericConfig} from './GenericNetworkDialog';

import * as React from 'react';
import GenericNetworkDialog from './GenericNetworkDialog';
import LoadingFillerBackdrop from '../LoadingFillerBackdrop';
import TextField from '@material-ui/core/TextField';

import MagmaAPI from '../../../api/MagmaAPI';
import useMagmaAPI from '../../../api/useMagmaAPI';
import {CwfNetwork} from '../../../generated-ts';
import {getErrorMessage} from '../../util/ErrorUtils';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useState} from 'react';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';

const useStyles = makeStyles(() => ({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
}));

type Props = {
  onClose: () => void;
  onSave: () => void;
  networkConfig: GenericConfig;
};

export default function CWFNetworkDialog(props: Props) {
  const classes = useStyles();
  const [fegNetworkID, setFegNetworkID] = useState('');
  const enqueueSnackbar = useEnqueueSnackbar();

  const {isLoading, response: cwfNetworkConfig} = useMagmaAPI(
    MagmaAPI.carrierWifiNetworks.cwfNetworkIdGet,
    {networkId: props.networkConfig.id},
    useCallback(
      (response: CwfNetwork) =>
        setFegNetworkID(response?.federation?.feg_network_id),
      [],
    ),
  );

  if (isLoading || !cwfNetworkConfig) {
    return <LoadingFillerBackdrop />;
  }

  const onSave = (genericFields: GenericConfig) => {
    MagmaAPI.carrierWifiNetworks
      .cwfNetworkIdPut({
        networkId: cwfNetworkConfig.id,
        cwfNetwork: {
          ...cwfNetworkConfig,
          name: genericFields.name,
          description: genericFields.description,
          federation: {
            ...cwfNetworkConfig.federation,
            feg_network_id: fegNetworkID,
          },
        },
      })
      .then(props.onSave)
      .catch(error => {
        enqueueSnackbar(
          getErrorMessage(error, "error: couldn't edit network"),
          {
            variant: 'error',
          },
        );
      });
  };

  return (
    <GenericNetworkDialog
      onSave={onSave}
      onClose={props.onClose}
      networkConfig={props.networkConfig}>
      <TextField
        name="fegNetworkID"
        label="Federation Network ID"
        className={classes.input}
        value={fegNetworkID}
        onChange={({target}) => setFegNetworkID(target.value)}
      />
    </GenericNetworkDialog>
  );
}
