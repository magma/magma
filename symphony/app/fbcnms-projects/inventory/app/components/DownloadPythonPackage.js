/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import Button from '@fbcnms/ui/components/design-system/Button';
import InventoryQueryRenderer from './InventoryQueryRenderer';
import React from 'react';
import axios from 'axios';
import {DocumentAPIUrls} from '../common/DocumentAPI';
import {graphql} from 'relay-runtime';
import {makeStyles} from '@material-ui/styles';

type PythonPackage = {
  version: string,
  whlFileKey: string,
};

const useStyles = makeStyles(_ => ({
  link: {
    padding: '20px',
  },
}));

const DownloadPythonPackageQuery = graphql`
  query DownloadPythonPackageQuery {
    latestPythonPackage {
      lastPythonPackage {
        version
        whlFileKey
      }
    }
  }
`;

const handleDownload = (key: string, fileName: string) => {
  const method = 'GET';
  const url = DocumentAPIUrls.download_url(key, fileName);

  axios
    .request({
      url,
      method,
      responseType: 'blob',
      headers: {
        'Is-Global': 'True',
      },
    })
    .then(({data}) => {
      const downloadUrl = window.URL.createObjectURL(new Blob([data]));
      const link = document.createElement('a');
      link.href = downloadUrl;
      link.setAttribute('download', fileName);
      if (document.body != null) {
        document.body.appendChild(link);
      }
      link.click();
      link.remove();
    });
};

const handleWhlDownload = (pythonPackage: PythonPackage) => {
  handleDownload(
    pythonPackage.whlFileKey,
    `pyinventory-${pythonPackage.version}-py3-none-any.whl`,
  );
};

const DownloadPythonPackage = () => {
  const classes = useStyles();
  return (
    <InventoryQueryRenderer
      query={DownloadPythonPackageQuery}
      variables={{}}
      render={props => {
        const pythonPackage = props.latestPythonPackage.lastPythonPackage;

        if (pythonPackage == null) {
          return <div />;
        }

        return (
          <div className={classes.link}>
            <Button
              variant="text"
              onClick={() => handleWhlDownload(pythonPackage)}>
              Python Package {pythonPackage.version}
            </Button>
          </div>
        );
      }}
    />
  );
};

export default DownloadPythonPackage;
