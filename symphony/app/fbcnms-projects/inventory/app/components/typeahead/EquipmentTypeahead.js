/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {EquipmentTypeahead_equipmentQueryResponse} from './__generated__/EquipmentTypeahead_equipmentQuery.graphql';
import type {Suggestion} from '@fbcnms/ui/components/Typeahead';

import * as React from 'react';
import RelayEnvironment from '../../common/RelayEnvironment.js';
import Typeahead from '@fbcnms/ui/components/Typeahead';
import {debounce} from 'lodash';
import {fetchQuery, graphql} from 'relay-runtime';

const equipmentTypeaheadQuery = graphql`
  query EquipmentTypeahead_equipmentQuery($filters: [EquipmentFilterInput!]!) {
    equipmentSearch(limit: 10, filters: $filters) {
      equipment {
        id
        name
      }
    }
  }
`;

const EQUIPMENT_SEARCH_DEBOUNCE_TIMEOUT_MS = 200;
const DEBOUNCE_CONFIG = {
  trailing: true,
  leading: true,
};

type Props = {
  className?: string,
  selectedEquipment?: ?{id: string, name: string},
  margin?: ?string,
  onEquipmentSelection: (?{id: string, name: string}) => void,
  headline?: ?string,
};

type State = {
  equipmentSuggestions: Array<Suggestion>,
};

class EquipmentTypeahead extends React.Component<Props, State> {
  static defaultProps = {
    headline: 'Equipment',
  };

  state = {
    equipmentSuggestions: [],
  };

  _debounceEquipmentFetchSuggestions = debounce(
    (searchTerm: string) => this.fetchNewEquipmentSuggestions(searchTerm),
    EQUIPMENT_SEARCH_DEBOUNCE_TIMEOUT_MS,
    DEBOUNCE_CONFIG,
  );

  fetchNewEquipmentSuggestions(searchTerm: string) {
    fetchQuery(RelayEnvironment, equipmentTypeaheadQuery, {
      filters: [
        {
          filterType: 'EQUIP_INST_NAME',
          operator: 'CONTAINS',
          stringValue: searchTerm,
        },
      ],
    }).then((response: ?EquipmentTypeahead_equipmentQueryResponse) => {
      if (!response || !response.equipmentSearch) {
        return;
      }
      this.setState({
        equipmentSuggestions: response.equipmentSearch.equipment
          .filter(Boolean)
          .map(e => ({
            name: e.name,
            entityId: e.id,
            entityType: '',
            type: 'equipment',
          })),
      });
    });
  }

  onEquipmentSuggestionsFetchRequested = (searchTerm: string) => {
    this._debounceEquipmentFetchSuggestions(searchTerm);
  };

  render() {
    const {
      selectedEquipment,
      className,
      headline,
      onEquipmentSelection,
      margin,
    } = this.props;
    const {equipmentSuggestions} = this.state;
    return (
      <div className={className}>
        <Typeahead
          className={className}
          margin={margin}
          required
          suggestions={equipmentSuggestions}
          onSuggestionsFetchRequested={
            this.onEquipmentSuggestionsFetchRequested
          }
          onEntitySelected={suggestion => {
            onEquipmentSelection({
              id: suggestion.entityId,
              name: suggestion.name,
            });
          }}
          onEntriesRequested={() => {}}
          onSuggestionsClearRequested={() => onEquipmentSelection(null)}
          headline={headline}
          value={
            selectedEquipment
              ? {
                  name: selectedEquipment.name,
                  entityId: selectedEquipment.id,
                  entityType: '',
                  type: 'equipment',
                }
              : null
          }
        />
      </div>
    );
  }
}

export default EquipmentTypeahead;
