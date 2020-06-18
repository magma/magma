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

export const DEACTIVATED_PAGE_PATH = '/deactivated';

export default function DeactivatedPage() {
  return (
    <>
      <div>
        Your user had been deactivated. Contact you system administrator.
      </div>
      <div>
        <a href="/user/logout">logout</a>
      </div>
    </>
  );
}
