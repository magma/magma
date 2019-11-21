/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {SurveyTemplateCategory} from '../../common/LocationType';

import * as React from 'react';
import Button from '@material-ui/core/Button';
import DeleteIcon from '@material-ui/icons/Delete';
import IconButton from '@material-ui/core/IconButton';
import SurveyTemplateQuestionsDialog from './SurveyTemplateQuestionsDialog';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import TextField from '@material-ui/core/TextField';

import inventoryTheme from '../../common/theme';
import nullthrows from '@fbcnms/util/nullthrows';
import {makeStyles} from '@material-ui/styles';
import {removeItem, updateItem} from '@fbcnms/util/arrays';
import {useState} from 'react';

const useStyles = makeStyles({
  table: {
    marginBottom: '12px',
  },
  input: {
    ...inventoryTheme.textField,
    marginTop: '0px',
    marginBottom: '0px',
  },
  cell: {
    ...inventoryTheme.textField,
    paddingLeft: '0px',
  },
  addButton: {
    marginBottom: '12px',
  },
});

type Props = {
  categories: SurveyTemplateCategory[],
  onCategoriesChanged: (SurveyTemplateCategory[]) => void,
};

export default function SurveyTemplateCategories(props: Props) {
  const classes = useStyles();
  const [editingCategory, setEditingCategory] = useState<?number>(null);
  const categories = props.categories || [];

  const onChange = (
    changedProp: 'categoryTitle' | 'categoryDescription',
    index,
  ) => event => {
    props.onCategoriesChanged(
      updateItem<SurveyTemplateCategory, typeof changedProp>(
        props.categories,
        index,
        changedProp,
        event.target.value,
      ),
    );
  };

  return (
    <div>
      <Table className={classes.table}>
        <TableHead>
          <TableRow>
            <TableCell className={classes.cell}>Category Title</TableCell>
            <TableCell className={classes.cell}>Category Description</TableCell>
            <TableCell className={classes.cell} />
            <TableCell align="right" />
          </TableRow>
        </TableHead>
        <TableBody>
          {categories.map((category, i) => (
            <TableRow key={category.id}>
              <TableCell className={classes.cell} scope="row">
                <TextField
                  placeholder="Title"
                  variant="outlined"
                  className={classes.input}
                  value={category.categoryTitle}
                  onChange={onChange('categoryTitle', i)}
                  margin="dense"
                />
              </TableCell>
              <TableCell className={classes.cell} scope="row">
                <TextField
                  placeholder="Description"
                  variant="outlined"
                  className={classes.input}
                  value={category.categoryDescription}
                  onChange={onChange('categoryDescription', i)}
                  margin="dense"
                />
              </TableCell>
              <TableCell className={classes.cell} scope="row">
                <Button
                  color="primary"
                  variant="outlined"
                  onClick={() => setEditingCategory(i)}>
                  {category.surveyTemplateQuestions?.length || 'No'}{' '}
                  {category.surveyTemplateQuestions?.length == 1
                    ? 'Question'
                    : 'Questions'}
                </Button>
              </TableCell>
              <TableCell align="right" className={classes.cell}>
                <IconButton
                  onClick={() =>
                    props.onCategoriesChanged(removeItem(props.categories, i))
                  }>
                  <DeleteIcon />
                </IconButton>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
      <Button
        className={classes.addButton}
        color="primary"
        variant="outlined"
        onClick={() =>
          props.onCategoriesChanged([
            ...props.categories,
            getEmptyCategory(categories.length),
          ])
        }>
        Add Category
      </Button>
      {editingCategory !== null && (
        <SurveyTemplateQuestionsDialog
          onClose={() => setEditingCategory(null)}
          onSave={questions => {
            props.onCategoriesChanged(
              updateItem<SurveyTemplateCategory, 'surveyTemplateQuestions'>(
                props.categories,
                nullthrows(editingCategory),
                'surveyTemplateQuestions',
                questions,
              ),
            );
            setEditingCategory(null);
          }}
          questions={getQuestions(categories, editingCategory)}
        />
      )}
    </div>
  );
}

function getEmptyCategory(index) {
  return {
    id: 'SurveyTemplateCategory@tmp' + index,
    categoryTitle: '',
    categoryDescription: '',
    surveyTemplateQuestions: [],
  };
}

function getQuestions(categories, index: ?number) {
  if (index == null) return [];
  return categories[index].surveyTemplateQuestions || [];
}
