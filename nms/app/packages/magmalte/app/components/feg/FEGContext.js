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
import FEGGatewayContext from '../context/FEGGatewayContext';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';

import type {federation_gateway, network_id, network_type} from '@fbcnms/magma-api';

import {SetGatewayState} from '../../state/feg/EquipmentState';

import {useContext, useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';

type Props = {
    networkId: network_id,
    networkType: network_type,
    children: React.Node,
};
import {FEG} from '@fbcnms/types/network';
export function FEGGatewayContextProvider(props: Props) {
    const {networkId} = props;
    const [fegGateways, setFegGateways] = useState<{[string]: federation_gateway}>({});
    const [isLoading, setIsLoading] = useState(true);
    const enqueueSnackbar = useEnqueueSnackbar();

    useEffect(() => {
        const fetchState = async () => {
        try {
            const fegGateways = await MagmaV1API.getFegByNetworkIdGateways({
            networkId,
            });
            setFegGateways(fegGateways);
        } catch (e) {
            enqueueSnackbar?.('failed fetching gateway information', {
            variant: 'error',
            });
        }
        setIsLoading(false);
        };
        fetchState();
    }, [networkId, enqueueSnackbar]);

    if (isLoading) {
        return <LoadingFiller />;
    }

    return (
        <FEGGatewayContext.Provider
        value={{
            state: fegGateways,
            setState: (key, value?, newState?) => {
                return SetGatewayState({
                    fegGateways,
                    setFegGateways,
                    networkId,
                    key,
                    value,
                    newState,
                });
            },
        }}>
        {props.children}
        </FEGGatewayContext.Provider>
    );
}

export function FEGContextProvider(props: Props) {
    const {networkId, networkType} = props;
    const fegNetwork = networkType === FEG;
    if (!fegNetwork) {
      return props.children;
    }

    return (
        <FEGGatewayContextProvider {...{networkId, networkType}}>
            {props.children}
        </FEGGatewayContextProvider>
    );
}
