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
 *
 * @flow
 * @format
 */

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {WithAlert} from '../../components/Alert/withAlert';
import type {call_trace} from '../../../generated/MagmaAPIBindings';

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import ActionTable from '../../components/ActionTable';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import CardTitleRow from '../../components/layout/CardTitleRow';
import CreateTraceButton from './TraceStartDialog';
import HistoryIcon from '@material-ui/icons/History';
import LineStyleIcon from '@material-ui/icons/LineStyle';
// $FlowFixMe migrated to typescript
import NetworkContext from '../../components/context/NetworkContext';
import React from 'react';
// $FlowFixMe migrated to typescript
import TopBar from '../../components/TopBar';
// $FlowFixMe migrated to typescript
import TraceContext from '../../components/context/TraceContext';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import withAlert from '../../components/Alert/withAlert';
import {makeStyles} from '@material-ui/styles';
import {useContext, useState} from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(5),
  },
}));

function TracingDashboard() {
  return (
    <>
      <TopBar
        header="Call Tracing"
        tabs={[
          {
            label: 'Call Traces',
            to: 'overview',
            icon: LineStyleIcon,
            filters: <div />,
          },
        ]}
      />
      <TracingTable />
    </>
  );
}

type TracingRowType = {
  traceID: string,
  state: 'COMPLETED' | 'ACTIVE',
  gatewayID: string,
  traceType: string,
};

function TracingTableRaw(_: WithAlert) {
  const [currRow, setCurrRow] = useState<TracingRowType>({});
  const classes = useStyles();
  const ctx = useContext(TraceContext);
  const {networkId} = useContext(NetworkContext);
  const enqueueSnackbar = useEnqueueSnackbar();
  const traceMap = ctx.state;
  const tableData = tracesToRows(traceMap);
  const tableColumns = [
    {title: 'Trace ID', field: 'traceID'},
    {title: 'State', field: 'state'},
    {title: 'Gateway ID', field: 'gatewayID'},
    {title: 'Trace Type', field: 'traceType'},
  ];

  const TraceFilter = () => {
    return <CreateTraceButton />;
  };

  const handleStop = async () => {
    if (currRow.state === 'COMPLETED') {
      enqueueSnackbar('Call trace ' + currRow.traceID + ' already stopped.', {
        variant: 'error',
      });
      return;
    }

    try {
      // $FlowFixMe[prop-missing]: Suppress type error, cannot refine type
      await ctx.setState?.(currRow.traceID, {
        requested_end: true,
      });
      enqueueSnackbar('Call trace ended successfully', {
        variant: 'success',
      });
    } catch (e) {
      const errMsg = e.response?.data?.message ?? e.message ?? e;
      enqueueSnackbar('Failed stopping call trace: ' + errMsg, {
        variant: 'error',
      });
    }
  };

  const handleDownload = async () => {
    if (currRow.state != 'COMPLETED') {
      enqueueSnackbar('Call trace ' + currRow.traceID + ' is still active', {
        variant: 'error',
      });
      return;
    }

    if (networkId) {
      // TODO(andreilee): Build download link based on generated API bindings
      window.location.href =
        '/nms/apicontroller/magma/v1/networks/' +
        networkId +
        '/tracing/' +
        currRow.traceID +
        '/download';
    }
  };

  return (
    <div className={classes.dashboardRoot}>
      <CardTitleRow
        icon={HistoryIcon}
        label={'Call Traces'}
        filter={TraceFilter}
      />
      <ActionTable
        data={tableData}
        columns={tableColumns}
        handleCurrRow={(row: TracingRowType) => setCurrRow(row)}
        menuItems={[
          {
            name: 'Download',
            // $FlowFixMe[incompatible-type]
            handleFunc: handleDownload,
          },
          {
            name: 'Stop',
            // $FlowFixMe[incompatible-type]
            handleFunc: handleStop,
          },
        ]}
        options={{
          actionsColumnIndex: -1,
          pageSize: 10,
          pageSizeOptions: [10, 20],
        }}
      />
    </div>
  );
}

function tracesToRows(traceMap: {[string]: call_trace}): Array<TracingRowType> {
  const rows = [];
  Object.keys(traceMap).map((traceID: string, _) => {
    const isTraceEnding: boolean = !!traceMap[traceID]?.state
      ?.call_trace_ending;

    rows.push({
      traceID: traceID,
      state: isTraceEnding ? 'COMPLETED' : 'ACTIVE',
      gatewayID: traceMap[traceID].config?.gateway_id || '',
      traceType: 'GATEWAY',
    });
  });
  return rows;
}

const TracingTable = withAlert(TracingTableRaw);

export default TracingDashboard;
