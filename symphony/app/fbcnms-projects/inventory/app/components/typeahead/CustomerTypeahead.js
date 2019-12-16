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
import {fetchQuery, graphql} from 'relay-runtime';

type Props = {
  className?: string,
  required?: boolean,
  headline?: string,
  selectedCustomer?: ?string,
  margin?: ?string,
  onCustomerSelection: (?{id: string, name: string}) => void,
};

type State = {
  customerSuggestions: Array<Suggestion>,
  customers: Array<Suggestion>,
};

const customerSearchQuery = graphql`
  query CustomerTypeahead_CustomersQuery($limit: Int) {
    customerSearch(limit: $limit) {
      id
      name
      externalId
    }
  }
`;

class CustomerTypeahead extends React.Component<Props, State> {
  state = {
    customerSuggestions: [],
    customers: [],
  };

  componentDidMount() {
    fetchQuery(RelayEnvironment, customerSearchQuery, {
      limit: 1000,
      filters: [],
    }).then(response => {
      if (!response || !response.customerSearch) {
        return;
      }
      this.setState({
        customers: response.customerSearch.map(p => ({
          name: p.name,
          entityId: p.id,
          entityType: 'customer',
          type: p?.type?.name,
        })),
      });
    });
  }

  fetchNewCustomerSuggestions = (searchTerm: string) => {
    const searchTermLC = searchTerm.toLowerCase();
    const customers = this.state.customers;
    const suggestions = customers.filter(e =>
      e.name.toLowerCase().includes(searchTermLC),
    );
    this.setState({
      customerSuggestions: suggestions,
    });
  };

  render() {
    const {selectedCustomer, headline, required, className} = this.props;
    const {customerSuggestions} = this.state;
    return (
      <div className={className}>
        <Typeahead
          margin={this.props.margin}
          required={!!required}
          suggestions={customerSuggestions}
          onSuggestionsFetchRequested={this.fetchNewCustomerSuggestions}
          onEntitySelected={suggestion =>
            this.props.onCustomerSelection({
              id: suggestion.entityId,
              name: suggestion.name,
            })
          }
          onEntriesRequested={emptyFunction}
          onSuggestionsClearRequested={() =>
            this.props.onCustomerSelection(null)
          }
          placeholder={headline}
          value={
            selectedCustomer
              ? {
                  name: selectedCustomer,
                  entityId: '1',
                  entityType: '',
                  type: '',
                }
              : null
          }
          variant="small"
          disabled={true}
        />
      </div>
    );
  }
}

export default CustomerTypeahead;
