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

import * as React from 'react';
import CWFNetworkDialog from './CWFNetworkDialog';
import FEGNetworkDialog from './FEGNetworkDialog';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import GenericNetworkDialog from './GenericNetworkDialog';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import LoadingFillerBackdrop from '../LoadingFillerBackdrop';
import MagmaV1API from '../../../generated/WebClient';

import useMagmaAPI from '../../../api/useMagmaAPIFlow';
// $FlowFixMe migrated to typescript
import {CWF, FEG} from '../../../shared/types/network';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';
import {useParams} from 'react-router-dom';

type Props = {
  onClose: () => void,
  onSave: () => void,
};

export default function NetworkDialog(props: Props) {
  const {networkID: editingNetworkID} = useParams();
  const enqueueSnackbar = useEnqueueSnackbar();

  const {response: networkConfig, isLoading} = useMagmaAPI(
    MagmaV1API.getNetworksByNetworkId,
    {
      networkId: editingNetworkID,
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
    MagmaV1API.putNetworksByNetworkId({
      networkId: networkConfig.id,
      network: networkConfig,
    })
      .then(props.onSave)
      .catch(error =>
        enqueueSnackbar(
          error?.response?.data?.message || "error: couldn't edit network",
          {
            variant: 'error',
          },
        ),
      );
  };

  return <GenericNetworkDialog {...dialogProps} onSave={onSave} />;
}
