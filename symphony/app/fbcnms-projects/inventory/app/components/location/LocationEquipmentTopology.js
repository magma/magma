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
import type {TopologyNetwork} from '../../common/NetworkTopology';

import * as React from 'react';
import ForceNetworkTopology from '../topology/ForceNetworkTopology';
import {createFragmentContainer, graphql} from 'react-relay';

type Props = {
  topology: TopologyNetwork,
  equipment: Array<Equipment>,
};

const LocationEquipmentTopology = (props: Props) => {
  const {topology, equipment} = props;
  return (
    <ForceNetworkTopology
      networkTopology={topology}
      rootIds={equipment.map(eq => eq.id)}
    />
  );
};

export default createFragmentContainer(LocationEquipmentTopology, {
  topology: graphql`
    fragment LocationEquipmentTopology_topology on NetworkTopology {
      nodes {
        id
        name
      }
      links {
        source
        target
      }
    }
  `,
  equipment: graphql`
    fragment LocationEquipmentTopology_equipment on Equipment
      @relay(plural: true) {
      id
    }
  `,
});
