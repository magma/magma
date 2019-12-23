/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Equipment} from '../../common/Equipment';
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
  terminationPoints: Array<Equipment>,
} & WithStyles<typeof styles>;

const styles = _ => ({
  card: {
    height: '100%',
    position: 'relative',
  },
});

const ServiceEquipmentTopology = (props: Props) => {
  const {topology, terminationPoints, classes} = props;
  const rootIds = terminationPoints.map(eq => eq.id);
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
    terminationPoints: graphql`
      fragment ServiceEquipmentTopology_terminationPoints on Equipment
        @relay(plural: true) {
        id
      }
    `,
  }),
);
