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
 *
 * @flow strict-local
 * @format
 */
import DeviceStatusCircle from '../../theme/design-system/DeviceStatusCircle';
import LoadingFiller from '../../components/LoadingFiller';
import RadioIcon from '@material-ui/icons/Radio';
import React, {useCallback, useContext, useMemo, useState} from 'react';
import withAlert from '../../components/Alert/withAlert';
import {isNumber} from 'lodash';
import type {WithAlert} from '../../components/Alert/withAlert';

import ActionTable from '../../components/ActionTable';
import CardTitleRow from '../../components/layout/CardTitleRow';
import CbsdContext from '../../components/context/CbsdContext';
import {CbsdAddEditDialog} from './CbsdEdit';
import {Query} from '@material-table/core';
import type {Cbsd} from '../../../generated-ts';

type CbsdRowType = {
  id: number;
  isActive: boolean;
  serialNumber: string;
  state: 'registered' | 'unregistered';
  userId: string;
  fccId: string;
  cbsdId?: string;
  maxEirp?: number;
  bandwidth?: number;
  frequency?: number;
  grantExpireTime?: string;
  transmitExpireTime?: string;
};

function CbsdsTable(props: WithAlert) {
  const [isEditDialodOpen, setIsEditDialogOpen] = useState(false);

  const ctx = useContext(CbsdContext);

  const [currentRow, setCurrentRow] = useState<CbsdRowType>({} as CbsdRowType);

  const data: Array<CbsdRowType> = useMemo(() => {
    return ctx.state.cbsds
      ? ctx.state.cbsds.map(
          (item: Cbsd): CbsdRowType => {
            return {
              id: item.id,
              isActive: item.is_active,
              serialNumber: item.serial_number,
              state: item.state,
              userId: item.user_id,
              fccId: item.fcc_id,
              cbsdId: item.cbsd_id,
              maxEirp: item.grant?.max_eirp,
              bandwidth: item.grant?.bandwidth_mhz,
              frequency: item.grant?.frequency_mhz,
              grantExpireTime: item.grant?.grant_expire_time,
              transmitExpireTime: item.grant?.transmit_expire_time,
            };
          },
        )
      : [];
  }, [ctx.state.cbsds]);

  const getDataFn = useCallback(
    (query: Query<CbsdRowType>) => {
      ctx.setPaginationOptions({
        page: query.page,
        pageSize: query.pageSize,
      });

      return Promise.resolve({
        data,
        page: ctx.state.page,
        totalCount: ctx.state.totalCount,
      });
    },
    [ctx, data],
  );

  const currentCbsd = useMemo(() => {
    return ctx.state.cbsds
      ? ctx.state.cbsds.find(({id}) => id === currentRow.id)
      : undefined;
  }, [ctx.state.cbsds, currentRow]);

  if (ctx.state.isLoading) return <LoadingFiller />;

  return (
    <>
      <CbsdAddEditDialog
        open={isEditDialodOpen}
        onClose={() => setIsEditDialogOpen(false)}
        cbsd={currentCbsd}
      />

      <CardTitleRow
        key="title"
        icon={RadioIcon}
        label={`CBSDs (${data.length})`}
      />
      <ActionTable
        data={getDataFn}
        columns={[
          {
            title: 'Active Status',
            field: 'isActive',
            render: currRow => (
              <DeviceStatusCircle isGrey={false} isActive={currRow.isActive} />
            ),
          },
          {
            title: 'Serial Number',
            field: 'serialNumber',
          },
          {
            title: 'State',
            field: 'state',
          },
          {
            title: 'User ID',
            field: 'userId',
          },
          {
            title: 'FCC ID',
            field: 'fccId',
          },
          {
            title: 'CBSD ID',
            field: 'cbsdId',
          },
          {
            title: 'MAX EIRP(dBm/MHz)',
            field: 'maxEirp',
          },
          {
            title: 'Bandwidth(MHz)',
            field: 'bandwidth',
          },
          {
            title: 'Frequency(MHz)',
            field: 'frequency',
          },
          {
            title: 'Grant Expire Time',
            field: 'grantExpireTime',
          },
          {
            title: 'Transmit Expire Time',
            field: 'transmitExpireTime',
          },
        ]}
        handleCurrRow={(row: CbsdRowType) => setCurrentRow(row)}
        menuItems={[
          {
            name: 'Deregister',
            handleFunc: () => {
              void props
                .confirm(
                  `Are you sure you want to deregister ${currentRow?.serialNumber}?`,
                )
                .then(confirmed => {
                  if (!confirmed && isNumber(currentRow?.id)) {
                    return;
                  }
                  void ctx.deregister(currentRow.id);
                });
            },
          },
          {
            name: 'Edit',
            handleFunc: () => setIsEditDialogOpen(true),
          },
          {
            name: 'Remove',
            handleFunc: () => {
              void props
                .confirm(
                  `Are you sure you want to delete ${currentRow?.serialNumber}?`,
                )
                .then(confirmed => {
                  if (!confirmed && isNumber(currentRow?.id)) {
                    return;
                  }
                  void ctx.remove(currentRow.id);
                });
            },
          },
        ]}
        options={{
          actionsColumnIndex: -1,
          pageSize: ctx.state.pageSize,
          pageSizeOptions: [10, 20, 40],
          search: false,
          sorting: false,
        }}
      />
    </>
  );
}

export default withAlert(CbsdsTable);
