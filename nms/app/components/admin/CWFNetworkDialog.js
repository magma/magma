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
 * @flow strict-local
 * @format
 */

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {GenericConfig} from './GenericNetworkDialog';

import * as React from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import GenericNetworkDialog from './GenericNetworkDialog';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import LoadingFillerBackdrop from '../LoadingFillerBackdrop';
import MagmaV1API from '../../../generated/WebClient';
import TextField from '@material-ui/core/TextField';

import useMagmaAPI from '../../../api/useMagmaAPIFlow';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useState} from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';

const useStyles = makeStyles(() => ({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
}));

type Props = {
  onClose: () => void,
  onSave: () => void,
  networkConfig: GenericConfig,
};

export default function CWFNetworkDialog(props: Props) {
  const classes = useStyles();
  const [fegNetworkID, setFegNetworkID] = useState('');
  const enqueueSnackbar = useEnqueueSnackbar();

  const {isLoading, response: cwfNetworkConfig} = useMagmaAPI(
    MagmaV1API.getCwfByNetworkId,
    {networkId: props.networkConfig.id},
    useCallback(
      response => setFegNetworkID(response?.federation?.feg_network_id),
      [],
    ),
  );

  if (isLoading || !cwfNetworkConfig) {
    return <LoadingFillerBackdrop />;
  }

  const onSave = genericFields => {
    MagmaV1API.putCwfByNetworkId({
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
          error.response?.data?.message || "error: couldn't edit network",
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
