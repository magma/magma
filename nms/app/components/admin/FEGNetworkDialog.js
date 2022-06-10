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

import type {GenericConfig} from './GenericNetworkDialog';

import * as React from 'react';
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

export default function FEGNetworkDialog(props: Props) {
  const classes = useStyles();
  const [servedNetworkIDs, setServedNetworkIDs] = useState('');
  const enqueueSnackbar = useEnqueueSnackbar();

  const {response: fegNetworkConfig, isLoading} = useMagmaAPI(
    MagmaV1API.getFegByNetworkId,
    {networkId: props.networkConfig.id},
    useCallback(
      response =>
        setServedNetworkIDs(
          (response?.federation?.served_network_ids || []).join(','),
        ),
      [],
    ),
  );

  if (isLoading || !fegNetworkConfig) {
    return <LoadingFillerBackdrop />;
  }

  const onSave = genericFields => {
    MagmaV1API.putFegByNetworkId({
      networkId: fegNetworkConfig.id,
      fegNetwork: {
        ...fegNetworkConfig,
        name: genericFields.name,
        description: genericFields.description,
        federation: {
          ...fegNetworkConfig.federation,
          served_network_ids: servedNetworkIDs.split(','),
        },
      },
    })
      .then(props.onSave)
      .catch(error =>
        enqueueSnackbar(
          error.response?.data?.message || "error: couldn't edit network",
          {
            variant: 'error',
          },
        ),
      );
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
        value={servedNetworkIDs}
        onChange={({target}) => setServedNetworkIDs(target.value)}
      />
    </GenericNetworkDialog>
  );
}
