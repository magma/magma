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
import DateTimeFormat from '../../common/DateTimeFormat';
import LocationLink from '../location/LocationLink';
import React, {useMemo} from 'react';
import Table from '@fbcnms/ui/components/design-system/Table/Table';
import fbt from 'fbt';
import nullthrows from '@fbcnms/util/nullthrows';
import {InventoryAPIUrls} from '../../common/InventoryAPI';
import {createFragmentContainer, graphql} from 'react-relay';
import {formatMultiSelectValue} from '@fbcnms/ui/utils/displayUtils';
import {statusValues} from '../../common/WorkOrder';
import {useHistory} from 'react-router';

type Props = {
  workOrder: WorkOrdersView_workOrder,
  onWorkOrderSelected: string => void,
};

const WorkOrdersView = (props: Props) => {
  const {workOrder, onWorkOrderSelected} = props;
  const history = useHistory();

  const data = useMemo(() => workOrder.map(wo => ({...wo, key: wo.id})), [
    workOrder,
  ]);

  if (workOrder.length === 0) {
    return <div />;
  }

  return (
    <Table
      data={data}
      columns={[
        {
          key: 'name',
          title: 'Name',
          getSortingValue: row => row.name,
          render: row => (
            <Button variant="text" onClick={() => onWorkOrderSelected(row.id)}>
              {row.name}
            </Button>
          ),
        },
        {
          key: 'type',
          title: `${fbt('Template', '')}`,
          getSortingValue: row => row.workOrderType?.name,
          render: row => row.workOrderType?.name ?? '',
        },
        {
          key: 'project',
          title: 'Project',
          getSortingValue: row => row.project?.name,
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
          getSortingValue: row => row.owner.email,
          render: row => row.owner.email ?? '',
        },
        {
          key: 'status',
          title: 'Status',
          getSortingValue: row => row.status,
          render: row => formatMultiSelectValue(statusValues, row.status) ?? '',
        },
        {
          key: 'creationDate',
          title: 'Creation Time',
          getSortingValue: row => row.creationDate,
          render: row => DateTimeFormat.dateTime(row.creationDate),
        },
        {
          key: 'dueDate',
          title: 'Due Date',
          getSortingValue: row => row.installDate,
          render: row => DateTimeFormat.dateOnly(row.installDate),
        },
        {
          key: 'location',
          title: 'Location',
          getSortingValue: row => row.location?.name,
          render: row =>
            row.location ? (
              <LocationLink title={row.location.name} id={row.location.id} />
            ) : null,
        },
        {
          key: 'assignee',
          title: 'Assignee',
          getSortingValue: row => row.assignedTo?.email,
          render: row => row.assignedTo?.email || null,
        },
        {
          key: 'closeDate',
          title: 'Close Time',
          getSortingValue: row => row.closeDate,
          render: row => DateTimeFormat.dateTime(row.closeDate),
        },
      ]}
    />
  );
};

export default createFragmentContainer(WorkOrdersView, {
  workOrder: graphql`
    fragment WorkOrdersView_workOrder on WorkOrder @relay(plural: true) {
      id
      name
      description
      owner {
        id
        email
      }
      creationDate
      installDate
      status
      assignedTo {
        id
        email
      }
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
      closeDate
    }
  `,
});
