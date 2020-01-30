/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {SiteSurveyPane_survey} from './__generated__/SiteSurveyPane_survey.graphql.js';

import DateTimeFormat from '../../common/DateTimeFormat.js';
import MenuItem from '@material-ui/core/MenuItem';
import React, {useMemo, useState} from 'react';
import Select from '@material-ui/core/Select';
import SiteSurveyField from './SiteSurveyField';
import SiteSurveyQuestionReply from './SiteSurveyQuestionReply';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextField from '@material-ui/core/TextField';
import Typography from '@material-ui/core/Typography';
import classNames from 'classnames';
import {createFragmentContainer, graphql} from 'react-relay';
import {gray6} from '@fbcnms/ui/theme/colors';
import {makeStyles} from '@material-ui/styles';
import {uniqBy} from 'lodash';

const useStyles = makeStyles(theme => ({
  root: {
    minHeight: '400px',
    minWidth: '720px',
    display: 'flex',
    flexDirection: 'column',
  },
  headerDivider: {
    borderTop: `1px solid ${theme.palette.grey[200]}`,
    margin: '12px 0px',
  },
  header: {
    display: 'flex',
    borderBottom: `1px solid ${gray6}`,
    padding: '24px',
    alignItems: 'center',
  },
  body: {
    display: 'flex',
    flexDirection: 'column',
    flexGrow: 1,
  },
  title: {
    flexGrow: 1,
  },
  surveyName: {
    fontWeight: 500,
    fontSize: '20px',
    lineHeight: '24px',
    color: theme.palette.blueGrayDark,
  },
  completionDate: {
    fontSize: '14px',
    lineHeight: '24px',
    color: theme.palette.blueGrayDark,
  },
  firstField: {
    paddingTop: '24px',
    minWidth: '184px',
  },
  lastField: {
    paddingBottom: '24px',
    minWidth: '184px',
  },
  lastFieldContainer: {
    display: 'flex',
    flexGrow: 1,
  },
  reply: {
    paddingBottom: '16px',
  },
  select: {
    width: '224px',
  },
}));

type Props = {
  survey: SiteSurveyPane_survey,
};

const SiteSurveyPane = (props: Props) => {
  const {name, completionTimestamp, surveyResponses} = props.survey;
  const classes = useStyles();
  const forms = useMemo(
    () =>
      uniqBy(
        surveyResponses
          .filter(Boolean)
          .map(r => ({
            id: r.formIndex,
            name: r.formName,
            index: r.formIndex || 0,
          }))
          .slice()
          .sort((formA, formB) => formA.index - formB.index),
        'name',
      ),
    [surveyResponses],
  );
  const [selectedForm, setSelectedForm] = useState(
    forms.length > 0 ? forms[0].id : null,
  );
  const responses = useMemo(
    () =>
      surveyResponses
        .filter(Boolean)
        .filter(r => r.formIndex === selectedForm)
        .slice()
        .sort(
          (responseA, responseB) =>
            (responseA.questionIndex || 0) - (responseB.questionIndex || 0),
        ),
    [selectedForm, surveyResponses],
  );

  return (
    <div className={classes.root}>
      <div className={classes.header}>
        <div className={classes.title}>
          <Text className={classes.surveyName}>{name}</Text>
          <Typography className={classes.completionDate}>
            {'Completed: '}
            {DateTimeFormat.dateTime(completionTimestamp * 1000)}
          </Typography>
        </div>
        <Select
          className={classes.select}
          value={selectedForm}
          onChange={({target}) => setSelectedForm(target.value)}
          input={
            <TextField
              select
              value={selectedForm ?? ''}
              margin="dense"
              variant="outlined"
            />
          }>
          {forms.map(f => (
            <MenuItem value={f.id}>{f.name}</MenuItem>
          ))}
        </Select>
      </div>
      <div className={classes.body}>
        {responses.map((q, i) => (
          <div
            className={classNames({
              [classes.lastFieldContainer]: i === responses.length - 1,
            })}>
            <SiteSurveyField
              className={
                i === 0
                  ? classes.firstField
                  : i === responses.length - 1
                  ? classes.lastField
                  : null
              }
              key={q.id}
              label={q.questionText}>
              <div className={classes.reply}>
                <SiteSurveyQuestionReply question={q} />
              </div>
            </SiteSurveyField>
          </div>
        ))}
      </div>
    </div>
  );
};

export default createFragmentContainer(SiteSurveyPane, {
  survey: graphql`
    fragment SiteSurveyPane_survey on Survey {
      name
      completionTimestamp
      surveyResponses {
        id
        questionText
        formName
        formIndex
        questionIndex
        ...SiteSurveyQuestionReply_question
      }
    }
  `,
});
