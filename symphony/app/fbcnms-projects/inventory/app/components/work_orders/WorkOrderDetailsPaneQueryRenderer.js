/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import InventoryQueryRenderer from '../InventoryQueryRenderer';
import React from 'react';
import WorkOrderDetailsPane from './WorkOrderDetailsPane';
import {graphql} from 'react-relay';

type Props = {
  workOrderId: string,
};

const workOrderDetailsQuery = graphql`
  query WorkOrderDetailsPaneQueryRendererQuery($workOrderId: ID!) {
    workOrder: node(id: $workOrderId) {
      ... on WorkOrder {
        ...WorkOrderDetailsPane_workOrder
      }
    }
  }
`;

class WorkOrderDetailsPaneQueryRenderer extends React.Component<Props> {
  render() {
    const {workOrderId} = this.props;
    return (
      <InventoryQueryRenderer
        query={workOrderDetailsQuery}
        variables={{workOrderId: workOrderId}}
        render={props => {
          return <WorkOrderDetailsPane workOrder={props.workOrder} />;
        }}
      />
    );
  }
}

export default WorkOrderDetailsPaneQueryRenderer;
