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

import ActionButton from '@fbcnms/ui/components/ActionButton';
import ExpandButtonContext from './context/ExpandButtonContext';
import FormAction from '@fbcnms/ui/components/design-system/Form/FormAction';
import InventoryQueryRenderer from '../components/InventoryQueryRenderer';
import InventoryTreeView from './InventoryTreeView';
import React, {useContext} from 'react';
import classNames from 'classnames';
import withInventoryErrorBoundary from '../common/withInventoryErrorBoundary';
import {graphql} from 'relay-runtime';
import {makeStyles} from '@material-ui/styles';
import {sortLexicographically} from '@fbcnms/ui/utils/displayUtils';

const useStyles = makeStyles({
  root: {
    display: 'flex',
    flexGrow: 1,
    flexDirection: 'column',
    height: '100%',
  },
  collapsedTree: {
    flexGrow: 0,
    overflow: 'hidden',
    width: 0,
    flexBasis: 0,
  },
  expandedTree: {
    overflow: 'auto',
    width: '25%',
    minWidth: '25%',
    flexGrow: 0,
  },
});

type Props = {
  selectedLocationId: ?string,
  onSelect: ?(locationId: ?string) => void,
  onAddLocation: (parentLocation: ?Location) => void,
};

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
    locations(first: 500, onlyTopLevel: true)
      @connection(key: "LocationsTree_locations") {
      edges {
        node {
          ...LocationsTree_location @relay(mask: false)
        }
      }
    }
  }
`;

const LocationsTree = ({
  selectedLocationId,
  onAddLocation,
  onSelect,
}: Props) => {
  const classes = useStyles();
  const {isExpanded, showExpandButton, hideExpandButton} = useContext(
    ExpandButtonContext,
  );
  return (
    <div
      className={classNames({
        [classes.root]: true,
        [classes.collapsedTree]: !isExpanded,
        [classes.expandedTree]: isExpanded,
      })}
      onMouseEnter={showExpandButton}
      onMouseLeave={hideExpandButton}>
      <InventoryQueryRenderer
        query={locationsTreeQuery}
        variables={{}}
        render={props => {
          return (
            <InventoryTreeView
              title="Locations"
              selectedId={selectedLocationId}
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
                <FormAction>
                  <ActionButton
                    action="add"
                    onClick={() => onAddLocation(location)}
                  />
                </FormAction>
              )}
              onClick={(locationId: string) => {
                if (onSelect) {
                  onSelect(locationId);
                }
              }}
            />
          );
        }}
      />
    </div>
  );
};

export default withInventoryErrorBoundary(LocationsTree);
