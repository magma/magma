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

import * as React from 'react';
import RelayEnvironment from '../../common/RelayEnvironment.js';
import Typeahead from '@fbcnms/ui/components/Typeahead';
import emptyFunction from '@fbcnms/util/emptyFunction';
import {debounce} from 'lodash';
import {fetchQuery, graphql} from 'relay-runtime';

const SERVICE_SEARCH_DEBOUNCE_TIMEOUT_MS = 200;
const DEBOUNCE_CONFIG = {
  trailing: true,
  leading: true,
};

type Props = {
  className?: string,
  required?: boolean,
  headline?: ?string,
  selectedService?: ?{id: string, name: string},
  margin?: ?string,
  onServiceSelection: (?{id: string, name: string}) => void,
};

type State = {
  serviceSuggestions: Array<Suggestion>,
};

const serviceSearchQuery = graphql`
  query ServiceTypeahead_ServicesQuery(
    $filters: [ServiceFilterInput!]!
    $limit: Int
  ) {
    serviceSearch(filters: $filters, limit: $limit) {
      services {
        id
        name
      }
    }
  }
`;

class ServiceTypeahead extends React.Component<Props, State> {
  static defaultProps = {
    headline: 'Service',
  };

  state = {
    serviceSuggestions: [],
  };

  _debounceServiceFetchSuggestions = debounce(
    (searchTerm: string) => this.fetchNewServiceSuggestions(searchTerm),
    SERVICE_SEARCH_DEBOUNCE_TIMEOUT_MS,
    DEBOUNCE_CONFIG,
  );

  fetchNewServiceSuggestions = (searchTerm: string) => {
    fetchQuery(RelayEnvironment, serviceSearchQuery, {
      filters: [
        {
          filterType: 'SERVICE_INST_NAME',
          operator: 'CONTAINS',
          stringValue: searchTerm,
        },
      ],
      limit: 1000,
    }).then(response => {
      if (!response || !response.serviceSearch) {
        return;
      }
      this.setState({
        serviceSuggestions: response.serviceSearch.services.map(p => ({
          name: p.name,
          entityId: p.id,
          entityType: 'service',
          type: p?.type?.name,
        })),
      });
    });
  };

  onServiceSuggestionsFetchRequested = (searchTerm: string) => {
    this._debounceServiceFetchSuggestions(searchTerm);
  };

  render() {
    const {selectedService, headline, required, className} = this.props;
    const {serviceSuggestions} = this.state;
    return (
      <div className={className}>
        <Typeahead
          margin={this.props.margin}
          required={!!required}
          suggestions={serviceSuggestions}
          onSuggestionsFetchRequested={this.onServiceSuggestionsFetchRequested}
          onEntitySelected={suggestion =>
            this.props.onServiceSelection({
              id: suggestion.entityId,
              name: suggestion.name,
            })
          }
          onEntriesRequested={emptyFunction}
          onSuggestionsClearRequested={() =>
            this.props.onServiceSelection(null)
          }
          placeholder={headline}
          value={
            selectedService
              ? {
                  name: selectedService.name,
                  entityId: selectedService.id,
                  entityType: '',
                  type: 'service',
                }
              : null
          }
          variant="small"
        />
      </div>
    );
  }
}

export default ServiceTypeahead;
