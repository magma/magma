/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import FieldValue from './FieldValue';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  table: {
    boxShadow: '0px 1px 4px 0px rgba(0,0,0,0.17)',
    borderRadius: '4px',
  },
  headerCell: {
    fontSize: '12px',
    lineHeight: '16px',
    color: '#8895ad',
    fontWeight: 400,
    borderBottom: 'none',
  },
  tableCell: {
    fontSize: '12px',
    lineHeight: '16px',
    color: theme.palette.blueGrayDark,
    fontWeight: 400,
    borderBottom: 'none',
  },
  row: {
    '&:nth-child(odd)': {
      backgroundColor: '#f5f7fc',
    },
  },
}));

type Props = {
  fields: Array<Object>,
};

const FanbtlNeame = (props: Props) => {
  const {fields} = props;
  const classes = useStyles();
  return (
    <Table className={classes.table}>
      <TableHead>
        <TableRow>
          <TableCell className={classes.headerCell}>Name</TableCell>
          <TableCell className={classes.headerCell}>Value</TableCell>
          <TableCell className={classes.headerCell}>Type</TableCell>
        </TableRow>
      </TableHead>
      <TableBody>
        {fields.map(field => (
          <TableRow key={field.name} className={classes.row}>
            <TableCell className={classes.tableCell}>{field.name}</TableCell>
            <TableCell className={classes.tableCell}>
              <FieldValue field={field} />
            </TableCell>
            <TableCell className={classes.tableCell}>{field.type}</TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
};

export default FanbtlNeame;
