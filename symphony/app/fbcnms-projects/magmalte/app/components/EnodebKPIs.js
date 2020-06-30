/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {KPIData} from './KPITray';
import type {enodeb, enodeb_state} from '@fbcnms/magma-api';

import KPITray from './KPITray';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import nullthrows from '@fbcnms/util/nullthrows';
import React, {useEffect, useState} from 'react';
import SettingsInputAntennaIcon from '@material-ui/icons/SettingsInputAntenna';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';

import {useRouter} from '@fbcnms/ui/hooks';

export default function EnodebKPIs() {
  const [enodebSt, setEnodebSt] = useState<{[string]: enodeb_state}>({});
  const {match} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);

  const {response} = useMagmaAPI(MagmaV1API.getLteByNetworkIdEnodebs, {
    networkId: networkId,
  });

  useEffect(() => {
    const fetchEnodebState = async () => {
      if (!response) {
        return;
      }

      const enodebInfo: {[string]: enodeb} = response;
      const requests = Object.keys(enodebInfo).map(async k => {
        const {serial} = enodebInfo[k];
        try {
          // eslint-disable-next-line max-len
          const enbSt = await MagmaV1API.getLteByNetworkIdEnodebsByEnodebSerialState(
            {
              networkId: networkId,
              enodebSerial: serial,
            },
          );
          return {serial, enbSt};
        } catch (error) {
          console.error('error getting enodeb status for ' + serial);
          return null;
        }
      });
      Promise.all(requests).then(allResponses => {
        const enodebState = {};
        allResponses.filter(Boolean).forEach(r => {
          enodebState[r.serial] = r.enbSt;
        });
        setEnodebSt(enodebState);
      });
    };
    fetchEnodebState();
  }, [networkId, response]);

  const [total, transmitting] = enodebStatus(enodebSt);
  const kpiData: KPIData[] = [
    {category: 'Severe Events', value: 0},
    {category: 'Total', value: total || 0},
    {category: 'Transmitting', value: transmitting || 0},
  ];
  return (
    <KPITray
      icon={SettingsInputAntennaIcon}
      description="eNodeBs"
      data={kpiData}
    />
  );
}

function enodebStatus(enodebSt: {[string]: enodeb_state}): [number, number] {
  let transmitCnt = 0;
  Object.keys(enodebSt)
    .map((k: string) => enodebSt[k])
    .map((e: enodeb_state) => {
      if (e.rf_tx_on) {
        transmitCnt++;
      }
    });
  return [Object.keys(enodebSt).length, transmitCnt];
}
