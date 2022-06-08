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
import Dialog from '@material-ui/core/Dialog';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import Link from '@material-ui/core/Link';
import PeopleIcon from '@material-ui/icons/People';
import React from 'react';
import ReactJson from 'react-json-view';
import SettingsIcon from '@material-ui/icons/Settings';
// $FlowFixMe migrated to typescript
import SubscriberContext from '../../components/context/SubscriberContext';
import SubscriberDetail from './SubscriberDetail';
import SubscriberStateTable from './SubscriberStateTable';
import SubscriberTable from './SubscriberTable';
// $FlowFixMe migrated to typescript
import TopBar from '../../components/TopBar';
import type {
  mutable_subscriber,
  subscriber,
} from '../../../generated/MagmaAPIBindings';

import {Navigate, Route, Routes, useNavigate} from 'react-router-dom';
import {useContext} from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {SubscriberRowType} from '../../state/lte/SubscriberState';

const TITLE = 'Subscribers';

export default function SubscriberDashboard() {
  return (
    <Routes>
      <Route
        path="/overview/config/:subscriberId/*"
        element={<SubscriberDetail />}
      />
      <Route
        path="/overview/sessions/:subscriberId/*"
        element={<SubscriberDetail />}
      />

      <Route path="/overview/*" element={<SubscribersOverview />} />
      <Route index element={<Navigate to="overview" replace />} />
    </Routes>
  );
}

type Props = {
  open: boolean,
  onClose?: () => void,
  imsi: string,
};

export function SubscribersOverview() {
  return (
    <>
      <TopBar
        header={TITLE}
        tabs={[
          {
            label: 'Config',
            to: 'config',
            icon: SettingsIcon,
          },
          {
            label: 'Sessions',
            to: 'sessions',
            icon: PeopleIcon,
          },
        ]}
      />
      <Routes>
        <Route path="/config" element={<SubscriberTable />} />
        <Route path="/sessions" element={<SubscriberStateTable />} />
        <Route index element={<Navigate to="config" replace />} />
      </Routes>
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
  const navigate = useNavigate();
  const {subscriberConfig, currRow} = props;
  const imsi = currRow.imsi;
  return (
    <div>
      <Link
        variant="body2"
        component="button"
        onClick={() => navigate(imsi + `${!subscriberConfig ? '/event' : ''}`)}>
        {imsi}
      </Link>
    </div>
  );
}
