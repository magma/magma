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
import type {WithStyles} from '@material-ui/core';

import * as React from 'react';
import ForceNetworkTopology from '../topology/ForceNetworkTopology';
import {createFragmentContainer, graphql} from 'react-relay';
import {withStyles} from '@material-ui/core/styles';

type Props = {
  topology: TopologyNetwork,
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
  return (
    <div className={classes.card}>
      <ForceNetworkTopology
        networkTopology={topology}
        rootIds={terminationPoints.map(eq => eq.id)}
        className={classes.card}
      />
    </div>
  );
};

export default withStyles(styles)(
  createFragmentContainer(ServiceEquipmentTopology, {
    topology: graphql`
      fragment ServiceEquipmentTopology_topology on NetworkTopology {
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
    terminationPoints: graphql`
      fragment ServiceEquipmentTopology_terminationPoints on Equipment
        @relay(plural: true) {
        id
      }
    `,
  }),
);
