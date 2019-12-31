/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ServiceEquipmentTopology_endpoints} from './__generated__/ServiceEquipmentTopology_endpoints.graphql';
import type {ServiceEquipmentTopology_topology} from './__generated__/ServiceEquipmentTopology_topology.graphql';
import type {WithStyles} from '@material-ui/core';

import ActiveConsumerEndpointIcon from '@fbcnms/ui/icons/ActiveConsumerEndpointIcon';
import ActiveEquipmentIcon from '@fbcnms/ui/icons/ActiveEquipmentIcon';
import ActiveProviderEndpointIcon from '@fbcnms/ui/icons/ActiveProviderEndpointIcon';
import ForceNetworkTopology from '../topology/ForceNetworkTopology';
import React, {useCallback} from 'react';
import TopologyTextBox from '../topology/TopologyTextBox';
import {createFragmentContainer, graphql} from 'react-relay';
import {withStyles} from '@material-ui/core/styles';

type Props = {
  topology: ServiceEquipmentTopology_topology,
  endpoints: ServiceEquipmentTopology_endpoints,
} & WithStyles<typeof styles>;

const styles = _ => ({
  card: {
    height: '100%',
    position: 'relative',
  },
});

const ServiceEquipmentTopology = (props: Props) => {
  const {topology, endpoints, classes} = props;

  const getEndpointTopLevelEquipment = endpoint => {
    const port = endpoint.port;
    const positionHierarchySize = port.parentEquipment.positionHierarchy.length;
    if (positionHierarchySize > 0) {
      return port.parentEquipment.positionHierarchy[positionHierarchySize - 1]
        .parentEquipment.id;
    }
    return port.parentEquipment.id;
  };

  const renderNode = useCallback(
    (id: string) => {
      const node = topology.nodes.find(node => node.id === id);
      const consumerIds = endpoints
        .filter(endpoint => endpoint.role == 'CONSUMER')
        .map(getEndpointTopLevelEquipment);
      const providerIds = endpoints
        .filter(endpoint => endpoint.role == 'PROVIDER')
        .map(getEndpointTopLevelEquipment);
      return consumerIds.includes(id) ? (
        <g transform="translate(-18 -18)">
          <ActiveConsumerEndpointIcon variant="large" />
          <TopologyTextBox transform="translate(20 65)" text={node?.name} />
        </g>
      ) : providerIds.includes(id) ? (
        <g transform="translate(-18 -18)">
          <ActiveProviderEndpointIcon variant="large" />
          <TopologyTextBox transform="translate(20 65)" text={node?.name} />
        </g>
      ) : (
        <g transform="translate(-8 -8)">
          <ActiveEquipmentIcon />
          <TopologyTextBox transform="translate(10 40)" text={node?.name} />
        </g>
      );
    },
    [endpoints, topology.nodes],
  );

  return (
    <div className={classes.card}>
      <ForceNetworkTopology
        topology={topology}
        className={classes.card}
        renderNode={renderNode}
      />
    </div>
  );
};

export default withStyles(styles)(
  createFragmentContainer(ServiceEquipmentTopology, {
    topology: graphql`
      fragment ServiceEquipmentTopology_topology on NetworkTopology {
        nodes {
          ... on Equipment {
            id
            name
          }
        }
        ...ForceNetworkTopology_topology
      }
    `,
    endpoints: graphql`
      fragment ServiceEquipmentTopology_endpoints on ServiceEndpoint
        @relay(plural: true) {
        role
        port {
          parentEquipment {
            id
            positionHierarchy {
              parentEquipment {
                id
              }
            }
          }
        }
      }
    `,
  }),
);
