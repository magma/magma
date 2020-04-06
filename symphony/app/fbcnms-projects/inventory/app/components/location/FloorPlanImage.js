/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import Button from '@fbcnms/ui/components/design-system/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormGroup from '@material-ui/core/FormGroup';
import React, {useEffect, useRef, useState} from 'react';
import TextField from '@material-ui/core/TextField';

import nullthrows from '@fbcnms/util/nullthrows';
import useSnackbar from '@fbcnms/ui/hooks/useSnackbar';
import {makeStyles} from '@material-ui/styles';
import {uploadFile} from '../FileUpload/FileUploadUtils';

const useStyles = makeStyles(() => ({
  img: {
    maxWidth: '500px',
    maxHeight: '500px',
  },
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
}));

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
  file: File,
  onUpload: (key: string, ReferencePoint, Scale) => void,
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
  useSnackbar(message, {variant: 'info'}, message != '', true);

  useEffect(() => {
    const reader = new FileReader();
    reader.onload = () => {
      if (typeof reader.result == 'string') {
        nullthrows(imageRef.current).src = reader.result;
      }
    };
    reader.readAsDataURL(props.file);
    setMessage('Click a point on the image to provide a lat/lon reference');
  }, [props.file]);

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
            const newScale = {...nullthrows(scale), scaleInMeters};
            setScale(newScale);
            uploadFile(props.file.name, props.file, (_, imgKey) =>
              props.onUpload(imgKey, nullthrows(referencePoint), newScale),
            );
          }}
          onClose={() => setScaleDialogShown(false)}
        />
      )}
      <img
        ref={imageRef}
        className={classes.img}
        onClick={e => {
          const box = e.target.getBoundingClientRect();
          const x = Math.round(e.pageX - box.x);
          const y = Math.round(e.pageY - box.y);
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
