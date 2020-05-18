/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import Paper from '@material-ui/core/Paper';
import React, {useState} from 'react';
import Tab from '@material-ui/core/Tab';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableContainer from '@material-ui/core/TableContainer';
import TableRow from '@material-ui/core/TableRow';
import Tabs from '@material-ui/core/Tabs';

export type RowData = {
  name: string,
  cols: Array<string>,
};

export type Props = {
  data: {
    [string]: Array<RowData>,
  },
};

type TabPanelProps = {
  currTabIndex: number,
  index: number,
  itemData: Array<RowData>,
};

function TabPanel(props: TabPanelProps) {
  const {currTabIndex, index, itemData} = props;

  return (
    <>
      {currTabIndex === index ? (
        <TableContainer component={Paper}>
          <Table>
            <TableBody>
              {itemData.map((rowItem, rowIdx) => {
                return (
                  <TableRow key={rowIdx} data-testid={'alertName' + rowIdx}>
                    <TableCell component="th" scope="row">
                      {rowItem.name}
                    </TableCell>
                    {rowItem.cols.map((cellItem, cellIdx) => {
                      return (
                        <TableCell key={rowIdx + '-' + cellIdx}>
                          {cellItem}
                        </TableCell>
                      );
                    })}
                  </TableRow>
                );
              })}
            </TableBody>
          </Table>
        </TableContainer>
      ) : null}
    </>
  );
}

export default function TabbedTable(props: Props) {
  const [currTabIndex, setCurrTabIndex] = useState<number>(0);
  const tabPanel = Object.keys(props.data).map((k: string, idx: number) => {
    return (
      <TabPanel
        key={idx}
        index={idx}
        currTabIndex={currTabIndex}
        itemData={props.data[k]}
      />
    );
  });
  return (
    <>
      <Tabs
        value={currTabIndex}
        onChange={(_, newIndex: number) => setCurrTabIndex(newIndex)}
        indicatorColor="primary"
        textColor="primary"
        variant="fullWidth">
        {Object.keys(props.data).map((k: string, idx: number) => {
          return <Tab key={idx} label={k} />;
        })}
      </Tabs>
      {tabPanel}
    </>
  );
}
