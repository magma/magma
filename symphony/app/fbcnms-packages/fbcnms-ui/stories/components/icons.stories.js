/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import * as React from 'react';
import ActivityIcon from '../../components/design-system/Icons/Communication/ActivityIcon';
import AddIcon from '../../components/design-system/Icons/Actions/AddIcon';
import ArrowDownIcon from '../../components/design-system/Icons/Navigation/ArrowDownIcon';
import ArrowLeftIcon from '../../components/design-system/Icons/Navigation/ArrowLeftIcon';
import ArrowRightIcon from '../../components/design-system/Icons/Navigation/ArrowRightIcon';
import ArrowUpIcon from '../../components/design-system/Icons/Navigation/ArrowUpIcon';
import AssignIcon from '../../components/design-system/Icons/Actions/AssignIcon';
import AttachmentIcon from '../../components/design-system/Icons/Communication/AttachmentIcon';
import BackArrowIcon from '../../components/design-system/Icons/Navigation/BackArrowIcon';
import CalendarIcon from '../../components/design-system/Icons/Indications/CalendarIcon';
import CloseIcon from '../../components/design-system/Icons/Navigation/CloseIcon';
import CommentIcon from '../../components/design-system/Icons/Communication/CommentIcon';
import DeleteIcon from '../../components/design-system/Icons/Actions/DeleteIcon';
import DownloadAttachmentIcon from '../../components/design-system/Icons/Communication/DownloadAttachmentIcon';
import DownloadIcon from '../../components/design-system/Icons/Actions/DownloadIcon';
import DuplicateIcon from '../../components/design-system/Icons/Actions/DuplicateIcon';
import EditIcon from '../../components/design-system/Icons/Actions/EditIcon';
import EmojiIcon from '../../components/design-system/Icons/Communication/EmojiIcon';
import FiltersIcon from '../../components/design-system/Icons/Actions/FiltersIcon';
import HierarchyArrowIcon from '../../components/design-system/Icons/Indications/HierarchyArrowIcon';
import InfoSmallIcon from '../../components/design-system/Icons/Indications/InfoSmallIcon';
import InfoTinyIcon from '../../components/design-system/Icons/Indications/InfoTinyIcon';
import LinkIcon from '../../components/design-system/Icons/Actions/LinkIcon';
import ListViewIcon from '../../components/design-system/Icons/Navigation/ListViewIcon';
import MapViewIcon from '../../components/design-system/Icons/Navigation/MapViewIcon';
import MessageIcon from '../../components/design-system/Icons/Indications/MessageIcon';
import NextArrowIcon from '../../components/design-system/Icons/Navigation/NextArrowIcon';
import PlannedIcon from '../../components/design-system/Icons/Indications/PlannedIcon';
import RemoveIcon from '../../components/design-system/Icons/Actions/RemoveIcon';
import SearchIcon from '../../components/design-system/Icons/Actions/SearchIcon';
import Text from '../../components/design-system/Text';
import ThreeDotsIcon from '../../components/design-system/Icons/Actions/ThreeDotsIcon';
import UploadIcon from '../../components/design-system/Icons/Actions/UploadIcon';
import WorkOrdersIcon from '../../components/design-system/Icons/Indications/WorkOrdersIcon';
import {STORY_CATEGORIES} from '../storybookUtils';
import {makeStyles} from '@material-ui/styles';
import {storiesOf} from '@storybook/react';

const useStyles = makeStyles(_theme => ({
  root: {
    width: '100%',
  },
  iconColors: {
    marginBottom: '16px',
  },
  iconRoot: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: '16px',
  },
  iconName: {
    marginLeft: '8px',
    color: '#374050',
  },
}));

type IconProps = {
  name: string,
  icon: React.Node,
};

const Icon = ({icon, name}: IconProps) => {
  const classes = useStyles();
  return (
    <div className={classes.iconRoot}>
      {icon}
      <Text className={classes.iconName} variant="body1">
        {name}
      </Text>
    </div>
  );
};

const IconsRoot = () => {
  const classes = useStyles();

  return (
    <div className={classes.root}>
      <div className={classes.iconColors}>
        <AddIcon />
        <AddIcon color="light" />
        <AddIcon color="primary" />
        <AddIcon color="error" />
        <AddIcon color="gray" />
      </div>
      <Icon icon={<AddIcon />} name="AddIcon" />
      <Icon icon={<AssignIcon />} name="AssignIcon" />
      <Icon icon={<DeleteIcon />} name="DeleteIcon" />
      <Icon icon={<DownloadIcon />} name="DownloadIcon" />
      <Icon icon={<DuplicateIcon />} name="DuplicateIcon" />
      <Icon icon={<EditIcon />} name="EditIcon" />
      <Icon icon={<FiltersIcon />} name="FiltersIcon" />
      <Icon icon={<LinkIcon />} name="LinkIcon" />
      <Icon icon={<RemoveIcon />} name="RemoveIcon" />
      <Icon icon={<SearchIcon />} name="SearchIcon" />
      <Icon icon={<ThreeDotsIcon />} name="ThreeDotsIcon" />
      <Icon icon={<UploadIcon />} name="UploadIcon" />
      <Icon icon={<ActivityIcon />} name="ActivityIcon" />
      <Icon icon={<AttachmentIcon />} name="AttachmentIcon" />
      <Icon icon={<CalendarIcon />} name="CalendarIcon" />
      <Icon icon={<CommentIcon />} name="CommentIcon" />
      <Icon icon={<DownloadAttachmentIcon />} name="DownloadAttachmentIcon" />
      <Icon icon={<EmojiIcon />} name="EmojiIcon" />
      <Icon icon={<HierarchyArrowIcon />} name="HierarchyArrowIcon" />
      <Icon icon={<InfoSmallIcon />} name="InfoSmallIcon" />
      <Icon icon={<InfoTinyIcon />} name="InfoTinyIcon" />
      <Icon icon={<MessageIcon />} name="MessageIcon" />
      <Icon icon={<PlannedIcon />} name="PlannedIcon" />
      <Icon icon={<WorkOrdersIcon />} name="WorkOrdersIcon" />
      <Icon icon={<ArrowDownIcon />} name="ArrowDownIcon" />
      <Icon icon={<ArrowLeftIcon />} name="ArrowLeftIcon" />
      <Icon icon={<ArrowRightIcon />} name="ArrowRightIcon" />
      <Icon icon={<ArrowUpIcon />} name="ArrowUpIcon" />
      <Icon icon={<BackArrowIcon />} name="BackArrowIcon" />
      <Icon icon={<CloseIcon />} name="CloseIcon" />
      <Icon icon={<ListViewIcon />} name="ListViewIcon" />
      <Icon icon={<MapViewIcon />} name="MapViewIcon" />
      <Icon icon={<NextArrowIcon />} name="NextArrowIcon" />
    </div>
  );
};

storiesOf(`${STORY_CATEGORIES.FOUNDATION}`, module).add('Icons', () => (
  <IconsRoot />
));
