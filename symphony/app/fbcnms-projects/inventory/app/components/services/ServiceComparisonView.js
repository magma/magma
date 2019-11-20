/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {FilterConfig} from '../comparison_view/ComparisonViewTypes';
import type {PropertyType} from '../../common/PropertyType';
import type {ServiceComparisonViewQueryRendererPropertiesQueryResponse} from './__generated__/ServiceComparisonViewQueryRendererPropertiesQuery.graphql.js';

import AddServiceCard from './AddServiceCard';
import AddServiceDialog from './AddServiceDialog';
import Button from '@fbcnms/ui/components/design-system/Button';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import CardFooter from '@fbcnms/ui/components/CardFooter';
import InventoryErrorBoundary from '../../common/InventoryErrorBoundary';
import PowerSearchBar from '../power_search/PowerSearchBar';
import React, {useCallback, useMemo, useState} from 'react';
import RelayEnvironment from '../../common/RelayEnvironment';
import ServiceCard from './ServiceCard';
import ServiceComparisonViewQueryRenderer from './ServiceComparisonViewQueryRenderer';
import symphony from '@fbcnms/ui/theme/symphony';
import useLocationTypes from '../comparison_view/hooks/locationTypesHook';
import useRouter from '@fbcnms/ui/hooks/useRouter';
import {SERVICE_PROPERTY_FILTER_NAME} from './PowerSearchServicePropertyFilter';
import {
  ServiceSearchConfig,
  buildServicePropertyFilterConfigs,
} from './ServiceSearchConfig';
import {extractEntityIdFromUrl} from '../../common/RouterUtils';
import {getInitialFilterValue} from '../comparison_view/FilterUtils';
import {graphql} from 'relay-runtime';
import {groupBy} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {useGraphQL} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(_ => ({
  cardRoot: {
    height: '100%',
    display: 'flex',
    flexDirection: 'column',
    paddingLeft: '0px',
    paddingRight: '0px',
  },
  cardContent: {
    paddingLeft: '0px',
    paddingRight: '0px',
    paddingTop: '0px',
    flexGrow: 1,
    width: '100%',
  },
  root: {
    display: 'flex',
    flexDirection: 'column',
    backgroundColor: symphony.palette.white,
    height: '100%',
  },
  searchResults: {
    flexGrow: 1,
  },
  bar: {
    display: 'flex',
    flexDirection: 'row',
    boxShadow: '0px 2px 2px 0px rgba(0, 0, 0, 0.1)',
  },
  searchBar: {
    flexGrow: 1,
  },
}));

const servicePropertiesQuery = graphql`
  query ServiceComparisonViewQueryRendererPropertiesQuery {
    possibleProperties(entityType: SERVICE) {
      name
      type
      stringValue
    }
  }
`;

function getPossibleProperties(
  data: ?ServiceComparisonViewQueryRendererPropertiesQueryResponse,
): Array<PropertyType> {
  if (data == null || data.possibleProperties == null) {
    return [];
  }
  const propertiesGroup: {[string]: Array<PropertyType>} = groupBy(
    data.possibleProperties
      .filter(prop => prop.type !== 'gps_location' && prop.type !== 'range')
      .map((prop, index) => ({
        id: prop.name + prop.type,
        type: prop.type,
        name: prop.name,
        index: index,
        stringValue: prop.stringValue,
      })),
    prop => prop.name + prop.type,
  );
  const supportedProperties: Array<PropertyType> = [];
  for (const k in propertiesGroup) {
    supportedProperties.push(propertiesGroup[k][0]);
  }
  return supportedProperties;
}

const ServiceComparisonView = () => {
  const {match, history, location} = useRouter();
  const [dialogKey, setDialogKey] = useState(1);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [serviceKey, setServiceKey] = useState(1);
  const [filters, setFilters] = useState([]);
  const classes = useStyles();

  const selectedServiceTypeId = useMemo(
    () => extractEntityIdFromUrl('serviceType', location.search),
    [location.search],
  );

  const selectedServiceCardId = useMemo(
    () => extractEntityIdFromUrl('service', location.search),
    [location],
  );

  const serviceDataResponse = useGraphQL(
    RelayEnvironment,
    servicePropertiesQuery,
    {},
  );
  const possibleProperties = getPossibleProperties(
    serviceDataResponse.response,
  );
  const servicePropertiesFilterConfigs = buildServicePropertyFilterConfigs(
    possibleProperties,
  );

  const locationTypesFilterConfigs = useLocationTypes();

  const filterConfigs = ServiceSearchConfig.map(ent => ent.filters)
    .reduce((allFilters, currentFilter) => allFilters.concat(currentFilter), [])
    .concat(servicePropertiesFilterConfigs)
    .concat(locationTypesFilterConfigs);

  const navigateToAddService = (selectedServiceTypeId: ?string) => {
    history.push(
      match.url +
        (selectedServiceTypeId ? `?serviceType=${selectedServiceTypeId}` : ''),
    );
  };

  const navigateToService = (selectedServiceId: ?string) => {
    history.push(
      match.url + (selectedServiceId ? `?service=${selectedServiceId}` : ''),
    );
  };

  const showDialog = useCallback(() => {
    setDialogOpen(true);
    setDialogKey(dialogKey + 1);
    setServiceKey(serviceKey + 1);
  }, [setDialogOpen, dialogKey, setDialogKey, serviceKey, setServiceKey]);

  const hideDialog = useCallback(() => setDialogOpen(false), [setDialogOpen]);

  if (selectedServiceTypeId != null) {
    return (
      <InventoryErrorBoundary>
        <AddServiceCard serviceTypeId={selectedServiceTypeId} />
      </InventoryErrorBoundary>
    );
  }

  if (selectedServiceCardId != null) {
    return (
      <InventoryErrorBoundary>
        <ServiceCard serviceId={selectedServiceCardId} />
      </InventoryErrorBoundary>
    );
  }

  return (
    <InventoryErrorBoundary>
      <Card className={classes.cardRoot}>
        <CardContent className={classes.cardContent}>
          <div className={classes.root}>
            <div className={classes.bar}>
              <div className={classes.searchBar}>
                <PowerSearchBar
                  placeholder="Filter services"
                  filterConfigs={filterConfigs}
                  searchConfig={ServiceSearchConfig}
                  getSelectedFilter={(filterConfig: FilterConfig) =>
                    getInitialFilterValue(
                      filterConfig.key,
                      filterConfig.name,
                      filterConfig.defaultOperator,
                      filterConfig.name === SERVICE_PROPERTY_FILTER_NAME
                        ? possibleProperties.find(
                            propDef => propDef.name === filterConfig.label,
                          )
                        : null,
                    )
                  }
                  onFiltersChanged={filters => setFilters(filters)}
                />
              </div>
            </div>
            <div className={classes.searchResults}>
              <ServiceComparisonViewQueryRenderer
                limit={50}
                filters={filters}
                onServiceSelected={selectedServiceCardId =>
                  navigateToService(selectedServiceCardId)
                }
                serviceKey={serviceKey}
              />
            </div>
          </div>
        </CardContent>
        <CardFooter alignItems="left">
          <Button onClick={showDialog}>Add Service</Button>
        </CardFooter>
        <AddServiceDialog
          key={`new_service_${dialogKey}`}
          open={dialogOpen}
          onClose={hideDialog}
          onServiceTypeSelected={typeId => {
            navigateToAddService(typeId);
            setDialogOpen(false);
          }}
        />
      </Card>
    </InventoryErrorBoundary>
  );
};

export default ServiceComparisonView;
