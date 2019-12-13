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
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FileAttachment from '../FileAttachment';
import FormGroup from '@material-ui/core/FormGroup';
import React, {useRef, useState} from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TextField from '@material-ui/core/TextField';

import nullthrows from '@fbcnms/util/nullthrows';
import useSnackbar from '@fbcnms/ui/hooks/useSnackbar';
import {FileUploadButton, uploadFile} from '../FileUpload';
import {graphql, useFragment} from 'react-relay/hooks';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles({
  img: {
    maxWidth: '500px',
    maxHeight: '500px',
  },
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
  table: {
    minWidth: 70,
    marginBottom: '12px',
  },
});

type ReferencePoint = {
  x: number,
  y: number,
  latitude?: number,
  longitude?: number,
};

type Scale = {
  x1: number,
  y1: number,
  x2?: number,
  y2?: number,
  scaleInMeters?: number,
};

type Props = {
  location: LocationFloorPlansTab_location$key,
};

export default function LocationFloorPlansTab(props: Props) {
  const imageRef = useRef();
  const classes = useStyles();
  const [referencePointDialogShown, setReferencePointDialogShown] = useState(
    false,
  );
  const [referencePoint, setReferencePoint] = useState<?ReferencePoint>(null);
  const [scaleDialogShown, setScaleDialogShown] = useState(false);
  const [scale, setScale] = useState<?Scale>(null);
  const [message, setMessage] = useState('');
  const [file, setFile] = useState<?File>();
  useSnackbar(message, {variant: 'info'}, message != '', true);

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

  const uploadFloorPlan = (imgKey, scaleInMeters) => {
    const file2 = nullthrows(file);
    const {x, y, latitude, longitude} = nullthrows(referencePoint);
    const {x1, y1, x2, y2} = nullthrows(scale);

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
          scaleInMeters: scaleInMeters,
        },
      },
      {
        onCompleted: () => setMessage('Uploaded successfully'),
        onError: () => setMessage('Error uploading image'),
      },
      store => {
        const newNode = store.getRootField('addFloorPlan');
        const entityProxy = store.get(location.id);
        const floorPlans = entityProxy.getLinkedRecords('floorPlans') || [];
        entityProxy.setLinkedRecords([...floorPlans, newNode], 'floorPlans');
        setFile(null);
      },
    );
  };

  const onFileChanged = event => {
    const reader = new FileReader();
    reader.onload = () => {
      if (typeof reader.result == 'string') {
        nullthrows(imageRef.current).src = reader.result;
      }
    };
    reader.readAsDataURL(event.currentTarget.files[0]);
    setFile(event.currentTarget.files[0]);
    setMessage('Click a point on the image to provide a lat/lon reference');
  };

  return (
    <>
      {referencePointDialogShown && (
        <ReferencePointDialog
          onSave={(latitude, longitude) => {
            setMessage(
              'Please select two points on the image to specify the scale',
            );
            setReferencePointDialogShown(false);
            setReferencePoint({
              ...nullthrows(referencePoint),
              latitude,
              longitude,
            });
          }}
          onClose={() => {
            setReferencePointDialogShown(false);
            setReferencePoint(null);
          }}
        />
      )}
      {scaleDialogShown && (
        <ScaleDialog
          onSave={scaleInMeters => {
            setMessage('Uploading...');
            setScaleDialogShown(false);
            setScale({...nullthrows(scale), scaleInMeters});
            uploadFile(nullthrows(file), (_, imgKey) =>
              uploadFloorPlan(imgKey, scaleInMeters),
            );
          }}
          onClose={() => setScaleDialogShown(false)}
        />
      )}
      {file ? (
        <img
          ref={imageRef}
          className={classes.img}
          onClick={e => {
            const box = e.target.getBoundingClientRect();
            const x = e.pageX - box.x;
            const y = e.pageY - box.y;
            if (!referencePoint) {
              setReferencePointDialogShown(true);
              setReferencePoint({x, y});
            } else {
              if (scale && scale.x2 === undefined) {
                setScale({...scale, x2: x, y2: y});
                setScaleDialogShown(true);
              } else {
                setScale({x1: x, y1: y});
              }
            }
          }}
        />
      ) : (
        <Card>
          <CardHeader
            rightContent={
              <FileUploadButton
                button={<Button>Upload</Button>}
                onFileChanged={onFileChanged}
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
                  onDocumentDeleted={() => null}
                />
              ))}
            </TableBody>
          </Table>
        </Card>
      )}
    </>
  );
}

const ReferencePointDialog = (props: {
  onSave: (number, number) => void,
  onClose: () => void,
}) => {
  const [lat, setLat] = useState('');
  const [lon, setLon] = useState('');
  const classes = useStyles();

  return (
    <Dialog maxWidth="sm" open={true} onClose={props.onClose}>
      <DialogTitle>Provide Latitude/Longitude</DialogTitle>
      <DialogContent>
        <FormGroup row>
          <TextField
            required
            className={classes.input}
            label="Latitude"
            margin="normal"
            value={lat}
            onChange={event => setLat(event.target.value)}
          />
          <TextField
            required
            className={classes.input}
            label="Longitude"
            margin="normal"
            value={lon}
            onChange={event => setLon(event.target.value)}
          />
        </FormGroup>
      </DialogContent>
      <DialogActions>
        <Button onClick={() => props.onSave(parseFloat(lat), parseFloat(lon))}>
          Save
        </Button>
      </DialogActions>
    </Dialog>
  );
};

const ScaleDialog = (props: {onSave: number => void, onClose: () => void}) => {
  const [scale, setScale] = useState('');
  const classes = useStyles();

  return (
    <Dialog maxWidth="sm" open={true} onClose={props.onClose}>
      <DialogTitle>Provide Scale</DialogTitle>
      <DialogContent>
        <FormGroup row>
          <TextField
            required
            className={classes.input}
            label="Scale (in meters)"
            margin="normal"
            value={scale}
            onChange={event => setScale(event.target.value)}
          />
        </FormGroup>
      </DialogContent>
      <DialogActions>
        <Button onClick={() => props.onSave(parseFloat(scale))}>Save</Button>
      </DialogActions>
    </Dialog>
  );
};
