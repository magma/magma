/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import CircularProgress from '@material-ui/core/CircularProgress';
import React, {useState} from 'react';
import RelayEnvironemnt from '../common/RelayEnvironment.js';
import {QueryRenderer} from 'react-relay';
import {makeStyles} from '@material-ui/styles';
import {useSnackbar} from '@fbcnms/ui/hooks';

type Props = {
  query: any,
  variables: Object,
  render: (props: Object) => React$Element<any> | null,
};

const useStyles = makeStyles(_theme => ({
  progress: {
    display: 'block',
    margin: '24px auto',
  },
}));

const InventoryQueryRenderer = (compProps: Props) => {
  const {query, variables} = compProps;
  const classes = useStyles();
  const [errorPresent, setErrorPresent] = useState(false);
  useSnackbar('Sorry, something went wrong', {variant: 'error'}, errorPresent);

  return (
    <QueryRenderer
      environment={RelayEnvironemnt}
      query={query}
      variables={variables}
      render={({error, props}) => {
        if (error) {
          setErrorPresent(true);
          return null;
        }

        if (!props) {
          return <CircularProgress className={classes.progress} />;
        }

        return compProps.render(props);
      }}
    />
  );
};

export default InventoryQueryRenderer;
