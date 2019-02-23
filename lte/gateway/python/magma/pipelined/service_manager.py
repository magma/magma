#!/usr/bin/env python3
"""
Copyright (c) 2019-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
# pylint: skip-file
# pylint does not play well with aioeventlet, as it uses asyncio.async which
# produces a parse error

import asyncio
from typing import List

import aioeventlet
from lte.protos.mconfig.mconfigs_pb2 import PipelineD
from lte.protos.meteringd_pb2_grpc import MeteringdRecordsControllerStub
from lte.protos.mobilityd_pb2_grpc import MobilityServiceStub
from lte.protos.session_manager_pb2_grpc import LocalSessionManagerStub
from magma.pipelined.app import of_rest_server
from magma.pipelined.app.access_control import AccessControlController
from magma.pipelined.app.arp import ArpController
from magma.pipelined.app.dpi import DPIController
from magma.pipelined.app.enforcement import EnforcementController
from magma.pipelined.app.enforcement_stats import EnforcementStatsController
from magma.pipelined.app.inout import EGRESS, INGRESS, InOutController
from magma.pipelined.app.meter import MeterController
from magma.pipelined.app.meter_stats import MeterStatsController
from magma.pipelined.app.subscriber import SubscriberController
from magma.pipelined.rule_mapper import RuleIDToNumMapper
from ryu.base.app_manager import AppManager

from magma.common.service import MagmaService
from magma.common.service_registry import ServiceRegistry
from magma.configuration import environment


class ServiceManager:
    """
    ServiceManager manages the service lifecycle and chaining of services for
    the Ryu apps. Ryu apps are loaded based on the services specified in the
    YAML config for static apps and mconfig for dynamic apps.
    ServiceManager also maintains a mapping between apps to the flow
    tables they use.

    Currently, its use cases include:
        - Starting all Ryu apps
        - Flow table number lookup for Ryu apps
    """

    INGRESS_TABLE_NUM = 1
    EGRESS_TABLE_NUM = 20
    START_TABLE_NUM = 2
    RYU_REST_APP_NAME = 'ryu_rest_app'

    # Mapping between services defined in mconfig and the names of the
    # corresponding Ryu apps in PipelineD with flow tables assigned.
    # Note that a service may require multiple apps.
    SERVICE_TO_APP_NAMES_DICT = {
        PipelineD.METERING: [MeterController.APP_NAME,
                             MeterStatsController.APP_NAME,
                             SubscriberController.APP_NAME, ],
        PipelineD.DPI: [DPIController.APP_NAME],
        PipelineD.ENFORCEMENT: [EnforcementController.APP_NAME,
                                EnforcementStatsController.APP_NAME, ],
    }
    # Mapping between services defined in mconfig and the module of the
    # corresponding Ryu apps in PipelineD. The module is used to for the Ryu
    # app manager to instantiate the app.
    # Note that a service may require multiple apps.
    SERVICE_TO_APP_MODULES_DICT = {
        PipelineD.METERING: [MeterController.__module__,
                             MeterStatsController.__module__,
                             SubscriberController.__module__, ],
        PipelineD.DPI: [DPIController.__module__],
        PipelineD.ENFORCEMENT: [EnforcementController.__module__,
                                EnforcementStatsController.__module__, ],
    }

    # Mapping between the app names defined in pipelined.yml and the module of
    # their corresponding Ryu apps in PipelineD.
    STATIC_APP_NAME_TO_MODULE_DICT = {
        ArpController.APP_NAME: ArpController.__module__,
        AccessControlController.APP_NAME: AccessControlController.__module__,
        RYU_REST_APP_NAME: 'ryu.app.ofctl_rest',
    }

    # Some apps do not use a table, so they need to be excluded from table
    # allocation.
    STATIC_APPS_WITH_NO_TABLE = [
        RYU_REST_APP_NAME,
    ]

    def __init__(self, magma_service: MagmaService):
        self._magma_service = magma_service
        # inout is a mandatory app and it occupies both table 1(for ingress)
        # and table 20(for egress).
        self._app_modules = [InOutController.__module__]
        self._app_table_num_dict = {
            INGRESS: self.INGRESS_TABLE_NUM,
            EGRESS: self.EGRESS_TABLE_NUM,
        }
        self._last_table_num = self.START_TABLE_NUM - 1

        static_apps = self._magma_service.config['static_apps']
        dynamic_services = self._magma_service.mconfig.services
        self._init_static_apps(static_apps)
        self._init_dynamic_apps(dynamic_services)

    def _init_static_apps(self, static_apps: List[str]):
        """
        _init_static_apps populates app modules and table number dict for
        static apps.
        """
        static_app_modules = [self.STATIC_APP_NAME_TO_MODULE_DICT[app]
                              for app in static_apps]
        self._app_modules.extend(static_app_modules)

        # Populate app to table num dict with the static apps. Filter out any
        # apps that do not need a table.
        static_app_names = [app_name for app_name in static_apps if
                            app_name not in self.STATIC_APPS_WITH_NO_TABLE]
        for app_name in static_app_names:
            self._app_table_num_dict[app_name] = self._last_table_num + 1
            self._last_table_num += 1

    def _init_dynamic_apps(self, dynamic_services: List[str]):
        """
        _init_dynamic_apps populates app modules and table number dict for
        dynamic apps. Note that each dynamic service can consist of multiple
        apps that use the same table.
        """
        dynamic_app_modules = [app for service in dynamic_services for app in
                               self.SERVICE_TO_APP_MODULES_DICT[service]]
        self._app_modules.extend(dynamic_app_modules)

        for tbl_num, service in enumerate(dynamic_services,
                                          start=self._last_table_num + 1):
            self._last_table_num = tbl_num
            for app in self.SERVICE_TO_APP_NAMES_DICT[service]:
                self._app_table_num_dict[app] = tbl_num

        if PipelineD.ENFORCEMENT in dynamic_services and \
                not self._magma_service.mconfig.relay_enabled:
            self._app_modules.remove(EnforcementStatsController.__module__)
            self._app_table_num_dict.pop(EnforcementStatsController.APP_NAME)

    def load(self):
        """
        Instantiates and schedules the Ryu app eventlets in the service
        eventloop.
        """
        manager = AppManager.get_instance()
        manager.load_apps(self._app_modules)
        contexts = manager.create_contexts()
        contexts['rule_id_mapper'] = RuleIDToNumMapper()
        contexts['app_futures'] = {}
        contexts['config'] = self._magma_service.config
        contexts['mconfig'] = self._magma_service.mconfig
        contexts['loop'] = self._magma_service.loop
        contexts['service_manager'] = self

        records_chan = ServiceRegistry.get_rpc_channel(
            'meteringd_records', ServiceRegistry.CLOUD)
        sessiond_chan = ServiceRegistry.get_rpc_channel(
            'sessiond', ServiceRegistry.LOCAL)
        mobilityd_chan = ServiceRegistry.get_rpc_channel(
            'mobilityd', ServiceRegistry.LOCAL)
        contexts['rpc_stubs'] = {
            'metering_cloud': MeteringdRecordsControllerStub(records_chan),
            'mobilityd': MobilityServiceStub(mobilityd_chan),
            'sessiond': LocalSessionManagerStub(sessiond_chan),
        }

        # Instantiate and schedule apps
        for app in manager.instantiate_apps(**contexts):
            # Wrap the eventlet in asyncio so it will stop when the loop is
            # stopped
            future = aioeventlet.wrap_greenthread(app,
                                                  self._magma_service.loop)

            # Schedule the eventlet for evaluation in service loop
            asyncio.ensure_future(future)

        # In development mode, run server so that
        if environment.is_dev_mode():
            server_thread = of_rest_server.start(manager)
            future = aioeventlet.wrap_greenthread(server_thread,
                                                  self._magma_service.loop)
            asyncio.ensure_future(future)

    def get_table_num(self, app_name: str) -> int:
        """
        Args:
            app_name: Name of the app
        Returns:
            The app's corresponding table number
        """
        return self._app_table_num_dict[app_name]

    def get_next_table_num(self, app_name: str) -> int:
        """
        Args:
            app_name: Name of the app
        Returns:
            The next table after the app's table number.
            If there are no more used tables after the current table,
            return the EGRESS table
        """
        if self._app_table_num_dict[app_name] < self._last_table_num:
            return self._app_table_num_dict[app_name] + 1
        return self.EGRESS_TABLE_NUM

    def is_app_enabled(self, app_name: str) -> bool:
        """
        Args:
             app_name: Name of the app
        Returns:
            Whether or not the app is enabled
        """
        return app_name in self._app_table_num_dict or \
            app_name == InOutController.APP_NAME
