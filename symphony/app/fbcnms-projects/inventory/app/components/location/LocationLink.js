/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {ContextRouter} from 'react-router-dom';

import Button from '@fbcnms/ui/components/design-system/Button';
import React from 'react';
import {InventoryAPIUrls} from '../../common/InventoryAPI';
import {withRouter} from 'react-router-dom';

type Props = {
  id: string,
  title: string,
  newTab?: boolean,
} & ContextRouter;

const LocationLink = (props: Props) => {
  const {id, newTab = false} = props;
  return (
    <Button
      variant="text"
      onClick={() =>
        newTab
          ? window.open(InventoryAPIUrls.location(id))
          : props.history.push(`/inventory/inventory?location=${id}`)
      }>
      {props.title}
    </Button>
  );
};

export default withRouter(LocationLink);
