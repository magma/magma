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

import * as React from 'react';
import ActiveEquipmentIcon from '@fbcnms/ui/icons/ActiveEquipmentIcon';
import ActiveEquipmentInLocationIcon from '@fbcnms/ui/icons/ActiveEquipmentInLocationIcon';
import ForceNetworkTopology from '../topology/ForceNetworkTopology';
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
  const rootIds = endpoints.map(endpoint => {
    const port = endpoint.port;
    const positionHierarchySize = port.parentEquipment.positionHierarchy.length;
    if (positionHierarchySize > 0) {
      return port.parentEquipment.positionHierarchy[positionHierarchySize - 1]
        .parentEquipment.id;
    }
    return port.parentEquipment.id;
  });
  return (
    <div className={classes.card}>
      <ForceNetworkTopology
        topology={topology}
        className={classes.card}
        renderNode={(id: string) =>
          rootIds.includes(id) ? (
            <ActiveEquipmentInLocationIcon />
          ) : (
            <ActiveEquipmentIcon />
          )
        }
        renderNodeName={(id: string) => {
          const nodes = topology.nodes.filter(node => node.id === id);
          return nodes[0].name;
        }}
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
