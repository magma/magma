/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import * as React from 'react';
import CWFNetworkDialog from './CWFNetworkDialog';
import FEGNetworkDialog from './FEGNetworkDialog';
import GenericNetworkDialog from './GenericNetworkDialog';
import LoadingFillerBackdrop from '@fbcnms/ui/components/LoadingFillerBackdrop';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';

import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';
import {CWF, FEG} from '@fbcnms/types/network';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

type Props = {
  onClose: () => void,
  onSave: () => void,
};

export default function NetworkDialog(props: Props) {
  const editingNetworkID = useRouter().match.params.networkID;
  const enqueueSnackbar = useEnqueueSnackbar();

  const {response: networkConfig, isLoading} = useMagmaAPI(
    MagmaV1API.getNetworksByNetworkId,
    {networkId: editingNetworkID},
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
        enqueueSnackbar(error.response?.data?.error || error, {
          variant: 'error',
        }),
      );
  };

  return <GenericNetworkDialog {...dialogProps} onSave={onSave} />;
}
