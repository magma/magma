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
  EquipmentDocumentsCard_equipment$data,
  EquipmentDocumentsCard_equipment$key,
} from './__generated__/EquipmentDocumentsCard_equipment.graphql';

import Card from '@fbcnms/ui/components/design-system/Card/Card';
import CardHeader from '@fbcnms/ui/components/design-system/Card/CardHeader';
import DocumentsAddButton from '../DocumentsAddButton';
import EntityDocumentsTable from '../EntityDocumentsTable';
import React, {useMemo} from 'react';
import classNames from 'classnames';
import {graphql, useFragment} from 'react-relay/hooks';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(_theme => ({
  cardHasNoContent: {
    marginBottom: '0px',
  },
}));

type Props = {
  className?: string,
  equipment: EquipmentDocumentsCard_equipment$key,
};

const EquipmentDocumentsCard = (props: Props) => {
  const {className, equipment} = props;
  const classes = useStyles();

  const data: EquipmentDocumentsCard_equipment$data = useFragment(
    graphql`
      fragment EquipmentDocumentsCard_equipment on Equipment {
        id
        images {
          ...EntityDocumentsTable_files
        }
        files {
          ...EntityDocumentsTable_files
        }
      }
    `,
    equipment,
  );

  const documents = useMemo(
    () => [...data.files.filter(Boolean), ...data.images.filter(Boolean)],
    [data],
  );

  return (
    <Card className={className}>
      <CardHeader
        className={classNames({
          [classes.cardHasNoContent]: documents.length === 0,
        })}
        rightContent={
          <DocumentsAddButton entityType="EQUIPMENT" entityId={data.id} />
        }>
        Documents
      </CardHeader>
      <EntityDocumentsTable
        entityType="EQUIPMENT"
        entityId={data.id}
        files={documents}
      />
    </Card>
  );
};

export default EquipmentDocumentsCard;
