/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {tier} from '@fbcnms/magma-api';

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
  gatewayID: string,
  onChange: (gatewayID: string, newTierID: string) => Promise<void>,
  tierID: ?string,
  networkUpgradeTiers: ?(tier[]),
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
        props.onChange(props.gatewayID, target.value);
      }}
      className={classes.select}>
      <MenuItem key={0} value="" disabled>
        <em>Not Specified</em>
      </MenuItem>
      {options}
    </Select>
  );
}
