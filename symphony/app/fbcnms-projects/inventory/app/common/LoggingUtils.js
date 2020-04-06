/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {User} from '@fbcnms/ui/context/AppContext';

import {ServerLog} from '@fbcnms/ui/utils/Logging';

export const LogEvents = {
  CLIENT_FATAL_ERROR: 'client_fatal_error',

  ADD_EQUIPMENT_TYPE_BUTTON_CLICKED: 'add_equipment_type_button_clicked',
  SAVE_EQUIPMENT_TYPE_BUTTON_CLICKED: 'save_equipment_type_button_clicked',
  ADD_LOCATION_TYPE_BUTTON_CLICKED: 'add_location_type_button_clicked',
  SAVE_LOCATION_TYPE_BUTTON_CLICKED: 'save_location_type_button_clicked',
  ADD_SERVICE_TYPE_BUTTON_CLICKED: 'add_service_type_button_clicked',
  SAVE_SERVICE_TYPE_BUTTON_CLICKED: 'save_service_type_button_clicked',
  ADD_EQUIPMENT_PORT_TYPE_BUTTON_CLICKED: 'add_port_type_button_clicked',
  SAVE_EQUIPMENT_PORT_TYPE_BUTTON_CLICKED: 'save_port_type_button_clicked',
  ADD_LOCATION_BUTTON_CLICKED: 'add_location_button_clicked',
  ADD_EQUIPMENT_BUTTON_CLICKED: 'add_equipment_button_clicked',
  NAVIGATE_TO_LOCATION: 'navigate_to_location',
  NAVIGATE_TO_EQUIPMENT: 'navigate_to_equipment',
  SAVE_EQUIPMENT_BUTTON_CLICKED: 'save_equipment_button_clicked',
  EDIT_EQUIPMENT_BUTTON_CLICKED: 'edit_equipment_button_clicked',
  EDIT_EQUIPMENT_PORT_BUTTON_CLICKED: 'edit_equipment_port_button_clicked',
  SAVE_EQUIPMENT_PORT_BUTTON_CLICKED: 'save_equipment_port_button_clicked',
  LOCATION_CARD_CANCEL_BUTTON_CLICKED: 'location_card_cancel_button_clicked',
  EDIT_LOCATION_BUTTON_CLICKED: 'edit_location_button_clicked',
  SAVE_LOCATION_BUTTON_CLICKED: 'save_location_button_clicked',
  DELETE_LOCATION_BUTTON_CLICKED: 'delete_location_button_clicked',
  LOCATION_CARD_TAB_CLICKED: 'location_card_tab_clicked',
  EQUIPMENT_CARD_TAB_CLICKED: 'equipment_card_tab_clicked',
  CONFIGURE_NAV_CLICKED: 'configure_nav_clicked',
  SERVICES_NAV_CLICKED: 'services_nav_clicked',
  INVENTORY_NAV_CLICKED: 'inventory_nav_clicked',
  MAP_NAV_CLICKED: 'map_nav_clicked',
  SEARCH_NAV_CLICKED: 'search_nav_clicked',
  WORK_ORDERS_NAV_CLICKED: 'work_orders_nav_clicked',
  CONFIGURE_TAB_NAVIGATION_CLICKED: 'configure_tab_navigation_clicked',
  EQUIPMENT_CARD_LOCATION_BREADCRUMB_CLICKED:
    'equipment_card_location_breadcrumb_clicked',
  EQUIPMENT_CARD_EQUIPMENT_BREADCRUMB_CLICKED:
    'equipment_card_equipment_breadcrumb_clicked',
  LOCATION_CARD_BREADCRUMB_CLICKED: 'location_card_breadcrumb_clicked',
  DELETE_EQUIPMENT_CLICKED: 'delete_equipment_clicked',
  EQUIPMENT_COMPARISON_VIEW_EQUIPMENT_CLICKED:
    'equipment_comparison_view_equipment_clicked',
  EQUIPMENT_COMPARISON_VIEW_FILTER_REMOVED:
    'equipment_comparison_view_filter_removed',
  EQUIPMENT_COMPARISON_VIEW_FILTER_SET: 'equipment_comparison_view_filter_set',
  LINK_COMPARISON_VIEW_FILTER_REMOVED: 'link_comparison_view_filter_removed',
  LINK_COMPARISON_VIEW_FILTER_SET: 'link_comparison_view_filter_set',
  PORT_COMPARISON_VIEW_FILTER_REMOVED: 'port_comparison_view_filter_removed',
  PORT_COMPARISON_VIEW_FILTER_SET: 'port_comparison_view_filter_set',
  LOCATION_COMPARISON_VIEW_FILTER_REMOVED:
    'location_comparison_view_filter_removed',
  LOCATION_COMPARISON_VIEW_FILTER_SET: 'location_comparison_view_filter_set',
  COMPARISON_VIEW_SUBJECT_CHANGED: 'comparison_view_subject_changed',
  COMPARISON_VIEW_FILTERS_CHANGED: 'comparison_view_filters_changed',
  LOCATIONS_MAP_POPUP_OPENED: 'locations_map_popup_opened',
  PROJECTS_MAP_POPUP_OPENED: 'projects_map_popup_opened',
  ATTACH_EQUIPMENT_TO_POSITION_CLICKED: 'attach_equipment_to_position_clicked',
  ADD_LINK_CLICKED: 'add_link_clicked',
  EDIT_LINK_CLICKED: 'add_link_clicked',
  SAVE_LINK_BUTTON_CLICKED: 'save_link_clicked',
  CONNECT_PORTS_CLICKED: 'connect_ports_clicked',
  DISCONNECT_PORTS_CLICKED: 'disconnect_ports_clicked',
  LOCATION_CARD_ADD_HYPERLINK_CLICKED: 'location_card_add_hyperlink_clicked',
  LOCATION_CARD_UPLOAD_FILE_CLICKED: 'location_card_upload_file_clicked',
  LOCATION_TYPE_REORDERED: 'location_type_reordered',
  DOCUMENTATION_LINK_CLICKED: 'documentation_link_clicked',
  DOCUMENTATION_LINK_CLICKED_FROM_EXPORT_DIALOG:
    'documentation_link_clicked_from_export_dialog',
  SAVED_SEARCH_LOADED: 'saved_search_loaded',
  SAVED_SEARCH_CREATED: 'saved_search_created',
  SAVED_SEARCH_DELETED: 'saved_search_deleted',
  SAVED_SEARCH_EDITED: 'saved_search_edited',
  //Work Orders Logs:
  ADD_WORK_ORDER_TYPE_BUTTON_CLICKED: 'add_work_order_template_button_clicked',
  DELETE_WORK_ORDER_BUTTON_CLICKED: 'delete_work_order_button_clicked',
  EXECUTE_WORK_ORDER_BUTTON_CLICKED: 'execute_work_order_button_clicked',
  SAVE_WORK_ORDER_BUTTON_CLICKED: 'save_work_order_button_clicked',
  SAVE_WORK_ORDER_TYPE_BUTTON_CLICKED:
    'save_work_order_template_button_clicked',
  WORK_ORDER_DETAILS_NAV_CLICKED: 'work_order_details_nav_clicked',
  WORK_ORDERS_CONFIGURE_NAV_CLICKED: 'work_orders_configure_nav_clicked',
  WORK_ORDERS_CONFIGURE_TAB_NAVIGATION_CLICKED:
    'work_orders_configure_tab_navigation_clicked',
  WORK_ORDERS_SEARCH_NAV_CLICKED: 'work_orders_search_nav_clicked',

  //Projects logs:
  ADD_PROJECT_BUTTON_CLICKED: 'add_project_button_clicked',
  ADD_PROJECT_TEMPLATE_BUTTON_CLICKED: 'add_project_template_button_clicked',
  DELETE_PROJECT_BUTTON_CLICKED: 'delete_project_button_clicked',
  DELETE_PROJECT_TYPE_BUTTON_CLICKED: 'delete_project_type_button_clicked',
  EDIT_PROJECT_TEMPLATE_BUTTON_CLICKED: 'edit_project_template_button_clicked',
  PROJECTS_SEARCH_NAV_CLICKED: 'projects_search_nav_clicked',
  SAVE_PROJECT_BUTTON_CLICKED: 'save_project_button_clicked',
  SAVE_PROJECT_TEMPLATE_BUTTON_CLICKED: 'save_project_type_button_clicked',

  //Service logs:
  SERVICES_SEARCH_NAV_CLICKED: 'services_search_nav_clicked',
  SAVE_SERVICE_BUTTON_CLICKED: 'save_service_button_clicked',
  DELETE_SERVICE_BUTTON_CLICKED: 'delete_service_button_clicked',
  VIEW_EQUIPMENT_SERVICE_BUTTON_CLICKED:
    'view_equipment_service_button_clicked',
  ADD_EQUIPMENT_LINK_BUTTON_CLICKED: 'add_equipment_link_button_clicked',
  DELETE_SERVICE_LINK_BUTTON_CLICKED: 'delete_service_link_button_clicked',
  ADD_CONSUMER_ENDPOINT_BUTTON_CLICKED: 'add_consumer_endpoint_button_clicked',
  ADD_PROVIDER_ENDPOINT_BUTTON_CLICKED: 'add_provider_endpoint_button_clicked',
  DELETE_SERVICE_ENDPOINT_BUTTON_CLICKED:
    'delete_service_endpoint_button_clicked',
};

export const ServerLogger = ServerLog('inventory');

export const setLoggerUser = (user: User) => {
  ServerLogger.addPayloadBuilder(payload => {
    const {timestamp: _, ...rest} = payload;
    return {
      user,
      data: rest,
    };
  });
};
