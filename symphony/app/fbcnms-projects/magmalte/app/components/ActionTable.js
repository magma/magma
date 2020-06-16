/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import type {ComponentType} from 'react';

import ArrowUpward from '@material-ui/icons/ArrowUpward';
import ChevronLeft from '@material-ui/icons/ChevronLeft';
import ChevronRight from '@material-ui/icons/ChevronRight';
import Clear from '@material-ui/icons/Clear';
import FilterList from '@material-ui/icons/FilterList';
import FirstPage from '@material-ui/icons/FirstPage';
import Grid from '@material-ui/core/Grid';
import LastPage from '@material-ui/icons/LastPage';
import MaterialTable from 'material-table';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import MoreVertIcon from '@material-ui/icons/MoreVert';
import React, {useState} from 'react';
import Remove from '@material-ui/icons/Remove';
import SaveAlt from '@material-ui/icons/SaveAlt';
import Search from '@material-ui/icons/Search';
import Text from '../theme/design-system/Text';

import {colors} from '../theme/default';
import {forwardRef} from 'react';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  cardTitleRow: {
    marginBottom: theme.spacing(1),
    minHeight: '36px',
  },
  cardTitleIcon: {
    fill: colors.primary.comet,
    marginRight: theme.spacing(1),
  },
}));

const tableIcons = {
  Export: forwardRef((props, ref) => <SaveAlt {...props} ref={ref} />),
  FirstPage: forwardRef((props, ref) => <FirstPage {...props} ref={ref} />),
  LastPage: forwardRef((props, ref) => <LastPage {...props} ref={ref} />),
  NextPage: forwardRef((props, ref) => <ChevronRight {...props} ref={ref} />),
  PreviousPage: forwardRef((props, ref) => (
    <ChevronLeft {...props} ref={ref} />
  )),
  ResetSearch: forwardRef((props, ref) => <Clear {...props} ref={ref} />),
  Search: forwardRef((props, ref) => <Search {...props} ref={ref} />),
  SortArrow: forwardRef((props, ref) => <ArrowUpward {...props} ref={ref} />),
  ThirdStateCheck: forwardRef((props, ref) => <Remove {...props} ref={ref} />),
  Filter: forwardRef((props, ref) => <FilterList {...props} ref={ref} />),
};
type ActionMenuItems = {
  name: string,
  handleFunc?: () => void,
};

type ColumnType =
  | 'boolean'
  | 'numeric'
  | 'date'
  | 'datetime'
  | 'time'
  | 'currency';

type ActionTableColumn = {
  title: string,
  type?: ColumnType,
  field: string,
};

type ActionTableOptions = {
  actionsColumnIndex: number,
  pageSizeOptions: Array<number>,
};

type ActionOrderType = {
  field: string,
  title: string,
  tableData: {},
};

export type ActionQuery = {
  filters: Array<string>,
  orderBy: ActionOrderType,
  orderDirection: string,
  page: number,
  pageSize: number,
  search: string,
  totalCount: number,
};

export type ActionTableProps<T> = {
  titleIcon?: ComponentType<SvgIconExports>,
  title: string,
  handleCurrRow?: T => void,
  columns: Array<ActionTableColumn>,
  menuItems?: Array<ActionMenuItems>,
  data: Array<T> | (ActionQuery => {}),
  options: ActionTableOptions,
};

export default function ActionTable<T>(props: ActionTableProps<T>) {
  const classes = useStyles();
  const [anchorEl, setAnchorEl] = useState(null);
  const actionTableJSX = [];

  const handleClick = (event, row: T) => {
    setAnchorEl(event.currentTarget);
    if (props.handleCurrRow) {
      props.handleCurrRow(row);
    }
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  if (props.titleIcon) {
    const TitleIcon = props.titleIcon;
    actionTableJSX.push(
      <Grid
        container
        alignItems="center"
        className={classes.cardTitleRow}
        key="title">
        <TitleIcon className={classes.cardTitleIcon} />
        <Text variant="body1">
          {props.title} ({props.data.length})
        </Text>
      </Grid>,
    );
  }
  if (props.menuItems) {
    const menuItems: Array<ActionMenuItems> = props.menuItems;
    actionTableJSX.push(
      <Menu
        key="menu"
        id="actions-menu"
        data-testid="actions-menu"
        anchorEl={anchorEl}
        keepMounted
        open={Boolean(anchorEl)}
        onClose={handleClose}>
        {menuItems.map(item => (
          <MenuItem
            key={item.name}
            onClick={() => {
              if (item.handleFunc) {
                item.handleFunc();
              }
            }}>
            {item.name}
          </MenuItem>
        ))}
      </Menu>,
    );
  }
  return (
    <>
      {actionTableJSX}
      {/* TODO: How do I modify this component??? Such as changine paper elevation, search placement (should be toggle open/closed), etc. */}
      <MaterialTable
        title=""
        columns={props.columns}
        icons={tableIcons}
        data={props.data}
        actions={
          props.menuItems
            ? [
                {
                  icon: () => <MoreVertIcon />,
                  tooltip: 'Actions',
                  onClick: handleClick,
                },
              ]
            : null
        }
        options={props.options}
      />
    </>
  );
}
