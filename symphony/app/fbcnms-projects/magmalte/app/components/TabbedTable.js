/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import React, {useState} from 'react';
import Tab from '@material-ui/core/Tab';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableContainer from '@material-ui/core/TableContainer';
import TableRow from '@material-ui/core/TableRow';
import Tabs from '@material-ui/core/Tabs';
import Text from '../theme/design-system/Text';

import {colors, typography} from '../theme/default';
import {makeStyles} from '@material-ui/styles';
import {withStyles} from '@material-ui/core/styles';

const useStyles = makeStyles(theme => ({
  tab: {
    backgroundColor: colors.primary.white,
    borderRadius: '4px 4px 0 0',
    boxShadow: `inset 0 -2px 0 0 ${colors.primary.concrete}`,
    '& + &': {
      marginLeft: '4px',
    },
  },
  emptyTable: {
    backgroundColor: colors.primary.white,
    padding: theme.spacing(4),
    minHeight: '96px',
  },
  emptyTableContent: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    color: colors.primary.comet,
  },
  rowTitle: {
    color: colors.primary.brightGray,
  },
  rowText: {
    color: colors.primary.comet,
  },
}));

const MagmaTabs = withStyles({
  indicator: {
    backgroundColor: colors.secondary.dodgerBlue,
  },
})(Tabs);

const MagmaTab = withStyles({
  root: {
    fontFamily: typography.body1.fontFamily,
    fontWeight: typography.body1.fontWeight,
    fontSize: typography.body1.fontSize,
    lineHeight: typography.body1.lineHeight,
    letterSpacing: typography.body1.letterSpacing,
    color: colors.primary.brightGray,
    textTransform: 'none',
  },
})(Tab);

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
  label: string,
};

function TabPanel(props: TabPanelProps) {
  const classes = useStyles();
  const {currTabIndex, index, itemData, label} = props;

  if (itemData.length === 0) {
    return currTabIndex === index ? (
      <Paper elevation={0}>
        <Grid
          container
          alignItems="center"
          justify="center"
          className={classes.emptyTable}>
          <Grid item xs={12} className={classes.emptyTableContent}>
            <Text variant="body2">You have 0 {label} Alerts</Text>
            <Text variant="body3">
              To add alert triggers click <a href="#">alert settings</a>.
            </Text>
          </Grid>
        </Grid>
      </Paper>
    ) : null;
  }

  return (
    <>
      {currTabIndex === index ? (
        <TableContainer component={Paper} elevation={0}>
          <Table>
            <TableBody>
              {itemData.map((rowItem, rowIdx) => {
                return (
                  <TableRow key={rowIdx} data-testid={'alertName' + rowIdx}>
                    <TableCell component="th" scope="row">
                      <Text variant="body3" className={classes.rowTitle}>
                        {rowItem.name}
                      </Text>
                    </TableCell>
                    {rowItem.cols.map((cellItem, cellIdx) => {
                      return (
                        <TableCell key={rowIdx + '-' + cellIdx}>
                          <Text variant="body3" className={classes.rowText}>
                            {cellItem}
                          </Text>
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
  const classes = useStyles();
  const [currTabIndex, setCurrTabIndex] = useState<number>(0);
  const tabPanel = Object.keys(props.data).map((k: string, idx: number) => {
    return (
      <TabPanel
        key={idx}
        index={idx}
        label={k}
        currTabIndex={currTabIndex}
        itemData={props.data[k]}
      />
    );
  });
  return (
    <>
      <MagmaTabs
        value={currTabIndex}
        onChange={(_, newIndex: number) => setCurrTabIndex(newIndex)}
        variant="fullWidth">
        {Object.keys(props.data).map((k: string, idx: number) => {
          return (
            <MagmaTab
              key={idx}
              label={`${props.data[k].length} ${k}`}
              className={classes.tab}
            />
          );
        })}
      </MagmaTabs>
      {tabPanel}
    </>
  );
}
