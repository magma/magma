/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {Equipment, EquipmentPort} from '../../common/Equipment';
import type {WithStyles} from '@material-ui/core';

import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  message: {
    margin: theme.spacing(2),
  },
});

type Props = {
  aSideEquipment: Equipment,
  aSidePort: EquipmentPort,
  zSideEquipment: Equipment,
  zSidePort: EquipmentPort,
} & WithStyles<typeof styles>;

class PortsConnectConfirmation extends React.Component<Props> {
  _formatPortEquipment(equipment: Equipment, port: EquipmentPort) {
    return `${port.definition.portType?.name || ''}
    ${port.definition.name} on
    ${equipment.equipmentType.name}
    ${equipment.name}`;
  }

  render() {
    const {aSideEquipment, aSidePort, zSideEquipment, zSidePort} = this.props;
    return (
      <Text className={this.props.classes.message} variant="body2">
        {`Are you sure you would like to connect port
        ${this._formatPortEquipment(aSideEquipment, aSidePort)}
        to port
        ${this._formatPortEquipment(zSideEquipment, zSidePort)}?`}
      </Text>
    );
  }
}

export default withStyles(styles)(PortsConnectConfirmation);
