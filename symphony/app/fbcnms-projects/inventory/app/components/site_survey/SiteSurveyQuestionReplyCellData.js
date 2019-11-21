/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {SiteSurveyQuestionReplyCellData_data} from './__generated__/SiteSurveyQuestionReplyCellData_data.graphql';

import * as React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';

import {createFragmentContainer, graphql} from 'react-relay';

type Props = {
  data: SiteSurveyQuestionReplyCellData_data,
};

function SiteSurveyQuestionReplyCellData(props: Props) {
  const {cellData} = props.data;
  const rows = (cellData || []).filter(Boolean).map(row => (
    <TableRow>
      <TableCell>{row.networkType}</TableCell>
      <TableCell>{row.signalStrength}</TableCell>
      <TableCell>{row.baseStationID}</TableCell>
      <TableCell>{row.cellID}</TableCell>
      <TableCell>{row.locationAreaCode}</TableCell>
      <TableCell>{row.mobileCountryCode}</TableCell>
      <TableCell>{row.mobileNetworkCode}</TableCell>
    </TableRow>
  ));
  return (
    <Table>
      <TableHead>
        <TableRow>
          <TableCell>Type</TableCell>
          <TableCell>Signal</TableCell>
          <TableCell>Base Station ID</TableCell>
          <TableCell>Cell ID</TableCell>
          <TableCell>LAC</TableCell>
          <TableCell>MCC</TableCell>
          <TableCell>MNC</TableCell>
        </TableRow>
      </TableHead>
      <TableBody>{rows}</TableBody>
    </Table>
  );
}

export default createFragmentContainer(SiteSurveyQuestionReplyCellData, {
  data: graphql`
    fragment SiteSurveyQuestionReplyCellData_data on SurveyQuestion {
      cellData {
        networkType
        signalStrength
        baseStationID
        cellID
        locationAreaCode
        mobileCountryCode
        mobileNetworkCode
      }
    }
  `,
});
