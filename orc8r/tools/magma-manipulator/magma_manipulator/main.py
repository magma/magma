"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import logging
import threading
import time
from queue import Queue

from kubernetes import client, config, watch

from magma_manipulator.config_parser import cfg as CONF
from magma_manipulator import gateways
from magma_manipulator import k8s_tools
from magma_manipulator import magma_api
from magma_manipulator import utils


LOG = logging.getLogger(__name__)
logging.basicConfig(
    level=logging.DEBUG,
    format='%(asctime)s %(name)-12s %(levelname)-8s %(message)s')

GWS_CFG_PULL_INTERVAL = 30
RETRY_ON_FAIL = 3

K8S_STARTED_REASON = ('Started',)
K8S_ADDED_TYPE = ('ADDED',)

events_queue = Queue()
INIT_QUEUE_TIMEOUT = 10
EVENT_MAX_TIMEOUT = 900


def watch_for_gateways(kubeconfig_path, kube_namespace, gw_names):
    config.load_kube_config(config_file=kubeconfig_path)
    v1 = client.CoreV1Api()
    w = watch.Watch()
    # infinity loop for k8s events
    for event in w.stream(v1.list_namespaced_event,
                          kube_namespace, timeout_seconds=0):
        pod_name_prefix = event['object'].involved_object.name.split('-')[0]
        if pod_name_prefix in gw_names:
            if event['type'] in K8S_ADDED_TYPE and \
               event['object'].reason in K8S_STARTED_REASON:
                LOG.info('Received event from k8s: {type} {name} {reason} '
                         '{timestamp} {msg}'.format(
                             type=event['type'],
                             name=event['object'].involved_object.name,
                             reason=event['object'].reason,
                             timestamp=event['object'].first_timestamp,
                             msg=event['object'].message))
                event = {
                    'pod_name': event['object'].involved_object.name,
                    'timeout': INIT_QUEUE_TIMEOUT,
                    'retry_on_fail': RETRY_ON_FAIL
                }
                events_queue.put(event)


def pull_gws_configs(gateways):
    while True:
        for gw_name, gw in gateways.items():
            gw_config = magma_api.get_gateway_config(
                CONF.orc8r_api_url,
                gw.network, gw.network_type,
                gw.id, CONF.magma_certs_path)

            config_path = utils.save_gateway_config(
                gw.id, CONF.gateways.configs_dir, gw_config)
            gw.config_path = config_path
            LOG.info('Pulled config for {gw_name} {gw_id}'.format(
                gw_name=gw.name, gw_id=gw.id))
        time.sleep(GWS_CFG_PULL_INTERVAL)


def start_periodic_tasks(gateways):
    LOG.info('Start watching for k8s events from gateways {gws}'.format(
        gws=(gateways.keys())))
    watch_thread = threading.Thread(
        target=watch_for_gateways,
        args=(CONF.k8s.kubeconfig_path,
              CONF.k8s.namespace,
              gateways.keys()))
    watch_thread.start()

    LOG.info('Pulling gateways config at {interval} second interval'.format(
        interval=GWS_CFG_PULL_INTERVAL))
    cfg_puller_thread = threading.Thread(
        target=pull_gws_configs,
        args=(gateways,))
    cfg_puller_thread.start()


def put_event_after_timeout(event):
    LOG.debug('Wait {sec} seconds for event {event}'.format(
        sec=event['timeout'], event=event['pod_name']))
    if event['timeout'] > EVENT_MAX_TIMEOUT:
        LOG.error('Can not handle event for pod {pod_name}. Timeout expired'
                  .format(pod_name=event['pod_name']))
        return
    t = threading.Timer(event['timeout'], lambda: events_queue.put(event))
    t.start()


def main():
    gws_manager = gateways.GatewaysManager()
    start_periodic_tasks(gws_manager.get_gateways())

    while True:
        try:
            if not events_queue.empty():
                event = events_queue.get()
                gw_pod_name = event['pod_name']
                gw = gws_manager.get_gateway(gw_pod_name)

                LOG.info('Handle event for {gw_pod_name}'.format(
                    gw_pod_name=gw_pod_name))

                if not k8s_tools.is_pod_ready(CONF.k8s.kubeconfig_path,
                                              CONF.k8s.namespace,
                                              gw_pod_name):
                    event['timeout'] *= 2
                    put_event_after_timeout(event)
                    continue

                if not utils.is_gw_reachable(gw.get_ip(gw_pod_name)):
                    event['timeout'] *= 2
                    put_event_after_timeout(event)
                    continue

                if not utils.is_cloud_init_done(
                        gw.get_ip(gw_pod_name),
                        CONF.gateways.username,
                        CONF.gateways.rsa_private_key_path):
                    event['timeout'] *= 2
                    put_event_after_timeout(event)
                    continue

                if magma_api.is_gateway_in_network(CONF.orc8r_api_url,
                                                   gw.network,
                                                   gw.id,
                                                   CONF.magma_certs_path):
                    magma_api.delete_gateway(CONF.orc8r_api_url,
                                             gw.network, gw.network_type,
                                             gw.id, CONF.magma_certs_path)

                # get gw hardware id and challenge key
                gw_uuid, gw_key = gw.get_uuid_and_key()

                magma_api.register_gateway(CONF.orc8r_api_url,
                                           gw.network, gw.network_type,
                                           gw.id, gw_uuid, gw_key,
                                           gw.name, CONF.magma_certs_path)

                magma_api.apply_gateway_config(CONF.orc8r_api_url,
                                               gw.network, gw.network_type,
                                               gw.id, gw.get_config(),
                                               CONF.magma_certs_path)
            time.sleep(1)
        except Exception as e:
            LOG.error(e)
            if event['retry_on_fail'] > 0:
                event['retry_on_fail'] -= 1
                put_event_after_timeout(event)
                LOG.warning('Try to register gateway {gw_name} {gw_id} one '
                            'more time. Remainig attempts {attempts}'.format(
                                gw_name=gw.name,
                                gw_id=gw.id,
                                attempts=event['retry_on_fail']))
