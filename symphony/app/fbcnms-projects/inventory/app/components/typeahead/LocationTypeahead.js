/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {LocationTypeahead_LocationsQuery} from './__generated__/LocationTypeahead_LocationsQuery.graphql';
import type {Suggestion} from '@fbcnms/ui/components/Typeahead';

import * as React from 'react';
import Breadcrumbs from '@fbcnms/ui/components/Breadcrumbs';
import RelayEnvironment from '../../common/RelayEnvironment.js';
import Typeahead from '@fbcnms/ui/components/Typeahead';
import {debounce} from 'lodash';
import {fetchQuery, graphql} from 'relay-runtime';

const inventoryEntitiesTypeaheadQuery = graphql`
  query LocationTypeahead_LocationsQuery($name: String!) {
    searchForNode(name: $name, first: 10) {
      edges {
        node {
          __typename
          ... on Location {
            id
            externalId
            name
            locationType {
              name
            }
            locationHierarchy {
              id
              name
              locationType {
                name
              }
            }
          }
        }
      }
    }
  }
`;

const LOCATION_SEARCH_DEBOUNCE_TIMEOUT_MS = 200;
const DEBOUNCE_CONFIG = {
  trailing: true,
  leading: true,
};

type Props = {
  className?: string,
  selectedLocation?: ?{id: string, name: string},
  margin?: ?string,
  onLocationSelection: (location: ?{id: string, name: string}) => void,
  headline?: ?string,
};

type State = {
  locationSuggestions: Array<Suggestion>,
};

class LocationTypeahead extends React.Component<Props, State> {
  static defaultProps = {
    headline: 'Location',
  };

  state = {
    locationSuggestions: [],
  };

  _debounceLocationFetchSuggestions = debounce(
    (searchTerm: string) => this.fetchNewLocationSuggestions(searchTerm),
    LOCATION_SEARCH_DEBOUNCE_TIMEOUT_MS,
    DEBOUNCE_CONFIG,
  );

  fetchNewLocationSuggestions(searchTerm: string) {
    fetchQuery<LocationTypeahead_LocationsQuery>(
      RelayEnvironment,
      inventoryEntitiesTypeaheadQuery,
      {
        name: searchTerm,
      },
    ).then(response => {
      if (!response || !response.searchForNode) {
        return;
      }
      this.setState({
        locationSuggestions:
          response.searchForNode?.edges
            ?.filter(Boolean)
            .map(edge => edge.node)
            .filter(Boolean)
            .map(node => {
              if (node.__typename !== 'Location') {
                return null;
              }
              return {
                entityId: node.id,
                entityType: 'location',
                name: node.name,
                type: node.locationType.name,
                render: () => {
                  const breadcrumbs = [...node.locationHierarchy, node].map(
                    l => ({
                      id: l.id,
                      name: l.name,
                      subtext: l.locationType.name,
                    }),
                  );
                  return <Breadcrumbs breadcrumbs={breadcrumbs} size="small" />;
                },
              };
            })
            .filter(Boolean) ?? ([]: Array<Suggestion>),
      });
    });
  }

  onLocationSuggestionsFetchRequested = (searchTerm: string) => {
    this._debounceLocationFetchSuggestions(searchTerm);
  };

  render() {
    const {selectedLocation, className, headline} = this.props;
    const {locationSuggestions} = this.state;
    return (
      <div className={className}>
        <Typeahead
          className={className}
          margin={this.props.margin}
          required
          suggestions={locationSuggestions}
          onSuggestionsFetchRequested={this.onLocationSuggestionsFetchRequested}
          onEntitySelected={suggestion =>
            this.props.onLocationSelection({
              id: suggestion.entityId,
              name: suggestion.name,
            })
          }
          onEntriesRequested={() => {}}
          onSuggestionsClearRequested={() =>
            this.props.onLocationSelection(null)
          }
          placeholder={headline}
          value={
            selectedLocation
              ? {
                  name: selectedLocation.name,
                  entityId: selectedLocation.id,
                  entityType: '',
                  type: 'location',
                }
              : null
          }
        />
      </div>
    );
  }
}

export default LocationTypeahead;
