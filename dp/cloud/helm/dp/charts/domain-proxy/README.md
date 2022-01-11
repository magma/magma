
Domain-proxy
===========

A Helm chart for magma orchestrator's domain-proxy module.

This Chart will deploy the following:

- 1 x Configuration controller
- 1 x protocol controller with 80/443 TCP port exposed with either ingress/contour ingress resource.
- 1 x radio controller with GRPC 50053/TCP port.
- All using Kubernetes Deployment
- 1 x job for Postgres database migration.

## Installation.

### From Source.

```console
$ git clone
$ cd dp/cloud/helm/dp
$ helm dep update domain-proxy
$ helm install --name myname --namespace mynamespace domain-proxy
```

### Certificates.

In order to work properly Domain proxy requires set of certificates. To enable chart to consume them place them inside.
`certificates` directory and setup proper certificate paths in your `values.yaml` file.

## Development.

### Local Deployment using Minikube.
If you're running locally in Minikube, see the `examples/minikube_values.yml` file.

To run local development environment:

```console
$ make
```

## Configuration

The following table lists the configurable parameters of the Domain-proxy chart and their default values.

| Parameter                | Description             | Default        |
| ------------------------ | ----------------------- | -------------- |
| `create` | Deploy Domain Proxy Chart. | `true` |
| `nameOverride` | Replaces the name of the chart in the `Chart.yaml` file. | `""` |
| `fullnameOverride` | Completely replaces the helm release generated name. | `""` |
| `configuration_controller.sasEndpointUrl` | Endpoint where sas request should be send. | `""` |
| `configuration_controller.requestProcessingInterval` | How often configuration controller will send requests to SAS. In seconds. | `"10"` |
| `configuration_controller.database` | Database configuration. | `{}` |
| `configuration_controller.nameOverride` | Replaces service part of the dp component deployment name. | `""` |
| `configuration_controller.fullnameOverride` | Completely replaces dp component deployment name. | `""` |
| `configuration_controller.enabled` | Enables deployment of the given service. | `true` |
| `configuration_controller.name` | Domain proxy component name. | `"configuration-controller"` |
| `configuration_controller.image.repository` | Docker image repository. | `"configuration-controller"` |
| `configuration_controller.image.pullPolicy` | Default the pull policy of all containers in that pod. | `"IfNotPresent"` |
| `configuration_controller.image.tag` | Overrides the image tag whose default is the chart appVersion. | `""` |
| `configuration_controller.replicaCount` | How many replicas of particular component should be created. | `1` |
| `configuration_controller.imagePullSecrets` | Name of the secret that contains container image registry keys | `[]` |
| `configuration_controller.serviceAccount.create` | Specifies whether a service account should be created | `false` |
| `configuration_controller.serviceAccount.annotations` | Annotations to add to the service account | `{}` |
| `configuration_controller.serviceAccount.name` | The name of the service account to use,If not set and create is true, a name is generated using the fullname template. | `""` |
| `configuration_controller.podAnnotations` | Additional pod annotations | `{}` |
| `configuration_controller.podSecurityContext` | Holds pod-level security attributes | `{}` |
| `configuration_controller.securityContext` | Holds security configuration that will be applied to a container. | `{}` |
| `configuration_controller.service.enable` | Whether to enable kubernetes service for dp component. | `true` |
| `configuration_controller.service.port` | Default port of enabled kubernetes service. | `8080` |
| `configuration_controller.tlsConfig` | tls configuration for communication with SAS. | `{}` |
| `configuration_controller.ingress.enabled` | Enable kubernetes ingress resource. | `false` |
| `configuration_controller.ingress.annotations` | Annotations to kubernetes ingress resource. | `{}` |
| `configuration_controller.ingress.hosts` |  | `[]` |
| `configuration_controller.ingress.tls` | Kubernetes secret name for tls termination on ingress kubernetes resource. | `[]` |
| `configuration_controller.resources` | Resource requests and limits of Pod. | `{}` |
| `configuration_controller.readinessProbe` | Readines probe definition. | `{}` |
| `configuration_controller.livenessProbe` | Livenes probe definition. | `{}` |
| `configuration_controller.autoscaling.enabled` | Enables horizontal pod autscaler kubernetes resource. | `false` |
| `configuration_controller.autoscaling.minReplicas` | Minimum number of dp component replicas. | `1` |
| `configuration_controller.autoscaling.maxReplicas` | Maximum number of dp component replicas. | `100` |
| `configuration_controller.autoscaling.targetCPUUtilizationPercentage` | Target CPU utilization threshold in perecents when new replica should be created | `80` |
| `configuration_controller.podDisruptionBudget.enabled` | Creates kubernetes podDisruptionBudget resource. | `false` |
| `configuration_controller.podDisruptionBudget.minAvailable` | Minimum available pods for dp component. | `1` |
| `configuration_controller.podDisruptionBudget.maxUnavailable` | Maximum unavailable pods for dp component. | `""` |
| `configuration_controller.nodeSelector` | Kubernetes node selection constraint. | `{}` |
| `configuration_controller.tolerations` | Allow the pods to schedule onto nodes with matching taints. | `[]` |
| `configuration_controller.affinity` | Constrain which nodes your pod is eligible to be scheduled on. | `{}` |
| `protocol_controller.nameOverride` | Replaces service part of the dp component deployment name. | `""` |
| `protocol_controller.fullnameOverride` | Completely replaces dp component deployment name. | `""` |
| `protocol_controller.enabled` | Enables deployment of the given dp component. | `true` |
| `protocol_controller.name` | Domain proxy component name. | `"protocol-controller"` |
| `protocol_controller.image.repository` | Docker image repository. | `"protocol-controller"` |
| `protocol_controller.image.tag` | Overrides the image tag whose default is the chart appVersion. | `""` |
| `protocol_controller.image.pullPolicy` | Default the pull policy of all containers in that pod. | `"IfNotPresent"` |
| `protocol_controller.replicaCount` | How many replicas of particular component should be created. | `1` |
| `protocol_controller.imagePullSecrets` | Name of the secret that contains container image registry keys. | `[]` |
| `protocol_controller.serviceAccount.create` | Specifies whether a service account should be created | `false` |
| `protocol_controller.serviceAccount.annotations` | Annotations to add to the service account | `{}` |
| `protocol_controller.serviceAccount.name` | The name of the service account to use,If not set and create is true, a name is generated using the fullname template. | `""` |
| `protocol_controller.podAnnotations` | Additional pod annotations. | `{}` |
| `protocol_controller.podSecurityContext` | Holds pod-level security attributes. | `{}` |
| `protocol_controller.securityContext` | Holds security configuration that will be applied to a container. | `{}` |
| `protocol_controller.service.enable` | Whether to enable kubernetes service for dp component. | `true` |
| `protocol_controller.service.port` | Default port of enabled kubernetes service. | `8080` |
| `protocol_controller.tlsConfig` |  | `{}` |
| `protocol_controller.apiPrefix` | Protocol controller URL API prefix. | `"/sas/v1"` |
| `protocol_controller.ingress.enabled` | Enable kubernetes ingress resource. | `false` |
| `protocol_controller.ingress.annotations` | Annotations to kubernetes ingress resource. | `{}` |
| `protocol_controller.ingress.hosts` |  | `[]` |
| `protocol_controller.ingress.tls` | Kubernetes secret name for tls termination on ingress kubernetes resource. | `[]` |
| `protocol_controller.httpproxy.enabled` | Enables contour httpproxy CRD. | `false` |
| `protocol_controller.httpproxy.annotations` | Aditional annotations. | `{}` |
| `protocol_controller.httpproxy.virtualhost` |  | `{}` |
| `protocol_controller.resources` | Resource requests and limits of Pod. | `{}` |
| `protocol_controller.readinessProbe` | Readines probe definition. | `{}` |
| `protocol_controller.livenessProbe` | Livenes probe definition. | `{}` |
| `protocol_controller.autoscaling.enabled` | Enables horizontal pod autscaler kubernetes resource. | `false` |
| `protocol_controller.autoscaling.minReplicas` | Minimum number of dp component replicas. | `1` |
| `protocol_controller.autoscaling.maxReplicas` | Maximum number of dp component replicas. | `100` |
| `protocol_controller.autoscaling.targetCPUUtilizationPercentage` | Target CPU utilization threshold in perecents when new replica should be created. | `80` |
| `protocol_controller.podDisruptionBudget.enabled` | Creates kubernetes podDisruptionBudget resource. | `false` |
| `protocol_controller.podDisruptionBudget.minAvailable` | Minimum available pods for dp component. | `1` |
| `protocol_controller.podDisruptionBudget.maxUnavailable` | Minimum available pods for dp component. | `""` |
| `protocol_controller.nodeSelector` | Kubernetes node selection constraint. | `{}` |
| `protocol_controller.tolerations` | Allow the pods to schedule onto nodes with matching taints. | `[]` |
| `protocol_controller.affinity` | Constrain which nodes your pod is eligible to be scheduled on. | `{}` |
| `radio_controller.database` |  | `{}` |
| `radio_controller.nameOverride` | Replaces service part of the dp component deployment name. | `""` |
| `radio_controller.fullnameOverride` | Completely replaces dp component deployment name. | `""` |
| `radio_controller.enabled` | Enables deployment of the given dp component. | `true` |
| `radio_controller.name` | Domain proxy component name. | `"radio-controller"` |
| `radio_controller.image.repository` | Docker image repository. | `"radio-controller"` |
| `radio_controller.image.tag` | Overrides the image tag whose default is the chart appVersion. | `""` |
| `radio_controller.image.pullPolicy` | Default the pull policy of all containers in that pod. | `"IfNotPresent"` |
| `radio_controller.replicaCount` | How many replicas of particular component should be created. | `1` |
| `radio_controller.imagePullSecrets` | Name of the secret that contains container image registry keys. | `[]` |
| `radio_controller.serviceAccount.create` | Specifies whether a service account should be created. | `false` |
| `radio_controller.serviceAccount.annotations` | Annotations to add to the service account. | `{}` |
| `radio_controller.serviceAccount.name` | The name of the service account to use,If not set and create is true, a name is generated using the fullname template. | `""` |
| `radio_controller.podAnnotations` | Additional pod annotations. | `{}` |
| `radio_controller.podSecurityContext` | Holds pod-level security attributes. | `{}` |
| `radio_controller.securityContext` | Holds security configuration that will be applied to a container. | `{}` |
| `radio_controller.resources` | Resource requests and limits of Pod. | `{}` |
| `radio_controller.readinessProbe` | Readines probe definition. | `{}` |
| `radio_controller.livenessProbe` | Livenes probe definition. | `{}` |
| `radio_controller.autoscaling.enabled` | Enables horizontal pod autscaler kubernetes resource. | `false` |
| `radio_controller.autoscaling.minReplicas` | Minimum number of dp component replicas. | `1` |
| `radio_controller.autoscaling.maxReplicas` | Maximum number of dp component replicas. | `100` |
| `radio_controller.autoscaling.targetCPUUtilizationPercentage` | Target CPU utilization threshold in perecents when new replica should be created | `80` |
| `radio_controller.podDisruptionBudget.enabled` | Creates kubernetes podDisruptionBudget resource. | `false` |
| `radio_controller.podDisruptionBudget.minAvailable` | Minimum available pods for dp component. | `1` |
| `radio_controller.podDisruptionBudget.maxUnavailable` | Maximum unavailable pods for dp component. | `""` |
| `radio_controller.nodeSelector` | Kubernetes node selection constraint. | `{}` |
| `radio_controller.tolerations` | Allow the pods to schedule onto nodes with matching taints. | `[]` |
| `radio_controller.affinity` | Constrain which nodes your pod is eligible to be scheduled on. | `{}` |
| `active_mode_controller.nameOverride` | Replaces service part of the dp component deployment name. | `""` |
| `active_mode_controller.fullnameOverride` | Completely replaces dp component deployment name. | `""` |
| `active_mode_controller.enabled` | Enables deployment of the given dp component. | `true` |
| `active_mode_controller.name` | Domain proxy component name. | `"active-mode-controller"` |
| `active_mode_controller.image.repository` | Docker image repository. | `"active-mode-controller"` |
| `active_mode_controller.image.tag` | Overrides the image tag whose default is the chart appVersion. | `""` |
| `active_mode_controller.image.pullPolicy` | Default the pull policy of all containers in that pod. | `"IfNotPresent"` |
| `active_mode_controller.replicaCount` | How many replicas of particular component should be created. | `1` |
| `active_mode_controller.imagePullSecrets` | Name of the secret that contains container image registry keys. | `[]` |
| `active_mode_controller.serviceAccount.create` | Specifies whether a service account should be created. | `false` |
| `active_mode_controller.serviceAccount.annotations` | Annotations to add to the service account. | `{}` |
| `active_mode_controller.serviceAccount.name` | The name of the service account to use,If not set and create is true, a name is generated using the fullname template. | `""` |
| `active_mode_controller.podAnnotations` | Additional pod annotations. | `{}` |
| `active_mode_controller.podSecurityContext` | Holds pod-level security attributes. | `{}` |
| `active_mode_controller.securityContext` | Holds security configuration that will be applied to a container. | `{}` |
| `active_mode_controller.resources` | Resource requests and limits of Pod. | `{}` |
| `active_mode_controller.readinessProbe` | Readines probe definition. | `{}` |
| `active_mode_controller.livenessProbe` | Livenes probe definition. | `{}` |
| `active_mode_controller.autoscaling.enabled` | Enables horizontal pod autscaler kubernetes resource. | `false` |
| `active_mode_controller.autoscaling.minReplicas` | Minimum number of dp component replicas. | `1` |
| `active_mode_controller.autoscaling.maxReplicas` | Maximum number of dp component replicas. | `100` |
| `active_mode_controller.autoscaling.targetCPUUtilizationPercentage` | Target CPU utilization threshold in perecents when new replica should be created | `80` |
| `active_mode_controller.podDisruptionBudget.enabled` | Creates kubernetes podDisruptionBudget resource. | `false` |
| `active_mode_controller.podDisruptionBudget.minAvailable` | Minimum available pods for dp component. | `1` |
| `active_mode_controller.podDisruptionBudget.maxUnavailable` | Maximum unavailable pods for dp component. | `""` |
| `active_mode_controller.nodeSelector` | Kubernetes node selection constraint. | `{}` |
| `active_mode_controller.tolerations` | Allow the pods to schedule onto nodes with matching taints. | `[]` |
| `active_mode_controller.affinity` | Constrain which nodes your pod is eligible to be scheduled on. | `{}` |
| `db_service.database` |  | `{}` |
| `db_service.enabled` | Enables deployment of the given service. | `true` |
| `db_service.nameOverride` | Replaces service part of the dp component deployment name. | `""` |
| `db_service.fullnameOverride` | Completely replaces dp component deployment name. | `""` |
| `db_service.name` | Domain proxy component name. | `"db-service"` |
| `db_service.image.repository` | Docker image repository. | `"db-service"` |
| `db_service.image.pullPolicy` | Default the pull policy of all containers in that pod. | `"IfNotPresent"` |
| `db_service.image.tag` | Overrides the image tag whose default is the chart appVersion. | `""` |
| `db_service.imagePullSecrets` | Name of the secret that contains container image registry keys | `[]` |
| `db_service.serviceAccount.create` | Specifies whether a service account should be created. | `false` |
| `db_service.serviceAccount.annotations` | Annotations to add to the service account. | `{}` |
| `db_service.serviceAccount.name` | The name of the service account to use,If not set and create is true, a name is generated using the fullname template. | `""` |



---
_Documentation generated by [Frigate](https://frigate.readthedocs.io)._

