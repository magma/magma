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
import DeviceStatusCircle from '../../theme/design-system/DeviceStatusCircle';
import LoadingFiller from '../../components/LoadingFiller';
import RadioIcon from '@mui/icons-material/Radio';
import React, {useCallback, useContext, useMemo, useState} from 'react';
import withAlert from '../../components/Alert/withAlert';
import {isNumber} from 'lodash';
import type {WithAlert} from '../../components/Alert/withAlert';

import ActionTable from '../../components/ActionTable';
import CardTitleRow from '../../components/layout/CardTitleRow';
import CbsdContext from '../../context/CbsdContext';
import {CbsdAddEditDialog} from './CbsdEdit';
import {Query} from '@material-table/core';
import type {Cbsd, Grant} from '../../../generated';

type CbsdRowType = {
  id: number;
  isActive: boolean;
  serialNumber: string;
  state: 'registered' | 'unregistered';
  userId: string;
  fccId: string;
  cbsdId?: string;
  grantState?: React.ReactNode;
  maxEirp?: React.ReactNode;
  bandwidth?: React.ReactNode;
  frequency?: React.ReactNode;
  grantExpireTime?: React.ReactNode;
  transmitExpireTime?: React.ReactNode;
};

type GrantFieldCellContentProps = {
  grants: Cbsd['grants'];
  field: keyof Grant;
};

/**
 * Quick & dirty UX solution to display CBSDs with multiple grants for v1.8.
 * Each grant has 5 fields, and we have columns for those fields.
 * Previously we had 1 grant per CBSD, so it was easy.
 * Now we have multiple grants, so each of 5 columns can have multiple values.
 * As a quick solution display each value in the same cell, with <hr /> separator.
 * The UX is planned to be improved after v1.8 release.
 */
function GrantFieldCellContent({grants, field}: GrantFieldCellContentProps) {
  if (!grants?.length) {
    return null;
  }
  const fieldValues = grants?.map(g => g[field]);

  // Using index for key, because grants have no ids,
  // and may have identical values for some fields.
  const fragments = fieldValues.map((value, idx) => (
    <React.Fragment key={idx}>
      {value}
      {idx !== fieldValues.length - 1 && <hr />}
    </React.Fragment>
  ));

  return <>{fragments}</>;
}

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
              grantState: (
                <GrantFieldCellContent grants={item.grants} field="state" />
              ),
              maxEirp: (
                <GrantFieldCellContent grants={item.grants} field="max_eirp" />
              ),
              bandwidth: (
                <GrantFieldCellContent
                  grants={item.grants}
                  field="bandwidth_mhz"
                />
              ),
              frequency: (
                <GrantFieldCellContent
                  grants={item.grants}
                  field="frequency_mhz"
                />
              ),
              grantExpireTime: (
                <GrantFieldCellContent
                  grants={item.grants}
                  field="grant_expire_time"
                />
              ),
              transmitExpireTime: (
                <GrantFieldCellContent
                  grants={item.grants}
                  field="transmit_expire_time"
                />
              ),
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
            title: 'Grant state',
            field: 'grantState',
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
            name: 'Relinquish',
            handleFunc: () => {
              void props
                .confirm(
                  `Are you sure you want to relinquish ${currentRow?.serialNumber}?`,
                )
                .then(confirmed => {
                  if (!confirmed && isNumber(currentRow?.id)) {
                    return;
                  }
                  void ctx.relinquish(currentRow.id);
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
