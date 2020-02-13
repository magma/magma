/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {SurveyQuestionType} from '../configure/__generated__/AddEditLocationTypeCard_editingLocationType.graphql.js';
import type {SurveyTemplateQuestion} from '../../common/LocationType';

import * as React from 'react';
import AddCircleOutline from '@material-ui/icons/AddCircleOutline';
import DeleteIcon from '@material-ui/icons/Delete';
import DraggableTableRow from '../draggable/DraggableTableRow';
import DroppableTableBody from '../draggable/DroppableTableBody';
import IconButton from '@material-ui/core/IconButton';
import MenuItem from '@material-ui/core/MenuItem';
import Table from '@material-ui/core/Table';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import TextField from '@material-ui/core/TextField';

import {assertEnum} from '@fbcnms/util/enums';
import {makeStyles} from '@material-ui/styles';
import {removeItem, updateItem} from '@fbcnms/util/arrays';
import {reorder} from '../draggable/DraggableUtils';

const useStyles = makeStyles({
  root: {
    marginBottom: '12px',
  },
  input: {
    width: '200px',
    marginTop: '0px',
    marginBottom: '0px',
  },
  cell: {
    width: '200px',
    paddingRight: '5px',
  },
});

const questionTypes: {[SurveyQuestionType]: string} = {
  BOOL: 'Yes/No',
  CELLULAR: 'Cellular Signal Scan',
  COORDS: 'Coordinates',
  DATE: 'Date',
  FLOAT: 'Float',
  EMAIL: 'E-mail',
  INTEGER: 'Integer',
  PHONE: 'Phone Number',
  PHOTO: 'Photo',
  TEXT: 'Text',
  TEXTAREA: 'Text Area',
  WIFI: 'Wi-Fi Signal Scan',
};

type Props = {
  questions: SurveyTemplateQuestion[],
  onQuestionsChanged: (SurveyTemplateQuestion[]) => void,
};

export default function SurveyTemplateQuestionsTable(props: Props) {
  const classes = useStyles();
  const {questions} = props;

  const onChangeQuestionType = index => event => {
    props.onQuestionsChanged(
      updateItem<SurveyTemplateQuestion, 'questionType'>(
        props.questions,
        index,
        'questionType',
        assertEnum(questionTypes, event.target.value),
      ),
    );
  };

  const onChange = (
    changedProp: 'questionTitle' | 'questionDescription',
    index,
  ) => event => {
    const value = event.target.value;
    props.onQuestionsChanged(
      updateItem<SurveyTemplateQuestion, typeof changedProp>(
        props.questions,
        index,
        changedProp,
        value,
      ),
    );
  };

  const onDragEnd = result => {
    if (!result.destination) {
      return;
    }

    const items: SurveyTemplateQuestion[] = reorder(
      props.questions,
      result.source.index,
      result.destination.index,
    );

    const newItems = items.map((question, i) => ({...question, index: i}));
    props.onQuestionsChanged(newItems);
  };

  return (
    <div>
      <Table component="div" className={classes.root}>
        <TableHead component="div">
          <TableRow component="div">
            <TableCell size="small" padding="none" component="div" />
            <TableCell component="div" className={classes.cell}>
              Question Title
            </TableCell>
            <TableCell component="div" className={classes.cell}>
              Question Description
            </TableCell>
            <TableCell component="div" className={classes.cell}>
              Question Type
            </TableCell>
            <TableCell component="div" align="right">
              <IconButton
                onClick={() =>
                  props.onQuestionsChanged([
                    ...props.questions,
                    getEmptyQuestion(questions.length),
                  ])
                }>
                <AddCircleOutline />
              </IconButton>
            </TableCell>
          </TableRow>
        </TableHead>
        <DroppableTableBody onDragEnd={onDragEnd}>
          {questions.map((question, i) => (
            <DraggableTableRow id={question.id} index={i} key={i}>
              <TableCell className={classes.cell} component="div" scope="row">
                <TextField
                  placeholder="Title"
                  variant="outlined"
                  className={classes.input}
                  value={question.questionTitle}
                  onChange={onChange('questionTitle', i)}
                  margin="dense"
                />
              </TableCell>
              <TableCell className={classes.cell} component="div" scope="row">
                <TextField
                  placeholder="Description"
                  variant="outlined"
                  className={classes.input}
                  value={question.questionDescription}
                  onChange={onChange('questionDescription', i)}
                  margin="dense"
                />
              </TableCell>
              <TableCell className={classes.cell} component="div" scope="row">
                <TextField
                  select
                  variant="outlined"
                  className={classes.input}
                  value={question.questionType}
                  onChange={onChangeQuestionType(i)}
                  SelectProps={{
                    MenuProps: {
                      className: classes.menu,
                    },
                  }}
                  margin="dense">
                  {Object.keys(questionTypes).map(type => (
                    <MenuItem key={type} value={type}>
                      {questionTypes[type]}
                    </MenuItem>
                  ))}
                </TextField>
              </TableCell>
              <TableCell align="right" component="div">
                <IconButton
                  onClick={() =>
                    props.onQuestionsChanged(removeItem(props.questions, i))
                  }>
                  <DeleteIcon />
                </IconButton>
              </TableCell>
            </DraggableTableRow>
          ))}
        </DroppableTableBody>
      </Table>
    </div>
  );
}

function getEmptyQuestion(index) {
  return {
    id: 'SurveyTemplateQuestion@tmp' + index,
    questionTitle: '',
    questionDescription: '',
    questionType: 'BOOL',
    index,
  };
}
