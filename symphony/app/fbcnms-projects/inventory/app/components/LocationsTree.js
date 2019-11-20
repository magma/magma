/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import type {Location} from '../common/Location';
import type {WithStyles} from '@material-ui/core';

import ActionButton from '@fbcnms/ui/components/ActionButton';
import InventoryQueryRenderer from '../components/InventoryQueryRenderer';
import InventoryTreeView from './InventoryTreeView';
import React from 'react';
import withInventoryErrorBoundary from '../common/withInventoryErrorBoundary';
import {graphql} from 'relay-runtime';
import {sortLexicographically} from '@fbcnms/ui/utils/displayUtils';
import {withStyles} from '@material-ui/core/styles';

type Props = WithStyles<typeof styles> & {
  selectedLocationId: ?string,
  onSelect: ?(locationId: ?string) => void,
  onAddLocation: (parentLocation: ?Location) => void,
};

const styles = theme => ({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
  button: {
    margin: theme.spacing(),
  },
});

graphql`
  fragment LocationsTree_location on Location @relay(mask: false) {
    id
    externalId
    name
    locationType {
      id
      name
    }
    numChildren
    siteSurveyNeeded
  }
`;

const locationsTreeQuery = graphql`
  query LocationsTreeQuery {
    locations(first: 50, onlyTopLevel: true)
      @connection(key: "LocationsTree_locations") {
      edges {
        node {
          ...LocationsTree_location @relay(mask: false)
        }
      }
    }
  }
`;

class LocationsTree extends React.Component<Props> {
  render() {
    return (
      <InventoryQueryRenderer
        query={locationsTreeQuery}
        variables={{}}
        render={props => {
          return (
            <InventoryTreeView
              title="Locations"
              selectedId={this.props.selectedLocationId}
              dummyRootTitle="Add top-level location"
              tree={
                props?.locations
                  ? props.locations.edges
                      .map(x => x.node)
                      .sort((x, y) =>
                        sortLexicographically(x.name ?? '', y.name ?? ''),
                      )
                  : []
              }
              titlePropertyGetter={(location: Location) => location.name}
              subtitlePropertyGetter={(location: Location) =>
                location.locationType?.name
              }
              childrenPropertyGetter={(location: Location) =>
                (location?.children ?? [])
                  .slice()
                  .filter(Boolean)
                  .sort((x, y) =>
                    sortLexicographically(x.name ?? '', y.name ?? ''),
                  )
              }
              getHoverRightContent={(location: ?Location) => (
                <ActionButton
                  action="add"
                  onClick={() => this.props.onAddLocation(location)}
                />
              )}
              onClick={(locationId: string) => {
                if (this.props.onSelect) {
                  this.props.onSelect(locationId);
                }
              }}
            />
          );
        }}
      />
    );
  }
}

export default withStyles(styles)(withInventoryErrorBoundary(LocationsTree));
