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
 */

import type {Tier} from '../../../generated-ts';

import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import Select from '@material-ui/core/Select';

import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  select: {
    width: 200,
  },
}));

type Props = {
  gatewayID: string;
  onChange: (gatewayID: string, newTierID: string) => Promise<void>;
  tierID: string | null | undefined;
  networkUpgradeTiers: Array<Tier> | null | undefined;
};

export default function UpgradeStatusTierID(props: Props) {
  const classes = useStyles();
  const {networkUpgradeTiers} = props;
  const tierID = props.tierID || '';

  let options;
  if (!networkUpgradeTiers) {
    options = [
      <MenuItem key={1} value="default" disabled>
        <em>Default</em>
      </MenuItem>,
    ];
  } else {
    options = networkUpgradeTiers.map((data, i) => (
      <MenuItem key={i + 1} value={data.id}>
        {data.name}
      </MenuItem>
    ));
  }

  return (
    <Select
      value={tierID}
      onChange={({target}) => {
        void props.onChange(props.gatewayID, target.value as string);
      }}
      className={classes.select}>
      <MenuItem key={0} value="" disabled>
        <em>Not Specified</em>
      </MenuItem>
      {options}
    </Select>
  );
}
