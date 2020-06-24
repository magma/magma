/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Location} from './Location';
import type {NamedNode} from './EntUtils';
import type {ProjectTemplateNodesQuery} from './__generated__/ProjectTemplateNodesQuery.graphql';
import type {Property} from './Property';
import type {PropertyType} from './PropertyType';
import type {WorkOrder} from './WorkOrder';

import {graphql} from 'relay-runtime';
import {useLazyLoadQuery} from 'react-relay/hooks';

export type ProjectType = {
  id: string,
  name: string,
  propertyTypes: Array<PropertyType>,
};

export type Project = {
  id: string,
  type: ?ProjectType,
  name: string,
  description: ?string,
  location: ?Location,
  creatorId: ?string,
  properties: Array<Property>,
  workOrders: Array<WorkOrder>,
  numberOfWorkOrders: number,
};

const projectTemplateNodesQuery = graphql`
  query ProjectTemplateNodesQuery {
    projectTypes {
      edges {
        node {
          id
          name
        }
      }
    }
  }
`;

export type ProjectTemplateNode = $Exact<NamedNode>;

export function useProjectTemplateNodes(): $ReadOnlyArray<ProjectTemplateNode> {
  const response = useLazyLoadQuery<ProjectTemplateNodesQuery>(
    projectTemplateNodesQuery,
  );
  const projectTemplatesData = response.projectTypes?.edges || [];
  const projectTemplates = projectTemplatesData
    .map(p => p.node)
    .filter(Boolean);
  return projectTemplates;
}
