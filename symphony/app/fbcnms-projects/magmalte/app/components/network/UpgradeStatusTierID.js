/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {WithStyles} from '@material-ui/core';
import type {tier} from '@fbcnms/magma-api';

import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import Select from '@material-ui/core/Select';

import {withStyles} from '@material-ui/core/styles';

const styles = {
  select: {
    width: 200,
  },
};

type Props = WithStyles<typeof styles> & {
  gatewayID: string,
  onChange: (gatewayID: string, newTierID: string) => Promise<void>,
  tierID: ?string,
  networkUpgradeTiers: ?(tier[]),
};

class UpgradeStatusTierID extends React.Component<Props> {
  handleChange = (newValue: SyntheticInputEvent<EventTarget>) => {
    this.props.onChange(this.props.gatewayID, newValue.target.value);
  };

  render() {
    const {classes, networkUpgradeTiers} = this.props;
    const tierID = this.props.tierID || '';

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
        onChange={this.handleChange}
        className={classes.select}>
        <MenuItem key={0} value="" disabled>
          <em>Not Specified</em>
        </MenuItem>
        {options}
      </Select>
    );
  }
}

export default withStyles(styles)(UpgradeStatusTierID);
