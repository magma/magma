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
import AddSubscriberButton, {
  BulkEditSubscriberButton,
  validateSubscriberInfo,
} from './SubscriberAddDialog';
import Dialog from '@material-ui/core/Dialog';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import Grid from '@material-ui/core/Grid';
import Link from '@material-ui/core/Link';
import PeopleIcon from '@material-ui/icons/People';
import React from 'react';
import ReactJson from 'react-json-view';
import SettingsIcon from '@material-ui/icons/Settings';
import SubscriberContext from '../../components/context/SubscriberContext';
import SubscriberDetail from './SubscriberDetail';
import SubscriberStateTable from './SubscriberStateTable';
import SubscriberTable from './SubscriberTable';
import TopBar from '../../components/TopBar';
import type {SubscriberInfo} from './SubscriberAddDialog';
import type {
  mutable_subscriber,
  mutable_subscribers,
  subscriber,
} from '../../../generated/MagmaAPIBindings';

import {Redirect, Route, Switch} from 'react-router-dom';
import {hexToBase64, isValidHex} from '@fbcnms/util/strings';
import {useContext, useRef, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

const TITLE = 'Subscribers';

export default function SubscriberDashboard() {
  const {relativePath, relativeUrl} = useRouter();
  return (
    <Switch>
      <Route
        path={relativePath('/overview/config/:subscriberId')}
        component={SubscriberDetail}
      />
      <Route
        path={relativePath('/overview/sessions/:subscriberId')}
        component={SubscriberDetail}
      />

      <Route path={relativePath('/overview')} component={SubscribersOverview} />
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
type refreshProps = {
  refreshContext: () => void,
};
const SUBSCRIBERS_CHUNK_SIZE = 1000;

function SubscriberActions(props: refreshProps) {
  const [_error, setError] = useState('');
  const enqueueSnackbar = useEnqueueSnackbar();
  const ctx = useContext(SubscriberContext);
  const successCountRef = useRef(0);
  const [open, setOpen] = useState(false);

  const saveSubscribers = async (
    subscribers: Array<SubscriberInfo>,
    edit: boolean,
  ) => {
    let addedSubscribers = [];
    let subscriberErrors = '';
    for (const [i, subscriber] of subscribers.entries()) {
      try {
        const err = validateSubscriberInfo(subscriber, ctx.state);
        if (err.length > 0) {
          throw err;
        }
        const authKey =
          subscriber.authKey && isValidHex(subscriber.authKey)
            ? hexToBase64(subscriber.authKey)
            : '';

        const authOpc =
          subscriber.authOpc !== undefined && isValidHex(subscriber.authOpc)
            ? hexToBase64(subscriber.authOpc)
            : '';
        const newSubscriber = {
          active_apns: subscriber.apns,
          active_policies: subscriber.policies,
          id: subscriber.imsi,
          name: subscriber.name,
          lte: {
            auth_algo: 'MILENAGE',
            auth_key: authKey,
            auth_opc: authOpc,
            state: subscriber.state,
            sub_profile: subscriber.dataPlan,
          },
        };
        if (edit) {
          ctx.setState?.(subscriber.imsi, newSubscriber);
        } else {
          addedSubscribers.push(newSubscriber);
          // bulk add chunked subscribers
          if (
            addedSubscribers.length == SUBSCRIBERS_CHUNK_SIZE ||
            i == subscribers.length - 1
          ) {
            const success = await bulkAdd(addedSubscribers, subscriberErrors);
            if (success) {
              successCountRef.current =
                successCountRef.current + addedSubscribers.length;
              addedSubscribers = [];
            } else {
              enqueueSnackbar('Saving subscribers to the api failed: ', {
                variant: 'error',
              });
              return;
            }
          }
        }
      } catch (e) {
        const errMsg = e.response?.data?.message ?? e.message ?? e;
        subscriberErrors +=
          'error saving ' + subscriber.imsi + ' : ' + errMsg + '\n';
        //report saved errors if we reach end of loop without calling bulkadd.
        if (i == subscribers.length - 1) {
          setError(subscriberErrors);
          enqueueSnackbar('Saving subscribers to the api failed: ', {
            variant: 'error',
          });
          return;
        }
      }
    }
    enqueueSnackbar(
      ` Subscriber${successCountRef.current > 1 ? 's' : ''} saved successfully`,
      {
        variant: 'success',
      },
    );
    setOpen(false);
    props.refreshContext();
  };
  const bulkAdd = async (
    addedSubscribers: mutable_subscribers,
    subscriberErrors: string,
  ) => {
    let success = true;
    try {
      if (subscriberErrors.length > 0) {
        setError(subscriberErrors);
        return false;
      }
      await ctx.setState?.('', addedSubscribers);
      return success;
    } catch (e) {
      const errMsg = e.response?.data?.message ?? e.message ?? e;
      setError('error saving subscribers: ' + errMsg);
      success = false;
      return success;
    }
  };
  return (
    <Grid container alignItems="center" spacing={1}>
      <Grid item>
        <BulkEditSubscriberButton
          onClose={() => {
            props.refreshContext();
          }}
          onSave={subscribers => {
            saveSubscribers(subscribers, true);
          }}
        />
      </Grid>
      <Grid item>
        <AddSubscriberButton
          onClose={() => {
            props.refreshContext();
          }}
          onSave={subscribers => {
            saveSubscribers(subscribers, false);
          }}
          isOpen={open}
          handleOpen={isOpen => {
            setOpen(isOpen);
          }}
        />
      </Grid>
    </Grid>
  );
}

export function SubscribersOverview() {
  const {relativePath, relativeUrl} = useRouter();
  const [refresh, setRefresh] = useState(false);

  return (
    <>
      <TopBar
        header={TITLE}
        tabs={[
          {
            label: 'Config',
            to: '/config',
            icon: SettingsIcon,
            filters: (
              <SubscriberActions refreshContext={() => setRefresh(!refresh)} />
            ),
          },
          {
            label: 'Sessions',
            to: '/sessions',
            icon: PeopleIcon,
          },
        ]}
      />
      <Switch>
        <Route
          path={relativePath('/config')}
          component={() => <SubscriberTable refresh />}
        />
        <Route
          path={relativePath('/sessions')}
          component={SubscriberStateTable}
        />
        <Redirect to={relativeUrl('/config')} />
      </Switch>
    </>
  );
}

export function JsonDialog(props: Props) {
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

type RenderLinkType = {
  subscriberConfig: subscriber,
  currRow: SubscriberRowType,
};

export function RenderLink(props: RenderLinkType) {
  const {relativeUrl, history} = useRouter();
  const {subscriberConfig, currRow} = props;
  const imsi = currRow.imsi;
  return (
    <div>
      <Link
        variant="body2"
        component="button"
        onClick={() =>
          history.push(
            relativeUrl('/' + imsi + `${!subscriberConfig ? '/event' : ''}`),
          )
        }>
        {imsi}
      </Link>
    </div>
  );
}
