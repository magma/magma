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
import Text from '@fbcnms/ui/components/design-system/Text';
import {makeStyles} from '@material-ui/styles';
import {pascalCaseGoStyle} from '../../common/EntUtils';

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
  comma: {
    display: 'inline',
  },
}));

type Props = {
  edges: Array<Object>,
};

const FanbtlNeame = (props: Props) => {
  const {edges} = props;
  const classes = useStyles();
  return (
    <Table className={classes.table}>
      <TableHead>
        <TableRow>
          <TableCell className={classes.headerCell}>Name</TableCell>
          <TableCell className={classes.headerCell}>IDs</TableCell>
        </TableRow>
      </TableHead>
      <TableBody>
        {edges.map(edge => (
          <TableRow key={edge.name} className={classes.row}>
            <TableCell className={classes.tableCell}>
              {pascalCaseGoStyle(edge.name)}
            </TableCell>
            <TableCell className={classes.tableCell}>
              <div>
                {edge.ids.map((edgeId, i) => (
                  <>
                    <FieldValue field={{type: 'ID', value: edgeId}} />
                    {i === edge.ids.length - 1 ? null : (
                      <Text className={classes.comma}>{', '}</Text>
                    )}
                  </>
                ))}
              </div>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
};

export default FanbtlNeame;
