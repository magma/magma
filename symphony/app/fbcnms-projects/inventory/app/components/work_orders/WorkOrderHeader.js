/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import Breadcrumbs from '@fbcnms/ui/components/Breadcrumbs';
import Button from '@fbcnms/ui/components/design-system/Button';
import React from 'react';
import WorkOrderDeleteButton from './WorkOrderDeleteButton';
import WorkOrderSaveButton from './WorkOrderSaveButton';
import nullthrows from '@fbcnms/util/nullthrows';
import {InventoryAPIUrls} from '../../common/InventoryAPI';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';
import type {
  // $FlowFixMe (T62907961) Relay flow types
  CheckListCategoryTable_list,
  WorkOrderDetails_workOrder,
} from './__generated__/WorkOrderDetails_workOrder.graphql.js';
import type {Property} from '../../common/Property';

const useStyles = makeStyles(_theme => ({
  nameHeader: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    paddingBottom: '24px',
    overflow: 'hidden',
    flexBasis: 'auto',
  },
  breadcrumbs: {
    flexGrow: 1,
    width: '10px',
  },
  actionButtons: {
    display: 'flex',
    flexDirection: 'row',
  },
  deleteButton: {
    marginRight: '8px',
  },
  cancelButton: {
    marginRight: '8px',
  },
}));

type Props = {
  workOrderName: string,
  workOrder: WorkOrderDetails_workOrder,
  properties: Array<Property>,
  checkListCategories: CheckListCategoryTable_list,
  locationId: ?string,
  onWorkOrderRemoved: () => void,
  onCancelClicked: () => void,
};

const WorkOrderHeader = (props: Props) => {
  const classes = useStyles();
  const {history} = useRouter();
  const {
    workOrderName,
    workOrder,
    properties,
    checkListCategories,
    locationId,
    onWorkOrderRemoved,
    onCancelClicked,
  } = props;
  return (
    <div className={classes.nameHeader}>
      <div className={classes.breadcrumbs}>
        <Breadcrumbs
          breadcrumbs={[
            {
              id: 'work_orders',
              name: 'Work Orders',
              onClick: onCancelClicked,
            },
            {
              id: workOrder.project?.id ?? '',
              name: workOrder.project?.name ?? '',
              subtext: workOrder.project?.type?.name,
              onClick: () =>
                history.push(
                  InventoryAPIUrls.project(nullthrows(workOrder.project?.id)),
                ),
            },
            {
              id: workOrder.id,
              name: workOrderName,
              subtext: workOrder.workOrderType.name,
            },
          ].filter(x => !!x.id)}
          size="large"
        />
      </div>
      <div className={classes.actionButtons}>
        <WorkOrderDeleteButton
          className={classes.deleteButton}
          workOrder={workOrder}
          onWorkOrderRemoved={onWorkOrderRemoved}
        />
        <Button
          className={classes.cancelButton}
          skin="regular"
          onClick={onCancelClicked}>
          Cancel
        </Button>
        <WorkOrderSaveButton
          workOrder={workOrder}
          properties={properties}
          checkListCategories={checkListCategories}
          locationId={locationId}
        />
      </div>
    </div>
  );
};

export default WorkOrderHeader;
