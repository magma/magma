/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Link} from '../../common/Equipment';

import * as React from 'react';
import ServiceLinkDetails from './ServiceLinkDetails';
import ServiceLinksView_links from './__generated__/ServiceLinksView_links.graphql';
import {createFragmentContainer, graphql} from 'react-relay';
import {makeStyles} from '@material-ui/styles';

type Props = {
  links: ServiceLinksView_links,
  onDeleteLink: (link: Link) => void,
};

const useStyles = makeStyles(_ => ({
  root: {
    display: 'flex',
    flexDirection: 'column',
  },
}));

const ServiceLinksView = (props: Props) => {
  const {links, onDeleteLink} = props;
  const classes = useStyles();

  return (
    <div className={classes.root}>
      {links.map(link => (
        <ServiceLinkDetails
          link={link}
          onDeleteLink={() => onDeleteLink(link)}
        />
      ))}
    </div>
  );
};

export default createFragmentContainer(ServiceLinksView, {
  links: graphql`
    fragment ServiceLinksView_links on Link @relay(plural: true) {
      id
      ports {
        parentEquipment {
          id
          name
        }
        definition {
          id
          name
        }
      }
    }
  `,
});
