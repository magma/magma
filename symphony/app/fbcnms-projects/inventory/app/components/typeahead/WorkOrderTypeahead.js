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
import type {WorkOrderTypeaheadQuery} from './__generated__/WorkOrderTypeaheadQuery.graphql';

import * as React from 'react';
import RelayEnvironment from '../../common/RelayEnvironment.js';
import Typeahead from '@fbcnms/ui/components/Typeahead';
import emptyFunction from '@fbcnms/util/emptyFunction';
import {debounce} from 'lodash';
import {fetchQuery, graphql} from 'relay-runtime';
import {useState} from 'react';

const SEARCH_DEBOUNCE_TIMEOUT_MS = 200;
const DEBOUNCE_CONFIG = {
  trailing: true,
  leading: true,
};

type Props = {
  className?: string,
  required?: boolean,
  headline?: ?string,
  selectedWorkOrder?: ?{id: string, name: string},
  margin?: ?string,
  onWorkOrderSelected: (?{id: string, name: string}) => void,
};

const workOrderSearchQuery = graphql`
  query WorkOrderTypeaheadQuery(
    $filters: [WorkOrderFilterInput!]!
    $limit: Int
  ) {
    workOrderSearch(filters: $filters, limit: $limit) {
      workOrders {
        id
        name
        workOrderType {
          name
        }
      }
    }
  }
`;

const WorkOrderTypeahead = ({
  selectedWorkOrder,
  onWorkOrderSelected,
  headline,
  required,
  className,
  margin,
}: Props) => {
  const [suggestions, setSuggestions] = useState<Array<Suggestion>>([]);

  const debounceFetchSuggestions = debounce(
    (searchTerm: string) => fetchNewSuggestions(searchTerm),
    SEARCH_DEBOUNCE_TIMEOUT_MS,
    DEBOUNCE_CONFIG,
  );

  const fetchNewSuggestions = (searchTerm: string) => {
    fetchQuery<WorkOrderTypeaheadQuery>(
      RelayEnvironment,
      workOrderSearchQuery,
      {
        filters: [
          {
            filterType: 'WORK_ORDER_NAME',
            operator: 'CONTAINS',
            stringValue: searchTerm,
          },
        ],
        limit: 10,
      },
    ).then(response => {
      if (!response || !response.workOrderSearch) {
        return;
      }
      setSuggestions(
        response.workOrderSearch.workOrders.filter(Boolean).map(wo => ({
          name: wo.name,
          entityId: wo.id,
          entityType: 'work_order',
          type: wo.workOrderType.name,
        })),
      );
    });
  };

  const onSuggestionsFetchRequested = (searchTerm: string) => {
    debounceFetchSuggestions(searchTerm);
  };

  return (
    <div className={className}>
      <Typeahead
        margin={margin}
        required={!!required}
        suggestions={suggestions}
        onSuggestionsFetchRequested={onSuggestionsFetchRequested}
        onEntitySelected={suggestion =>
          onWorkOrderSelected({
            id: suggestion.entityId,
            name: suggestion.name,
          })
        }
        onEntriesRequested={emptyFunction}
        onSuggestionsClearRequested={() => onWorkOrderSelected(null)}
        placeholder={headline}
        value={
          selectedWorkOrder
            ? {
                name: selectedWorkOrder.name,
                entityId: selectedWorkOrder.id,
                entityType: '',
                type: 'work_order',
              }
            : null
        }
        variant="small"
      />
    </div>
  );
};

export default WorkOrderTypeahead;
