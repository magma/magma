/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {CellScanCheckListItemData} from '../checkListCategory/ChecklistItemsDialogMutateState';
import type {CheckListItemFillingProps} from './CheckListItemFilling';
import type {TableRowDataType} from '@fbcnms/ui/components/design-system/Table/Table';

import * as React from 'react';
import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import Table from '@fbcnms/ui/components/design-system/Table/Table';
import fbt from 'fbt';

type CellData = $NonMaybeType<
  $ElementType<CellScanCheckListItemData, 'cellData'>,
>;

const CellScanCheckListItemFilling = ({
  item,
  onChange: _onChange,
}: CheckListItemFillingProps): React.Node => {
  const cellData: CellData = item.cellData ?? [];
  const data: Array<
    TableRowDataType<{|
      data: $ElementType<CellData, number>,
    |}>,
  > = cellData.map(cell => ({
    key: cell.id,
    data: cell,
  }));

  return (
    <FormField>
      <Table
        variant="embedded"
        dataRowsSeparator="border"
        data={data}
        columns={[
          {
            key: 'network_type',
            title: <fbt desc="">Type</fbt>,
            render: row => row.data.networkType,
          },
          {
            key: 'signal',
            title: <fbt desc="">Signal</fbt>,
            render: row => row.data.signalStrength,
          },
          {
            key: 'base_station_id',
            title: <fbt desc="">Base Station ID</fbt>,
            render: row => row.data.baseStationID ?? '',
          },
          {
            key: 'cell_id',
            title: <fbt desc="">Cell ID</fbt>,
            render: row => row.data.cellID ?? '',
          },
          {
            key: 'lac',
            title: <fbt desc="">LAC</fbt>,
            render: row => row.data.locationAreaCode ?? '',
          },
          {
            key: 'mcc',
            title: <fbt desc="">MCC</fbt>,
            render: row => row.data.mobileCountryCode ?? '',
          },
          {
            key: 'mnc',
            title: <fbt desc="">MNC</fbt>,
            render: row => row.data.mobileNetworkCode ?? '',
          },
        ]}
      />
    </FormField>
  );
};

export default CellScanCheckListItemFilling;
