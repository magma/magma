/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import ActionTable from '../../components/ActionTable';
import React from 'react';

import type {EnodebInfo} from '../../components/lte/EnodebUtils';
import type {lte_gateway} from '@fbcnms/magma-api';

import {isEnodebHealthy} from '../../components/lte/EnodebUtils';
import {useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';

type EnodebRowType = {
  name: string,
  id: string,
  health: string,
};

export default function GatewayDetailEnodebs(
  {gwInfo}: {gwInfo: lte_gateway},
  enbInfo: {[string]: EnodebInfo},
) {
  const {history, match} = useRouter();
  const [currRow, setCurrRow] = useState<EnodebRowType>({});

  const enbRows: Array<EnodebRowType> = gwInfo.connected_enodeb_serials.map(
    (serialNum: string) => {
      const enbInf = enbInfo[serialNum];
      return {
        name: enbInf.enb.name,
        id: serialNum,
        health: isEnodebHealthy(enbInf) ? 'Good' : 'Bad',
      };
    },
  );

  return (
    <ActionTable
      title=""
      data={enbRows}
      columns={[
        {title: 'Name', field: 'name'},
        {title: 'Serial Number', field: 'id'},
        {title: 'Health', field: 'health'},
      ]}
      handleCurrRow={(row: EnodebRowType) => setCurrRow(row)}
      menuItems={[
        {
          name: 'View',
          handleFunc: () => {
            history.push(
              match.url.replace(`gateway/${gwInfo.id}`, `enodeb/${currRow.id}`),
            );
          },
        },
        {name: 'Edit'},
        {name: 'Remove'},
        {name: 'Deactivate'},
        {name: 'Reboot'},
      ]}
      options={{
        actionsColumnIndex: -1,
        pageSizeOptions: [5],
        toolbar: false,
      }}
    />
  );
}
