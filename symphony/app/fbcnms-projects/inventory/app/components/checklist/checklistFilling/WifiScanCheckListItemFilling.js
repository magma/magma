/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {CheckListItemFillingProps} from './CheckListItemFilling';
import type {TableRowDataType} from '@fbcnms/ui/components/design-system/Table/Table';
import type {WifiScanCheckListItemData} from '../checkListCategory/ChecklistItemsDialogMutateState';

import * as React from 'react';
import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import Table from '@fbcnms/ui/components/design-system/Table/Table';
import fbt from 'fbt';

type WifiData = $NonMaybeType<
  $ElementType<WifiScanCheckListItemData, 'wifiData'>,
>;

const WifiScanCheckListItemFilling = ({
  item,
  onChange: _onChange,
}: CheckListItemFillingProps): React.Node => {
  const wifiData: WifiData = item.wifiData ?? [];
  const data: Array<
    TableRowDataType<{|
      data: $ElementType<WifiData, number>,
    |}>,
  > = wifiData.map(wifi => ({
    key: wifi.id,
    data: wifi,
  }));

  return (
    <FormField>
      <Table
        variant="embedded"
        dataRowsSeparator="border"
        data={data}
        columns={[
          {
            key: 'ssid',
            title: <fbt desc="">SSID</fbt>,
            render: row => row.data.ssid ?? '',
          },
          {
            key: 'bssid',
            title: <fbt desc="">BSSID</fbt>,
            render: row => row.data.bssid,
          },
          {
            key: 'frequency',
            title: <fbt desc="">Frequency</fbt>,
            render: row => row.data.frequency,
          },
          {
            key: 'channel',
            title: <fbt desc="">Channel</fbt>,
            render: row => row.data.channel,
          },
          {
            key: 'band',
            title: <fbt desc="">Band</fbt>,
            render: row => row.data.band ?? '',
          },
          {
            key: 'signal',
            title: <fbt desc="">Signal</fbt>,
            render: row => row.data.strength ?? '',
          },
        ]}
      />
    </FormField>
  );
};

export default WifiScanCheckListItemFilling;
