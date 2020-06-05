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
import FirstPage from '@material-ui/icons/FirstPage';
import LastPage from '@material-ui/icons/LastPage';
import MaterialTable from 'material-table';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import MoreVertIcon from '@material-ui/icons/MoreVert';
import React, {useState} from 'react';
import Remove from '@material-ui/icons/Remove';
import Search from '@material-ui/icons/Search';
import Text from '@fbcnms/ui/components/design-system/Text';

import {forwardRef} from 'react';

const tableIcons = {
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

export type ActionTableProps<T> = {
  titleIcon?: ComponentType<SvgIconExports>,
  title: string,
  handleCurrRow?: T => void,
  columns: Array<ActionTableColumn>,
  menuItems?: Array<ActionMenuItems>,
  data: Array<T>,
  options: ActionTableOptions,
};

export default function ActionTable<T>(props: ActionTableProps<T>) {
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
      <Text key="title">
        <TitleIcon /> {props.title} ({props.data.length})
      </Text>,
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
