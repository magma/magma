import os

from kubernetes import client, config


def set_kubeconfig_environ(kubeconfig):
    os.environ["KUBECONFIG"] = kubeconfig
    os.environ["K8S_AUTH_KUBECONFIG"] = kubeconfig


def get_all_namespaces(kubeconfig) -> list[str]:
    config.load_kube_config(kubeconfig)
    v1 = client.CoreV1Api()
    response = v1.list_namespace()
    all_namespaces = [item.metadata.name for item in response.items]
    return all_namespaces
