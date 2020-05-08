/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import KPITray from './KPITray';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React, {useEffect, useState} from 'react';
import SettingsInputAntennaIcon from '@material-ui/icons/SettingsInputAntenna';
import nullthrows from '@fbcnms/util/nullthrows';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';
import type {KPIData} from './KPITray';
import type {enodeb, enodeb_state} from '@fbcnms/magma-api';

export default function EnodebKPIs() {
  const [enodebInfo, setEnodebInfo] = useState<{[string]: enodeb}>({});
  const [enodebSt, setEnodebSt] = useState<{[string]: enodeb_state}>({});
  const {match} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);
  const enqueueSnackbar = useEnqueueSnackbar();

  useEffect(() => {
    const fetchEnodeInfo = async () => {
      try {
        const enbInfo = await MagmaV1API.getLteByNetworkIdEnodebs({
          networkId: networkId,
        });
        setEnodebInfo(enbInfo);
      } catch (error) {
        enqueueSnackbar('Error getting enodeb information', {
          variant: 'error',
        });
      }
    };
    fetchEnodeInfo();
  }, [networkId, enqueueSnackbar]);

  useEffect(() => {
    const fetchEnodebState = async () => {
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
  }, [networkId, enqueueSnackbar, enodebInfo]);

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
