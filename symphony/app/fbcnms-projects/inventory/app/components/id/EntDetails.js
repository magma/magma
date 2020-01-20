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
import EdgesTable from './EdgesTable';
import EntGraph from './EntGraph';
import EntViewToggleButton from './EntViewToggleButton';
import ErrorIcon from '@material-ui/icons/Error';
import FieldsTable from './FieldsTable';
import React, {useState} from 'react';
import RelayEnvironment from '../../common/RelayEnvironment';
import Text from '@fbcnms/ui/components/design-system/Text';
import nullthrows from '@fbcnms/util/nullthrows';
import {QueryRenderer} from 'react-relay';
import {graphql} from 'relay-runtime';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  root: {
    padding: '32px',
  },
  titleContainer: {
    marginBottom: '16px',
    display: 'flex',
  },
  entTypeTitle: {
    fontSize: '20px',
    lineHeight: '24px',
    color: theme.palette.blueGrayDark,
    fontWeight: 500,
  },
  subtitle: {
    fontSize: '16px',
    lineHeight: '24px',
    color: theme.palette.blueGrayDark,
  },
  subtitleContainer: {
    margin: '16px 0px',
  },
  progressContainer: {
    flexGrow: 1,
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    color: theme.palette.blueGrayDark,
  },
  entId: {
    fontSize: '12px',
    lineHeight: '26px',
    color: theme.palette.grey[600],
    marginLeft: '8px',
  },
  errorIcon: {
    marginRight: '8px',
  },
  section: {
    marginBottom: '32px',
  },
  titleText: {
    display: 'flex',
    flexGrow: 1,
  },
}));

type Props = {};

const idQuery = graphql`
  query EntDetailsQuery($id: ID!) {
    vertex(id: $id) {
      id
      type
      fields {
        name
        value
        type
      }
      edges {
        name
        type
        ids
      }
    }
  }
`;

const EntDetails = (_props: Props) => {
  const {match} = useRouter();
  const id = nullthrows(match.params.id);
  const classes = useStyles();
  const [selectedView, setSelectedView] = useState<'details' | 'graph'>(
    'details',
  );

  return (
    <QueryRenderer
      environment={RelayEnvironment}
      query={idQuery}
      variables={{id}}
      render={({props, error}) => {
        if (!props) {
          return (
            <div className={classes.progressContainer}>
              <CircularProgress className={classes.progress} />
            </div>
          );
        }

        const {vertex} = props;
        if (error || !vertex) {
          return (
            <div className={classes.progressContainer}>
              <ErrorIcon className={classes.errorIcon} />
              <Text>Ent with the ID {id} was not found</Text>
            </div>
          );
        }

        return (
          <div className={classes.root}>
            <div className={classes.titleContainer}>
              <div className={classes.titleText}>
                <Text className={classes.entTypeTitle}>{vertex.type}</Text>
                <Text className={classes.entId}>({id})</Text>
              </div>
              <EntViewToggleButton
                selectedView={selectedView}
                onViewSelected={view => setSelectedView(view)}
              />
            </div>
            <div>
              {selectedView === 'details' ? (
                <>
                  <div className={classes.subtitleContainer}>
                    <Text className={classes.subtitle}>Fields</Text>
                  </div>
                  <div className={classes.section}>
                    <FieldsTable fields={vertex.fields} />
                  </div>
                  <div className={classes.subtitleContainer}>
                    <Text className={classes.subtitle}>Edges</Text>
                  </div>
                  <div className={classes.section}>
                    <EdgesTable edges={vertex.edges} />
                  </div>
                </>
              ) : null}
              {selectedView === 'graph' ? (
                <>
                  <div className={classes.subtitleContainer}>
                    <Text className={classes.subtitle}>Graph</Text>
                  </div>
                  <div className={classes.section}>
                    <EntGraph rootNode={vertex} />
                  </div>
                </>
              ) : null}
            </div>
          </div>
        );
      }}
    />
  );
};

export default EntDetails;
