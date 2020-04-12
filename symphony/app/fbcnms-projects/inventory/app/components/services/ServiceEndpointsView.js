/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ServiceEndpoint} from '../../common/Service';

import * as React from 'react';
import ServiceEndpointDetails from './ServiceEndpointDetails';
import ServiceEndpointsView_endpoints from './__generated__/ServiceEndpointsView_endpoints.graphql';
import {createFragmentContainer, graphql} from 'react-relay';
import {sortLexicographically} from '@fbcnms/ui/utils/displayUtils';

type Props = {
  // $FlowFixMe (T62907961) Relay flow types
  endpoints: ServiceEndpointsView_endpoints,
  onDeleteEndpoint: (endpoint: ServiceEndpoint) => void,
};

const ServiceEndpointsView = (props: Props) => {
  const {endpoints, onDeleteEndpoint} = props;

  return (
    <div>
      {endpoints
        .sort((e1, e2) =>
          sortLexicographically(
            e1.port.parentEquipment.name,
            e2.port.parentEquipment.name,
          ),
        )
        .sort((e1, e2) => sortLexicographically(e1.role, e2.role))
        .map(endpoint => (
          <ServiceEndpointDetails
            endpoint={endpoint}
            onDeleteEndpoint={() => onDeleteEndpoint(endpoint)}
          />
        ))}
    </div>
  );
};

export default createFragmentContainer(ServiceEndpointsView, {
  endpoints: graphql`
    fragment ServiceEndpointsView_endpoints on ServiceEndpoint
      @relay(plural: true) {
      id
      port {
        parentEquipment {
          name
          ...EquipmentBreadcrumbs_equipment
        }
        definition {
          id
          name
        }
      }
      definition {
        role
      }
    }
  `,
});
