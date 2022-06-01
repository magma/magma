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
 */

import * as React from 'react';
import AddBox from '@material-ui/icons/AddBox';
import ArrowDownward from '@material-ui/icons/ArrowDownward';
import Check from '@material-ui/icons/Check';
import ChevronLeft from '@material-ui/icons/ChevronLeft';
import ChevronRight from '@material-ui/icons/ChevronRight';
import Chip from '@material-ui/core/Chip';
import Clear from '@material-ui/icons/Clear';
import DeleteOutline from '@material-ui/icons/DeleteOutline';
import Edit from '@material-ui/icons/Edit';
import FilterList from '@material-ui/icons/FilterList';
import FirstPage from '@material-ui/icons/FirstPage';
import LastPage from '@material-ui/icons/LastPage';
import MaterialTable, {MaterialTableProps} from '@material-table/core';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import MoreVertIcon from '@material-ui/icons/MoreVert';
import RefreshIcon from '@material-ui/icons/Refresh';
import Remove from '@material-ui/icons/Remove';
import SaveAlt from '@material-ui/icons/SaveAlt';
import Search from '@material-ui/icons/Search';
import {colors} from '../../../../theme/default';
import {forwardRef} from 'react';
import {makeStyles} from '@material-ui/styles';
import {useState} from 'react';

const useStyles = makeStyles({
  labelChip: {
    backgroundColor: colors.primary.mercury,
    color: colors.primary.brightGray,
    margin: '5px',
  },
  ellipsisChip: {
    display: 'block',
    maxWidth: 256,
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
  },
});

type ActionMenuItems = {
  name: string;
  handleFunc?: () => any;
};

type ActionTableProps<T extends object> = {
  onRowClick?: (rowData: T) => void;
  columnStruct: MaterialTableProps<T>['columns'];
  menuItems?: Array<ActionMenuItems>;
  actions?: MaterialTableProps<T>['actions'];
  tableData: MaterialTableProps<T>['data'];
  dataTestId?: string;
};

const renderLabelValue = (labelValue: LabelVal) => {
  if (typeof labelValue === 'boolean') {
    return labelValue ? 'true' : 'false';
  }
  if (typeof labelValue === 'string' && labelValue.trim() === '') {
    return null;
  }
  return labelValue;
};

type CellProps<TValue> = {
  value: TValue;
};
type LabelVal = string | number | boolean;
type Labels = Record<string, LabelVal>;
export function LabelsCell({
  value,
}: CellProps<Labels> & {hideFields?: Array<string>}) {
  const classes = useStyles();
  const labels = value;
  return (
    <div>
      {Object.keys(labels).map(keyName => {
        const val = renderLabelValue(labels[keyName]);
        return (
          <Chip
            key={keyName}
            classes={{label: classes.ellipsisChip}}
            className={classes.labelChip}
            label={
              <span>
                <em>{keyName}</em>
                {val !== null && typeof val !== 'undefined' ? '=' : null}
                {val}
              </span>
            }
            size="small"
          />
        );
      })}
    </div>
  );
}
type GroupsList = Array<Labels>;

export function MultiGroupsCell({value}: CellProps<GroupsList>) {
  const classes = useStyles();
  return (
    <>
      {value.map((cellValue, idx) => (
        <div key={idx}>
          {Object.keys(cellValue).map(keyName => (
            <Chip
              key={keyName}
              classes={{label: classes.ellipsisChip}}
              className={classes.labelChip}
              label={
                <span>
                  <em>{keyName}</em>={renderLabelValue(cellValue[keyName])}
                </span>
              }
              size="small"
            />
          ))}
        </div>
      ))}
    </>
  );
}

const tableIcons = {
  Add: forwardRef<SVGSVGElement>((props, ref) => (
    <AddBox {...props} ref={ref} />
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

export default function SimpleTable<T extends object>(
  props: ActionTableProps<T>,
) {
  const {columnStruct, tableData, onRowClick} = props;
  const actionTableJSX = [];
  const [anchorEl, setAnchorEl] = useState<(EventTarget & Element) | null>(
    null,
  );
  const handleClick = (event: React.MouseEvent, row: T | Array<T>) => {
    setAnchorEl(event.currentTarget);
    if (onRowClick) {
      onRowClick(row as T);
    }
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  if (props.menuItems && anchorEl) {
    // Actions menu
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
        data-testid={props.dataTestId}
        columns={columnStruct}
        data={tableData || ([] as Array<T>)}
        icons={tableIcons}
        onRowClick={(event, rowData) =>
          onRowClick ? onRowClick(rowData!) : null
        }
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
        options={{
          actionsColumnIndex: -1,
          filtering: true,
          // hide table title and toolbar
          toolbar: false,
        }}
        localization={{
          // hide 'Actions' in table header
          header: {actions: ''},
        }}
      />
    </>
  );
}

export function toLabels(obj: Record<string, any>): Labels {
  if (!obj) {
    return {};
  }
  return Object.keys(obj).reduce((map, key) => {
    map[key] = obj[key] as string;
    return map;
  }, {} as Labels);
}
