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

import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import {useAxios} from '@fbcnms/ui/hooks';

export type Props = {|
  endpoint: string,
|};

export default function FBCMobileAppQRCode(props: Props) {
  const {isLoading, response} = useAxios({
    method: 'GET',
    url: props.endpoint,
  });

  if (isLoading || !response) {
    return <LoadingFiller />;
  }

  return (
    <img
      src={response.data}
      style={{
        height: 250,
        width: 250,
      }}
    />
  );
}
