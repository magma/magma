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
  LocationFloorPlansTab_location$data,
  LocationFloorPlansTab_location$key,
} from './__generated__/LocationFloorPlansTab_location.graphql';

import AddFloorPlanMutation from '../../mutations/AddFloorPlanMutation';
import Button from '@fbcnms/ui/components/design-system/Button';
import Card from '@fbcnms/ui/components/design-system/Card/Card';
import CardHeader from '@fbcnms/ui/components/design-system/Card/CardHeader';
import DeleteFloorPlanMutation from '../../mutations/DeleteFloorPlanMutation';
import FileAttachment from '../FileAttachment';
import FloorPlanImage from './FloorPlanImage';
import React, {useState} from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import axios from 'axios';

import nullthrows from '@fbcnms/util/nullthrows';
import {DocumentAPIUrls} from '../../common/DocumentAPI';
import {FileUploadButton} from '../FileUpload';
import {graphql, useFragment} from 'react-relay/hooks';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';

const useStyles = makeStyles(() => ({
  table: {
    minWidth: 70,
    marginBottom: '12px',
  },
}));

type Props = {
  location: LocationFloorPlansTab_location$key,
};

const FLOOR_PLANS_KEY = 'floorPlans';

export default function LocationFloorPlansTab(props: Props) {
  const classes = useStyles();
  const [file, setFile] = useState<?File>();
  const enqueueSnackbar = useEnqueueSnackbar();

  const location: LocationFloorPlansTab_location$data = useFragment(
    graphql`
      fragment LocationFloorPlansTab_location on Location {
        id
        floorPlans {
          id
          name
          image {
            ...FileAttachment_file
          }
        }
      }
    `,
    props.location,
  );

  const uploadFloorPlan = (imgKey, referencePoint, scale) => {
    const file2 = nullthrows(file);
    const {x, y, latitude, longitude} = referencePoint;
    const {x1, y1, x2, y2, scaleInMeters} = scale;

    AddFloorPlanMutation(
      {
        input: {
          name: '', // TODO expose name field
          locationID: location.id,
          image: {
            entityType: 'LOCATION',
            entityId: '', // we are not using this field here
            imgKey,
            fileName: file2.name,
            fileSize: file2.size,
            modified: new Date(file2.lastModified).toISOString(),
            contentType: file2.type,
          },
          referenceX: x,
          referenceY: y,
          latitude: nullthrows(latitude),
          longitude: nullthrows(longitude),
          referencePoint1X: x1,
          referencePoint1Y: y1,
          referencePoint2X: nullthrows(x2),
          referencePoint2Y: nullthrows(y2),
          scaleInMeters: nullthrows(scaleInMeters),
        },
      },
      {
        onCompleted: () => {
          enqueueSnackbar('Uploaded successfully', {variant: 'success'});
        },
        onError: () => {
          enqueueSnackbar('Error uploading image', {variant: 'error'});
        },
      },
      store => {
        const newNode = store.getRootField('addFloorPlan');
        const entityProxy = store.get(location.id);
        const floorPlans = entityProxy.getLinkedRecords(FLOOR_PLANS_KEY) || [];
        entityProxy.setLinkedRecords([...floorPlans, newNode], FLOOR_PLANS_KEY);
        setFile(null);
      },
    );
  };

  if (file) {
    return <FloorPlanImage file={file} onUpload={uploadFloorPlan} />;
  }

  return (
    <Card>
      <CardHeader
        rightContent={
          <FileUploadButton
            button={<Button>Upload</Button>}
            onFileChanged={event => setFile(event.currentTarget.files[0])}
          />
        }>
        Floor Plans
      </CardHeader>
      <Table className={classes.table}>
        <TableBody>
          {location.floorPlans.filter(Boolean).map(floorPlan => (
            <FileAttachment
              key={floorPlan.id}
              file={floorPlan.image}
              onDocumentDeleted={() =>
                DeleteFloorPlanMutation(
                  {id: floorPlan.id},
                  {
                    onCompleted: () => {
                      enqueueSnackbar('Floor Plan deleted successfully', {
                        variant: 'success',
                      });
                    },
                  },
                  store => {
                    const proxy = store.get(location.id);
                    const records = proxy
                      .getLinkedRecords(FLOOR_PLANS_KEY)
                      .filter(f => f && f.id !== floorPlan.id);
                    proxy.setLinkedRecords(records, FLOOR_PLANS_KEY);
                    store.delete(floorPlan.id);
                    axios.delete(DocumentAPIUrls.delete_url(floorPlan.id));
                  },
                )
              }
            />
          ))}
        </TableBody>
      </Table>
    </Card>
  );
}
