/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Suggestion} from '@fbcnms/ui/components/Typeahead';
import type {Theme, WithStyles} from '@material-ui/core';

import * as React from 'react';
import RelayEnvironment from '../common/RelayEnvironment.js';
import Typeahead from '@fbcnms/ui/components/Typeahead';
import {debounce} from 'lodash';
import {fetchQuery, graphql} from 'relay-runtime';
import {withStyles} from '@material-ui/core/styles';

const inventoryEntitiesTypeaheadQuery = graphql`
  query InventoryEntitiesTypeaheadQuery($name: String!) {
    searchForEntity(name: $name, first: 10) {
      edges {
        node {
          entityId
          entityType
          name
          type
          externalId
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
    searchTerm => this.fetchNewSuggestions(searchTerm),
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
    fetchQuery(RelayEnvironment, inventoryEntitiesTypeaheadQuery, {
      name: searchTerm,
    }).then(response => {
      if (!response || !response.searchForEntity) {
        return;
      }

      const suggestions = response.searchForEntity.edges
        .filter(Boolean)
        .map(edge => ({
          ...edge.node,
        }));
      suggestions.forEach(node => {
        if (!!node.externalId) {
          node.type = `${node.type} - ${node.externalId}`;
        }
      });

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
          getSuggestionValue={suggestion => suggestion.name}
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
