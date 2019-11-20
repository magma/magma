/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import * as React from 'react';
import InventoryQueryRenderer from '../../components/InventoryQueryRenderer';
import ServiceEquipmentTopology from './ServiceEquipmentTopology';
import {graphql} from 'react-relay';

type Props = {
  serviceId: string,
};

const networkTopologyQuery = graphql`
  query ServiceNetworkMapTabQuery($serviceId: ID!) {
    service(id: $serviceId) {
      terminationPoints {
        ...ServiceEquipmentTopology_terminationPoints
      }
      topology {
        ...ServiceEquipmentTopology_topology
      }
    }
  }
`;

const ServiceNetworkMap = (props: Props) => {
  const {serviceId} = props;
  return (
    <InventoryQueryRenderer
      query={networkTopologyQuery}
      variables={{
        serviceId: serviceId,
      }}
      render={props => {
        const service = props.service;
        return (
          <ServiceEquipmentTopology
            topology={service.topology}
            terminationPoints={service.terminationPoints}
          />
        );
      }}
    />
  );
};

export default ServiceNetworkMap;
