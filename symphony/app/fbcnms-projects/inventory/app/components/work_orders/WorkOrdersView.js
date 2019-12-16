/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {WorkOrdersView_workOrder} from './__generated__/WorkOrdersView_workOrder.graphql';

import Button from '@fbcnms/ui/components/design-system/Button';
import LocationLink from '../location/LocationLink';
import React, {useMemo, useState} from 'react';
import Table from '@fbcnms/ui/components/design-system/Table/Table';
import classNames from 'classnames';
import nullthrows from '@fbcnms/util/nullthrows';
import {InventoryAPIUrls} from '../../common/InventoryAPI';
import {createFragmentContainer, graphql} from 'react-relay';
import {formatMultiSelectValue} from '@fbcnms/ui/utils/displayUtils';
import {statusValues} from '../../common/WorkOrder';
import {useRouter} from '@fbcnms/ui/hooks';

type Props = {
  workOrder: WorkOrdersView_workOrder,
  onWorkOrderSelected: string => void,
  className?: string,
};

const WorkOrdersView = (props: Props) => {
  const {className, workOrder, onWorkOrderSelected} = props;
  const [sortDirection, setSortDirection] = useState('desc');
  const [sortColumn, setSortColumn] = useState('name');
  const {history} = useRouter();

  const sortedWorkOrders = useMemo(
    () =>
      workOrder
        .slice()
        .sort(
          (wo1, wo2) =>
            wo1[sortColumn].localeCompare(wo2[sortColumn]) *
            (sortDirection === 'asc' ? -1 : 1),
        )
        .map(wo => ({...wo, key: wo.id})),
    [sortColumn, sortDirection, workOrder],
  );

  if (workOrder.length === 0) {
    return <div />;
  }

  return (
    <div className={classNames(className)}>
      <Table
        data={sortedWorkOrders}
        onSortClicked={col => {
          if (sortColumn === col) {
            setSortDirection(sortDirection === 'asc' ? 'desc' : 'asc');
          } else {
            setSortColumn(col);
            setSortDirection('desc');
          }
        }}
        columns={[
          {
            key: 'name',
            title: 'Name',
            render: row => (
              <Button
                variant="text"
                onClick={() => onWorkOrderSelected(row.id)}>
                {row.name}
              </Button>
            ),
            sortable: true,
            sortDirection: sortColumn === 'name' ? sortDirection : undefined,
          },
          {
            key: 'type',
            title: 'Type',
            render: row => row.workOrderType?.name ?? '',
          },
          {
            key: 'project',
            title: 'Project',
            render: row =>
              row.project ? (
                <Button
                  variant="text"
                  onClick={() =>
                    history.push(
                      InventoryAPIUrls.project(nullthrows(row.project).id),
                    )
                  }>
                  {row.project?.name ?? ''}
                </Button>
              ) : null,
          },
          {
            key: 'owner',
            title: 'Owner',
            render: row => row.ownerName ?? '',
          },
          {
            key: 'status',
            title: 'Status',
            render: row =>
              formatMultiSelectValue(statusValues, row.status) ?? '',
          },
          {
            key: 'creationDate',
            title: 'Creation Date',
            render: row => new Date(row.creationDate).toLocaleDateString(),
          },
          {
            key: 'dueDate',
            title: 'Due Date',
            render: row =>
              !!row.installDate
                ? new Date(row.installDate).toLocaleDateString()
                : '',
          },
          {
            key: 'location',
            title: 'Location',
            render: row =>
              row.location ? (
                <LocationLink title={row.location.name} id={row.location.id} />
              ) : null,
          },
          {
            key: 'assignee',
            title: 'Assignee',
            render: row => row.assignee || null,
          },
        ]}
      />
    </div>
  );
};

export default createFragmentContainer(WorkOrdersView, {
  workOrder: graphql`
    fragment WorkOrdersView_workOrder on WorkOrder @relay(plural: true) {
      id
      name
      description
      ownerName
      creationDate
      installDate
      status
      assignee
      location {
        id
        name
      }
      workOrderType {
        id
        name
      }
      project {
        id
        name
      }
    }
  `,
});
