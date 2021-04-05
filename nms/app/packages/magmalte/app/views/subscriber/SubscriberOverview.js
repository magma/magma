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
 * @flow strict-local
 * @format
 */
import type {ActionQuery} from '../../components/ActionTable';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {
  mutable_subscriber,
  subscriber,
  subscriber_state,
} from '@fbcnms/magma-api';

import ActionTable from '../../components/ActionTable';
import AddSubscriberButton from './SubscriberAddDialog';
import CardTitleRow from '../../components/layout/CardTitleRow';
import Dialog from '@material-ui/core/Dialog';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import Link from '@material-ui/core/Link';
import NetworkContext from '../../components/context/NetworkContext';
import PeopleIcon from '@material-ui/icons/People';
import React from 'react';
import ReactJson from 'react-json-view';
import SubscriberContext from '../../components/context/SubscriberContext';
import SubscriberDetail from './SubscriberDetail';
import TopBar from '../../components/TopBar';
import nullthrows from '@fbcnms/util/nullthrows';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';

import {FEG_LTE} from '@fbcnms/types/network';
import {Redirect, Route, Switch} from 'react-router-dom';
import {colors, typography} from '../../theme/default';
import {handleSubscriberQuery} from '../../state/lte/SubscriberState';
import {makeStyles} from '@material-ui/styles';
import {useContext, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

const TITLE = 'Subscribers';
const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(5),
  },
  appBarBtn: {
    color: colors.primary.white,
    background: colors.primary.comet,
    fontFamily: typography.button.fontFamily,
    fontWeight: typography.button.fontWeight,
    fontSize: typography.button.fontSize,
    lineHeight: typography.button.lineHeight,
    letterSpacing: typography.button.letterSpacing,

    '&:hover': {
      background: colors.primary.mirage,
    },
  },
  appBarBtnSecondary: {
    color: colors.primary.white,
  },
  cardTitleRow: {
    marginBottom: theme.spacing(1),
    minHeight: '36px',
  },
  cardTitleIcon: {
    fill: colors.primary.comet,
    marginRight: theme.spacing(1),
  },
}));

export default function SubscriberDashboard() {
  const {relativePath, relativeUrl} = useRouter();
  return (
    <Switch>
      <Route
        path={relativePath('/overview/:subscriberId')}
        component={SubscriberDetail}
      />

      <Route path={relativePath('/overview')} component={SubscriberOverview} />
      <Redirect to={relativeUrl('/overview')} />
    </Switch>
  );
}

export type SubscriberRowType = {
  name: string,
  imsi: string,
  activeApns?: string,
  ipAddresses?: string,
  activeSessions?: number,
  service: string,
  currentUsage: string,
  dailyAvg: string,
  lastReportedTime: Date | string,
};

type Props = {
  open: boolean,
  onClose?: () => void,
  imsi: string,
};

function JsonDialog(props: Props) {
  const ctx = useContext(SubscriberContext);
  const sessionState = ctx.sessionState[props.imsi] || {};
  const configuredSubscriberState = ctx.state[props.imsi];
  const subscriber: mutable_subscriber = {
    ...configuredSubscriberState,
    state: sessionState,
  };
  return (
    <Dialog open={props.open} onClose={props.onClose} fullWidth={true}>
      <DialogTitle>{props.imsi}</DialogTitle>
      <DialogContent>
        <ReactJson
          src={subscriber}
          enableClipboard={false}
          displayDataTypes={false}
        />
      </DialogContent>
    </Dialog>
  );
}

function SubscriberInternal(props: WithAlert) {
  const {history, match, relativeUrl} = useRouter();
  const [currRow, setCurrRow] = useState<SubscriberRowType>({});
  const networkId: string = nullthrows(match.params.networkId);
  const enqueueSnackbar = useEnqueueSnackbar();
  const ctx = useContext(SubscriberContext);
  const classes = useStyles();
  const state = ctx;
  const networkCtx = useContext(NetworkContext);
  // $FlowIgnore
  const subscriberMap: {[string]: subscriber} = state.state;
  // $FlowIgnore
  const sessionState: {[string]: subscriber_state} = state.sessionState;
  const subscriberMetrics = ctx.metrics;
  const [jsonDialog, setJsonDialog] = useState(false);
  const [pageSize, setPageSize] = useState(10);
  // first token (page 1) is an empty string
  const [tokenList, setTokenList] = useState(['']);
  const onClose = () => setJsonDialog(false);
  const tableRef = React.useRef();
  const tableColumns = [
    {title: 'Name', field: 'name'},
    {
      title: 'IMSI',
      field: 'imsi',
      render: currRow => {
        const subscriberConfig = subscriberMap[currRow.imsi];
        return (
          <Link
            variant="body2"
            component="button"
            onClick={() =>
              // Link to event tab if FEG_LTE network
              history.push(
                relativeUrl(
                  '/' +
                    currRow.imsi +
                    `${
                      networkCtx.networkType === FEG_LTE && !subscriberConfig
                        ? '/event'
                        : ''
                    }`,
                ),
              )
            }>
            {currRow.imsi}
          </Link>
        );
      },
    },
    {title: 'Service', field: 'service', width: 100},
    {title: 'Current Usage', field: 'currentUsage', width: 175},
    {title: 'Daily Average', field: 'dailyAvg', width: 175},
    {
      title: 'Last Reported Time',
      field: 'lastReportedTime',
      type: 'datetime',
      width: 200,
    },
  ];
  return (
    <>
      <TopBar
        header={TITLE}
        tabs={[
          {
            label: 'Subscribers',
            to: '/subscribersv2',
            icon: PeopleIcon,
            filters: <AddSubscriberButton onClose={() => {}} />,
          },
        ]}
      />
      <div className={classes.dashboardRoot}>
        <CardTitleRow key="title" icon={PeopleIcon} label={TITLE} />
        {subscriberMap || sessionState ? (
          <div>
            <JsonDialog
              open={jsonDialog}
              onClose={onClose}
              imsi={currRow.imsi}
            />
            <ActionTable
              tableRef={tableRef}
              data={(query: ActionQuery) => {
                return handleSubscriberQuery({
                  networkId,
                  query,
                  ctx,
                  tokenList,
                  setTokenList,
                  pageSize,
                  subscriberMetrics,
                });
              }}
              columns={tableColumns}
              handleCurrRow={(row: SubscriberRowType) => setCurrRow(row)}
              menuItems={
                networkCtx.networkType === FEG_LTE
                  ? [
                      {
                        name: 'View JSON',
                        handleFunc: () => {
                          setJsonDialog(true);
                        },
                      },
                    ]
                  : [
                      {
                        name: 'View JSON',
                        handleFunc: () => {
                          setJsonDialog(true);
                        },
                      },
                      {
                        name: 'View',
                        handleFunc: () => {
                          history.push(relativeUrl('/' + currRow.imsi));
                        },
                      },
                      {
                        name: 'Edit',
                        handleFunc: () => {
                          history.push(
                            relativeUrl('/' + currRow.imsi + '/config'),
                          );
                        },
                      },
                      {
                        name: 'Remove',
                        handleFunc: () => {
                          props
                            .confirm(
                              `Are you sure you want to delete ${currRow.imsi}?`,
                            )
                            .then(async confirmed => {
                              if (!confirmed) {
                                return;
                              }

                              try {
                                await ctx.setState?.(currRow.imsi);
                                tableRef.current?.onQueryChange();
                              } catch (e) {
                                enqueueSnackbar(
                                  'failed deleting subscriber ' + currRow.imsi,
                                  {
                                    variant: 'error',
                                  },
                                );
                              }
                            });
                        },
                      },
                    ]
              }
              onChangeRowsPerPage={newPageSize => {
                setPageSize(newPageSize);
                tableRef.current?.onQueryChange();
              }}
              options={{
                actionsColumnIndex: -1,
                pageSize: pageSize,
                pageSizeOptions: [10, 20],
                showFirstLastPageButtons: false,
              }}
            />
          </div>
        ) : (
          '<Text>No Subscribers Found</Text>'
        )}
      </div>
    </>
  );
}

const SubscriberOverview = withAlert(SubscriberInternal);
