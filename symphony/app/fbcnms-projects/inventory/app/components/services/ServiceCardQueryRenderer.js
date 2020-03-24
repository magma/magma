/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import InventoryQueryRenderer from '../InventoryQueryRenderer';
import React from 'react';
import ServiceCard from './ServiceCard';
import {graphql} from 'react-relay';

type Props = {
  serviceId: ?string,
};

const serviceQuery = graphql`
  query ServiceCardQueryRendererQuery($serviceId: ID!) {
    node(id: $serviceId) {
      ... on Service {
        ...ServiceCard_service
      }
    }
  }
`;

const ServiceCardQueryRenderer = (props: Props) => {
  const {serviceId} = props;

  return (
    <InventoryQueryRenderer
      query={serviceQuery}
      variables={{
        serviceId,
      }}
      render={props => {
        const {node} = props;
        return <ServiceCard service={node} />;
      }}
    />
  );
};

export default ServiceCardQueryRenderer;
