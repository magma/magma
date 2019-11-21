/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Equipment} from './Equipment';
import type {FileAttachmentType} from './FileAttachment.js';
import type {ImageAttachmentType} from './ImageAttachment.js';
import type {LocationSiteSurveyTab_location} from './../components/location/__generated__/LocationSiteSurveyTab_location.graphql.js';
import type {LocationType} from './LocationType';
import type {Property} from './Property';
import type {TopologyNetwork} from './NetworkTopology';

// TODO: Usage of the Location type should eventually be replaced by the
// generated Relay type.
export type Location = {
  id: string,
  externalId: ?string,
  name: string,
  locationType: LocationType,
  parentLocation: ?Location,
  children: Array<Location>,
  numChildren: number,
  equipments: Array<Equipment>,
  latitude: number,
  longitude: number,
  properties: Array<Property>,
  images: Array<ImageAttachmentType>,
  files: Array<FileAttachmentType>,
  siteSurveyNeeded: boolean,
  topology: TopologyNetwork,
  locationHierarchy: Array<Location>,
  surveys: $PropertyType<LocationSiteSurveyTab_location, 'surveys'>,
};
