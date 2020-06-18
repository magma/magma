/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {Service} from '../../common/Service';

import Breadcrumbs from '@fbcnms/ui/components/Breadcrumbs';
import Button from '@fbcnms/ui/components/design-system/Button';
import FormAction from '@fbcnms/ui/components/design-system/Form/FormAction';
import React from 'react';
import ServiceDeleteButton from './ServiceDeleteButton';
import symphony from '@fbcnms/ui/theme/symphony';
import {makeStyles} from '@material-ui/styles';
import {useFormContext} from '../../common/FormContext';

const useStyles = makeStyles(_ => ({
  nameHeader: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'start',
    marginBottom: '24px',
  },
  breadcrumbs: {
    flexGrow: 1,
  },
  deleteButton: {
    cursor: 'pointer',
    color: symphony.palette.D400,
    width: '32px',
    height: '32px',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    marginRight: '12px',
  },
}));

type Props = {
  service: Service,
  onBackClicked: () => void,
  onServiceRemoved: () => void,
};

const ServiceHeader = (props: Props) => {
  const classes = useStyles();
  const {service, onBackClicked, onServiceRemoved} = props;
  const form = useFormContext();
  return (
    <div className={classes.nameHeader}>
      <div className={classes.breadcrumbs}>
        <Breadcrumbs
          breadcrumbs={[
            {
              id: 'services',
              name: 'Services',
              onClick: onBackClicked,
            },
            {
              id: service.id,
              name: service.name,
            },
          ]}
          size="large"
        />
      </div>
      <ServiceDeleteButton
        className={classes.deleteButton}
        service={service}
        onServiceRemoved={onServiceRemoved}
      />
      <FormAction ignorePermissions={true}>
        <Button onClick={onBackClicked} disabled={form.alerts.error.detected}>
          Done
        </Button>
      </FormAction>
    </div>
  );
};

export default ServiceHeader;
