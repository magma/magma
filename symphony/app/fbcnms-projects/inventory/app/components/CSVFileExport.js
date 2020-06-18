/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {FiltersQuery} from './comparison_view/ComparisonViewTypes';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithStyles} from '@material-ui/core';

import Button from '@fbcnms/ui/components/design-system/Button';
import CircularProgress from '@material-ui/core/CircularProgress';
import React, {useState} from 'react';
import axios from 'axios';
import classNames from 'classnames';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {withStyles} from '@material-ui/core/styles';

const styles = {
  exportButton: {
    paddingLeft: '8px',
    paddingRight: '8px',
  },
  exportButtonContainer: {
    display: 'flex',
  },
  exportButtonContent: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    '& $hiddenContent': {
      maxHeight: '0px',
      overflowY: 'hidden',
    },
  },
  hiddenContent: {},
};
const PATH_PREFIX = '/graph/export';

type Props = {
  exportPath: string,
  title: string,
  filters: ?FiltersQuery,
} & WithStyles<typeof styles> &
  WithAlert;

const CSVFileExport = (props: Props) => {
  const {classes, title, exportPath} = props;
  const [isDownloading, setIsDownloading] = useState(false);

  const filters = props.filters?.map(f => {
    if (f.name == 'property') {
      const property = f.propertyValue;
      f.propertyValue = property;
    }
    return f;
  });

  const onClick = async () => {
    const path = PATH_PREFIX + exportPath;
    const fileName = exportPath.replace('/', '').replace(/\//g, '_') + '.csv';
    setIsDownloading(true);
    try {
      await axios
        .get(path, {
          params: {
            filters: JSON.stringify(filters),
          },
          responseType: 'blob',
        })
        .then(response => {
          setIsDownloading(false);
          const url = window.URL.createObjectURL(new Blob([response.data]));
          const link = document.createElement('a');
          link.href = url;
          link.setAttribute('download', fileName);
          link.click();
        });
    } catch (error) {
      props.alert(error.response?.data?.error || error);
      setIsDownloading(false);
    }
  };
  return (
    <div className={classes.exportButtonContainer}>
      <Button className={classes.exportButton} variant="text" onClick={onClick}>
        <div className={classes.exportButtonContent}>
          <span
            className={classNames({
              [classes.hiddenContent]: isDownloading,
            })}>
            {title}
          </span>
          <CircularProgress
            size={20}
            color="inherit"
            className={classNames({
              [classes.hiddenContent]: !isDownloading,
            })}
          />
        </div>
      </Button>
    </div>
  );
};

export default withStyles(styles)(withAlert(CSVFileExport));
