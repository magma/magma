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
import Button from '@fbcnms/ui/components/design-system/Button';
import DialogActions from '@material-ui/core/DialogActions';

type Props = {
  onAbort: () => void,
  onUpload: () => void,
};

const UploadAnywayDialog = (props: Props) => {
  return (
    <div>
      <DialogActions>
        <Button onClick={props.onAbort} skin="regular">
          Cancel
        </Button>
        <Button onClick={props.onUpload}>Upload Anyway</Button>
      </DialogActions>
    </div>
  );
};

export default UploadAnywayDialog;
