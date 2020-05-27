/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {SurveyTemplateQuestion} from '../../common/LocationType';

import * as React from 'react';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import SurveyTemplateQuestionsTable from '../survey_templates/SurveyTemplateQuestionsTable';

import {useState} from 'react';

type Props = {
  onClose: () => void,
  onSave: (SurveyTemplateQuestion[]) => void,
  questions: SurveyTemplateQuestion[],
};

export default function QuestionsDialog(props: Props) {
  const [questions, setQuestions] = useState(props.questions);
  return (
    <Dialog onClose={props.onClose} open={true} maxWidth={false}>
      <DialogContent>
        <SurveyTemplateQuestionsTable
          questions={questions}
          onQuestionsChanged={q => setQuestions(q)}
        />
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose} color="primary">
          Cancel
        </Button>
        <Button
          onClick={() => props.onSave(questions)}
          color="primary"
          variant="contained">
          Save
        </Button>
      </DialogActions>
    </Dialog>
  );
}
