/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import * as React from 'react';
import ServiceLinksTable from './ServiceLinksTable';
import ServiceLinksView_links from './__generated__/ServiceLinksView_links.graphql';
import {createFragmentContainer, graphql} from 'react-relay';

type Props = {
  links: ServiceLinksView_links,
};

const ServiceLinksView = (props: Props) => {
  const {links} = props;
  return <ServiceLinksTable links={links} />;
};

export default createFragmentContainer(ServiceLinksView, {
  links: graphql`
    fragment ServiceLinksView_links on Link @relay(plural: true) {
      id
      ports {
        parentEquipment {
          id
          name
        }
        definition {
          id
          name
        }
      }
    }
  `,
});
