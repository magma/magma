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

type Props = {
  // $FlowFixMe (T62907961) Relay flow types
  endpoints: ServiceEndpointsView_endpoints,
  onDeleteEndpoint: ?(endpoint: ServiceEndpoint) => void,
};

const ServiceEndpointsView = (props: Props) => {
  const {endpoints, onDeleteEndpoint} = props;

  return (
    <div>
      {endpoints
        .sort((e1, e2) => e1.definition.index - e2.definition.index)
        .map(endpoint => (
          <ServiceEndpointDetails
            endpoint={endpoint}
            onDeleteEndpoint={
              onDeleteEndpoint ? () => onDeleteEndpoint(endpoint) : null
            }
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
      equipment {
        name
        ...EquipmentBreadcrumbs_equipment
      }
      definition {
        name
        role
      }
    }
  `,
});
