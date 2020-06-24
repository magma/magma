/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {TokenizerDisplayProps} from '@fbcnms/ui/components/design-system/Token/Tokenizer';

import * as React from 'react';
import StaticNamedNodesTokenizer from './StaticNamedNodesTokenizer';
import {useProjectTemplateNodes} from './Project';

type ProjectTemplatesTokenizerProps = $ReadOnly<{|
  ...TokenizerDisplayProps,
  selectedProjectTemplateIds?: ?$ReadOnlyArray<string>,
  onSelectedProjectTemplateIdsChange?: ($ReadOnlyArray<string>) => void,
|}>;

function ProjectTemplatesTokenizer(props: ProjectTemplatesTokenizerProps) {
  const {
    selectedProjectTemplateIds,
    onSelectedProjectTemplateIdsChange,
    ...rest
  } = props;
  const projectTemplates = useProjectTemplateNodes();
  return (
    <StaticNamedNodesTokenizer
      allNamedNodes={projectTemplates}
      selectedNodeIds={selectedProjectTemplateIds}
      onSelectedNodeIdsChange={onSelectedProjectTemplateIdsChange}
      {...rest}
    />
  );
}

export default ProjectTemplatesTokenizer;
