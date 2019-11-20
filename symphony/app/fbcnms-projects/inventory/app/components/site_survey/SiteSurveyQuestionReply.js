/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {
  SiteSurveyQuestionReply_question,
  SurveyQuestionType,
} from './__generated__/SiteSurveyQuestionReply_question.graphql.js';

import * as React from 'react';
import MapView from 'inventory/app/components/map/MapView';
import SiteSurveyQuestionReplyCellData from './SiteSurveyQuestionReplyCellData';
import SiteSurveyQuestionReplyWifiData from './SiteSurveyQuestionReplyWifiData';
import Text from '@fbcnms/ui/components/design-system/Text';

import nullthrows from '@fbcnms/util/nullthrows';
import {DocumentAPIUrls} from '../../common/DocumentAPI';
import {coordsToGeoJson} from '../map/MapUtil';
import {createFragmentContainer, graphql} from 'react-relay';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  questionImg: {
    maxWidth: '400px',
    maxHeight: '200px',
  },
  mapContainer: {
    height: '200px',
  },
  text: {
    color: theme.palette.blueGrayDark,
    fontSize: '14px',
    lineHeight: '24px',
    fontWeight: 500,
  },
}));

type Props = {
  question: SiteSurveyQuestionReply_question,
};

function SiteSurveyQuestionReply(props: Props) {
  const {question} = props;
  const questionFormat: ?SurveyQuestionType = question.questionFormat;
  const classes = useStyles();

  switch (questionFormat) {
    case 'BOOL':
      return (
        <Text className={classes.text}>{question.boolData ? 'Yes' : 'No'}</Text>
      );
    case 'EMAIL':
      return <Text className={classes.text}>{question.emailData}</Text>;
    case 'PHONE':
      return <Text className={classes.text}>{question.phoneData}</Text>;
    case 'TEXT':
    case 'TEXTAREA':
      return <Text className={classes.text}>{question.textData}</Text>;
    case 'PHOTO':
      const storeKey = question.photoData?.storeKey;
      return storeKey ? (
        <img
          className={classes.questionImg}
          src={DocumentAPIUrls.get_url(storeKey)}
        />
      ) : null;
    case 'COORDS':
      return (
        <div className={classes.mapContainer}>
          <MapView
            id="mapView"
            center={{lat: question.latitude, lng: question.longitude}}
            mode="streets"
            zoomLevel="14"
            markers={coordsToGeoJson(
              nullthrows(question.latitude),
              nullthrows(question.longitude),
            )}
          />
        </div>
      );
    case 'WIFI':
      return <SiteSurveyQuestionReplyWifiData data={question} />;
    case 'CELLULAR':
      return <SiteSurveyQuestionReplyCellData data={question} />;
    case 'FLOAT':
      return <Text className={classes.text}>{question.floatData}</Text>;
    case 'INTEGER':
      return <Text className={classes.text}>{question.intData}</Text>;
    case 'DATE':
      return (
        <Text className={classes.text}>
          {question.dateData != null &&
            new Intl.DateTimeFormat('en-US').format(
              new Date(question.dateData * 1000),
            )}
        </Text>
      );
  }

  return null;
}

export default createFragmentContainer(SiteSurveyQuestionReply, {
  question: graphql`
    fragment SiteSurveyQuestionReply_question on SurveyQuestion {
      questionFormat
      longitude
      latitude
      boolData
      textData
      emailData
      phoneData
      floatData
      intData
      dateData
      photoData {
        storeKey
      }
      ...SiteSurveyQuestionReplyWifiData_data
      ...SiteSurveyQuestionReplyCellData_data
    }
  `,
});
