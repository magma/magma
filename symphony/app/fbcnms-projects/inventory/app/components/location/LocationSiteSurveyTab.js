/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {LocationSiteSurveyTab_location} from './__generated__/LocationSiteSurveyTab_location.graphql.js';
import type {WithStyles} from '@material-ui/core';

import Card from '@fbcnms/ui/components/design-system/Card/Card';
import CardHeader from '@fbcnms/ui/components/design-system/Card/CardHeader';
import Dialog from '@material-ui/core/Dialog';
import DialogContent from '@material-ui/core/DialogContent';
import LocationSiteSurvey from './LocationSiteSurvey';
import React from 'react';
import RequestSiteSurveyLocationButton from './RequestSiteSurveyLocationButton';
import SiteSurveyPane from '../site_survey/SiteSurveyPane';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import classNames from 'classnames';
import nullthrows from '@fbcnms/util/nullthrows';
import {createFragmentContainer, graphql} from 'react-relay';
import {makeStyles} from '@material-ui/styles';
import {useState} from 'react';

type Props = {
  location: LocationSiteSurveyTab_location,
} & WithStyles<typeof styles>;

const styles = {
  content: {
    padding: '0px 24px 24px 24px',
  },
  dialogContent: {
    '&&': {padding: '0px 0px 0px 0px'},
  },
  cardHasNoContent: {
    marginBottom: '0px',
  },
};
const useStyles = makeStyles(styles);

function LocationSiteSurveyTab(props: Props) {
  const {location} = props;
  const {surveys} = location;
  const classes = useStyles();
  const [selectedSurvey, setSelectedSurvey] = useState(null);
  return (
    <Card>
      <CardHeader
        className={classNames({
          [classes.cardHasNoContent]: surveys.filter(Boolean).length === 0,
        })}
        rightContent={
          <RequestSiteSurveyLocationButton
            location={{
              id: location.id,
              siteSurveyNeeded: location.siteSurveyNeeded,
            }}
          />
        }>
        Site Surveys
      </CardHeader>
      {surveys.filter(Boolean).length > 0 ? (
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Name</TableCell>
              <TableCell>Owner</TableCell>
              <TableCell>Status</TableCell>
              <TableCell>Site Survey</TableCell>
              <TableCell />
            </TableRow>
          </TableHead>
          <TableBody>
            {surveys
              .filter(Boolean)
              .slice()
              .sort(
                (surveyA, surveyB) =>
                  surveyB.completionTimestamp - surveyA.completionTimestamp,
              )
              .map(survey => (
                <LocationSiteSurvey
                  survey={survey}
                  onSurveySelected={() => setSelectedSurvey(survey)}
                />
              ))}
          </TableBody>
        </Table>
      ) : null}
      {selectedSurvey !== null ? (
        <Dialog
          maxWidth="md"
          onClose={() => setSelectedSurvey(null)}
          open={true}>
          <DialogContent className={classes.dialogContent}>
            <SiteSurveyPane survey={nullthrows(selectedSurvey)} />
          </DialogContent>
        </Dialog>
      ) : null}
    </Card>
  );
}

export default createFragmentContainer(LocationSiteSurveyTab, {
  location: graphql`
    fragment LocationSiteSurveyTab_location on Location {
      id
      siteSurveyNeeded
      surveys {
        id
        completionTimestamp
        name
        ownerName
        sourceFile {
          id
          fileName
          storeKey
        }
        ...SiteSurveyPane_survey
      }
    }
  `,
});
