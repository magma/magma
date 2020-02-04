/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {EquipmentPosition} from '../../common/Equipment';
import type {EquipmentType} from '../../common/EquipmentType';

import AppContext from '@fbcnms/ui/context/AppContext';
import Button from '@fbcnms/ui/components/design-system/Button';
import Card from '@material-ui/core/Card';
import EditIcon from '@material-ui/icons/Edit';
import EquipmentBreadcrumbs from './EquipmentBreadcrumbs';
import EquipmentDetails from './EquipmentDetails';
import EquipmentDocumentsCard from './EquipmentDocumentsCard';
import EquipmentPortsTable from './EquipmentPortsTable';
import EquipmentServicesTable from './EquipmentServicesTable';
import ErrorMessage from '@fbcnms/ui/components/ErrorMessage';
import FormAction from '@fbcnms/ui/components/design-system/Form/FormAction';
import InventoryQueryRenderer from '../InventoryQueryRenderer';
import PerfectScrollbar from 'react-perfect-scrollbar';
import React, {useContext, useState} from 'react';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import {FormValidationContextProvider} from '@fbcnms/ui/components/design-system/Form/FormValidationContext';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {graphql} from 'react-relay';
import {makeStyles} from '@material-ui/styles';

import 'react-perfect-scrollbar/dist/css/styles.css';

graphql`
  fragment EquipmentPropertiesCard_baseEquipmentProps on Equipment
    @relay(mask: false) {
    id
    name
    futureState
    parentLocation {
      id
      name
    }
  }
`;

graphql`
  fragment EquipmentPropertiesCard_position on EquipmentPosition
    @relay(mask: false) {
    id
    definition {
      id
      name
      index
      visibleLabel
    }
    attachedEquipment {
      id
      name
      futureState
      workOrder {
        id
        status
      }
    }
  }
`;

const equipmentsPropertiesCardQuery = graphql`
  query EquipmentPropertiesCardQuery($equipmentId: ID!) {
    equipment: node(id: $equipmentId) {
      ... on Equipment {
        id
        name
        ...EquipmentPortsTable_equipment
        equipmentType {
          id
          name
          propertyTypes {
            ...PropertyTypeFormField_propertyType
            ...DynamicPropertiesGrid_propertyTypes
          }
          positionDefinitions {
            id
            ...PositionDefinitionsTable_positionDefinitions
          }
          portDefinitions {
            id
          }
        }
        ...EquipmentBreadcrumbs_equipment
        parentLocation {
          id
          name
        }
        parentPosition {
          parentEquipment {
            parentLocation {
              id
            }
          }
        }
        ...EquipmentPositionsGrid_equipment
        positions {
          parentEquipment {
            id
          }
        }
        properties {
          ...PropertyFormField_property
          ...DynamicPropertiesGrid_properties
        }
        services {
          id
          name
          externalId
          customer {
            name
          }
          serviceType {
            id
            name
          }
        }
        ...EquipmentDocumentsCard_equipment
      }
    }
  }
`;

type Props = {
  equipmentId: string,
  workOrderId: ?string,
  onAttachingEquipmentToPosition: (
    equipmentType: EquipmentType,
    position: EquipmentPosition,
  ) => void,
  onEquipmentClicked: (equipmentId: string) => void,
  onParentLocationClicked: (locationId: string) => void,
  onEdit: () => void,
  onWorkOrderSelected: (workOrderId: string) => void,
};

const useStyles = makeStyles(theme => ({
  cardRoot: {
    height: '100%',
  },
  root: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'flex-start',
    justifyContent: 'space-around',
    height: '100%',
  },
  tabs: {
    backgroundColor: 'white',
  },
  tabContainer: {
    width: 'auto',
  },
  cardContent: {
    paddingLeft: '0px',
    paddingRight: '0px',
    paddingTop: '0px',
    overflow: 'hidden',
    height: '100%',
    width: '100%',
    backgroundColor: theme.palette.common.white,
  },
  equipmentDetails: {
    marginTop: '20px',
    paddingLeft: theme.spacing(3),
    paddingRight: theme.spacing(3),
  },
  cardHeader: {
    padding: '0px',
    width: '100%',
    display: 'block',
  },
  equipmentBreadcrumbs: {
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'space-between',
    padding: '20px 20px 24px 20px',
  },
  iconButton: {
    padding: '0px',
    marginLeft: theme.spacing(),
  },
  tabsContainer: {
    marginBottom: '16px',
  },
  documentsTable: {
    padding: '24px',
  },
}));

const EquipmentPropertiesCard = (props: Props) => {
  const {
    equipmentId,
    workOrderId,
    onEquipmentClicked,
    onParentLocationClicked,
  } = props;
  const classes = useStyles();
  const [selectedTab, setSelectedTab] = useState('details');
  const servicesEnabled = useContext(AppContext).isFeatureEnabled('services');
  return (
    <>
      <div className={classes.cardRoot}>
        <InventoryQueryRenderer
          query={equipmentsPropertiesCardQuery}
          variables={{
            equipmentId,
          }}
          render={innerProps => {
            const {equipment} = innerProps;
            if (!equipment) {
              return (
                <Card className={classes.cardRoot}>
                  <ErrorMessage message="It appears this equipment does not exist" />
                </Card>
              );
            }

            return (
              <div className={classes.root}>
                <FormValidationContextProvider>
                  <div className={classes.cardHeader}>
                    <div>
                      <div className={classes.equipmentBreadcrumbs}>
                        <EquipmentBreadcrumbs
                          onEquipmentClicked={onEquipmentClicked}
                          onParentLocationClicked={onParentLocationClicked}
                          equipment={equipment}
                        />
                        <FormAction>
                          <Button
                            variant="text"
                            skin="primary"
                            onClick={props.onEdit}>
                            <EditIcon />
                          </Button>
                        </FormAction>
                      </div>
                      {equipment.equipmentType.portDefinitions.length > 0 ||
                      equipment.positions.length > 0 ||
                      equipment.services.length > 0 ? (
                        <div className={classes.tabsContainer}>
                          <Tabs
                            className={classes.tabs}
                            value={selectedTab}
                            onChange={(_e, selectedTab) => {
                              ServerLogger.info(
                                LogEvents.EQUIPMENT_CARD_TAB_CLICKED,
                                {tab: selectedTab},
                              );
                              setSelectedTab(selectedTab);
                            }}
                            indicatorColor="primary"
                            textColor="primary">
                            <Tab
                              classes={{root: classes.tabContainer}}
                              label="Details"
                              value="details"
                            />
                            <Tab
                              classes={{root: classes.tabContainer}}
                              label="Documents"
                              value="documents"
                            />
                            <Tab
                              classes={{root: classes.tabContainer}}
                              label="Ports"
                              value="ports"
                            />
                            {servicesEnabled && (
                              <Tab
                                classes={{root: classes.tabContainer}}
                                label="Services"
                                value="services"
                              />
                            )}
                          </Tabs>
                        </div>
                      ) : null}
                    </div>
                  </div>
                  <div className={classes.cardContent}>
                    <PerfectScrollbar>
                      {selectedTab === 'details' ? (
                        <div className={classes.equipmentDetails}>
                          <EquipmentDetails
                            equipment={equipment}
                            workOrderId={workOrderId}
                            onAttachingEquipmentToPosition={
                              props.onAttachingEquipmentToPosition
                            }
                            onEquipmentClicked={props.onEquipmentClicked}
                            onWorkOrderSelected={props.onWorkOrderSelected}
                          />
                        </div>
                      ) : null}
                      {selectedTab === 'ports' ? (
                        <EquipmentPortsTable
                          equipment={equipment}
                          workOrderId={workOrderId}
                          onPortEquipmentClicked={props.onEquipmentClicked}
                          onParentLocationClicked={
                            props.onParentLocationClicked
                          }
                          onWorkOrderSelected={props.onWorkOrderSelected}
                        />
                      ) : null}
                      {selectedTab === 'services' ? (
                        <EquipmentServicesTable services={equipment.services} />
                      ) : null}
                      {selectedTab === 'documents' ? (
                        <EquipmentDocumentsCard
                          className={classes.documentsTable}
                          equipment={equipment}
                        />
                      ) : null}
                    </PerfectScrollbar>
                  </div>
                </FormValidationContextProvider>
              </div>
            );
          }}
        />
      </div>
    </>
  );
};

export default EquipmentPropertiesCard;
