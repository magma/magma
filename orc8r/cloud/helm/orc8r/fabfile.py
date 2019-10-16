"""
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
"""

import time

import os
from fabric.api import env, local
from fabric.context_managers import lcd
from fabric.operations import prompt

MAGMA_ROOT = os.path.realpath('../../..')
HELM_ROOT = os.path.join(MAGMA_ROOT, 'orc8r/cloud/helm/orc8r')
SECRETS_CHART_ROOT = os.path.join(HELM_ROOT, 'charts/secrets')


def deploy(vals_file: str,
           kubeconfig: str,
           secrets_s3_bucket: str,
           install: bool=False):
    """
    Upgrade or initiate a Helm release for Orchestrator.

    Args:
        vals_file: Full path to the vals.yml file for the Helm deployment
        kubeconfig: Full path to the kubeconfig file for the target k8s cluster
        secrets_s3_bucket: Name of the s3 bucket where secrets are stored
        install: Set to True to do a fresh helm install. False by default.
    """
    print("You're initiating a Helm release upgrade. Have you updated the "
          f'image tags in {vals_file}?')
    print('Think about it for 3 seconds...')
    time.sleep(3)
    sure = prompt('Are you ready to continue?',
                  default='no', validate='^(yes|no)$')
    if sure != 'yes':
        exit()

    _copy_secrets(secrets_s3_bucket)
    os.environ['KUBECONFIG'] = kubeconfig

    # template secrets and kubectl apply
    env.release_success = False
    try:
        with lcd(HELM_ROOT):
            if not install:
                local(f'helm upgrade orc8r . --values={vals_file}')
            else:
                local(f'helm install --name orc8r --namespace magma . '
                      f'--values={vals_file}')
            env.release_success = True
    except Exception as e:
        print(e)
    finally:
        local(f"rm -rf {os.path.join(SECRETS_CHART_ROOT, '.secrets')}")
        if env.release_success:
            text = f'Upgrade Successful!' \
                   'Use `kubectl -n magma get pods -w` ' \
                   'to monitor the health of the release.'
            print(text)
        else:
            print('Failed to upgrade release')
            exit(1)


def _copy_secrets(s3_bucket_name: str):
    new_secrets_dir = os.path.join(SECRETS_CHART_ROOT, '.secrets-deploy')
    local(f'aws s3 cp s3://{s3_bucket_name} {new_secrets_dir} --recursive')

    final_secrets_dir = os.path.join(SECRETS_CHART_ROOT, '.secrets')
    local(f'rm -rf {final_secrets_dir}')
    local(f'mv {new_secrets_dir} {final_secrets_dir}')
