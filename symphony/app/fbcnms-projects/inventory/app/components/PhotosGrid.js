/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {PhotosGrid_images} from './__generated__/PhotosGrid_images.graphql';
import type {WithStyles} from '@material-ui/core';

import CardSection from './CardSection';
import GridList from '@material-ui/core/GridList';
import GridListTile from '@material-ui/core/GridListTile';
import ImageAttachment from './ImageAttachment';
import React from 'react';
import {createFragmentContainer, graphql} from 'react-relay';
import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  gridListRoot: {
    marginTop: '6px',
  },
  gridList: {
    transform: 'translateZ(0)',
  },
  newImage: {
    alignItems: 'center',
    border: `1px solid ${theme.palette.grey[200]}`,
    borderRadius: '2px',
    display: 'flex',
    flexDirection: 'column',
    height: '100%',
    justifyContent: 'center',
    textAlign: 'center',
  },
  buttonText: {
    marginTop: '6px',
    lineHeight: '100%',
  },
  addFabIcon: {
    fontSize: '30px',
  },
  uploadButton: {
    padding: '8px',
    borderRadius: '2px',
    '&:hover': {
      backgroundColor: theme.palette.grey[50],
    },
  },
});

type Props = {
  images: PhotosGrid_images,
  onImageDeleted: (img: PhotosGrid_images) => void,
} & WithStyles<typeof styles>;

class PhotosGrid extends React.Component<Props> {
  render() {
    const {classes, images, onImageDeleted} = this.props;
    return (
      <div>
        <CardSection title="Photos">
          <div className={classes.gridListRoot}>
            <GridList
              className={classes.gridList}
              cellHeight={150}
              cols={6.5}
              rows={2}>
              {images.map(img => (
                <GridListTile key={img.id}>
                  <ImageAttachment img={img} onImageDeleted={onImageDeleted} />
                </GridListTile>
              ))}
            </GridList>
          </div>
        </CardSection>
      </div>
    );
  }
}

export default withStyles(styles)(
  createFragmentContainer(PhotosGrid, {
    images: graphql`
      fragment PhotosGrid_images on File @relay(plural: true) {
        id
        ...ImageAttachment_img
      }
    `,
  }),
);
