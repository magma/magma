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

import Dialog from '@material-ui/core/Dialog';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import Link from '@material-ui/core/Link';
import React from 'react';
import ReactJson from 'react-json-view';
import SubscriberContext from '../../components/context/SubscriberContext';
import {useContext} from 'react';
import {useNavigate} from 'react-router-dom';
import type {MutableSubscriber, Subscriber} from '../../../generated-ts';
import type {SubscriberActionType, SubscriberInfo} from './SubscriberUtils';
import type {SubscriberRowType} from '../../state/lte/SubscriberState';

export type subscriberStaticIpsRowType = {
  apnName: string;
  staticIp: string;
};

export type subscriberForbiddenNetworkTypes = {
  nwTypes: string;
};

export type EditSubscriberProps = {
  subscriberState: Subscriber;
  onSubscriberChange: (key: string, val: string | number | undefined) => void;
  inputClass: string;
  onTrafficPolicyChange: (
    key: string,
    val: string | number | undefined,
    index: number,
  ) => void;
  onDeleteApn: (apn: undefined) => void;
  onAddApnStaticIP: () => void;
  subProfiles: undefined;
  subscriberStaticIPRows: Array<subscriberStaticIpsRowType>;
  forbiddenNetworkTypes: Array<subscriberForbiddenNetworkTypes>;
  authKey: string;
  authOpc: string;
  setAuthKey: (key: string) => void;
  setAuthOpc: (key: string) => void;
};

export type SubscribersDialogDetailProps = {
  // Subscribers to add, edit or delete
  setSubscribers: (arg0: Array<SubscriberInfo>) => void;
  subscribers: Array<SubscriberInfo>;
  // Formatting error (eg: field missing, wrong IMSI format)
  setAddError: (arg0: Array<string>) => void;
  addError: Array<string>;
  // Display dropzone if set to true
  setUpload: (arg0: boolean) => void;
  upload: boolean;
  onClose: () => void;
  // Add, edit or delete subscribers
  onSave: (
    arg0: Array<SubscriberInfo>,
    selectedSubscribers?: Array<string>,
  ) => void;
  error?: string;
  // Row added with the Add New Row button
  rowAdd: boolean;
  setRowAdd: (arg0: boolean) => void;
  // Delete, Edit or Add subscriber
  subscriberAction: SubscriberActionType;
};

type JsonProps = {
  open: boolean;
  onClose?: () => void;
  imsi: string;
};

export function JsonDialog(props: JsonProps) {
  const ctx = useContext(SubscriberContext);
  const sessionState = ctx.sessionState[props.imsi] || {};
  const configuredSubscriberState = ctx.state[props.imsi];
  // TODO[ts-migration] The state composition needs attention in more detail
  const subscriber: MutableSubscriber = {
    ...configuredSubscriberState,
    state: sessionState,
  } as MutableSubscriber;
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
  subscriberConfig: Subscriber;
  currRow: SubscriberRowType;
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
