/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import type {magmad_gateway} from '@fbcnms/magma-api';

import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableRow from '@material-ui/core/TableRow';
import Typography from '@material-ui/core/Typography';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles({
  markerContainer: {
    paddingTop: '8px',
  },
});

type Props = {
  gateway: ?magmad_gateway,
};

export default function GatewayMapMarkerPopup(props: Props) {
  const classes = useStyles();
  const {gateway} = props;
  const meta = gateway?.status?.meta;
  if (!meta) {
    return 'No data';
  }
  return (
    <div className={classes.markerContainer}>
      <Typography variant="h6" id="tableTitle">
        Gateway: {gateway?.id}
      </Typography>
      <Table>
        <TableBody>
          <TableRow key="enodeb_connected">
            <TableCell component="th" scope="row">
              Connected
            </TableCell>
            <TableCell>{meta.enodeb_connected ? 'yes' : 'no'}</TableCell>
          </TableRow>
          <TableRow key="rf_tx_on">
            <TableCell component="th" scope="row">
              RF TX On
            </TableCell>
            <TableCell>{meta.rf_tx_on ? 'yes' : 'no'}</TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>
  );
}
