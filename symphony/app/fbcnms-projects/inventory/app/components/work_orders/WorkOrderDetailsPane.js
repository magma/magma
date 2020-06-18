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
import WorkOrderDetailsPaneEquipmentItem from './WorkOrderDetailsPaneEquipmentItem';
import WorkOrderDetailsPaneItem from './WorkOrderDetailsPaneItem';
import WorkOrderDetailsPaneLinkItem from './WorkOrderDetailsPaneLinkItem';
import {createFragmentContainer, graphql} from 'react-relay';
import {withStyles} from '@material-ui/core/styles';
import type {WithStyles} from '@material-ui/core';
import type {WorkOrderDetailsPane_workOrder} from './__generated__/WorkOrderDetailsPane_workOrder.graphql.js';

type Props = WithStyles<typeof styles> & {
  workOrder: WorkOrderDetailsPane_workOrder,
};

const styles = theme => ({
  root: {
    backgroundColor: theme.palette.background.paper,
    minWidth: '200px',
  },
});

class WorkOrderDetailsPane extends React.Component<Props> {
  render() {
    const {classes} = this.props;
    const workOrder = this.props.workOrder;
    if (
      workOrder.equipmentToAdd.length == 0 &&
      workOrder.equipmentToRemove.length == 0 &&
      workOrder.linksToAdd.length == 0 &&
      workOrder.linksToRemove.length == 0
    ) {
      return (
        <div dense="true" className={classes.root}>
          <WorkOrderDetailsPaneItem
            key="placeholder"
            text="No items in this work order"
          />
        </div>
      );
    }

    return (
      <div className={classes.root}>
        {workOrder.equipmentToAdd.filter(Boolean).map(equipment => (
          <WorkOrderDetailsPaneEquipmentItem
            key={`add_${equipment.id}`}
            equipment={equipment}
            futureState="INSTALL"
          />
        ))}
        {workOrder.equipmentToRemove.filter(Boolean).map(equipment => (
          <WorkOrderDetailsPaneEquipmentItem
            key={`remove_${equipment.id}`}
            equipment={equipment}
            futureState="REMOVE"
          />
        ))}
        {workOrder.linksToAdd.filter(Boolean).map(link => (
          <WorkOrderDetailsPaneLinkItem
            key={`connect_${link.id}`}
            link={link}
            futureState="INSTALL"
          />
        ))}
        {workOrder.linksToRemove.filter(Boolean).map(link => (
          <WorkOrderDetailsPaneLinkItem
            key={`disconnect_${link.id}`}
            link={link}
            futureState="REMOVE"
          />
        ))}
      </div>
    );
  }
}

export default withStyles(styles)(
  createFragmentContainer(WorkOrderDetailsPane, {
    workOrder: graphql`
      fragment WorkOrderDetailsPane_workOrder on WorkOrder {
        id
        name
        equipmentToAdd {
          id
          ...WorkOrderDetailsPaneEquipmentItem_equipment
        }
        equipmentToRemove {
          id
          ...WorkOrderDetailsPaneEquipmentItem_equipment
        }
        linksToAdd {
          id
          ...WorkOrderDetailsPaneLinkItem_link
        }
        linksToRemove {
          id
          ...WorkOrderDetailsPaneLinkItem_link
        }
      }
    `,
  }),
);
