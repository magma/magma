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

import Checkbox from '@material-ui/core/Checkbox';
import DialogContent from '@material-ui/core/DialogContent';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableContainer from '@material-ui/core/TableContainer';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Text from '../../theme/design-system/Text';

import type {policy_rule, subscriber} from '@fbcnms/magma-api';

type Props = {
  policyRule: policy_rule,
  onChange: policy_rule => void,
  subscribers: {[string]: subscriber},
};

export default function PolicySubscribersEdit(props: Props) {
  const {policyRule, subscribers} = props;
  const allIMSI = Object.keys(subscribers);
  const assignedSubscribers = new Set(policyRule.assigned_subscribers ?? []);
  const rows: {[string]: boolean} = {};
  Object.keys(subscribers).forEach((imsi: string) => {
    rows[imsi] = assignedSubscribers.has(imsi);
  });

  const handleChange = (imsi: string, checked: boolean) => {
    if (checked) {
      props.onChange({
        ...policyRule,
        assigned_subscribers: Array.from(assignedSubscribers).concat([imsi]),
      });
    } else {
      assignedSubscribers.delete(imsi);
      props.onChange({
        ...policyRule,
        assigned_subscribers: Array.from(assignedSubscribers),
      });
    }
  };

  return (
    <>
      <DialogContent data-testid="policySubscribersEdit">
        <List>
          <Text
            weight="medium"
            variant="subtitle1"
            display="block"
            gutterBottom>
            {'Assigned subscribers'}
          </Text>
          <ListItem dense disableGutters />
          <TableContainer component={Paper}>
            <Table size="small" aria-label="a dense table">
              <TableHead>
                <TableRow>
                  <TableCell>{'Assigned'}</TableCell>
                  <TableCell align="right">{'IMSI'}</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {allIMSI.map((imsi: string) => (
                  <TableRow key={imsi}>
                    <TableCell component="th" scope="row">
                      <Checkbox
                        checked={assignedSubscribers.has(imsi)}
                        onChange={({target}) =>
                          handleChange(imsi, target.checked)
                        }
                      />
                    </TableCell>
                    <TableCell align="right">{imsi}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        </List>
      </DialogContent>
    </>
  );
}
