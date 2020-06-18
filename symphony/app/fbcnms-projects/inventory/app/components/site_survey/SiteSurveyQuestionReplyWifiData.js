/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {SiteSurveyQuestionReplyWifiData_data} from './__generated__/SiteSurveyQuestionReplyWifiData_data.graphql';

import * as React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';

import {createFragmentContainer, graphql} from 'react-relay';

type Props = {
  data: SiteSurveyQuestionReplyWifiData_data,
};

function SiteSurveyQuestionReplyWifiData(props: Props) {
  const {wifiData} = props.data;
  const rows = (wifiData || []).filter(Boolean).map(row => (
    <TableRow>
      <TableCell>{row.ssid}</TableCell>
      <TableCell>{row.bssid}</TableCell>
      <TableCell>{row.frequency}</TableCell>
      <TableCell>{row.channel}</TableCell>
      <TableCell>{row.band}</TableCell>
      <TableCell>{row.strength}</TableCell>
    </TableRow>
  ));
  return (
    <Table>
      <TableHead>
        <TableRow>
          <TableCell>SSID</TableCell>
          <TableCell>BSSID</TableCell>
          <TableCell>Frequency</TableCell>
          <TableCell>Channel</TableCell>
          <TableCell>Band</TableCell>
          <TableCell>Signal</TableCell>
        </TableRow>
      </TableHead>
      <TableBody>{rows}</TableBody>
    </Table>
  );
}

export default createFragmentContainer(SiteSurveyQuestionReplyWifiData, {
  data: graphql`
    fragment SiteSurveyQuestionReplyWifiData_data on SurveyQuestion {
      wifiData {
        band
        bssid
        channel
        frequency
        strength
        ssid
      }
    }
  `,
});
