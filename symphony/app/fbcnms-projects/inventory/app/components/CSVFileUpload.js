/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {WithStyles} from '@material-ui/core';

import * as React from 'react';
import {withStyles} from '@material-ui/core/styles';

const styles = {
  hiddenInput: {
    width: '0px',
    height: '0px',
    opacity: 0,
    overflow: 'hidden',
    position: 'absolute',
    zIndex: -1,
  },
  fileButton: {
    cursor: 'pointer',
    width: 'fit-content',
  },
};

type Props = {
  button: React.Element<any>,
  uploadPath: string,
  entity: ?string,
  onFileChanged: (SyntheticInputEvent<HTMLInputElement>, ?string) => void,
} & WithStyles<typeof styles>;

class CSVFileUpload extends React.Component<Props> {
  input = null;

  constructor(props: Props) {
    super(props);
    this.input = null;
  }

  render() {
    const {button, classes, entity} = this.props;
    return (
      <>
        <input
          className={classes.hiddenInput}
          type="file"
          accept=".csv"
          onChange={e => {
            this.props.onFileChanged(e, entity);
          }}
          multiple
          ref={ref => {
            this.input = ref;
          }}
        />
        <div className={classes.fileButton} onClick={this.openFileDialog}>
          {button}
        </div>
      </>
    );
  }

  openFileDialog = () => {
    if (this.input) {
      this.input.click();
    }
  };
}

export default withStyles(styles)(CSVFileUpload);
