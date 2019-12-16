/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ContextRouter} from 'react-router-dom';

import Button from '@fbcnms/ui/components/design-system/Button';
import React from 'react';
import {withRouter} from 'react-router-dom';

type Props = {
  id: string,
  title: string,
} & ContextRouter;

const LocationLink = (props: Props) => {
  return (
    <Button
      variant="text"
      onClick={() =>
        props.history.push(`/inventory/inventory?location=${props.id}`)
      }>
      {props.title}
    </Button>
  );
};

export default withRouter(LocationLink);
