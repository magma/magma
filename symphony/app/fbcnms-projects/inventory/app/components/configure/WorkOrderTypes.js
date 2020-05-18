/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {
  WorkOrderTypesQuery,
  WorkOrderTypesQueryResponse,
} from './__generated__/WorkOrderTypesQuery.graphql';

import AddEditWorkOrderTypeCard from './AddEditWorkOrderTypeCard';
import Button from '@fbcnms/ui/components/design-system/Button';
import InventoryView from '../InventoryViewContainer';
import React, {useMemo, useState} from 'react';
import Table from '@fbcnms/ui/components/design-system/Table/Table';
import fbt from 'fbt';
import withInventoryErrorBoundary from '../../common/withInventoryErrorBoundary';
import {ButtonAction} from '@fbcnms/ui/components/design-system/View/ViewHeaderActions';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {TABLE_SORT_ORDER} from '@fbcnms/ui/components/design-system/Table/TableContext';
import {graphql} from 'relay-runtime';
import {makeStyles} from '@material-ui/styles';
import {useLazyLoadQuery} from 'react-relay/hooks';

const useStyles = makeStyles(() => ({
  paper: {
    flexGrow: 1,
    overflowY: 'hidden',
  },
}));

const workOrderTypesQuery = graphql`
  query WorkOrderTypesQuery {
    workOrderTypes(first: 500) @connection(key: "Configure_workOrderTypes") {
      edges {
        node {
          id
          name
          description
          ...AddEditWorkOrderTypeCard_workOrderType
        }
      }
    }
  }
`;

type WorkOrderTypeEdge = $ElementType<
  $ElementType<
    $NonMaybeType<$ElementType<WorkOrderTypesQueryResponse, 'workOrderTypes'>>,
    'edges',
  >,
  number,
>;

type WorkOrderTypeNode = $NonMaybeType<$ElementType<WorkOrderTypeEdge, 'node'>>;

const WorkOrderTypes = () => {
  const classes = useStyles();
  const {
    workOrderTypes,
  }: WorkOrderTypesQueryResponse = useLazyLoadQuery<WorkOrderTypesQuery>(
    workOrderTypesQuery,
  );
  const [dialogKey, setDialogKey] = useState(0);
  const [showAddEditCard, setShowAddEditCard] = useState(false);
  const [
    editingWorkOrderType,
    setEditingWorkOrderType,
  ] = useState<?WorkOrderTypeNode>(null);

  const tableData: Array<WorkOrderTypeNode> = useMemo(
    () =>
      (workOrderTypes?.edges ?? [])
        .map((edge: WorkOrderTypeEdge) => edge.node)
        .filter(Boolean),
    [workOrderTypes],
  );

  const onClose = () => {
    setEditingWorkOrderType(null);
    setDialogKey(key => key + 1);
    setShowAddEditCard(false);
  };

  const saveWorkOrder = () => {
    ServerLogger.info(LogEvents.SAVE_WORK_ORDER_TYPE_BUTTON_CLICKED);
    onClose();
  };

  const showAddEditWorkOrderTypeCard = (woType: ?WorkOrderTypeNode) => {
    ServerLogger.info(LogEvents.ADD_WORK_ORDER_TYPE_BUTTON_CLICKED);
    setEditingWorkOrderType(woType);
    setShowAddEditCard(true);
  };

  if (showAddEditCard) {
    return (
      <div className={classes.paper}>
        <AddEditWorkOrderTypeCard
          key={'new_work_order_type@' + dialogKey}
          open={showAddEditCard}
          onClose={onClose}
          onSave={saveWorkOrder}
          workOrderType={editingWorkOrderType}
        />
      </div>
    );
  }
  return (
    <InventoryView
      header={{
        title: <fbt desc="">Work Order Templates</fbt>,
        subtitle: <fbt desc="">Create and manage reusable work orders.</fbt>,
        actionButtons: [
          <ButtonAction action={() => showAddEditWorkOrderTypeCard(null)}>
            <fbt desc="">Create Work Order Template</fbt>
          </ButtonAction>,
        ],
      }}
      permissions={{
        entity: 'workorderTemplate',
      }}>
      <Table
        data={tableData}
        columns={[
          {
            key: 'name',
            title: 'Work order template',
            render: (row: WorkOrderTypeNode) => (
              <Button
                useEllipsis={true}
                variant="text"
                onClick={() => showAddEditWorkOrderTypeCard(row)}>
                {row.name}
              </Button>
            ),
            getSortingValue: (row: WorkOrderTypeNode) => row.name,
          },
          {
            key: 'description',
            title: 'Description',
            render: (row: WorkOrderTypeNode) => row.description ?? '',
          },
        ]}
        sortSettings={{
          columnKey: 'name',
          order: TABLE_SORT_ORDER.ascending,
        }}
      />
    </InventoryView>
  );
};

export default withInventoryErrorBoundary(WorkOrderTypes);
