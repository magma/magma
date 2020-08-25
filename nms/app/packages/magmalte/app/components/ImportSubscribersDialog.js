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
import type {Match} from 'react-router-dom';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {subscriber} from '@fbcnms/magma-api';

import Button from '@material-ui/core/Button';
import Chip from '@material-ui/core/Chip';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import DownloadIcon from '@material-ui/icons/GetApp';
import FormLabel from '@material-ui/core/FormLabel';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React from 'react';

import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {hexToBase64, isValidHex} from '@fbcnms/util/strings';
import {last} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';

const useStyles = makeStyles(() => ({
  content: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    height: '150px',
    borderStyle: 'dashed',
    borderColor: 'gray',
    padding: '10px',
  },
  error: {
    marginBottom: '10px',
  },
  icon: {
    marginLeft: '5px',
  },
}));

const IMSI = 'imsi';
const LTE_STATE = 'lte_state';
const LTE_AUTH_KEY = 'lte_auth_key';
const LTE_AUTH_OPC = 'lte_auth_opc';
const SUB_PROFILE = 'sub_profile';
const APN_LIST = 'active_apns';

const CSV_TEMPLATE_DATA: Array<Array<string>> = [
  [IMSI, LTE_STATE, LTE_AUTH_KEY, LTE_AUTH_OPC, SUB_PROFILE, APN_LIST],
  [
    '"""200056789012345"""',
    'ACTIVE',
    '20000000001234567890ABCDEFABCDEF',
    '21111111111234567890ABCDEFABCDEF',
    'low rate 1',
    'supernet;inet-1',
  ],
  [
    '"""200056789012346"""',
    'INACTIVE',
    '20000000001234567890ABCDEFABCDEF',
    '21111111111234567890ABCDEFABCDEF',
    'default',
  ],
];

const CSV_DOWNLOAD_LINK = `data:text/csv;charset=utf-8,${CSV_TEMPLATE_DATA.map(
  row => row.join(','),
).join(`\n`)}`;

const SUBSCRIBER_UPLOAD_LIMIT = 250;

type Props = WithAlert & {
  open: boolean,
  onClose: () => void,
  onSave: (successIDs: Array<string>) => void | Promise<void>,
  onSaveError: (failureIDs: Array<string>) => void,
};

export async function parseFileAndSave(
  csvText: null | string | ArrayBuffer,
  setErrorMsg: string => void,
  match: Match,
  props: {
    onSave: (successIDs: Array<string>) => void | Promise<void>,
    onSaveError: (failureIDs: Array<string>) => void,
  },
) {
  if (csvText == null) {
    setErrorMsg('Failed to get CSV, is it empty? Please see template.');
    return;
  } else if (typeof csvText != 'string') {
    setErrorMsg('Failed to get CSV, is it a binary? Please see template.');
    return;
  }
  // We only support windows style and the more popular unix style
  const newlineSeparator = csvText.includes('\r\n') ? '\r\n' : '\n';
  const rows: string[][] = csvText
    .split(newlineSeparator)
    .map(row => row.split(','));
  const columnNames = rows.shift(); // first row in csv are column names
  if (JSON.stringify(columnNames) !== JSON.stringify(CSV_TEMPLATE_DATA[0])) {
    setErrorMsg(
      'CSV column names are not properly formatted, please see template.',
    );
    return;
  }
  // Remove last line if it's blank
  const lastRow = last(rows);
  if (lastRow?.length === 1 && !lastRow?.[0]) {
    rows.pop();
  }
  if (rows.length > SUBSCRIBER_UPLOAD_LIMIT) {
    setErrorMsg(
      `Upload limit exceeded! Please limit the file to ${SUBSCRIBER_UPLOAD_LIMIT} subscribers.`,
    );
    return;
  }
  const subs = rows.map(row => getSubscriberFromRow(row));
  const validSubs = subs.filter(Boolean);
  if (validSubs.length < subs.length) {
    setErrorMsg('At least one row is incomplete, please fill in all fields.');
    return;
  }
  const failureIDs = [];
  const results = await Promise.all(
    validSubs.map(subscriber => {
      const {config: _, ...mutableSubscriber} = subscriber;
      return MagmaV1API.postLteByNetworkIdSubscribers({
        networkId: match.params.networkId || '',
        subscriber: mutableSubscriber,
      }).catch(e => {
        failureIDs.push(subscriber.id);
        return e;
      });
    }),
  );
  const successIDs = [];
  results.forEach(result => {
    if (!(result instanceof Error)) {
      successIDs.push(result);
    }
  });
  props.onSave(successIDs);
  if (failureIDs.length > 0) {
    props.onSaveError(failureIDs);
  }
}

function getSubscriberFromRow(row: Array<string>): ?subscriber {
  const data: {[string]: string} = CSV_TEMPLATE_DATA[0].reduce(
    (accumulator, colName, idx) => {
      // $FlowFixMe Set state for each field
      return {...accumulator, [colName]: row[idx]};
    },
    {},
  );
  if (!data[IMSI] || !data[LTE_AUTH_KEY]) {
    return null;
  }
  const state = data[LTE_STATE];
  if (state !== 'ACTIVE' && state !== 'INACTIVE') {
    return;
  }

  const lteValue = {
    state,
    auth_algo: 'MILENAGE', // default auth algo,
    auth_key: isValidHex(data[LTE_AUTH_KEY])
      ? hexToBase64(data[LTE_AUTH_KEY])
      : data[LTE_AUTH_KEY],
    auth_opc:
      data[LTE_AUTH_OPC] && isValidHex(data[LTE_AUTH_OPC])
        ? hexToBase64(data[LTE_AUTH_OPC])
        : data[LTE_AUTH_OPC],
    sub_profile: data[SUB_PROFILE] || 'default',
  };

  return {
    id: 'IMSI' + data[IMSI].replace(/"/g, ''), // strip surrounding quotes
    lte: lteValue,
    active_apns: data[APN_LIST]?.split(';'),
    config: {
      lte: lteValue,
    },
  };
}

function ImportSubscribersDialog(props: Props) {
  const [errorMsg, setErrorMsg] = useState();
  const [file, setFile] = useState();

  const classes = useStyles();
  const {match} = useRouter();

  async function onFileUpload() {
    const confirmed = await props.confirm('Upload file and add subscribers?');
    if (!confirmed) {
      return;
    }
    if (file) {
      const reader = new FileReader();
      reader.onload = () =>
        parseFileAndSave(reader.result, setErrorMsg, match, props);
      reader.readAsText(file);
    } else {
      setErrorMsg("Sorry, we couldn't find your file, please try again.");
    }
  }

  return (
    <Dialog open={props.open} onClose={props.onClose}>
      <DialogTitle>Upload Subscribers</DialogTitle>
      <DialogContent>
        <div className={classes.content}>
          {errorMsg && (
            <FormLabel className={classes.error} error>
              {errorMsg}
            </FormLabel>
          )}
          {file ? (
            <Chip
              label={file.name}
              onDelete={() => {
                setFile(null);
                setErrorMsg(null);
              }}
            />
          ) : (
            <Button variant="contained" component="label">
              Select File
              <input
                type="file"
                accept=".csv"
                style={{display: 'none'}}
                onChange={e => setFile(e.target.files[0])}
              />
            </Button>
          )}
        </div>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose} color="primary">
          Cancel
        </Button>
        <Button
          color="primary"
          variant="outlined"
          href={CSV_DOWNLOAD_LINK}
          download="MagmaSubscribersTemplate.csv">
          Get Subscribers Template
          <DownloadIcon className={classes.icon} />
        </Button>
        <Button
          onClick={onFileUpload}
          color="primary"
          variant="contained"
          disabled={file ? false : true}>
          Upload
        </Button>
      </DialogActions>
    </Dialog>
  );
}

export default withAlert(ImportSubscribersDialog);
