/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {SurveyTemplateCategory} from '../../common/LocationType';

import * as React from 'react';
import Button from '@fbcnms/ui/components/design-system/Button';
import DeleteIcon from '@fbcnms/ui/components/design-system/Icons/Actions/DeleteIcon';
import FormAction from '@fbcnms/ui/components/design-system/Form/FormAction';
import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import IconButton from '@fbcnms/ui/components/design-system/IconButton';
import PlusIcon from '@fbcnms/ui/components/design-system/Icons/Actions/PlusIcon';
import SurveyTemplateQuestionsDialog from './SurveyTemplateQuestionsDialog';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import fbt from 'fbt';

import inventoryTheme from '../../common/theme';
import nullthrows from '@fbcnms/util/nullthrows';
import {makeStyles} from '@material-ui/styles';
import {removeItem, updateItem} from '@fbcnms/util/arrays';
import {useState} from 'react';

const useStyles = makeStyles(() => ({
  table: inventoryTheme.table,
  input: {
    ...inventoryTheme.textField,
    marginTop: '0px',
    marginBottom: '0px',
  },
  cell: {
    paddingLeft: '0px',
  },
  addButton: {
    marginBottom: '12px',
  },
}));

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
                <FormField>
                  <TextInput
                    placeholder={`${fbt('Title', '')}`}
                    variant="outlined"
                    className={classes.input}
                    value={category.categoryTitle}
                    onChange={onChange('categoryTitle', i)}
                  />
                </FormField>
              </TableCell>
              <TableCell className={classes.cell} scope="row">
                <FormField>
                  <TextInput
                    placeholder={`${fbt('Description', '')}`}
                    variant="outlined"
                    className={classes.input}
                    value={category.categoryDescription}
                    onChange={onChange('categoryDescription', i)}
                  />
                </FormField>
              </TableCell>
              <TableCell className={classes.cell} scope="row">
                <FormAction>
                  <Button onClick={() => setEditingCategory(i)}>
                    {category.surveyTemplateQuestions?.length || 'No'}{' '}
                    {category.surveyTemplateQuestions?.length == 1
                      ? 'Question'
                      : 'Questions'}
                  </Button>
                </FormAction>
              </TableCell>
              <TableCell align="right" className={classes.cell}>
                <FormAction>
                  <IconButton
                    icon={DeleteIcon}
                    onClick={() =>
                      props.onCategoriesChanged(removeItem(props.categories, i))
                    }
                  />
                </FormAction>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
      <FormAction>
        <Button
          leftIcon={PlusIcon}
          variant="text"
          className={classes.addButton}
          onClick={() =>
            props.onCategoriesChanged([
              ...props.categories,
              getEmptyCategory(categories.length),
            ])
          }>
          Add Category
        </Button>
      </FormAction>
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
