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
import React, {useContext} from 'react';
import ReactJson from 'react-json-view';
import SubscriberContext from '../../components/context/SubscriberContext';
import {MutableSubscriber, Subscriber} from '../../../generated-ts';
import {SubscriberRowType} from '../../state/lte/SubscriberState';
import {isValidHex} from '../../util/strings';
import {useNavigate} from 'react-router-dom';
import type {
  PromqlReturnObject,
  SubscriberForbiddenNetworkTypesEnum,
} from '../../../generated-ts';

const mBIT = 1000000;
const kBIT = 1000;
export function getLabelUnit(val: number) {
  if (val > mBIT) {
    return [(val / mBIT).toFixed(2), 'mb'];
  } else if (val > kBIT) {
    return [(val / kBIT).toFixed(2), 'kb'];
  }

  return [val.toFixed(2), 'bytes'];
}

/**
 * Converts bits to megabits
 * @param {number} val The value in bits to be converted
 * @returns {string} Megabits value of the number passed in
 */
export function convertBitToMbit(val: number) {
  return (val / mBIT).toFixed(2);
}

export const CoreNetworkTypes = Object.freeze({
  NT_EPC: 'EPC',
  NT_5GC: '5GC',
});

export function getPromValue(resp: PromqlReturnObject) {
  const respArr = resp?.data?.result
    ?.map(item => {
      const value = item?.value?.[1];
      return value ? parseFloat(value) : undefined;
    })
    .filter(Boolean) as Array<number>;
  return respArr?.length ? respArr[0] : 0;
}

// default subscriber count in get subscriber query
export const DEFAULT_PAGE_SIZE = 25;

// susbcriber export colums title
export const SUBSCRIBER_EXPORT_COLUMNS = [
  {title: 'Name', field: 'name'},
  {title: 'IMSI', field: 'id'},
  {title: 'Auth Key', field: 'auth_key'},
  {title: 'Auth OPC', field: 'auth_opc'},
  {title: 'Service', field: 'state'},
  {title: 'Forbidden Network Types', field: 'forbidden_network_types'},
  {title: 'Data Plan', field: 'sub_profile'},
  {title: 'Active APNs', field: 'active_apns'},
];
export const SUBSCRIBER_ADD_ERRORS = Object.freeze({
  INVALID_IMSI:
    'The IMSI should be a string IMSI followed by a number with 10-15 digits',
  INVALID_AUTH_KEY:
    'Auth key is not a valid hex (example: 000102030405060708090A0B0C0D0E0F)',
  INVALID_AUTH_OPC:
    'Auth opc is not a valid hex (example: 000102030405060708090A0B0C0D0E0F)',
  REQUIRED_SUB_PROFILE: 'Please select a data plan',
  DUPLICATE_IMSI: 'The IMSI is duplicated',
  REQUIRED_AUTH_KEY: 'Auth key is required',
});
export const SUBSCRIBER_ACTION_TYPE = Object.freeze({
  ADD: 'add',
  EDIT: 'edit',
  DELETE: 'delete',
});
export type SubscriberActionType = typeof SUBSCRIBER_ACTION_TYPE[keyof typeof SUBSCRIBER_ACTION_TYPE];
export const REFRESH_TIMEOUT = 1000;

export type SubscriberInfo = {
  name: string;
  imsi: string;
  authKey: string;
  authOpc: string;
  state: 'INACTIVE' | 'ACTIVE';
  forbiddenNetworkTypes: Array<SubscriberForbiddenNetworkTypesEnum>;
  dataPlan: string;
  apns: Array<string>;
  policies?: Array<string>;
};
type SubscriberErrorKey = keyof typeof SUBSCRIBER_ADD_ERRORS;

/**
 * Checks subscriber fields format
 *
 * @param {Array<SubscriberInfo>} subscribers Array of subcribers to validate
 * @returns {Array<string>} Returns array of errors
 */
export function validateSubscribers(
  subscribers: Array<SubscriberInfo>,
  action: SubscriberActionType,
) {
  const errors: Record<string, Array<number>> = {};
  const imsiList: Array<string> = [];

  Object.keys(SUBSCRIBER_ADD_ERRORS).map(error => {
    const subscriberError = SUBSCRIBER_ADD_ERRORS[error as SubscriberErrorKey];
    errors[subscriberError] = [];
  });
  subscribers.forEach((info, i) => {
    if (!(action === 'delete')) {
      if (!info.authKey) {
        errors[SUBSCRIBER_ADD_ERRORS['REQUIRED_AUTH_KEY']].push(i + 1);
      }

      if (!info.dataPlan) {
        errors[SUBSCRIBER_ADD_ERRORS['REQUIRED_SUB_PROFILE']].push(i + 1);
      }

      if (imsiList.includes(info.imsi)) {
        errors[SUBSCRIBER_ADD_ERRORS['DUPLICATE_IMSI']].push(i + 1);
      }
    }

    if (!imsiList.includes(info.imsi)) {
      imsiList.push(info.imsi);
    }

    if (!info?.imsi?.match(/^(IMSI\d{10,15})$/)) {
      errors[SUBSCRIBER_ADD_ERRORS['INVALID_IMSI']].push(i + 1);
    }

    if (info.authKey && !isValidHex(info.authKey)) {
      errors[SUBSCRIBER_ADD_ERRORS['INVALID_AUTH_KEY']].push(i + 1);
    }

    if (info.authOpc && !isValidHex(info.authOpc)) {
      errors[SUBSCRIBER_ADD_ERRORS['INVALID_AUTH_OPC']].push(i + 1);
    }
  });

  const errorList: Array<string> = Object.keys(SUBSCRIBER_ADD_ERRORS)
    .map(error => SUBSCRIBER_ADD_ERRORS[error as SubscriberErrorKey])
    .reduce((res: Array<string>, errorMessage) => {
      if (errors[errorMessage].length > 0) {
        res.push(
          `${errorMessage} : Row ${errors[errorMessage].sort().join(', ')}`,
        );
      }

      return res;
    }, []);

  return errorList;
}

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
