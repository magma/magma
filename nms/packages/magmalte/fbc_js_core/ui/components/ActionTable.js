/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @flow
 * @format
 */

import typeof SvgIcon from '@material-ui/core/@@SvgIcon';

import AddCircleOutlined from '@material-ui/icons/AddCircleOutlined';
import ArrowDownward from '@material-ui/icons/ArrowDownward';
import Button from '@material-ui/core/Button';
import Check from '@material-ui/icons/Check';
import ChevronLeft from '@material-ui/icons/ChevronLeft';
import ChevronRight from '@material-ui/icons/ChevronRight';
import Clear from '@material-ui/icons/Clear';
import DeleteOutline from '@material-ui/icons/DeleteOutline';
import Edit from '@material-ui/icons/Edit';
import FilterList from '@material-ui/icons/FilterList';
import FirstPage from '@material-ui/icons/FirstPage';
import Grid from '@material-ui/core/Grid';
import LastPage from '@material-ui/icons/LastPage';
import MaterialTable from '@material-table/core';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import MoreVertIcon from '@material-ui/icons/MoreVert';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import RefreshIcon from '@material-ui/icons/Refresh';
import Remove from '@material-ui/icons/Remove';
import SaveAlt from '@material-ui/icons/SaveAlt';
import Search from '@material-ui/icons/Search';
import Text from './design-system/Text';

import {forwardRef} from 'react';
import {makeStyles} from '@material-ui/styles';
import {useState} from 'react';

const useStyles = makeStyles(theme => ({
  inputRoot: {
    '&.MuiOutlinedInput-root': {
      padding: 0,
    },
  },
  cardTitleRow: {
    marginBottom: theme.spacing(1),
    minHeight: '36px',
  },
  cardTitleIcon: {
    fill: '#545F77',
    marginRight: theme.spacing(1),
  },
}));

const tableIcons = {
  Add: forwardRef((props, ref) => (
    <Button
      startIcon={<AddCircleOutlined {...props} ref={ref} />}
      variant="outlined"
      color="primary">
      {'Add New Row'}
    </Button>
  )),
  Check: forwardRef((props, ref) => <Check {...props} ref={ref} />),

  Clear: forwardRef((props, ref) => <Clear {...props} ref={ref} />),
  Edit: forwardRef((props, ref) => <Edit {...props} ref={ref} />),
  Delete: forwardRef((props, ref) => <DeleteOutline {...props} ref={ref} />),

  Export: forwardRef((props, ref) => <SaveAlt {...props} ref={ref} />),
  FirstPage: forwardRef((props, ref) => <FirstPage {...props} ref={ref} />),
  LastPage: forwardRef((props, ref) => <LastPage {...props} ref={ref} />),
  NextPage: forwardRef((props, ref) => <ChevronRight {...props} ref={ref} />),
  PreviousPage: forwardRef((props, ref) => (
    <ChevronLeft {...props} ref={ref} />
  )),
  ResetSearch: forwardRef((props, ref) => <Clear {...props} ref={ref} />),
  Retry: forwardRef((props, ref) => <RefreshIcon {...props} ref={ref} />),
  Search: forwardRef((props, ref) => <Search {...props} ref={ref} />),
  SortArrow: forwardRef((props, ref) => <ArrowDownward {...props} ref={ref} />),
  ThirdStateCheck: forwardRef((props, ref) => <Remove {...props} ref={ref} />),
  Filter: forwardRef((props, ref) => <FilterList {...props} ref={ref} />),
};

type ActionMenuItems = {
  name: string,
  handleFunc?: () => void | (() => Promise<void>),
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
  // Order of actions column
  actionsColumnIndex: number,
  // Number of rows that would be rendered on every page
  pageSize?: number,
  // Page size options that could be selected by user
  pageSizeOptions: Array<number>,
  // Css style to be applied rows
  rowStyle?: {},
  // Header cell style for all headers
  headerStyle?: {},
  // Flag for showing toolbar
  toolbar?: boolean,
};

type ActionOrderType = {
  field: string,
  title: string,
  tableData: {},
};

type ActionFilter = {
  column: ActionTableColumn,
  value: string,
};

type ActionQuery = {
  filters: Array<ActionFilter>,
  orderBy: ActionOrderType,
  orderDirection: string,
  page: number,
  pageSize: number,
  search: string,
  totalCount: number,
};

type ActionTableProps<T> = {
  titleIcon?: SvgIcon,
  tableRef?: {},
  toolbar?: {},
  editable?: {},
  // Change/translate default texts of datatable (Eg: toolbar placeholder)
  localization?: {},
  // Table title
  title?: string,
  handleCurrRow?: T => void,
  columns: Array<ActionTableColumn>,
  // action list item
  menuItems?: Array<ActionMenuItems>,
  // Action list. An icon button will be rendered for each actions
  actions?: Array<{}>,
  // Data to be rendered
  data: Array<T> | (ActionQuery => {}),
  options: ActionTableOptions,
  // Component(s) to be rendered on detail panel
  detailPanel?: Array<{}>,
};

function PaperComponent(props: {}) {
  return <Paper {...props} elevation={0} />;
}

export default function ActionTable<T>(props: ActionTableProps<T>) {
  const actionTableJSX = [];
  const [anchorEl, setAnchorEl] = useState(null);
  const classes = useStyles();

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
      <Grid container alignItems="center" className={classes.cardTitleRow}>
        <Grid item xs>
          <Grid container alignItems="center">
            <TitleIcon className={classes.cardTitleIcon} />
            <Text variant="body1">{`${props.title || ''} (${
              props.data.length
            })`}</Text>
          </Grid>
        </Grid>
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
              handleClose();
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
        localization={props.localization}
        toolbar={props.toolbar}
        tableRef={props.tableRef}
        editable={props.editable}
        components={{
          Container: PaperComponent,
        }}
        title=""
        columns={props.columns}
        icons={tableIcons}
        data={props.data}
        actions={
          props.menuItems?.length
            ? [
                ...(props.actions ? props.actions : []),
                {
                  icon: () => <MoreVertIcon />,
                  tooltip: 'Actions',
                  onClick: handleClick,
                },
              ]
            : props.actions
        }
        options={props.options}
        detailPanel={props.detailPanel}
      />
    </>
  );
}
