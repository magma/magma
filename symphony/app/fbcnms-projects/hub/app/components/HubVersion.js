/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import * as React from 'react';
import RelayEnvironment from '../common/RelayEnvironment';
import Text from '@fbcnms/ui/components/design-system/Text';
import type {HubVersionQueryResponse} from './__generated__/HubVersionQuery.graphql';
// flowlint untyped-import:warn
import {QueryRenderer, graphql} from 'react-relay';

const HubVersionQuery = graphql`
  query HubVersionQuery {
    version {
      string
    }
  }
`;

function HubVersion() {
  return (
    <QueryRenderer
      environment={RelayEnvironment}
      query={HubVersionQuery}
      variables={{}}
      render={({
        error,
        props,
      }: {
        error?: Error,
        props?: HubVersionQueryResponse,
      }) => {
        if (error) {
          return <Text>Couldn't load hub version!</Text>;
        }
        if (!props) {
          return <Text>Loading hub version...</Text>;
        }
        return <Text>Hub version: {props.version.string}</Text>;
      }}
    />
  );
}

export default HubVersion;
