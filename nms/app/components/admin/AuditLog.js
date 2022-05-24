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
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import ActionTable from '../ActionTable';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import CardTitleRow from '../../components/layout/CardTitleRow';
import DeviceStatusCircle from '../../theme/design-system/DeviceStatusCircle';
import Dialog from '@material-ui/core/Dialog';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import Grid from '@material-ui/core/Grid';
import ListAltIcon from '@material-ui/icons/ListAlt';
// $FlowFixMe migrated to typescript
import LoadingFiller from '../LoadingFiller';
import React from 'react';
import ReactJson from 'react-json-view';

import {makeStyles} from '@material-ui/styles';
import {useAxios} from '../../../app/hooks';
import {useState} from 'react';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(5),
  },
}));

export type AuditLogRowType = {
  id: string,
  updatedAt: Date,
  ipAddress: string,
  url: string,
  actingUserEmail: string,
  mutationType: string,
  objectType: string,
  objectId: string,
  mutationData: {},
};

const DEFAULT_PAGE_SIZE = 25;

/**
 * AuditLog functional component to display audit logs in a tabular format
 */
function AuditLog() {
  const classes = useStyles();
  const {response, error, isLoading} = useAxios({
    url: '/admin/auditlog/async',
    method: 'get',
  });
  const [currRow, setCurrRow] = useState<AuditLogRowType>({});
  const onClose = () => setJsonDialog(false);
  const [jsonDialog, setJsonDialog] = useState(false);

  if (error || isLoading || !response || !response.data) {
    return <LoadingFiller />;
  }
  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={4}>
        <Grid item xs={12}>
          <CardTitleRow
            icon={ListAltIcon}
            label={`Audit Logs (${response?.data?.length ?? 0})`}
          />
        </Grid>
        <Grid item xs={12}>
          <JsonDialog open={jsonDialog} onClose={onClose} row={currRow} />
          <ActionTable
            data={response.data}
            columns={[
              {title: 'Time', field: 'updatedAt', type: 'datetime'},
              {title: 'IP Address', field: 'ipAddress'},
              {title: 'User', field: 'actingUserEmail'},
              {title: 'Action', field: 'mutationType'},
              {title: 'Object Type', field: 'objectType'},
              {title: 'Object ID', field: 'objectId'},
              {
                title: 'Status',
                field: 'status',
                render: currRow => (
                  <>
                    <DeviceStatusCircle
                      isGrey={false}
                      isActive={currRow.status === 'SUCCESS'}
                    />
                    {currRow.status}
                  </>
                ),
              },
            ]}
            options={{
              actionsColumnIndex: -1,
              pageSize: DEFAULT_PAGE_SIZE,
              pageSizeOptions: [DEFAULT_PAGE_SIZE],
            }}
            handleCurrRow={row => setCurrRow(row)}
            menuItems={[
              {
                name: 'View JSON',
                handleFunc: () => {
                  setJsonDialog(true);
                },
              },
            ]}
          />
        </Grid>
      </Grid>
    </div>
  );
}

type DialogProps = {
  row: AuditLogRowType,
  open: boolean,
  onClose: () => void,
};

/**
 * JSONDialog functional component is used to display the audit log data
 * in a JSON view
 * @param {*} props
 */
function JsonDialog(props: DialogProps) {
  return (
    <Dialog
      open={props.open}
      onClose={props.onClose}
      maxWidth="lg"
      fullWidth={true}>
      <DialogTitle>{props.row.url}</DialogTitle>
      <DialogContent>
        <ReactJson
          src={props.row.mutationData}
          enableClipboard={false}
          displayDataTypes={false}
        />
      </DialogContent>
    </Dialog>
  );
}

export default AuditLog;
