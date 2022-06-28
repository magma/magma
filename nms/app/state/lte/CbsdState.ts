/*
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
import type {OptionsObject} from 'notistack';

import MagmaAPI from '../../../api/MagmaAPI';
import type {MutableCbsd, PaginatedCbsds} from '../../../generated-ts';

type FetchProps = {
  networkId: string;
  page: number;
  pageSize: number;
  setIsLoading: (value: boolean) => void;
  setFetchResponse: (response: PaginatedCbsds) => void;
  enqueueSnackbar?: (
    msg: string,
    cfg: OptionsObject,
  ) => string | number | null | undefined;
};

export async function fetch(props: FetchProps) {
  const {
    networkId,
    page,
    pageSize,
    setIsLoading,
    setFetchResponse,
    enqueueSnackbar,
  } = props;
  if (networkId == null) return;

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
  } catch {
    enqueueSnackbar?.('failed fetching CBSDs information', {
      variant: 'error',
    });
  } finally {
    setIsLoading(false);
  }
}

type CreateProps = {
  networkId: string;
  newCbsd: MutableCbsd;
};

export async function create(props: CreateProps) {
  const {networkId, newCbsd} = props;
  if (networkId == null) return;

  await MagmaAPI.cbsds.dpNetworkIdCbsdsPost({
    networkId,
    cbsd: newCbsd,
  });
}

type UpdateProps = {
  networkId: string;
  id: number;
  cbsd: MutableCbsd;
};

export async function update(props: UpdateProps) {
  const {networkId, id, cbsd} = props;
  if (networkId == null) return;

  await MagmaAPI.cbsds.dpNetworkIdCbsdsCbsdIdPut({
    networkId,
    cbsdId: id,
    cbsd,
  });
}

type DeregisterProps = {
  networkId: string;
  id: number;
};

export async function deregister(props: DeregisterProps) {
  const {networkId, id} = props;
  if (networkId == null) return;

  await MagmaAPI.cbsds.dpNetworkIdCbsdsCbsdIdDeregisterPost({
    networkId,
    cbsdId: id,
  });
}

type RemoveProps = {
  networkId: string;
  cbsdId: number;
};

export async function remove(props: RemoveProps) {
  const {networkId, cbsdId} = props;
  if (networkId == null) return;

  await MagmaAPI.cbsds.dpNetworkIdCbsdsCbsdIdDelete({
    networkId,
    cbsdId,
  });
}
