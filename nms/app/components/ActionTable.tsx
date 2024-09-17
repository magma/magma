/*
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
 */

import AddCircleOutlined from '@mui/icons-material/AddCircleOutlined';
import ArrowDownward from '@mui/icons-material/ArrowDownward';
import Autocomplete from '@mui/material/Autocomplete';
import Button from '@mui/material/Button';
import CardTitleRow from './layout/CardTitleRow';
import Check from '@mui/icons-material/Check';
import ChevronLeft from '@mui/icons-material/ChevronLeft';
import ChevronRight from '@mui/icons-material/ChevronRight';
import Clear from '@mui/icons-material/Clear';
import DeleteOutline from '@mui/icons-material/DeleteOutline';
import Edit from '@mui/icons-material/Edit';
import FilterList from '@mui/icons-material/FilterList';
import FirstPage from '@mui/icons-material/FirstPage';
import FormControl from '@mui/material/FormControl';
import LastPage from '@mui/icons-material/LastPage';
import MaterialTable, {MaterialTableProps, Query} from '@material-table/core';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import OutlinedInput from '@mui/material/OutlinedInput';
import Paper, {PaperProps} from '@mui/material/Paper';
import React from 'react';
import RefreshIcon from '@mui/icons-material/Refresh';
import Remove from '@mui/icons-material/Remove';
import SaveAlt from '@mui/icons-material/SaveAlt';
import Search from '@mui/icons-material/Search';
import Select from '@mui/material/Select';
import SvgIcon from '@mui/material/SvgIcon';
import TextField from '@mui/material/TextField';
import {forwardRef, useState} from 'react';
import {makeStyles} from '@mui/styles';

const useStyles = makeStyles({
  inputRoot: {
    '&.MuiOutlinedInput-root': {
      padding: 0,
    },
  },
});

const tableIcons = {
  Add: forwardRef<SVGSVGElement>((props, ref) => (
    <Button
      startIcon={<AddCircleOutlined {...props} ref={ref} />}
      component="span"
      variant="outlined"
      color="primary">
      {'Add New Row'}
    </Button>
  )),
  Check: forwardRef<SVGSVGElement>((props, ref) => (
    <Check {...props} ref={ref} />
  )),

  Clear: forwardRef<SVGSVGElement>((props, ref) => (
    <Clear {...props} ref={ref} />
  )),
  Edit: forwardRef<SVGSVGElement>((props, ref) => (
    <Edit {...props} ref={ref} />
  )),
  Delete: forwardRef<SVGSVGElement>((props, ref) => (
    <DeleteOutline {...props} ref={ref} />
  )),

  Export: forwardRef<SVGSVGElement>((props, ref) => (
    <SaveAlt {...props} ref={ref} />
  )),
  FirstPage: forwardRef<SVGSVGElement>((props, ref) => (
    <FirstPage {...props} ref={ref} />
  )),
  LastPage: forwardRef<SVGSVGElement>((props, ref) => (
    <LastPage {...props} ref={ref} />
  )),
  NextPage: forwardRef<SVGSVGElement>((props, ref) => (
    <ChevronRight {...props} ref={ref} />
  )),
  PreviousPage: forwardRef<SVGSVGElement>((props, ref) => (
    <ChevronLeft {...props} ref={ref} />
  )),
  ResetSearch: forwardRef<SVGSVGElement>((props, ref) => (
    <Clear {...props} ref={ref} />
  )),
  Retry: forwardRef<SVGSVGElement>((props, ref) => (
    <RefreshIcon {...props} ref={ref} />
  )),
  Search: forwardRef<SVGSVGElement>((props, ref) => (
    <Search {...props} ref={ref} />
  )),
  SortArrow: forwardRef<SVGSVGElement>((props, ref) => (
    <ArrowDownward {...props} ref={ref} />
  )),
  ThirdStateCheck: forwardRef<SVGSVGElement>((props, ref) => (
    <Remove {...props} ref={ref} />
  )),
  Filter: forwardRef<SVGSVGElement>((props, ref) => (
    <FilterList {...props} ref={ref} />
  )),
};

type ActionMenuItems = {
  name: string;
  handleFunc?: () => void | (() => Promise<void>);
};

export type ActionQuery = Query<any>;
export type TableRef = React.MutableRefObject<
  {onQueryChange: VoidFunction} | undefined
>;

export type ActionTableProps<T extends object> = {
  titleIcon?: typeof SvgIcon;
  tableRef?: TableRef;
  editable?: MaterialTableProps<T>['editable'];
  localization?: MaterialTableProps<T>['localization'];
  title?: string;
  handleCurrRow?: (currentRow: T) => void;
  columns: MaterialTableProps<T>['columns'];
  menuItems?: Array<ActionMenuItems>;
  actions?: MaterialTableProps<T>['actions'];
  data: MaterialTableProps<T>['data'];
  options: MaterialTableProps<T>['options'];
  detailPanel?: MaterialTableProps<T>['detailPanel'];
  onSelectionChange?: MaterialTableProps<T>['onSelectionChange'];
};

export function PaperComponent(props: PaperProps) {
  return <Paper {...props} elevation={0} />;
}

type SelectProps = {
  content: Array<string>;
  defaultValue?: string;
  value: string;
  onChange: (value: string) => void;
  testId?: string;
};

export function SelectEditComponent(props: SelectProps) {
  if (props.value === undefined || props.value === null) {
    if (props.defaultValue !== undefined) {
      props.onChange(props.defaultValue);
      return null;
    }
  }

  return (
    <FormControl>
      <Select
        data-testid={props.testId ?? ''}
        value={props.value}
        onChange={({target}) => props.onChange(target.value)}
        input={<OutlinedInput />}>
        {props.content.map((k: string, idx: number) => (
          <MenuItem key={idx} value={k}>
            {k}
          </MenuItem>
        ))}
      </Select>
    </FormControl>
  );
}

export function AutoCompleteEditComponent(props: SelectProps) {
  const classes = useStyles();

  return (
    <Autocomplete
      disableClearable
      options={props.content}
      freeSolo
      value={props.value}
      classes={{
        inputRoot: classes.inputRoot,
      }}
      onChange={(_, newValue) => {
        props.onChange(newValue);
      }}
      inputValue={props.value}
      onInputChange={(_, newInputValue) => {
        props.onChange(newInputValue);
      }}
      renderInput={params => <TextField {...params} variant="outlined" />}
    />
  );
}

export default function ActionTable<T extends object>(
  props: ActionTableProps<T>,
) {
  const actionTableJSX = [];
  const [anchorEl, setAnchorEl] = useState<(EventTarget & Element) | null>(
    null,
  );

  const handleClick = (event: React.MouseEvent, row: T | Array<T>) => {
    setAnchorEl(event.currentTarget);

    if (props.handleCurrRow) {
      props.handleCurrRow(row as T);
    }
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  if (props.titleIcon) {
    const TitleIcon = props.titleIcon;
    actionTableJSX.push(
      <CardTitleRow
        key="title"
        icon={TitleIcon}
        label={`${props.title || ''} (${props.data.length})`}
      />,
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
        onSelectionChange={props.onSelectionChange}
      />
    </>
  );
}
