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
import CircularProgress from '@mui/material/CircularProgress';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import React from 'react';
import SimpleTable, {LabelsCell, toLabels} from '../table/SimpleTable';
import TableActionDialog from '../table/TableActionDialog';
import {getErrorMessage} from '../../../../util/ErrorUtils';
import {makeStyles} from '@mui/styles';
import {useAlarmContext} from '../AlarmContext';
import {useNetworkId} from '../hooks';
import {useSnackbars} from '../../../../hooks/useSnackbar';
import {useState} from 'react';

const useStyles = makeStyles(() => ({
  loading: {
    display: 'flex',
    height: '100%',
    alignItems: 'center',
    justifyContent: 'center',
  },
}));
export default function Routes() {
  const {apiUtil} = useAlarmContext();
  const [menuAnchorEl, setMenuAnchorEl] = useState<HTMLElement | null>(null);
  const [currentRow, setCurrentRow] = useState({});
  const [showDialog, setShowDialog] = useState<'view' | null>(null);
  const [lastRefreshTime] = useState<string>(new Date().toLocaleString());
  const classes = useStyles();
  const snackbars = useSnackbars();

  const onDialogAction = (args: 'view' | null) => {
    setShowDialog(args);
    setMenuAnchorEl(null);
  };

  const networkId = useNetworkId();
  const {isLoading, error, response} = apiUtil.useAlarmsApi(
    apiUtil.getRouteTree,
    {networkId},
    lastRefreshTime,
  );

  if (error) {
    snackbars.error(`Unable to load receivers: ${getErrorMessage(error)}`);
  }

  const routesList = response?.routes || [];

  return (
    <>
      <SimpleTable
        onRowClick={row => setCurrentRow(row)}
        columnStruct={[
          {title: 'Name', field: 'receiver'},
          {
            title: 'Group By',
            field: 'group_by',
            render: row => row.group_by?.join(','),
          },
          {
            title: 'Match',
            field: 'match',
            render: row => {
              const labels = toLabels(row.match);
              return <LabelsCell value={labels} />;
            },
          },
        ]}
        tableData={routesList || []}
        dataTestId="routes"
        menuItems={[
          {
            name: 'View',
            handleFunc: () => onDialogAction('view'),
          },
        ]}
      />
      {isLoading && routesList.length === 0 && (
        <div className={classes.loading}>
          <CircularProgress />
        </div>
      )}
      <Menu
        anchorEl={menuAnchorEl}
        keepMounted
        open={Boolean(menuAnchorEl)}
        onClose={() => setMenuAnchorEl(null)}>
        <MenuItem onClick={() => onDialogAction('view')}>View</MenuItem>
      </Menu>
      <TableActionDialog
        open={showDialog != null}
        onClose={() => onDialogAction(null)}
        title={'View Alert'}
        row={currentRow || {}}
        showCopyButton={true}
        showDeleteButton={false}
      />
    </>
  );
}
