#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from pysymphony import SymphonyClient

from ..common.data_class import WorkOrder, WorkOrderType
from ..graphql.input.add_work_order import AddWorkOrderInput
from ..graphql.input.add_work_order_type import AddWorkOrderTypeInput
from ..graphql.mutation.add_workorder import AddWorkOrderMutation
from ..graphql.mutation.add_workorder_type import AddWorkOrderTypeMutation


def add_workorder_type(client: SymphonyClient, name: str) -> WorkOrderType:
    """This function creates work order type with the given name

        Args:
            name (str): work order type name

        Returns:
            `pyworkforce.common.data_class.WorkOrderType`

        Example:
            ```
            client.add_workorder_type("Deployment work order")
            ```
    """
    result = AddWorkOrderTypeMutation.execute(
        client, AddWorkOrderTypeInput(name=name, checkListCategories=[])
    )
    return WorkOrderType(id=result.id)


def add_workorder(
    client: SymphonyClient, name: str, workorder_type: WorkOrderType
) -> WorkOrder:
    """This function creates work order of with the given name and type

        Args:
            name (str): work order name
            workorder_type (`pyworkforce.common.data_class.WorkOrderType`): work order type

        Returns:
            `pyworkforce.common.data_class.WorkOrder`

        Example:
            ```
            workorder_type = client.add_workorder_type("Deployment work order")
            client.add_workorder_type("new work order", workorder_type)
            ```
    """
    result = AddWorkOrderMutation.execute(
        client,
        AddWorkOrderInput(
            name=name,
            workOrderTypeId=workorder_type.id,
            properties=[],
            checkList=[],
            checkListCategories=[],
        ),
    )
    return WorkOrder(id=result.id, name=result.name)
