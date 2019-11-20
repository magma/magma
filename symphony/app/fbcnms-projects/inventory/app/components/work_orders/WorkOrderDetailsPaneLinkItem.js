/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import React from 'react';
import WorkOrderDetailsPaneItem from './WorkOrderDetailsPaneItem';
import nullthrows from '@fbcnms/util/nullthrows';
import {createFragmentContainer, graphql} from 'react-relay';
import {withStyles} from '@material-ui/core/styles';
import type {FutureState} from '../../common/WorkOrder';
import type {WithStyles} from '@material-ui/core';
import type {WorkOrderDetailsPaneLinkItem_link} from './__generated__/WorkOrderDetailsPaneLinkItem_link.graphql.js';

type Props = WithStyles<typeof styles> & {
  equipment: WorkOrderDetailsPaneLinkItem_link,
  futureState: FutureState,
};

const styles = theme => ({
  root: {
    backgroundColor: theme.palette.background.paper,
    minWidth: '200px',
  },
});

class WorkOrderDetailsPaneEquipmentItem extends React.Component<Props> {
  render() {
    const {equipment} = this.props;
    return (
      <WorkOrderDetailsPaneItem
        text={this._getLinkDescription(equipment, 'INSTALL')}
      />
    );
  }

  _getLinkDescription(
    link: WorkOrderDetailsPaneLinkItem_link,
    futureState: FutureState,
  ): string {
    const aSidePort = nullthrows(link.ports[0]);
    const zSidePort = nullthrows(link.ports[1]);
    const aSide = `${aSidePort.parentEquipment.name} - Port: ${aSidePort.definition.name}`;
    const zSide = `${zSidePort.parentEquipment.name} - Port: ${zSidePort.definition.name}`;
    return futureState === 'INSTALL'
      ? `Connect ${aSide} to ${zSide}`
      : `Disconnect ${aSide} from ${zSide}`;
  }
}

export default withStyles(styles)(
  createFragmentContainer(WorkOrderDetailsPaneEquipmentItem, {
    link: graphql`
      fragment WorkOrderDetailsPaneLinkItem_link on Link {
        ...EquipmentPortsTable_link @relay(mask: false)
      }
    `,
  }),
);
