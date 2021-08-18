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

import type {gateway_device} from '@fbcnms/magma-api';

import React from 'react';
import TextField from '@material-ui/core/TextField';

import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
}));

type Props = {
  record: gateway_device,
};

export default function WifiDeviceHardwareFields(props: Props) {
  const classes = useStyles();
  return (
    <>
      <TextField
        label="HW ID"
        className={classes.input}
        value={props.record.hardware_id}
        disabled={true}
      />
    </>
  );
}
