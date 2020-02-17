/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import Button from '@material-ui/core/Button';
import CircularProgress from '@material-ui/core/CircularProgress';
import React from 'react';
import RelayEnvironment from '../../common/RelayEnvironment';
import Text from '@fbcnms/ui/components/design-system/Text';
import Typography from '@material-ui/core/Typography';
import {QueryRenderer} from 'react-relay';
import {graphql} from 'relay-runtime';
import {makeStyles} from '@material-ui/styles';
import {useHistory} from 'react-router';

const useStyles = makeStyles(theme => ({
  root: {
    display: 'flex',
    flexDirection: 'column',
    flexGrow: 1,
  },
  vertexTitle: {
    fontSize: '16px',
    lineHeight: '24px',
    color: theme.palette.blueGrayDark,
    fontWeight: 500,
  },
  buttonsContainer: {
    display: 'flex',
    flexDirection: 'column',
    marginTop: '16px',
  },
  button: {
    width: '100%',
    background: '#3984FF',
    color: 'white',
    fontSize: '14px',
    fontWeight: 500,
    lineHeight: '16px',
    padding: '8px 16px',
    '&:hover': {
      background: '#3984FF',
    },
  },
  progressContainer: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    padding: '16px',
  },
  field: {
    color: '#303846',
    fontSize: '12px',
    fontWeight: 400,
    lineHeight: '20px',
  },
  detailsContainer: {
    flexGrow: 1,
  },
  fieldName: {
    fontWeight: 500,
  },
  vertexType: {
    color: '#606770',
    fontSize: '12px',
    marginBottom: '16px',
  },
}));

const vertexQuery = graphql`
  query GraphVertexDetailsQuery($id: ID!) {
    vertex(id: $id) {
      id
      type
      fields {
        name
        value
      }
    }
  }
`;

type Props = {
  vertexId: string,
};

const GraphVertexDetails = ({vertexId}: Props) => {
  const classes = useStyles();
  const history = useHistory();
  return (
    <div className={classes.root}>
      <div className={classes.detailsContainer}>
        <QueryRenderer
          environment={RelayEnvironment}
          query={vertexQuery}
          variables={{id: vertexId}}
          render={({props}) => {
            if (!props || !props.vertex) {
              return (
                <div className={classes.progressContainer}>
                  <CircularProgress />
                </div>
              );
            }

            return (
              <div>
                <Text className={classes.vertexTitle}>Vertex {vertexId}</Text>
                <Text className={classes.vertexType}>{props.vertex.type}</Text>
                {props.vertex.fields.map(field => (
                  <div>
                    <Typography className={classes.field}>
                      <span className={classes.fieldName}>{field.name}:</span>{' '}
                      {field.value}
                    </Typography>
                  </div>
                ))}
              </div>
            );
          }}
        />
      </div>
      <div className={classes.buttonsContainer}>
        <Button
          variant="outlined"
          className={classes.button}
          onClick={() => history.push(`/id/${vertexId}`)}>
          View
        </Button>
      </div>
    </div>
  );
};

export default GraphVertexDetails;
