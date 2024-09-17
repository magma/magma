/**
 * Copyright 2022 The Magma Authors.
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
import {AxiosError} from 'axios';
import {useCallback, useEffect, useMemo, useState} from 'react';

import MagmaAPI from '../api/MagmaAPI';
import {NetworkId} from '../../shared/types/network';
import {OptionsObject} from 'notistack';
import {useEnqueueSnackbar} from '../hooks/useSnackbar';
import type {Cbsd, MutableCbsd, PaginatedCbsds} from '../../generated';

export type CbsdContextType = {
  state: {
    isLoading: boolean;
    cbsds: Array<Cbsd>;
    totalCount: number;
    page: number;
    pageSize: number;
  };
  setPaginationOptions: (options: {page: number; pageSize: number}) => void;
  refetch: () => Promise<void>;
  create: (newCbsd: MutableCbsd) => Promise<void>;
  update: (id: number, cbsd: MutableCbsd) => Promise<void>;
  deregister: (id: number) => Promise<void>;
  relinquish: (id: number) => Promise<void>;
  remove: (id: number) => Promise<void>;
};

const CbsdContext = React.createContext<CbsdContextType>({} as CbsdContextType);

async function fetch(params: {
  networkId: string;
  page: number;
  pageSize: number;
  setIsLoading: (value: boolean) => void;
  setFetchResponse: (response: PaginatedCbsds) => void;
  enqueueSnackbar?: (
    msg: string,
    cfg: OptionsObject,
  ) => string | number | null | undefined;
}) {
  const {
    networkId,
    page,
    pageSize,
    setIsLoading,
    setFetchResponse,
    enqueueSnackbar,
  } = params;
  if (networkId == null) return;

  /**
   * Temporary hotfix for https://github.com/magma/domain-proxy/issues/469.
   * Ignore error when domain proxy module is disabled,
   * so that an error is not shown on the main page.
   */
  const handleError = (error: unknown) => {
    const status = (error as AxiosError)?.response?.status;

    if (status === 404) {
      console.error(
        'CBSD endpoint not found. Is Domain Proxy module enabled in orc8r deployment?',
      );
      return;
    }

    if (status === 502) {
      console.error(
        'CBSD endpoint returns 502. Is Domain Proxy module up and running?',
      );
      return;
    }

    enqueueSnackbar?.('failed fetching CBSDs information', {
      variant: 'error',
    });
  };

  try {
    setIsLoading(true);

    const response = (
      await MagmaAPI.cbsds.dpNetworkIdCbsdsGet({
        networkId,
        offset: page * pageSize,
        limit: pageSize,
      })
    ).data;
    setFetchResponse(response);
  } catch (error) {
    handleError(error);
  } finally {
    setIsLoading(false);
  }
}

async function create(params: {networkId: string; newCbsd: MutableCbsd}) {
  const {networkId, newCbsd} = params;
  if (networkId == null) return;

  await MagmaAPI.cbsds.dpNetworkIdCbsdsPost({
    networkId,
    cbsd: newCbsd,
  });
}

async function update(params: {
  networkId: string;
  id: number;
  cbsd: MutableCbsd;
}) {
  const {networkId, id, cbsd} = params;
  if (networkId == null) return;

  await MagmaAPI.cbsds.dpNetworkIdCbsdsCbsdIdPut({
    networkId,
    cbsdId: id,
    cbsd,
  });
}

export async function deregister(params: {networkId: string; id: number}) {
  const {networkId, id} = params;
  if (networkId == null) return;

  await MagmaAPI.cbsds.dpNetworkIdCbsdsCbsdIdDeregisterPost({
    networkId,
    cbsdId: id,
  });
}

export async function relinquish(params: {networkId: string; id: number}) {
  const {networkId, id} = params;
  if (networkId == null) return;

  await MagmaAPI.cbsds.dpNetworkIdCbsdsCbsdIdRelinquishPost({
    networkId,
    cbsdId: id,
  });
}

export async function remove(params: {networkId: string; cbsdId: number}) {
  const {networkId, cbsdId} = params;
  if (networkId == null) return;

  await MagmaAPI.cbsds.dpNetworkIdCbsdsCbsdIdDelete({
    networkId,
    cbsdId,
  });
}

export function CbsdContextProvider({
  networkId,
  children,
}: {
  networkId: NetworkId;
  children: React.ReactNode;
}) {
  const enqueueSnackbar = useEnqueueSnackbar();

  const [isLoading, setIsLoading] = useState(false);
  const [fetchResponse, setFetchResponse] = useState<PaginatedCbsds>({
    cbsds: [],
    total_count: 0,
  });
  const [paginationOptions, setPaginationOptions] = useState<{
    page: number;
    pageSize: number;
  }>({
    page: 0,
    pageSize: 10,
  });

  const refetch = useCallback(() => {
    return fetch({
      networkId,
      page: paginationOptions.page,
      pageSize: paginationOptions.pageSize,
      setIsLoading,
      setFetchResponse,
      enqueueSnackbar,
    });
  }, [
    networkId,
    paginationOptions.page,
    paginationOptions.pageSize,
    setIsLoading,
    setFetchResponse,
    enqueueSnackbar,
  ]);

  useEffect(() => {
    void refetch();
  }, [refetch, paginationOptions.page, paginationOptions.pageSize]);

  const state = useMemo(() => {
    return {
      isLoading,
      cbsds: fetchResponse.cbsds,
      totalCount: fetchResponse.total_count,
      page: paginationOptions.page,
      pageSize: paginationOptions.pageSize,
    };
  }, [
    isLoading,
    fetchResponse.cbsds,
    fetchResponse.total_count,
    paginationOptions.page,
    paginationOptions.pageSize,
  ]);

  return (
    <CbsdContext.Provider
      value={{
        state,
        setPaginationOptions,
        refetch,
        create: (newCbsd: MutableCbsd) => {
          return create({
            networkId,
            newCbsd,
          })
            .catch(e => {
              enqueueSnackbar?.('failed to create CBSD', {
                variant: 'error',
              });
              throw e as Error;
            })
            .then(() => {
              void refetch();
            });
        },
        update: (id: number, cbsd: MutableCbsd) => {
          return update({
            networkId,
            id,
            cbsd,
          })
            .catch(e => {
              enqueueSnackbar?.('failed to update CBSD', {
                variant: 'error',
              });
              throw e as Error;
            })
            .then(() => {
              void refetch();
            });
        },
        deregister: (id: number) => {
          return deregister({
            networkId,
            id,
          })
            .catch(() => {
              enqueueSnackbar?.('failed to deregister CBSD', {
                variant: 'error',
              });
            })
            .then(() => {
              void refetch();
            });
        },
        relinquish: (id: number) => {
          return relinquish({
            networkId,
            id,
          })
            .catch(() => {
              enqueueSnackbar?.('failed to relinquish CBSD', {
                variant: 'error',
              });
            })
            .then(() => {
              void refetch();
            });
        },
        remove: (id: number) => {
          return remove({
            networkId,
            cbsdId: id,
          })
            .catch(() => {
              enqueueSnackbar?.('failed to remove CBSD', {
                variant: 'error',
              });
            })
            .then(() => {
              void refetch();
            });
        },
      }}>
      {children}
    </CbsdContext.Provider>
  );
}

export default CbsdContext;
