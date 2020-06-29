/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {
  Entries,
  TokenizerDisplayProps,
} from '@fbcnms/ui/components/design-system/Token/Tokenizer';
import type {NamedNode} from './EntUtils';

import * as React from 'react';
import Tokenizer from '@fbcnms/ui/components/design-system/Token/Tokenizer';
import withSuspense from './withSuspense';
import {useCallback, useMemo, useState} from 'react';

export type ExactNamedNode = $Exact<NamedNode>;

const wrapAsEntries = (items: $ReadOnlyArray<ExactNamedNode>) =>
  (items || []).map(item => ({
    key: item.id,
    label: item.name,
    ...item,
  }));

type StaticNamedNodesTokenizerProps = $ReadOnly<{|
  ...TokenizerDisplayProps,
  allNamedNodes: $ReadOnlyArray<ExactNamedNode>,
  selectedNodeIds?: ?$ReadOnlyArray<string>,
  onSelectedNodeIdsChange?: ($ReadOnlyArray<string>) => void,
|}>;

function StaticNamedNodesTokenizer(props: StaticNamedNodesTokenizerProps) {
  const {
    allNamedNodes,
    onSelectedNodeIdsChange,
    selectedNodeIds,
    disabled: disabledProp,
    ...tokenizerDisplayProps
  } = props;

  const namedNodesEntries = useMemo(() => wrapAsEntries(allNamedNodes), [
    allNamedNodes,
  ]);

  const selectedNamedNodes = useMemo(() => {
    const ids: $ReadOnlyArray<string> = selectedNodeIds || [];
    const validIds = ids
      .map(id => allNamedNodes.find(node => node.id === id))
      .filter(Boolean);
    return wrapAsEntries(validIds);
  }, [allNamedNodes, selectedNodeIds]);

  const callOnSelectedNodeIdsChange = useCallback(
    (newEntries: Entries<ExactNamedNode>) => {
      if (!onSelectedNodeIdsChange) {
        return;
      }
      onSelectedNodeIdsChange(newEntries.map(lte => lte.id));
    },
    [onSelectedNodeIdsChange],
  );

  const [queryString, setQueryString] = useState('');

  const disabled = disabledProp || onSelectedNodeIdsChange == null;

  return (
    <Tokenizer
      {...tokenizerDisplayProps}
      disabled={disabled}
      tokens={selectedNamedNodes}
      onTokensChange={callOnSelectedNodeIdsChange}
      queryString={queryString}
      onQueryStringChange={setQueryString}
      dataSource={{
        fetchNetwork: () => Promise.resolve(namedNodesEntries),
      }}
    />
  );
}

export default withSuspense(StaticNamedNodesTokenizer);
