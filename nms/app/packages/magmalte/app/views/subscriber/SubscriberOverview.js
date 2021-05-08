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
import AddSubscriberButton from './SubscriberAddDialog';
import Dialog from '@material-ui/core/Dialog';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
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

import {Redirect, Route, Switch} from 'react-router-dom';
import {useContext, useState} from 'react';
import {useRouter} from '@fbcnms/ui/hooks';
import type {mutable_subscriber, subscriber} from '@fbcnms/magma-api';

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
              <AddSubscriberButton
                onClose={() => {
                  setRefresh(!refresh);
                }}
              />
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
