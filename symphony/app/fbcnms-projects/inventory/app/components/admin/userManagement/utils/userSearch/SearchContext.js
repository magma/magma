/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 *
 * @flow strict-local
 * @format
 */

import * as React from 'react';
import emptyFunction from '@fbcnms/util/emptyFunction';
import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from 'react';
import {debounce} from 'lodash';

export const NO_SEARCH_VALUE = '';
function getEmptyResults<T>(): Array<T> {
  return ([]: Array<T>);
}

type SearchContextValueType<T> = $ReadOnly<{|
  searchTerm: string,
  results: Array<T>,
  setSearchTerm: string => void,
  clearSearch: () => void,
  isSearchInProgress: boolean,
  isEmptySearchTerm: boolean,
|}>;

function useSearchManagerBuilder<M, R>(
  queryMetadata: ?M,
  searchCallback: (string, ?M) => Promise<Array<R>>,
): SearchContextValueType<R> {
  const EMPTY_SEARCH_RESULTS = useMemo(() => getEmptyResults<R>(), []);
  const [lastSearchedTerm, setLastSearchedTerm] = useState(NO_SEARCH_VALUE);
  const [searchTerm, setSearchTerm] = useState(NO_SEARCH_VALUE);
  const [results, setResults] = useState<Array<R>>(EMPTY_SEARCH_RESULTS);
  const [isSearchInProgress, setIsSearchInProgress] = useState(false);
  const [isEmptySearchTerm, setIsEmptySearchTerm] = useState(true);

  const runSearch = useCallback(
    debounce(currentSearchTerm => {
      setIsSearchInProgress(true);
      searchCallback(currentSearchTerm, queryMetadata)
        .then(setResults)
        .finally(() => setIsSearchInProgress(false));
    }, 200),
    [],
  );

  useEffect(() => {
    const actualSearchTerm = searchTerm.trim();
    if (actualSearchTerm == NO_SEARCH_VALUE) {
      setIsSearchInProgress(false);
      setResults(EMPTY_SEARCH_RESULTS);
      setLastSearchedTerm(NO_SEARCH_VALUE);
    } else if (actualSearchTerm == lastSearchedTerm) {
      setIsSearchInProgress(false);
    } else {
      runSearch(actualSearchTerm);
      setIsSearchInProgress(true);
      setLastSearchedTerm(actualSearchTerm);
    }
    setIsEmptySearchTerm(actualSearchTerm == NO_SEARCH_VALUE);
  }, [EMPTY_SEARCH_RESULTS, lastSearchedTerm, runSearch, searchTerm]);

  const clearSearch = useCallback(() => setSearchTerm(NO_SEARCH_VALUE), []);

  return {
    searchTerm,
    results,
    setSearchTerm,
    clearSearch,
    isSearchInProgress,
    isEmptySearchTerm,
  };
}

type ContextProviderProps<T> = $ReadOnly<{|
  children: React.Node,
  queryMetadata: ?T,
|}>;

export default function createSearchContext<
  QUERY_METADATA_TYPE,
  RESULT_VALUE_TYPE,
>(
  searchCallback: (
    string,
    ?QUERY_METADATA_TYPE,
  ) => Promise<Array<RESULT_VALUE_TYPE>>,
) {
  const SearchContext = createContext<
    $Exact<SearchContextValueType<RESULT_VALUE_TYPE>>,
  >({
    searchTerm: NO_SEARCH_VALUE,
    results: getEmptyResults<RESULT_VALUE_TYPE>(),
    setSearchTerm: emptyFunction,
    clearSearch: emptyFunction,
    isSearchInProgress: false,
    isEmptySearchTerm: false,
  });

  const useSearch = (queryMetadata?: ?QUERY_METADATA_TYPE) =>
    useSearchManagerBuilder<QUERY_METADATA_TYPE, RESULT_VALUE_TYPE>(
      queryMetadata,
      searchCallback,
    );

  const SearchContextProvider = (
    props: ContextProviderProps<QUERY_METADATA_TYPE>,
  ) => {
    const {children, queryMetadata} = props;

    return (
      <SearchContext.Provider value={useSearch(queryMetadata)}>
        {children}
      </SearchContext.Provider>
    );
  };

  const useSearchContext = () => {
    return useContext(SearchContext);
  };

  return {SearchContext, SearchContextProvider, useSearchContext, useSearch};
}
