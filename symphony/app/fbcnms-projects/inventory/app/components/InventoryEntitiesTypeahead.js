/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {InventoryEntitiesTypeaheadQuery} from './__generated__/InventoryEntitiesTypeaheadQuery.graphql';
import type {Suggestion} from '@fbcnms/ui/components/Typeahead';
import type {Theme, WithStyles} from '@material-ui/core';

import * as React from 'react';
import EquipmentBreadcrumbs from './equipment/EquipmentBreadcrumbs';
import LocationBreadcrumbsTitle from './location/LocationBreadcrumbsTitle';
import RelayEnvironment from '../common/RelayEnvironment.js';
import Text from '@fbcnms/ui/components/design-system/Text';
import Typeahead from '@fbcnms/ui/components/Typeahead';
import {debounce} from 'lodash';
import {fetchQuery, graphql} from 'relay-runtime';
import {withStyles} from '@material-ui/core/styles';

const inventoryEntitiesTypeaheadQuery = graphql`
  query InventoryEntitiesTypeaheadQuery($name: String!) {
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
          ... on Equipment {
            id
            externalId
            name
            equipmentType {
              name
            }
            ...EquipmentBreadcrumbs_equipment
          }
        }
      }
    }
  }
`;

const styles = (theme: Theme) => ({
  container: {
    minWidth: '250px',
  },
  suggestionRoot: {
    display: 'flex',
  },
  suggestionType: {
    color: theme.palette.text.secondary,
    fontSize: theme.typography.pxToRem(13),
    lineHeight: '21px',
    marginLeft: theme.spacing(),
  },
  breadcrumbsContainer: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    maxWidth: '600px',
    overflow: 'hidden',
  },
  externalId: {
    marginLeft: '6px',
  },
});

type EntityType = 'location' | 'equipment';

type Props = {
  onEntitySelected: (entityId: string, entityType: EntityType) => void,
} & WithStyles<typeof styles>;

type State = {
  suggestions: Array<Suggestion>,
};

const SEARCH_DEBOUNCE_TIMEOUT_MS = 200;

class InventoryEntitiesTypeahead extends React.Component<Props, State> {
  _debounceFetchSuggestions = debounce(
    (searchTerm: string) => this.fetchNewSuggestions(searchTerm),
    SEARCH_DEBOUNCE_TIMEOUT_MS,
    {
      trailing: true,
      leading: true,
    },
  );

  state = {
    suggestions: [],
  };

  fetchNewSuggestions(searchTerm: string) {
    const {classes} = this.props;
    fetchQuery<InventoryEntitiesTypeaheadQuery>(
      RelayEnvironment,
      inventoryEntitiesTypeaheadQuery,
      {
        name: searchTerm,
      },
    ).then(response => {
      if (!response || !response.searchForNode) {
        return;
      }

      const mapToSuggestion = (node): ?Suggestion => {
        if (node.__typename === 'Equipment') {
          return {
            entityId: node.id,
            entityType: 'equipment',
            name: node.name,
            type: node.equipmentType.name,
            render: () => {
              return (
                <div className={classes.breadcrumbsContainer}>
                  <EquipmentBreadcrumbs equipment={node} size="small" />
                </div>
              );
            },
          };
        } else if (node.__typename === 'Location') {
          return {
            entityId: node.id,
            entityType: 'location',
            name: node.name,
            type: node.locationType.name,
            render: () => (
              <div className={classes.breadcrumbsContainer}>
                <LocationBreadcrumbsTitle
                  locationDetails={node}
                  size="small"
                  hideTypes={true}
                  navigateOnClick={false}
                />
                {node.externalId && (
                  <Text
                    variant="caption"
                    color="gray"
                    className={classes.externalId}>
                    ({node.externalId})
                  </Text>
                )}
              </div>
            ),
          };
        }
        return (null: ?Suggestion);
      };
      const suggestions: Array<Suggestion> = (
        response.searchForNode.edges?.map(edge => {
          if (edge.node == null) {
            return null;
          }
          return mapToSuggestion(edge.node);
        }) ?? []
      ).filter(Boolean);
      this.setState({suggestions});
    });
  }

  onSuggestionsFetchRequested = searchTerm => {
    this._debounceFetchSuggestions(searchTerm);
  };

  render() {
    const {classes, onEntitySelected} = this.props;
    const {suggestions} = this.state;
    return (
      <div className={classes.container}>
        <Typeahead
          required={false}
          suggestions={suggestions}
          getSuggestionValue={(suggestion: Suggestion) => suggestion.name}
          onSuggestionsFetchRequested={this.onSuggestionsFetchRequested}
          onEntitySelected={suggestion => {
            const entityType =
              suggestion.entityType === 'location' ? 'location' : 'equipment';
            onEntitySelected(suggestion.entityId, entityType);
          }}
        />
      </div>
    );
  }
}

export default withStyles(styles)(InventoryEntitiesTypeahead);
