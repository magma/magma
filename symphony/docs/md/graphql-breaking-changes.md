---
id: graphql-breaking-changes
title: Graphql API Breaking Changes
---

[//]: <> (@generated This file was created by cli/extract_graphql_deprecations.pydo not change it manually)

## Deprecated Queries
* `equipmentSearch` - Use `equipments` instead. Will be removed on 2020-09-01
* `workOrderSearch` - Use `workOrders` instead. Will be removed on 2020-09-01
* `linkSearch` - Use `links` instead. Will be removed on 2020-09-01
* `portSearch` - Use `equipmentPorts` instead. Will be removed on 2020-09-01
* `locationSearch` - Use `locations` instead. Will be removed on 2020-09-01
* `projectSearch` - Use `projects` instead. Will be removed on 2020-09-01
* `customerSearch` - Use `customers` instead. Will be removed on 2020-09-01
* `serviceSearch` - Use `services` instead. Will be removed on 2020-09-01
* `userSearch` - Use `users` instead. Will be removed on 2020-09-01
* `permissionsPolicySearch` - Use `permissionsPolicies` instead. Will be removed on 2020-09-01
* `usersGroupSearch` - Use `usersGroups` instead. Will be removed on 2020-09-01

## Deprecated Mutations


## Deprecated Fields
* `WorkOrder.ownerName` - Use `WorkOrder.owner.email` instead. Will be removed on 2020-05-01
* `WorkOrder.assignee` - Use `WorkOrder.assignedTo.email` instead. Will be removed on 2020-05-01
* `Project.creator` - Use `Project.createdBy.email` instead. Will be removed on 2020-05-01

## Deprecated Input Fields
* `AddWorkOrderInput.ownerName` - Use `AddWorkOrderInput.ownerId` instead. Will be removed on 2020-05-01. You cannot use `AddWorkOrderInput.ownerName` and `AddWorkOrderInput.ownerId` together
* `AddWorkOrderInput.assignee` - Use `AddWorkOrderInput.assigneeId` instead. Will be removed on 2020-05-01. You cannot use `AddWorkOrderInput.assignee` and `AddWorkOrderInput.assigneeId` together
* `EditWorkOrderInput.ownerName` - Use `EditWorkOrderInput.ownerId` instead. Will be removed on 2020-05-01. You cannot use `EditWorkOrderInput.ownerName` and `EditWorkOrderInput.ownerId` together
* `EditWorkOrderInput.assignee` - Use `EditWorkOrderInput.assigneeId` instead. Will be removed on 2020-05-01. You cannot use `EditWorkOrderInput.assignee` and `EditWorkOrderInput.assigneeId` together
* `AddProjectInput.creator` - Use `AddProjectInput.creatorId` instead. Will be removed on 2020-05-01. You cannot use `AddProjectInput.creator` and `AddProjectInput.creatorId` together
* `EditProjectInput.creator` - Use `EditProjectInput.creatorId` instead. Will be removed on 2020-05-01. You cannot use `EditProjectInput.creator` and `EditProjectInput.creatorId` together

## Deprecated Enums


