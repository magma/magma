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
import TextField from '@mui/material/TextField';

import MagmaAPI from '../../api/MagmaAPI';
import useMagmaAPI from '../../api/useMagmaAPI';
import {AltFormField} from '../FormField';
import {CwfNetwork} from '../../../generated';
import {getErrorMessage} from '../../util/ErrorUtils';
import {useCallback, useState} from 'react';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';

type Props = {
  onClose: () => void;
  onSave: () => void;
  networkConfig: GenericConfig;
};

export default function CWFNetworkDialog(props: Props) {
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
      <AltFormField label="Federation Network ID">
        <TextField
          name="fegNetworkID"
          fullWidth
          value={fegNetworkID}
          onChange={({target}) => setFegNetworkID(target.value)}
        />
      </AltFormField>
    </GenericNetworkDialog>
  );
}
