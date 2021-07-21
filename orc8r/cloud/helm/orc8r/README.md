# Orchestrator Helm Deployment

## Configuration

The following table list the configurable parameters of the orchestrator chart and their default values.

| Parameter        | Description     | Default   |
| ---              | ---             | ---       |
| `imagePullSecrets` | Reference to one or more secrets to be used when pulling images. | `[]` |
| `secrets.create` | Create orchestrator secrets. See charts/secrets subchart. | `false` |
| `secret.certs` | Secret name containing orchestrator certs. | `orc8r-secrets-certs` |
| `secret.configs` | Secret name containing orchestrator configs. | `orc8r-secrets-configs` |
| `secret.envdir` | Secret name containing orchestrator envdir. | `orc8r-secrets-envdir` |
| `proxy.podDisruptionBudget.enabled` | Enables creation of a PodDisruptionBudget for proxy. | `false` |
| `proxy.podDisruptionBudget.minAvailable` | Minimum number / percentage of pods that should remain scheduled. | `1` |
| `proxy.podDisruptionBudget.maxUnavailable` | Maximum number / percentage of pods that may be made unavailable. | `""` |
| `proxy.service.enabled` | Enables proxy service. | `true` |
| `proxy.service.lagacyEnabled` | Enables proxy legacy service. | `true` |
| `proxy.service.annotations` | Annotations to be added to the proxy service. | `{}` |
| `proxy.service.extraAnnotations.bootstrapLagacy` | Extra annotations to be added to the bootstrap-legacy proxy service. | `{}` |
| `proxy.service.extraAnnotations.clientcertLegacy` | Extra annotations to be added to the clientcert-legacy proxy service. | `{}` |
| `proxy.service.labels` | Proxy service labels. | `{}` |
| `proxy.service.type` | Proxy service type. | `ClusterIP` |
| `proxy.service.port.clientcert.port` | Proxy client certificate service external port. | `9443` |
| `proxy.service.port.clientcert.targetPort` | Proxy client certificate service internal port. | `9443` |
| `proxy.service.port.clientcert.nodePort` | Proxy client certificate service node port. | `nil` |
| `proxy.service.port.open.port` | Proxy open service external port. | `9444` |
| `proxy.service.port.open.targetPort` | Proxy open service internal port. | `9444` |
| `proxy.service.port.open.nodePort` | Proxy open service node port. | `nil` |
| `proxy.image.repository` | Repository for orchestrator proxy image. | `nil` |
| `proxy.image.tag` | Tag for orchestrator proxy image. | `latest` |
| `proxy.image.pullPolicy` | Pull policy for orchestrator proxy image. | `IfNotPresent` |
| `proxy.spec.hostname` | Magma controller domain name. | `""` |
| `proxy.replicas` | Number of instances to deploy for orchestrator proxy. | `1` |
| `proxy.resources` | Define resources requests and limits for Pods. | `{}` |
| `proxy.nodeSelector` | Define which Nodes the Pods are scheduled on. | `{}` |
| `proxy.tolerations` | If specified, the pod's tolerations. | `[]` |
| `proxy.affinity` | Assign the orchestrator proxy to run on specific nodes. | `{}` |
| `controller.podDisruptionBudget.enabled` | Enables creation of a PodDisruptionBudget for proxy. | `false` |
| `controller.podDisruptionBudget.minAvailable` | Minimum number / percentage of pods that should remain scheduled. | `1` |
| `controller.podDisruptionBudget.maxUnavailable` | Maximum number / percentage of pods that may be made unavailable. | `""` |
| `controller.service.annotations` | Annotations to be added to the controller service. | `{}` |
| `controller.service.labels` | Controller service labels. | `{}` |
| `controller.service.type` | Controller service type. | `ClusterIP` |
| `controller.service.port` | Controller web service external port. | `8080` |
| `controller.service.targetPort` | Controller web service internal port. | `8080` |
| `controller.service.portStart` | Controller service port range start. | `9079` |
| `controller.service.portEnd` | Controller service inclusive port range end. | `9108` |
| `controller.image.repository` | Repository for orchestrator controller image. | `nil` |
| `controller.image.tag` | Tag for orchestrator controller image. | `latest` |
| `controller.image.pullPolicy` | Pull policy for orchestrator controller image. | `IfNotPresent` |
| `controller.spec.database.driver` | orc8r database name. | `mysql/postgres` |
| `controller.spec.database.sql_dialect` | database dialect name. | `maria/psql` |
| `controller.spec.database.db` | orc8r database name. | `magma` |
| `controller.spec.database.host` | database host. | `postgresql` |
| `controller.spec.database.port` | database port. | `5432` |
| `controller.spec.database.user` | Database username. | `postgres` |
| `controller.spec.database.pass` | Database password. | `postgres` |
| `controller.replicas` | Number of instances to deploy for orchestrator controller. | `1` |
| `controller.resources` | Define resources requests and limits for Pods. | `{}` |
| `controller.nodeSelector` | Define which Nodes the Pods are scheduled on. | `{}` |
| `controller.tolerations` | If specified, the pod's tolerations. | `[]` |
| `controller.affinity` | Assign the orchestrator proxy to run on specific nodes. | `{}` |
| `nms.enabled` | If true, deploy the nms sub-chart | `true` |
| `nms.nginx.create` | Enable nms nginx service. | `true` |
| `nms.magmalte.create` | Enable nms magmalte app. | `true` |
| `nms.rbac` | Enable rbac for nginx and magmalte app. | `false` |
| `logging.enabled` | If true, deploy the logging sub-chart | `true` |
