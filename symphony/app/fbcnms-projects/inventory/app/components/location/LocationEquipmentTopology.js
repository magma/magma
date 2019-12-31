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
import type {LocationEquipmentTopology_topology} from './__generated__/LocationEquipmentTopology_topology.graphql';

import * as React from 'react';
import ActiveEquipmentIcon from '@fbcnms/ui/icons/ActiveEquipmentIcon';
import ActiveEquipmentInLocationIcon from '@fbcnms/ui/icons/ActiveEquipmentInLocationIcon';
import ForceNetworkTopology from '../topology/ForceNetworkTopology';
import TopologyTextBox from '../topology/TopologyTextBox';
import {createFragmentContainer, graphql} from 'react-relay';

type Props = {
  topology: LocationEquipmentTopology_topology,
  equipment: Array<Equipment>,
};

const LocationEquipmentTopology = (props: Props) => {
  const {topology, equipment} = props;
  const rootIds = equipment.map(eq => eq.id);
  return (
    <ForceNetworkTopology
      topology={topology}
      renderNode={(id: string) => {
        const node = topology.nodes.find(node => node.id === id);
        return (
          <g transform="translate(-8 -8)">
            {rootIds.includes(id) ? (
              <ActiveEquipmentInLocationIcon />
            ) : (
              <ActiveEquipmentIcon />
            )}
            <TopologyTextBox transform="translate(8 40)" text={node?.name} />
          </g>
        );
      }}
    />
  );
};

export default createFragmentContainer(LocationEquipmentTopology, {
  topology: graphql`
    fragment LocationEquipmentTopology_topology on NetworkTopology {
      nodes {
        ... on Equipment {
          id
          name
        }
      }
      ...ForceNetworkTopology_topology
    }
  `,
  equipment: graphql`
    fragment LocationEquipmentTopology_equipment on Equipment
      @relay(plural: true) {
      id
    }
  `,
});
