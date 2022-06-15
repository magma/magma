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

import * as React from 'react';
import CWFNetworkDialog from './CWFNetworkDialog';
import FEGNetworkDialog from './FEGNetworkDialog';
import GenericNetworkDialog from './GenericNetworkDialog';
import LoadingFillerBackdrop from '../LoadingFillerBackdrop';

import MagmaAPI from '../../../api/MagmaAPI';
import useMagmaAPI from '../../../api/useMagmaAPI';
import {CWF, FEG} from '../../../shared/types/network';
import {getErrorMessage} from '../../util/ErrorUtils';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';
import {useParams} from 'react-router-dom';

type Props = {
  onClose: () => void;
  onSave: () => void;
};

export default function NetworkDialog(props: Props) {
  const {networkID: editingNetworkID} = useParams();
  const enqueueSnackbar = useEnqueueSnackbar();

  const {response: networkConfig, isLoading} = useMagmaAPI(
    MagmaAPI.networks.networksNetworkIdGet,
    {
      networkId: editingNetworkID!,
    },
  );

  if (!networkConfig || isLoading) {
    return <LoadingFillerBackdrop />;
  }

  const dialogProps = {
    onSave: props.onSave,
    onClose: props.onClose,
    networkConfig,
  };

  switch (networkConfig.type) {
    case FEG:
      return <FEGNetworkDialog {...dialogProps} />;
    case CWF:
      return <CWFNetworkDialog {...dialogProps} />;
  }

  const onSave = () => {
    MagmaAPI.networks
      .networksNetworkIdPut({
        networkId: networkConfig.id,
        network: networkConfig,
      })
      .then(props.onSave)
      .catch(error =>
        enqueueSnackbar(
          getErrorMessage(error, "error: couldn't edit network"),
          {
            variant: 'error',
          },
        ),
      );
  };

  return <GenericNetworkDialog {...dialogProps} onSave={onSave} />;
}
