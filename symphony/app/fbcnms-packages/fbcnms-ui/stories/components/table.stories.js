/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import type {RowsSeparationTypes} from '../../components/design-system/Table/TableContent';
import type {
  TableRowDataType,
  TableVariantTypes,
} from '../../components/design-system/Table/Table';

import Button from '../../components/design-system/Button';
import Checkbox from '../../components/design-system/Checkbox/Checkbox';
import RadioGroup from '../../components/design-system/RadioGroup/RadioGroup';
import React, {useMemo, useState} from 'react';
import Table from '../../components/design-system/Table/Table';
import Text from '../../components/design-system/Text';
import ThreeDotsVerticalIcon from '../../components/design-system/Icons/Actions/ThreeDotsVerticalIcon';
import {ROW_SEPARATOR_TYPES} from '../../components/design-system/Table/TableContent';
import {STORY_CATEGORIES} from '../storybookUtils';
import {TABLE_SORT_ORDER} from '../../components/design-system/Table/TableContext';
import {TABLE_VARIANT_TYPES} from '../../components/design-system/Table/Table';
import {makeStyles} from '@material-ui/styles';
import {storiesOf} from '@storybook/react';

type DataType = {|
  title?: string,
  firstName: string,
  lastName: string,
  startingDate: Date,
  age: number,
  city: string,
|};
type RowDataType = TableRowDataType<DataType>;

const DATA: Array<RowDataType> = [
  {
    key: '1',
    firstName: 'Meghan',
    lastName: 'Bishop',
    age: 32,
    startingDate: new Date('Febuary 13, 2020'),
    city: 'Tel Aviv',
  },
  {
    key: '2',
    title: 'Dr.',
    firstName: 'Sara',
    lastName: 'Porter',
    age: 21,
    startingDate: new Date('Febuary 28, 1999'),
    city: 'Raanana',
  },
  {
    key: '3',
    title: 'Don',
    firstName: 'Dolev',
    lastName: 'Hadar',
    age: 22,
    startingDate: new Date('May 02, 1990'),
    city: 'Tel Aviv',
  },
  {
    key: '4',
    title: 'Mr.',
    firstName: 'Walter',
    lastName: 'Jenning',
    age: 76,
    startingDate: new Date('July 11, 2001'),
    city: 'Ramat Gan',
  },
];

const useStyles = makeStyles(_theme => ({
  root: {
    width: '100%',
  },
  table: {
    marginBottom: '24px',
  },
  optionsContainer: {
    display: 'flex',
    flexDirection: 'column',
    marginTop: '32px',
  },
  displayOption: {
    marginTop: '4px',
    display: 'flex',
    alignItems: 'center',
  },
  displayMenuOption: {
    marginTop: '4px',
    display: 'flex',
    alignItems: 'top',
  },
  optionCheckbox: {
    marginRight: '8px',
  },
  iconColumn: {
    width: '36px',
  },
}));

const TablesRoot = () => {
  const classes = useStyles();
  const [showSelection, setShowSelection] = useState(false);
  const [selectedIds, setSelectedIds] = useState([]);

  const [showSorting, setShowSorting] = useState(false);

  const columns = useMemo(
    () => [
      {
        key: 'title',
        title: 'Title',
        render: row => row.title || '',
        getSortingValue: showSorting ? row => row.title : undefined,
      },
      {
        key: 'firstName',
        title: 'First Name',
        render: row => row.firstName,
        getSortingValue: showSorting ? row => row.firstName : undefined,
      },
      {
        key: 'lastName',
        title: 'Last Name',
        render: row => row.lastName,
      },
      {
        key: 'age',
        title: 'Age',
        render: row => row.age,
        getSortingValue: showSorting ? row => row.age : undefined,
      },
      {
        key: 'startingDate',
        title: 'Starting Date',
        render: row => Intl.DateTimeFormat('default').format(row.startingDate),
        getSortingValue: showSorting
          ? row => row.startingDate.getTime()
          : undefined,
      },
      {
        key: 'city',
        title: 'City',
        render: row => (
          <Button variant="text" onClick={() => alert(`clicked ${row.city}`)}>
            {row.city}
          </Button>
        ),
        getSortingValue: showSorting ? row => row.city : undefined,
      },
      {
        key: 'menu_icon',
        title: '',
        titleClassName: classes.iconColumn,
        className: classes.iconColumn,
        render: _row => (
          <Button variant="text" onClick={() => alert(`menu opening`)}>
            <ThreeDotsVerticalIcon color="gray" />
          </Button>
        ),
      },
    ],
    [classes.iconColumn, showSorting],
  );

  const [rowsSeparator, setRowsSeparator] = useState<RowsSeparationTypes>(
    ROW_SEPARATOR_TYPES.bands,
  );
  const [tableVariant, setTableVariant] = useState<TableVariantTypes>(
    TABLE_VARIANT_TYPES.standalone,
  );

  const [showActiveRow, setShowActiveRow] = useState(false);
  const [activeRowId, setActiveRowId] = useState(null);

  const [showDetailsCard, setShowDetailsCard] = useState(false);

  const tableProps = useMemo(
    () => ({
      data: DATA,
      columns: columns,
      variant: tableVariant,
      dataRowsSeparator: rowsSeparator,
      sortSettings: showSorting
        ? {
            columnKey: 'title',
            order: TABLE_SORT_ORDER.ascending,
          }
        : undefined,
      showSelection: showSelection,
      selectedIds: showSelection ? selectedIds : undefined,
      onSelectionChanged: showSelection ? setSelectedIds : undefined,
      activeRowId: showActiveRow ? activeRowId : undefined,
      onActiveRowIdChanged: showActiveRow ? setActiveRowId : undefined,
      detailsCard: showDetailsCard ? (
        <div>
          <div>
            <Text variant="h6">Here you can show some intersting details</Text>
          </div>
          <div>
            <Text variant="subtitle2">Usually be used with 'activeRow'</Text>
          </div>
        </div>
      ) : (
        undefined
      ),
    }),
    [
      activeRowId,
      columns,
      rowsSeparator,
      selectedIds,
      showActiveRow,
      showDetailsCard,
      showSelection,
      showSorting,
      tableVariant,
    ],
  );

  return (
    <div className={classes.root}>
      <div className={classes.table}>
        <Table {...tableProps} />
      </div>
      <div className={classes.optionsContainer}>
        <div className={classes.displayOption}>
          <Checkbox
            className={classes.optionCheckbox}
            checked={showSorting}
            onChange={selection =>
              setShowSorting(selection === 'checked' ? true : false)
            }
          />
          <Text>With Sorting</Text>
        </div>
        <div className={classes.displayOption}>
          <Checkbox
            className={classes.optionCheckbox}
            checked={showSelection}
            onChange={selection =>
              setShowSelection(selection === 'checked' ? true : false)
            }
          />
          <Text>With Selection</Text>
        </div>
        <div className={classes.displayOption}>
          <Checkbox
            className={classes.optionCheckbox}
            checked={showActiveRow}
            onChange={selection =>
              setShowActiveRow(selection === 'checked' ? true : false)
            }
          />
          <Text>Row can be active (clickable)</Text>
        </div>
        <div className={classes.displayOption}>
          <Checkbox
            className={classes.optionCheckbox}
            checked={showDetailsCard}
            onChange={selection =>
              setShowDetailsCard(selection === 'checked' ? true : false)
            }
          />
          <Text>Details Card Shown</Text>
        </div>
        <div className={classes.displayMenuOption}>
          <div>
            <Text>Row Separation Type: </Text>
          </div>
          <RadioGroup
            options={[
              {
                value: 'bands',
                label: `'${ROW_SEPARATOR_TYPES.bands}'`,
                details: 'Rows are banded with stripes',
              },
              {
                value: 'border',
                label: `'${ROW_SEPARATOR_TYPES.border}'`,
                details: 'Rows have light border in between',
              },
              {
                value: 'none',
                label: `'${ROW_SEPARATOR_TYPES.none}'`,
                details: 'Rows have no visual separation',
              },
            ]}
            value={rowsSeparator}
            onChange={value => setRowsSeparator(ROW_SEPARATOR_TYPES[value])}
          />
        </div>
        <div className={classes.displayMenuOption}>
          <div>
            <Text>Table Variant: </Text>
          </div>
          <RadioGroup
            options={[
              {
                value: 'standalone',
                label: `'${TABLE_VARIANT_TYPES.standalone}'`,
                details: 'Table is shown elevated',
              },
              {
                value: 'embedded',
                label: `'${TABLE_VARIANT_TYPES.embedded}'`,
                details: 'No elevation and no inner padding',
              },
            ]}
            value={tableVariant}
            onChange={value => setTableVariant(TABLE_VARIANT_TYPES[value])}
          />
        </div>
      </div>
    </div>
  );
};

storiesOf(`${STORY_CATEGORIES.COMPONENTS}`, module).add('Table', () => (
  <TablesRoot />
));
