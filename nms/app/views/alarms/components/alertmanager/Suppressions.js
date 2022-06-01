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
 * @flow strict-local
 * @format
 */
import CircularProgress from '@material-ui/core/CircularProgress';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import SimpleTable, {MultiGroupsCell, toLabels} from '../table/SimpleTable';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import TableActionDialog from '../table/TableActionDialog';
import {makeStyles} from '@material-ui/styles';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useAlarmContext} from '../AlarmContext';
import {useNetworkId} from '../hooks';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useSnackbars} from '../../../../hooks/useSnackbar';

import {useState} from 'react';

const useStyles = makeStyles(theme => ({
  addButton: {
    position: 'fixed',
    bottom: 0,
    right: 0,
    margin: theme.spacing(2),
  },
  loading: {
    display: 'flex',
    height: '100%',
    alignItems: 'center',
    justifyContent: 'center',
  },
}));

export default function Suppressions() {
  const {apiUtil} = useAlarmContext();
  const [menuAnchorEl, setMenuAnchorEl] = useState<?HTMLElement>(null);
  const [currentRow, setCurrentRow] = useState<{}>({});
  const [showDialog, setShowDialog] = useState(false);
  const [lastRefreshTime, _setLastRefreshTime] = useState<string>(
    new Date().toLocaleString(),
  );
  const [_isAddEditAlert, _setIsAddEditAlert] = useState<boolean>(false);
  const classes = useStyles();
  const snackbars = useSnackbars();
  const networkId = useNetworkId();
  const {isLoading, error, response} = apiUtil.useAlarmsApi(
    apiUtil.getSuppressions,
    {networkId},
    lastRefreshTime,
  );

  if (error) {
    snackbars.error(
      `Unable to load suppressions: ${
        error.response ? error.response.data.message : error.message
      }`,
    );
  }

  const silencesList = response || [];

  return (
    <>
      <SimpleTable
        onRowClick={row => setCurrentRow(row)}
        columnStruct={[
          {title: 'Name', field: 'comment'},
          {title: 'Active', field: 'status.state'},
          {title: 'Created By', field: 'createdBy'},
          {
            title: 'Matchers',
            field: 'matchers',
            render: row => {
              const value = row.matchers
                ? row.matchers.map(matcher => toLabels(matcher))
                : [];
              return <MultiGroupsCell value={value} />;
            },
          },
        ]}
        tableData={silencesList || []}
        dataTestId="suppressions"
        menuItems={[
          {
            name: 'View',
            handleFunc: () => setShowDialog(true),
          },
        ]}
      />
      {isLoading && silencesList.length === 0 && (
        <div className={classes.loading}>
          <CircularProgress />
        </div>
      )}
      <Menu
        anchorEl={menuAnchorEl}
        keepMounted
        open={Boolean(menuAnchorEl)}
        onClose={() => setMenuAnchorEl(null)}>
        <MenuItem onClick={() => setShowDialog(true)}>View</MenuItem>
      </Menu>
      <TableActionDialog
        open={showDialog}
        onClose={() => setShowDialog(false)}
        title={'View Suppression'}
        row={currentRow || {}}
        showCopyButton={true}
        showDeleteButton={false}
      />
    </>
  );
}
