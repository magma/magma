---
id: graphql-breaking-changes
title: Graphql API Breaking Changes
---

[//]: <> (@generated This file was created by cli/extract_graphql_deprecations.pydo not change it manually)

## Deprecated Queries
* `searchForEntity` - Use `searchForNode` instead. Will be removed on 2020-05-01

## Deprecated Mutations


## Deprecated Fields
* `Viewer.email` - Use `Viewer.user.email` instead. Will be removed on 2020-05-01
* `WorkOrder.ownerName` - Use `WorkOrder.owner.email` instead. Will be removed on 2020-05-01
* `WorkOrder.assignee` - Use `WorkOrder.assignedTo.email` instead. Will be removed on 2020-05-01
* `Comment.authorName` - Use `Comment.author.email` instead. Will be removed on 2020-05-01
* `Project.creator` - Use `Project.createdBy.email` instead. Will be removed on 2020-05-01

## Deprecated Input Fields
* `AddWorkOrderInput.ownerName` - Use `AddWorkOrderInput.ownerId` instead. Will be removed on 2020-05-01. You cannot use `AddWorkOrderInput.ownerName` and `AddWorkOrderInput.ownerId` together
* `AddWorkOrderInput.assignee` - Use `AddWorkOrderInput.assigneeId` instead. Will be removed on 2020-05-01. You cannot use `AddWorkOrderInput.assignee` and `AddWorkOrderInput.assigneeId` together
* `EditWorkOrderInput.ownerName` - Use `EditWorkOrderInput.ownerId` instead. Will be removed on 2020-05-01. You cannot use `EditWorkOrderInput.ownerName` and `EditWorkOrderInput.ownerId` together
* `EditWorkOrderInput.assignee` - Use `EditWorkOrderInput.assigneeId` instead. Will be removed on 2020-05-01. You cannot use `EditWorkOrderInput.assignee` and `EditWorkOrderInput.assigneeId` together
* `AddProjectInput.creator` - Use `AddProjectInput.creatorId` instead. Will be removed on 2020-05-01. You cannot use `AddProjectInput.creator` and `AddProjectInput.creatorId` together
* `EditProjectInput.creator` - Use `EditProjectInput.creatorId` instead. Will be removed on 2020-05-01. You cannot use `EditProjectInput.creator` and `EditProjectInput.creatorId` together

## Deprecated Enums
* `WorkOrderFilterType.WORK_ORDER_OWNER` - Use `WorkOrderFilterType.WORK_ORDER_OWNED_BY` instead. Will be removed on 2020-05-01
* `WorkOrderFilterType.WORK_ORDER_ASSIGNEE` - Use `WorkOrderFilterType.WORK_ORDER_ASSIGNED_TO` instead. Will be removed on 2020-05-01

